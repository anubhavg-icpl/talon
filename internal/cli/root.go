// Package cli implements the production operator CLI for Talon
// (cmd/talon). It is a thin client over talon-core's HTTP control plane.
package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// GlobalFlags are shared across all commands.
type GlobalFlags struct {
	CoreURL string
	Output  string
	Timeout time.Duration
	Quiet   bool
}

// RootOptions holds runtime state for command handlers.
type RootOptions struct {
	Flags   GlobalFlags
	Client  *Client
	Printer *Printer
}

// NewRootCommand builds the root `talon` command and all subcommands.
func NewRootCommand() *cobra.Command {
	opts := &RootOptions{}

	root := &cobra.Command{
		Use:   "talon",
		Short: "Talon operator CLI — manage runs, HITL gates, and stack health",
		Long: `Talon is the operator surface for the Talon penetration-testing platform.

Talks to talon-core over HTTP (default http://localhost:8000). Configure
with --core-url or the TALON_CORE_URL environment variable.

Examples:
  talon status
  talon run start --ip 127.0.0.1 --cve CVE-2011-2523 --lhost 192.168.0.176
  talon run status <run_id>
  talon run watch <run_id>
  talon run approve <run_id>
  talon run tools <run_id> --output json`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip client wiring for pure-local commands.
			if cmd.Name() == "version" || cmd.Name() == "help" || cmd.Name() == "completion" {
				return nil
			}
			// completion is nested under root; also skip when generating docs.
			if cmd.Parent() != nil && cmd.Parent().Name() == "completion" {
				return nil
			}

			format, err := ParseOutputFormat(opts.Flags.Output)
			if err != nil {
				return err
			}
			opts.Printer = NewPrinter(format)

			coreURL := opts.Flags.CoreURL
			if coreURL == "" {
				coreURL = os.Getenv("TALON_CORE_URL")
			}
			if coreURL == "" {
				coreURL = "http://localhost:8000"
			}
			opts.Flags.CoreURL = coreURL

			client, err := NewClient(coreURL, opts.Flags.Timeout)
			if err != nil {
				return err
			}
			opts.Client = client
			return nil
		},
	}

	pf := root.PersistentFlags()
	pf.StringVar(&opts.Flags.CoreURL, "core-url", "", "talon-core base URL (env TALON_CORE_URL, default http://localhost:8000)")
	pf.StringVarP(&opts.Flags.Output, "output", "o", "table", "output format: table|json|raw")
	pf.DurationVar(&opts.Flags.Timeout, "timeout", 30*time.Second, "HTTP client timeout for control-plane requests")
	pf.BoolVarP(&opts.Flags.Quiet, "quiet", "q", false, "minimal output (ids and statuses only)")

	root.AddCommand(newVersionCmd())
	root.AddCommand(newStatusCmd(opts))
	root.AddCommand(newRunCmd(opts))
	root.AddCommand(newCompletionCmd())

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
