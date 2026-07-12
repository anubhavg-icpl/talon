package cli

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type serviceProbe struct {
	Name    string `json:"name"`
	Target  string `json:"target"`
	Status  string `json:"status"` // ok | fail | skip
	Detail  string `json:"detail,omitempty"`
	Latency string `json:"latency,omitempty"`
}

type statusReport struct {
	CoreURL  string         `json:"core_url"`
	Overall  string         `json:"overall"` // healthy | degraded | down
	Services []serviceProbe `json:"services"`
	Checked  string         `json:"checked_at"`
}

func newStatusCmd(opts *RootOptions) *cobra.Command {
	var (
		arsenalURL string
		msfAddr    string
		amqpAddr   string
		skipExtra  bool
	)

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Probe talon-core and optional stack endpoints",
		Long: `Check reachability of the Talon control plane and common stack ports.

Always probes talon-core. Optionally probes:
  --arsenal-url   Arsenal engine health (default http://localhost:8888/health)
  --msf           Metasploit RPC TCP (default localhost:5554)
  --amqp          RabbitMQ AMQP TCP (default localhost:5672)

Overall is healthy when core is up; degraded if core is up but a side
service fails; down if core is unreachable.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			report := statusReport{
				CoreURL: opts.Client.BaseURL(),
				Checked: time.Now().UTC().Format(time.RFC3339),
			}

			// --- core ---
			core := serviceProbe{Name: "talon-core", Target: opts.Client.BaseURL()}
			start := time.Now()
			if err := opts.Client.ProbeCore(ctx); err != nil {
				core.Status = "fail"
				core.Detail = err.Error()
			} else {
				core.Status = "ok"
				core.Detail = "listening"
			}
			core.Latency = time.Since(start).Round(time.Millisecond).String()
			report.Services = append(report.Services, core)

			if !skipExtra {
				// --- arsenal ---
				if arsenalURL == "" {
					arsenalURL = envOr("TALON_ARSENAL_URL", "http://localhost:8888/health")
				}
				report.Services = append(report.Services, probeHTTP(ctx, "arsenal-engine", arsenalURL, opts.Flags.Timeout))

				// --- msf ---
				if msfAddr == "" {
					host := envOr("MSF_SERVER", "localhost")
					port := envOr("MSF_PORT", "5554")
					msfAddr = net.JoinHostPort(host, port)
				}
				report.Services = append(report.Services, probeTCP(ctx, "metasploit-rpc", msfAddr, 3*time.Second))

				// --- amqp ---
				if amqpAddr == "" {
					if u := os.Getenv("AMQP_URL"); u != "" {
						if host, port, ok := amqpHostPort(u); ok {
							amqpAddr = net.JoinHostPort(host, port)
						}
					}
					if amqpAddr == "" {
						amqpAddr = "localhost:5672"
					}
				}
				report.Services = append(report.Services, probeTCP(ctx, "rabbitmq", amqpAddr, 3*time.Second))
			}

			report.Overall = overallFrom(report.Services)

			err := opts.Printer.PrintValue(report, func(w io.Writer) error {
				if !opts.Flags.Quiet {
					fmt.Fprintf(w, "Talon stack status  overall=%s  core=%s\n\n", report.Overall, report.CoreURL)
				}
				rows := make([][2]string, 0, len(report.Services)+1)
				for _, s := range report.Services {
					val := s.Status
					if s.Latency != "" {
						val += " (" + s.Latency + ")"
					}
					if s.Detail != "" {
						val += " — " + s.Detail
					}
					rows = append(rows, [2]string{s.Name, val})
				}
				return KeyValueTable(w, rows)
			})
			if err != nil {
				return err
			}
			if report.Overall == "down" {
				return withExitCode(ExitError, "talon-core is unreachable at %s", report.CoreURL)
			}
			if report.Overall == "degraded" {
				// Degraded is informational: exit 0 so scripts can still proceed,
				// but print to stderr when not quiet.
				if !opts.Flags.Quiet {
					fmt.Fprintln(opts.Printer.Err, "warning: stack is degraded (core is up; a side service failed)")
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&arsenalURL, "arsenal-url", "", "arsenal health URL (env TALON_ARSENAL_URL)")
	cmd.Flags().StringVar(&msfAddr, "msf", "", "metasploit RPC host:port (default localhost:5554)")
	cmd.Flags().StringVar(&amqpAddr, "amqp", "", "rabbitmq host:port (default localhost:5672)")
	cmd.Flags().BoolVar(&skipExtra, "core-only", false, "only probe talon-core")
	return cmd
}

func overallFrom(services []serviceProbe) string {
	coreOK := false
	sideFail := false
	for _, s := range services {
		if s.Name == "talon-core" {
			coreOK = s.Status == "ok"
			continue
		}
		if s.Status == "fail" {
			sideFail = true
		}
	}
	if !coreOK {
		return "down"
	}
	if sideFail {
		return "degraded"
	}
	return "healthy"
}

func probeHTTP(ctx context.Context, name, rawURL string, timeout time.Duration) serviceProbe {
	p := serviceProbe{Name: name, Target: rawURL}
	start := time.Now()
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, rawURL, nil)
	if err != nil {
		p.Status = "fail"
		p.Detail = err.Error()
		return p
	}
	resp, err := http.DefaultClient.Do(req)
	p.Latency = time.Since(start).Round(time.Millisecond).String()
	if err != nil {
		p.Status = "fail"
		p.Detail = err.Error()
		return p
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		p.Status = "ok"
		p.Detail = "HTTP " + strconv.Itoa(resp.StatusCode)
	} else {
		p.Status = "fail"
		p.Detail = "HTTP " + strconv.Itoa(resp.StatusCode)
	}
	return p
}

func probeTCP(ctx context.Context, name, addr string, timeout time.Duration) serviceProbe {
	p := serviceProbe{Name: name, Target: addr}
	start := time.Now()
	d := net.Dialer{Timeout: timeout}
	conn, err := d.DialContext(ctx, "tcp", addr)
	p.Latency = time.Since(start).Round(time.Millisecond).String()
	if err != nil {
		p.Status = "fail"
		p.Detail = err.Error()
		return p
	}
	_ = conn.Close()
	p.Status = "ok"
	p.Detail = "tcp open"
	return p
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func amqpHostPort(raw string) (host, port string, ok bool) {
	// Accept amqp://user:pass@host:5672/vhost and plain host:port.
	if strings.Contains(raw, "://") {
		u, err := url.Parse(raw)
		if err != nil {
			return "", "", false
		}
		h := u.Hostname()
		p := u.Port()
		if p == "" {
			p = "5672"
		}
		if h == "" {
			return "", "", false
		}
		return h, p, true
	}
	h, p, err := net.SplitHostPort(raw)
	if err != nil {
		return "", "", false
	}
	return h, p, true
}
