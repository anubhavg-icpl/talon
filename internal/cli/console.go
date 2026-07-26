package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// newConsoleCmd is the Athesis-inspired black/red operator console TUI.
func newConsoleCmd(opts *RootOptions) *cobra.Command {
	var (
		runID       string
		interval    time.Duration
		autoApprove bool
	)
	cmd := &cobra.Command{
		Use:   "console",
		Short: "Interactive black/red operator console (Athesis-style TUI)",
		Long: `Live operator surface for Talon — OLED black panels, chatak-laal accents,
rectangle boxes, monospace layout. Inspired by the Athesis AGORA console.

Polls stack health and an optional run (tools + status). Keys:
  r  refresh now
  a  approve pending HITL (when a run is selected)
  j/k  scroll tool feed
  q  quit`,
		Example: `  talon console
  talon console --run <run_id>
  talon console --run <run_id> --auto-approve`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Client == nil {
				return withExitCode(ExitError, "console needs a talon-core client")
			}
			m := newConsoleModel(opts, runID, interval, autoApprove)
			p := tea.NewProgram(m, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				return withExitCode(ExitError, "console: %v", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&runID, "run", "", "run_id to follow (optional)")
	cmd.Flags().DurationVar(&interval, "interval", 2*time.Second, "poll interval")
	cmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "auto-approve HITL interrupts for the followed run")
	return cmd
}

type consoleModel struct {
	opts        *RootOptions
	theme       Theme
	runID       string
	interval    time.Duration
	autoApprove bool

	width  int
	height int

	report   *statusReport
	status   *StatusResponse
	tools    []ToolCallRecord
	history  []string
	errMsg   string
	lastTick time.Time
	feed     viewport.Model
	ready    bool
}

type tickMsg time.Time
type snapshotMsg struct {
	report  *statusReport
	status  *StatusResponse
	tools   []ToolCallRecord
	history []string
	err     error
}

func newConsoleModel(opts *RootOptions, runID string, interval time.Duration, autoApprove bool) consoleModel {
	if interval <= 0 {
		interval = 2 * time.Second
	}
	return consoleModel{
		opts:        opts,
		theme:       NewTheme(opts.Printer.Out),
		runID:       strings.TrimSpace(runID),
		interval:    interval,
		autoApprove: autoApprove,
		feed:        viewport.New(40, 10),
	}
}

func (m consoleModel) Init() tea.Cmd {
	return tea.Batch(m.fetchSnapshot(), tickCmd(m.interval))
}

func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m consoleModel) fetchSnapshot() tea.Cmd {
	opts := m.opts
	runID := m.runID
	auto := m.autoApprove
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()

		// Build status report the same way as `talon status`.
		report := &statusReport{
			CoreURL: opts.Client.BaseURL(),
			Checked: time.Now().UTC().Format(time.RFC3339),
		}
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

		arsenalURL := opts.Resolved.ArsenalURL
		if arsenalURL == "" {
			arsenalURL = envOr("TALON_ARSENAL_URL", "http://localhost:8888/health")
		}
		report.Services = append(report.Services, probeHTTP(ctx, "arsenal-engine", arsenalURL, 5*time.Second))
		msfHost := envOr("MSF_SERVER", "localhost")
		msfPort := envOr("MSF_PORT", "5554")
		report.Services = append(report.Services, probeTCP(ctx, "metasploit-rpc", msfHost+":"+msfPort, 2*time.Second))
		amqp := "localhost:5672"
		if u := envOr("AMQP_URL", ""); u != "" {
			if h, p, ok := amqpHostPort(u); ok {
				amqp = h + ":" + p
			}
		}
		report.Services = append(report.Services, probeTCP(ctx, "rabbitmq", amqp, 2*time.Second))
		report.Overall = overallFrom(report.Services)

		msg := snapshotMsg{report: report}

		if runID == "" {
			return msg
		}

		st, err := opts.Client.Status(ctx, runID)
		if err != nil {
			msg.err = err
			return msg
		}
		msg.status = st

		if auto && st.Status == "awaiting_approval" && st.Interrupt != nil {
			_, _ = opts.Client.Resume(ctx, runID, ResumeRequest{Decision: "approve"})
			// re-fetch after approve
			if st2, err2 := opts.Client.Status(ctx, runID); err2 == nil {
				msg.status = st2
			}
		}

		if tr, err := opts.Client.Tools(ctx, runID); err == nil && tr != nil {
			msg.tools = tr.ToolLog
		}
		if th, err := opts.Client.Traces(ctx, runID); err == nil && th != nil {
			msg.history = th.History
		}
		return msg
	}
}

func (m consoleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layoutFeed()
		m.ready = true
		return m, nil
	case tickMsg:
		m.lastTick = time.Time(msg)
		return m, tea.Batch(m.fetchSnapshot(), tickCmd(m.interval))
	case snapshotMsg:
		if msg.err != nil {
			m.errMsg = msg.err.Error()
		} else {
			m.errMsg = ""
		}
		if msg.report != nil {
			m.report = msg.report
		}
		if msg.status != nil {
			m.status = msg.status
		}
		if msg.tools != nil {
			m.tools = msg.tools
		}
		if msg.history != nil {
			m.history = msg.history
		}
		m.refreshFeedContent()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "r":
			return m, m.fetchSnapshot()
		case "a":
			if m.runID == "" {
				m.errMsg = "no --run set; cannot approve"
				return m, nil
			}
			return m, m.approveCmd()
		case "j", "down":
			m.feed.LineDown(1)
			return m, nil
		case "k", "up":
			m.feed.LineUp(1)
			return m, nil
		case "g":
			m.feed.GotoTop()
			return m, nil
		case "G":
			m.feed.GotoBottom()
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.feed, cmd = m.feed.Update(msg)
	return m, cmd
}

