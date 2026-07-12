package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Known compose service names in docker-compose.yml.
var composeServices = map[string]string{
	"core":     "talon-core",
	"relay":    "talon-relay",
	"arsenal":  "arsenal-engine",
	"msf":      "metasploit",
	"rabbitmq": "rabbitmq",
	"vuln":     "vuln-target",
	// aliases
	"talon-core":     "talon-core",
	"talon-relay":    "talon-relay",
	"arsenal-engine": "arsenal-engine",
	"metasploit":     "metasploit",
	"vuln-target":    "vuln-target",
}

func newLogsCmd(opts *RootOptions) *cobra.Command {
	var (
		follow bool
		tail   string
		since  string
	)

	cmd := &cobra.Command{
		Use:   "logs [service]",
		Short: "Show logs for a stack service (via docker compose)",
		Long: `Fetch logs for a Talon stack service using docker compose.

Services: core, relay, arsenal, msf, rabbitmq, vuln
(also accepts compose names: talon-core, arsenal-engine, …)

Examples:
  talon logs core
  talon logs arsenal --follow
  talon logs msf --tail 100
  talon logs rabbitmq --since 10m

Requires docker compose and a project dir (config project_dir, env
TALON_PROJECT_DIR, or the current working directory containing
docker-compose.yml).`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc := "talon-core"
			if len(args) == 1 {
				key := strings.ToLower(args[0])
				mapped, ok := composeServices[key]
				if !ok {
					return withExitCode(ExitUsage, "unknown service %q (want: core|relay|arsenal|msf|rabbitmq|vuln)", args[0])
				}
				svc = mapped
			}

			dir := opts.Resolved.ProjectDir
			if dir == "" {
				dir = findComposeDir(opts.Resolved.ComposeFile)
			}
			if dir == "" {
				return withExitCode(ExitError, "could not find compose project (set project_dir in config or TALON_PROJECT_DIR)")
			}

			composeFile := opts.Resolved.ComposeFile
			if composeFile == "" {
				composeFile = "docker-compose.yml"
			}
			if !filepath.IsAbs(composeFile) {
				composeFile = filepath.Join(dir, composeFile)
			}

			dargs := []string{"compose", "-f", composeFile, "logs", "--no-color"}
			if follow {
				dargs = append(dargs, "--follow")
			}
			if tail != "" {
				dargs = append(dargs, "--tail", tail)
			} else {
				dargs = append(dargs, "--tail", "200")
			}
			if since != "" {
				dargs = append(dargs, "--since", since)
			}
			dargs = append(dargs, svc)

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			// Follow has no timeout; one-shot gets a bound.
			if !follow {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, 2*time.Minute)
				defer cancel()
			}

			c := exec.CommandContext(ctx, "docker", dargs...)
			c.Dir = dir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				// Fallback: docker logs by container name conventions
				if fallbackErr := dockerLogsByName(ctx, svc, follow, tail, since, os.Stdout, os.Stderr); fallbackErr == nil {
					return nil
				}
				return withExitCode(ExitError, "docker compose logs failed: %v (cwd=%s file=%s)", err, dir, composeFile)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "stream logs")
	cmd.Flags().StringVar(&tail, "tail", "200", "number of lines to show from the end")
	cmd.Flags().StringVar(&since, "since", "", "show logs since timestamp (e.g. 10m, 2024-01-01T00:00:00)")
	// logs talks to docker, not core — skip core client requirement via name check in root
	return cmd
}

func findComposeDir(composeFile string) string {
	name := composeFile
	if name == "" {
		name = "docker-compose.yml"
	}
	// walk up from cwd
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	dir := cwd
	for i := 0; i < 8; i++ {
		candidate := filepath.Join(dir, filepath.Base(name))
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			return dir
		}
		// also try compose.yaml
		if st, err := os.Stat(filepath.Join(dir, "compose.yaml")); err == nil && !st.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// dockerLogsByName falls back to container_name mappings used in compose.
func dockerLogsByName(ctx context.Context, service string, follow bool, tail, since string, stdout, stderr io.Writer) error {
	container := map[string]string{
		"talon-core":     "talon_core",
		"talon-relay":    "talon_relay",
		"arsenal-engine": "arsenal_engine",
		"metasploit":     "msf_rpc",
		"rabbitmq":       "talon_rabbitmq",
		"vuln-target":    "talon_vuln_target",
	}[service]
	if container == "" {
		container = service
	}
	args := []string{"logs", "--timestamps"}
	if follow {
		args = append(args, "--follow")
	}
	if tail != "" {
		args = append(args, "--tail", tail)
	}
	if since != "" {
		args = append(args, "--since", since)
	}
	args = append(args, container)
	c := exec.CommandContext(ctx, "docker", args...)
	c.Stdout = stdout
	c.Stderr = stderr
	return c.Run()
}

// streamReader is unused helper kept for possible future local log files.
func streamReader(r io.Reader, w io.Writer) error {
	sc := bufio.NewScanner(r)
	// raise limit for long docker lines
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)
	for sc.Scan() {
		if _, err := fmt.Fprintln(w, sc.Text()); err != nil {
			return err
		}
	}
	return sc.Err()
}
