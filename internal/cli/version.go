package cli

import (
	"fmt"
	"io"
	"runtime"

	"github.com/spf13/cobra"
)

// Version is set at link time via -ldflags "-X github.com/anubhavg-icpl/talon/internal/cli.Version=..."
// when building release artifacts. Dev builds fall back to "dev".
var Version = "dev"

// Commit is the short git SHA when set via ldflags; empty in plain local builds.
var Commit = ""

// BuildDate is the UTC build timestamp when set via ldflags.
var BuildDate = ""

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print CLI version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			info := map[string]string{
				"version":    Version,
				"commit":     Commit,
				"build_date": BuildDate,
				"go":         runtime.Version(),
				"os/arch":    runtime.GOOS + "/" + runtime.GOARCH,
			}
			// version should work without PersistentPreRun client setup —
			// but Printer may be nil if someone calls with weird parent wiring.
			format, _ := ParseOutputFormat(cmd.Flag("output").Value.String())
			p := NewPrinter(format)
			return p.PrintValue(info, func(w io.Writer) error {
				fmt.Fprintf(w, "talon %s", Version)
				if Commit != "" {
					fmt.Fprintf(w, " (%s)", Commit)
				}
				fmt.Fprintln(w)
				if BuildDate != "" {
					fmt.Fprintf(w, "built:  %s\n", BuildDate)
				}
				fmt.Fprintf(w, "go:     %s\n", runtime.Version())
				fmt.Fprintf(w, "os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
				return nil
			})
		},
	}
}

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate completion scripts for your shell.

  # bash
  source <(talon completion bash)

  # zsh
  talon completion zsh > "${fpath[1]}/_talon"

  # fish
  talon completion fish | source`,
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			default:
				return fmt.Errorf("unsupported shell %q", args[0])
			}
		},
	}
	// Avoid PersistentPreRun trying to dial core for completion generation.
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error { return nil }
	return cmd
}