func (m consoleModel) approveCmd() tea.Cmd {
	opts := m.opts
	runID := m.runID
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if _, err := opts.Client.Resume(ctx, runID, ResumeRequest{Decision: "approve"}); err != nil {
			return snapshotMsg{err: err}
		}
		// chain a refresh
		return m.fetchSnapshot()()
	}
}

func (m *consoleModel) layoutFeed() {
	// Reserve space for header + two panels + footer.
	h := m.height - 16
	if h < 6 {
		h = 6
	}
	w := m.width - 4
	if w < 40 {
		w = 40
	}
	m.feed.Width = w
	m.feed.Height = h
}

func (m *consoleModel) refreshFeedContent() {
	var b strings.Builder
	if len(m.tools) == 0 && len(m.history) == 0 {
		b.WriteString(m.theme.Mute.Render("waiting for tool activity…"))
		m.feed.SetContent(b.String())
		return
	}
	// Prefer recent history lines then tool rows.
	if len(m.history) > 0 {
		start := 0
		if len(m.history) > 40 {
			start = len(m.history) - 40
		}
		for _, line := range m.history[start:] {
			b.WriteString(m.theme.Mute.Render(">_ "))
			b.WriteString(m.theme.Dim.Render(line))
			b.WriteByte('\n')
		}
		b.WriteByte('\n')
	}
	for _, t := range m.tools {
		preview := strings.ReplaceAll(t.Output, "\n", " ")
		if len(preview) > 90 {
			preview = preview[:87] + "..."
		}
		idx := m.theme.Mute.Render(fmt.Sprintf("[%d]", t.Index))
		name := m.theme.Accent.Render(t.ToolName)
		out := m.theme.Dim.Render(preview)
		b.WriteString(fmt.Sprintf("%s %s  %s\n", idx, name, out))
	}
	m.feed.SetContent(b.String())
	m.feed.GotoBottom()
}

func (m consoleModel) View() string {
	if !m.ready {
		return m.theme.Mute.Render(" booting console…")
	}
	w := m.width
	if w < 60 {
		w = 60
	}
	half := (w - 3) / 2
	if half < 28 {
		half = 28
	}

	// Header
	live := m.theme.Mute.Render("offline")
	if m.report != nil && m.report.Overall == "healthy" {
		live = m.theme.Live.Render("● LIVE")
	} else if m.report != nil && m.report.Overall == "degraded" {
		live = m.theme.Warn.Render("● DEGRADED")
	} else if m.report != nil {
		live = m.theme.Danger.Render("● DOWN")
	}
	core := ""
	if m.opts != nil && m.opts.Client != nil {
		core = m.theme.Mute.Render(m.opts.Client.BaseURL())
	}
	header := lipgloss.JoinHorizontal(lipgloss.Center,
		m.theme.BrandBar("operator console"),
		"  ",
		live,
		"  ",
		core,
	)

	// Services panel
	svcBody := m.theme.Mute.Render("probing…")
	if m.report != nil {
		svcBody = m.theme.ServiceRowsBox("services", m.report.Services, half)
	} else {
		svcBody = m.theme.BoxTitle("services", svcBody, half)
	}

	// Run panel
	runLines := []string{
		m.theme.Label.Render("run_id") + "  " + m.theme.Value.Render(orDash(m.runID)),
	}
	if m.status != nil {
		runLines = append(runLines,
			m.theme.Label.Render("status")+"  "+m.theme.StatusDot(m.status.Status)+" "+m.theme.StatusLabel(m.status.Status),
		)
		if m.status.JudgeVerdict != nil {
			jv := "false"
			if *m.status.JudgeVerdict {
				jv = "true"
			}
			runLines = append(runLines, m.theme.Label.Render("judge")+"   "+m.theme.Accent.Render(jv))
		}
		if m.status.Interrupt != nil {
			args, _ := json.Marshal(m.status.Interrupt.Args)
			runLines = append(runLines,
				m.theme.Label.Render("hitl")+"   "+m.theme.Warn.Render(m.status.Interrupt.ToolName),
				m.theme.Mute.Render(string(args)),
			)
		}
		if m.status.Output != "" {
			out := m.status.Output
			if len(out) > 180 {
				out = out[:177] + "..."
			}
			runLines = append(runLines, m.theme.Mute.Render(strings.ReplaceAll(out, "\n", " ")))
		}
	} else if m.runID == "" {
		runLines = append(runLines, m.theme.Mute.Render("pass --run <id> to follow a validation"))
	}
	runBody := m.theme.BoxTitle("active run", strings.Join(runLines, "\n"), half)

	top := lipgloss.JoinHorizontal(lipgloss.Top, svcBody, " ", runBody)

	// Feed panel
	feedTitle := m.theme.BoxTitle("live feed · tools / traces", m.feed.View(), w-2)

	// Footer
	keys := m.theme.Footer.Render(" r refresh · a approve · j/k scroll · g/G top/end · q quit ")
	if m.errMsg != "" {
		keys = m.theme.Danger.Render(" ! "+m.errMsg+" ") + "  " + keys
	}
	if !m.lastTick.IsZero() {
		keys += m.theme.Mute.Render("  polled " + m.lastTick.Local().Format("15:04:05"))
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		"",
		header,
		"",
		top,
		"",
		feedTitle,
		"",
		keys,
		"",
	)
}

func orDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}
