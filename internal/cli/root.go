// Package cli implements the production operator CLI for Talon
// (cmd/talon). It is a thin client over talon-core's HTTP control plane.
package cli

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// GlobalFlags are shared across all commands.
type GlobalFlags struct {
	CoreURL    string
	Output     string
	Timeout    time.Duration
	Quiet      bool
	ConfigPath string
}

// RootOptions holds runtime state for command handlers.
type RootOptions struct {
	Flags    GlobalFlags
	Client   *Client
	Printer  *Printer
	Resolved ResolvedConfig
	File     FileConfig
}

// NewRootCommand builds the root `talon` command and all subcommands.
func NewRootCommand() *cobra.Command {
	opts := &RootOptions{}

	root := &cobra.Command{
		Use:   "talon",
		Short: "Talon operator CLI — manage runs, HITL gates, and stack health",
		Long: `Talon is the operator surface for the Talon penetration-testing platform.

Talks to talon-core over HTTP (default http://localhost:8000).

Configuration precedence (highest wins):
  1. CLI flags (--core-url, --output, --timeout)
  2. Environment (TALON_CORE_URL, TALON_OUTPUT, TALON_PROJECT_DIR, …)
  3. Config file (~/.config/talon/config.yaml or $TALON_CONFIG)
  4. Built-in defaults

Examples:
  talon status
  talon run start --ip 127.0.0.1 --cve CVE-2011-2523 --lhost 192.168.0.176
  talon run status <run_id>
  talon run watch <run_id>
  talon run approve <run_id>
  talon run tools <run_id> --output json
  talon logs core --tail 100
  talon logs arsenal -f`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			name := cmd.Name()
			// Pure-local commands skip remote client setup.
			localOnly := name == "version" || name == "help" || name == "completion" || name == "logs" || name == "config"
			if cmd.Parent() != nil && cmd.Parent().Name() == "completion" {
				localOnly = true
			}

			cfgPath := opts.Flags.ConfigPath
			if cfgPath == "" {
				cfgPath = DefaultConfigPath()
			}
			file, err := LoadFileConfig(cfgPath)
			if err != nil {
				return err
			}
			opts.File = file

			// Detect whether output flag was explicitly changed.
			outFlag := cmd.Flags().Lookup("output")
			if outFlag == nil {
				outFlag = cmd.InheritedFlags().Lookup("output")
			}
			flagOutput := opts.Flags.Output
			if outFlag != nil && !outFlag.Changed {
				// Prefer file/env default over cobra default "table".
				flagOutput = ""
			}
			flagCore := opts.Flags.CoreURL // empty unless user set --core-url

			timeoutFlag := cmd.Flags().Lookup("timeout")
			if timeoutFlag == nil {
				timeoutFlag = cmd.InheritedFlags().Lookup("timeout")
			}
			flagTimeout := opts.Flags.Timeout
			if timeoutFlag != nil && !timeoutFlag.Changed {
				flagTimeout = 0 // let ResolveConfig use file/env/default
			}

			opts.Resolved = ResolveConfig(file, flagCore, flagOutput, flagTimeout, cfgPath)
			// If output still empty, default table
			if opts.Resolved.Output == "" {
				opts.Resolved.Output = "table"
			}
			// Keep Flags in sync for commands that read them.
			opts.Flags.CoreURL = opts.Resolved.CoreURL
			opts.Flags.Output = opts.Resolved.Output
			opts.Flags.Timeout = opts.Resolved.Timeout

			format, err := ParseOutputFormat(opts.Resolved.Output)
			if err != nil {
				return err
			}
			opts.Printer = NewPrinter(format)

			if localOnly {
				return nil
			}

			client, err := NewClient(opts.Resolved.CoreURL, opts.Resolved.Timeout)
			if err != nil {
				return err
			}
			opts.Client = client
			return nil
		},
	}

	pf := root.PersistentFlags()
	pf.StringVar(&opts.Flags.CoreURL, "core-url", "", "talon-core base URL (env TALON_CORE_URL)")
	pf.StringVarP(&opts.Flags.Output, "output", "o", "table", "output format: table|json|raw")
	pf.DurationVar(&opts.Flags.Timeout, "timeout", 30*time.Second, "HTTP client timeout for control-plane requests")
	pf.BoolVarP(&opts.Flags.Quiet, "quiet", "q", false, "minimal output (ids and statuses only)")
	pf.StringVar(&opts.Flags.ConfigPath, "config", "", "config file path (env TALON_CONFIG, default ~/.config/talon/config.yaml)")

	root.AddCommand(newVersionCmd())
	root.AddCommand(newStatusCmd(opts))
	root.AddCommand(newRunCmd(opts))
	root.AddCommand(newLogsCmd(opts))
	root.AddCommand(newCompletionCmd())
	root.AddCommand(newConfigCmd(opts))

	return root
}

// Execute runs the root command and maps errors to exit codes.
func Execute() int {
	root := NewRootCommand()
	if err := root.Execute(); err != nil {
		if code, ok := err.(exitError); ok {
			if code.Msg != "" {
				fmt.Fprintln(os.Stderr, "error:", code.Msg)
			}
			return code.Code
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		return ExitError
	}
	return ExitOK
}

// exitError carries a process exit code from a command handler.
type exitError struct {
	Code int
	Msg  string
}

func (e exitError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return fmt.Sprintf("exit %d", e.Code)
}

func withExitCode(code int, format string, args ...any) error {
	return exitError{Code: code, Msg: fmt.Sprintf(format, args...)}
}

func newConfigCmd(opts *RootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Show resolved CLI configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			payload := map[string]any{
				"config_path":  opts.Resolved.ConfigPath,
				"core_url":     opts.Resolved.CoreURL,
				"arsenal_url":  opts.Resolved.ArsenalURL,
				"msf":          opts.Resolved.MSF,
				"amqp":         opts.Resolved.AMQP,
				"timeout":      opts.Resolved.Timeout.String(),
				"output":       opts.Resolved.Output,
				"compose_file": opts.Resolved.ComposeFile,
				"project_dir":  opts.Resolved.ProjectDir,
			}
			return opts.Printer.PrintValue(payload, func(w io.Writer) error {
				return KeyValueTable(w, [][2]string{
					{"config_path", opts.Resolved.ConfigPath},
					{"core_url", opts.Resolved.CoreURL},
					{"arsenal_url", opts.Resolved.ArsenalURL},
					{"msf", opts.Resolved.MSF},
					{"amqp", opts.Resolved.AMQP},
					{"timeout", opts.Resolved.Timeout.String()},
					{"output", opts.Resolved.Output},
					{"compose_file", opts.Resolved.ComposeFile},
					{"project_dir", opts.Resolved.ProjectDir},
				})
			})
		},
	}
}
