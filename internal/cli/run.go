package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func newRunCmd(opts *RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start, inspect, and control validation runs",
		Long:  "Manage agent validation runs against talon-core (start, status, watch, HITL decisions, tools, traces).",
	}
	cmd.AddCommand(newRunStartCmd(opts))
	cmd.AddCommand(newRunStatusCmd(opts))
	cmd.AddCommand(newRunWatchCmd(opts))
	cmd.AddCommand(newRunApproveCmd(opts))
	cmd.AddCommand(newRunRejectCmd(opts))
	cmd.AddCommand(newRunEditCmd(opts))
	cmd.AddCommand(newRunToolsCmd(opts))
	cmd.AddCommand(newRunTracesCmd(opts))
	return cmd
}

func newRunStartCmd(opts *RootOptions) *cobra.Command {
	var (
		ip, cve, service, description, lhost string
		lport                                int
		watch                                bool
		interval                             time.Duration
		autoApprove                          bool
	)

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a new validation run",
		Example: `  talon run start --ip 127.0.0.1 --cve CVE-2011-2523 --lhost 192.168.0.176
  talon run start --ip 10.0.0.5 --service "vsftpd 2.3.4" --watch --auto-approve`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(ip) == "" {
				return withExitCode(ExitUsage, "--ip is required")
			}
			ctx := cmd.Context()
			req := StartRequest{
				IP:          ip,
				CVEID:       cve,
				ServiceName: service,
				Description: description,
				LHOST:       lhost,
				LPORT:       lport,
			}
			resp, err := opts.Client.Start(ctx, req)
			if err != nil {
				return err
			}

			if err := opts.Printer.PrintValue(resp, func(w io.Writer) error {
				if opts.Flags.Quiet {
					fmt.Fprintln(w, resp.RunID)
					return nil
				}
				return KeyValueTable(w, [][2]string{
					{"run_id", resp.RunID},
					{"message", resp.Message},
				})
			}); err != nil {
				return err
			}

			if watch {
				return watchRun(ctx, opts, resp.RunID, interval, autoApprove)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&ip, "ip", "", "target IP or hostname (required)")
	cmd.Flags().StringVar(&cve, "cve", "", "CVE identifier (e.g. CVE-2011-2523)")
	cmd.Flags().StringVar(&service, "service", "", "service fingerprint name")
	cmd.Flags().StringVar(&description, "description", "", "free-text target description")
	cmd.Flags().StringVar(&lhost, "lhost", "", "attacker LHOST for reverse payloads")
	cmd.Flags().IntVar(&lport, "lport", 0, "attacker LPORT (default 4444 on server)")
	cmd.Flags().BoolVar(&watch, "watch", false, "poll status until the run finishes")
	cmd.Flags().DurationVar(&interval, "interval", 3*time.Second, "poll interval when --watch is set")
	cmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "auto-approve HITL interrupts while watching (authorized lab use only)")
	_ = cmd.MarkFlagRequired("ip")
	return cmd
}

func newRunStatusCmd(opts *RootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "status <run_id>",
		Short: "Show status of a run",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := opts.Client.Status(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if err := printStatus(opts, args[0], st); err != nil {
				return err
			}
			code := ExitCodeForStatus(st.Status)
			if code == ExitOK {
				return nil
			}
			return withExitCode(code, "run %s status=%s", args[0], st.Status)
		},
	}
}

func newRunWatchCmd(opts *RootOptions) *cobra.Command {
	var (
		interval    time.Duration
		autoApprove bool
		maxWait     time.Duration
	)
	cmd := &cobra.Command{
		Use:   "watch <run_id>",
		Short: "Poll a run until it completes or errors",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if maxWait > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, maxWait)
				defer cancel()
			}
			return watchRun(ctx, opts, args[0], interval, autoApprove)
		},
	}
	cmd.Flags().DurationVar(&interval, "interval", 3*time.Second, "poll interval")
	cmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "auto-approve HITL interrupts (authorized lab use only)")
	cmd.Flags().DurationVar(&maxWait, "max-wait", 0, "give up after this duration (0 = no limit)")
	return cmd
}

func newRunApproveCmd(opts *RootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "approve <run_id>",
		Short: "Approve a pending HITL interrupt (e.g. nmap_scan)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resumeDecision(cmd.Context(), opts, args[0], ResumeRequest{Decision: "approve"})
		},
	}
}

func newRunRejectCmd(opts *RootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "reject <run_id>",
		Short: "Reject a pending HITL interrupt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resumeDecision(cmd.Context(), opts, args[0], ResumeRequest{Decision: "reject"})
		},
	}
}

func newRunEditCmd(opts *RootOptions) *cobra.Command {
	var argsJSON string
	cmd := &cobra.Command{
		Use:   "edit <run_id>",
		Short: "Approve a HITL interrupt with edited tool arguments",
		Long: `Resume a paused run with decision=edit and a JSON object of tool args.

Example:
  talon run edit <run_id> --args '{"target":"127.0.0.1","ports":"21,6200","scan_type":"-sT -Pn"}'`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(argsJSON) == "" {
				return withExitCode(ExitUsage, "--args is required (JSON object)")
			}
			var edited map[string]any
			if err := json.Unmarshal([]byte(argsJSON), &edited); err != nil {
				return withExitCode(ExitUsage, "invalid --args JSON: %v", err)
			}
			return resumeDecision(cmd.Context(), opts, args[0], ResumeRequest{
				Decision:   "edit",
				EditedArgs: edited,
			})
		},
	}
	cmd.Flags().StringVar(&argsJSON, "args", "", "JSON object of edited tool arguments (required)")
	_ = cmd.MarkFlagRequired("args")
	return cmd
}

