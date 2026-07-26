package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// loginResponse is POST /auth/login's body.
type loginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

// meResponse is GET /auth/me's body.
type meResponse struct {
	Username string `json:"username"`
	Auth     string `json:"auth"`
}

func newAuthCmd(opts *RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with talon-core (login, logout, status)",
	}
	cmd.AddCommand(newAuthLoginCmd(opts))
	cmd.AddCommand(newAuthLogoutCmd(opts))
	cmd.AddCommand(newAuthStatusCmd(opts))
	return cmd
}

func newAuthLoginCmd(opts *RootOptions) *cobra.Command {
	var username, password string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in and save a session token to the config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := bufio.NewReader(os.Stdin)
			if username == "" {
				fmt.Fprint(opts.Printer.Out, "username: ")
				line, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				username = strings.TrimSpace(line)
			}
			if password == "" {
				fmt.Fprint(opts.Printer.Out, "password: ")
				if term.IsTerminal(int(os.Stdin.Fd())) {
					raw, err := term.ReadPassword(int(os.Stdin.Fd()))
					fmt.Fprintln(opts.Printer.Out)
					if err != nil {
						return err
					}
					password = string(raw)
				} else {
					line, err := reader.ReadString('\n')
					if err != nil {
						return err
					}
					password = strings.TrimSpace(line)
				}
			}
			if username == "" || password == "" {
				return fmt.Errorf("username and password are required")
			}

			var out loginResponse
			err := opts.Client.doJSON(context.Background(), "POST", "/auth/login",
				map[string]string{"username": username, "password": password}, &out)
			if err != nil {
				return err
			}
			if out.Token == "" {
				return fmt.Errorf("login succeeded but no token returned")
			}
			if err := SaveConfigToken(opts.Resolved.ConfigPath, out.Token); err != nil {
				return fmt.Errorf("save token: %w", err)
			}
			fmt.Fprintf(opts.Printer.Out, "logged in as %s — token saved to %s\n", out.Username, opts.Resolved.ConfigPath)
			return nil
		},
	}
	cmd.Flags().StringVarP(&username, "username", "u", "", "username (prompts when omitted)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "password (prompts when omitted; avoid in scripts, prefer TALON_TOKEN)")
	return cmd
}

func newAuthLogoutCmd(opts *RootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Invalidate the session and remove the saved token",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Best-effort server-side invalidation; local token always removed.
			_ = opts.Client.doJSON(context.Background(), "POST", "/auth/logout", nil, nil)
			if err := SaveConfigToken(opts.Resolved.ConfigPath, ""); err != nil {
				return fmt.Errorf("clear token: %w", err)
			}
			fmt.Fprintln(opts.Printer.Out, "logged out")
			return nil
		},
	}
}

func newAuthStatusCmd(opts *RootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the current session user",
		RunE: func(cmd *cobra.Command, args []string) error {
			var out meResponse
			if err := opts.Client.doJSON(context.Background(), "GET", "/auth/me", nil, &out); err != nil {
				return err
			}
			if out.Auth == "disabled" {
				fmt.Fprintln(opts.Printer.Out, "auth disabled on this core — no login required")
				return nil
			}
			fmt.Fprintf(opts.Printer.Out, "authenticated as %s\n", out.Username)
			return nil
		},
	}
}