func newRunToolsCmd(opts *RootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "tools <run_id>",
		Short: "Show the tool-call log for a run",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := opts.Client.Tools(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			log := resp.ToolLog
			if log == nil {
				log = []ToolCallRecord{}
			}
			return opts.Printer.PrintValue(map[string]any{
				"run_id":   args[0],
				"count":    len(log),
				"tool_log": log,
			}, func(w io.Writer) error {
				if !opts.Flags.Quiet {
					fmt.Fprintf(w, "run_id=%s  tools=%d\n\n", args[0], len(log))
				}
				if len(log) == 0 {
					fmt.Fprintln(w, "(no tool calls recorded yet)")
					return nil
				}
				return ToolsTable(w, log)
			})
		},
	}
}

func newRunTracesCmd(opts *RootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "traces <run_id>",
		Short: "Show stored history/traces for a run",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := opts.Client.Traces(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			hist := resp.History
			if hist == nil {
				hist = []string{}
			}
			return opts.Printer.PrintValue(map[string]any{
				"run_id":  args[0],
				"history": hist,
			}, func(w io.Writer) error {
				if len(hist) == 0 {
					fmt.Fprintln(w, "(no history recorded)")
					return nil
				}
				for i, line := range hist {
					fmt.Fprintf(w, "[%d] %s\n", i, line)
				}
				return nil
			})
		},
	}
}

func resumeDecision(ctx context.Context, opts *RootOptions, runID string, body ResumeRequest) error {
	resp, err := opts.Client.Resume(ctx, runID, body)
	if err != nil {
		return err
	}
	return opts.Printer.PrintValue(map[string]any{
		"run_id":   runID,
		"decision": body.Decision,
		"message":  resp.Message,
	}, func(w io.Writer) error {
		if opts.Flags.Quiet {
			fmt.Fprintln(w, "ok")
			return nil
		}
		return KeyValueTable(w, [][2]string{
			{"run_id", runID},
			{"decision", body.Decision},
			{"message", resp.Message},
		})
	})
}

func printStatus(opts *RootOptions, runID string, st *StatusResponse) error {
	payload := map[string]any{
		"run_id":    runID,
		"status":    st.Status,
		"output":    st.Output,
		"interrupt": st.Interrupt,
	}
	return opts.Printer.PrintValue(payload, func(w io.Writer) error {
		rows := [][2]string{
			{"run_id", runID},
			{"status", st.Status},
		}
		if st.Interrupt != nil {
			rows = append(rows, [2]string{"interrupt.tool", st.Interrupt.ToolName})
			if b, err := json.Marshal(st.Interrupt.Args); err == nil {
				rows = append(rows, [2]string{"interrupt.args", string(b)})
			}
		}
		if st.Output != "" && !opts.Flags.Quiet {
			out := st.Output
			if len(out) > 400 {
				out = out[:397] + "..."
			}
			rows = append(rows, [2]string{"output", strings.ReplaceAll(out, "\n", " ")})
		}
		return KeyValueTable(w, rows)
	})
}

func watchRun(ctx context.Context, opts *RootOptions, runID string, interval time.Duration, autoApprove bool) error {
	if interval <= 0 {
		interval = 3 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	handle := func() (done bool, err error) {
		done, err = pollOnce(ctx, opts, runID, autoApprove)
		if err != nil {
			if _, ok := err.(exitError); ok {
				return true, err
			}
			if !opts.Flags.Quiet {
				fmt.Fprintln(opts.Printer.Err, "warn:", err)
			}
			return false, nil // transient
		}
		return done, nil
	}

	if done, err := handle(); done {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return withExitCode(ExitError, "watch cancelled: %v", ctx.Err())
		case <-ticker.C:
			if done, err := handle(); done {
				return err
			}
		}
	}
}

// pollOnce returns done=true when the run has reached a terminal state or
// is blocked on HITL without auto-approve. Transient API errors return
// done=false with a non-exit error so the watch loop can retry.
func pollOnce(ctx context.Context, opts *RootOptions, runID string, autoApprove bool) (done bool, err error) {
	st, err := opts.Client.Status(ctx, runID)
	if err != nil {
		return false, err
	}

	// Human-readable progress on stderr so stdout stays clean for JSON pipes.
	if opts.Printer.Format == OutputTable && !opts.Flags.Quiet {
		line := fmt.Sprintf("%s  status=%s", time.Now().Format("15:04:05"), st.Status)
		if st.Interrupt != nil {
			line += "  interrupt=" + st.Interrupt.ToolName
		}
		fmt.Fprintln(opts.Printer.Err, line)
	}

	if st.Status == "awaiting_approval" && st.Interrupt != nil {
		if autoApprove {
			if !opts.Flags.Quiet {
				fmt.Fprintf(opts.Printer.Err, "auto-approving %s\n", st.Interrupt.ToolName)
			}
			if _, err := opts.Client.Resume(ctx, runID, ResumeRequest{Decision: "approve"}); err != nil {
				return false, err
			}
			return false, nil
		}
		_ = printStatus(opts, runID, st)
		return true, withExitCode(ExitAwaitingApproval,
			"run %s awaiting approval for tool %s (use: talon run approve %s)",
			runID, st.Interrupt.ToolName, runID)
	}

	if IsTerminalStatus(st.Status) {
		if err := printStatus(opts, runID, st); err != nil {
			return true, err
		}
		code := ExitCodeForStatus(st.Status)
		if code != ExitOK {
			return true, withExitCode(code, "run %s finished with status=%s", runID, st.Status)
		}
		return true, nil
	}
	return false, nil
}
