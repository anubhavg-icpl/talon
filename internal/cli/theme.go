package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
)

// Athesis-inspired OLED black + chatak-laal palette for the operator TUI.
// Warm off-white body text; no blue/purple defaults.
var (
	ColorBg       = lipgloss.Color("#000000")
	ColorSurface  = lipgloss.Color("#0a0a0b")
	ColorSurface2 = lipgloss.Color("#121215")
	ColorBorder   = lipgloss.Color("#26262c")
	ColorBorderHi = lipgloss.Color("#33333b")
	ColorText     = lipgloss.Color("#e8e8ea")
	ColorDim      = lipgloss.Color("#9a9aa2")
	ColorMute     = lipgloss.Color("#6b6b73")
	// Accent is chatak laal with a slightly modern coral edge for legibility on black.
	ColorAccent     = lipgloss.Color("#ff2a44")
	ColorAccentSoft = lipgloss.Color("#3a0a12")
	ColorLive       = lipgloss.Color("#2ee06a")
	ColorWarn       = lipgloss.Color("#f5b23a")
	ColorDanger     = lipgloss.Color("#ff0033")
)

// Theme styles shared by boxed table output and the interactive console.
type Theme struct {
	Enabled bool
	Base    lipgloss.Style
	Title   lipgloss.Style
	Label   lipgloss.Style
	Value   lipgloss.Style
	Dim     lipgloss.Style
	Mute    lipgloss.Style
	Accent  lipgloss.Style
	Live    lipgloss.Style
	Warn    lipgloss.Style
	Danger  lipgloss.Style
	OK      lipgloss.Style
	Box     lipgloss.Style
	Header  lipgloss.Style
	Footer  lipgloss.Style
}

// NewTheme builds styles. Color is enabled when out is a TTY and NO_COLOR is unset.
func NewTheme(out io.Writer) Theme {
	enabled := colorEnabled(out)
	t := Theme{Enabled: enabled}
	if !enabled {
		t.Base = lipgloss.NewStyle()
		t.Title = lipgloss.NewStyle().Bold(true)
		t.Label = lipgloss.NewStyle()
		t.Value = lipgloss.NewStyle()
		t.Dim = lipgloss.NewStyle()
		t.Mute = lipgloss.NewStyle()
		t.Accent = lipgloss.NewStyle().Bold(true)
		t.Live = lipgloss.NewStyle()
		t.Warn = lipgloss.NewStyle()
		t.Danger = lipgloss.NewStyle()
		t.OK = lipgloss.NewStyle()
		t.Box = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
		t.Header = lipgloss.NewStyle().Bold(true)
		t.Footer = lipgloss.NewStyle()
		return t
	}
	t.Base = lipgloss.NewStyle().Foreground(ColorText).Background(ColorBg)
	t.Title = lipgloss.NewStyle().Foreground(ColorText).Bold(true).Background(ColorBg)
	t.Label = lipgloss.NewStyle().Foreground(ColorMute).Background(ColorBg)
	t.Value = lipgloss.NewStyle().Foreground(ColorText).Background(ColorBg)
	t.Dim = lipgloss.NewStyle().Foreground(ColorDim).Background(ColorBg)
	t.Mute = lipgloss.NewStyle().Foreground(ColorMute).Background(ColorBg)
	t.Accent = lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Background(ColorBg)
	t.Live = lipgloss.NewStyle().Foreground(ColorLive).Background(ColorBg)
	t.Warn = lipgloss.NewStyle().Foreground(ColorWarn).Background(ColorBg)
	t.Danger = lipgloss.NewStyle().Foreground(ColorDanger).Bold(true).Background(ColorBg)
	t.OK = lipgloss.NewStyle().Foreground(ColorLive).Background(ColorBg)
	t.Box = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(ColorBorder).
		Foreground(ColorText).
		Background(ColorSurface).
		Padding(0, 1)
	t.Header = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Background(ColorBg)
	t.Footer = lipgloss.NewStyle().Foreground(ColorMute).Background(ColorBg)
	return t
}

func colorEnabled(out io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if f, ok := out.(*os.File); ok {
		return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
	}
	return false
}

// BoxTitle draws a titled panel (rectangle) with body lines.
func (t Theme) BoxTitle(title string, body string, width int) string {
	if width < 20 {
		width = 20
	}
	inner := width - 4
	if inner < 8 {
		inner = 8
	}
	head := t.Accent.Render(" " + strings.ToUpper(title) + " ")
	content := body
	if content == "" {
		content = t.Mute.Render("—")
	}
	// Wrap long lines to panel width.
	var wrapped []string
	for _, line := range strings.Split(content, "\n") {
		wrapped = append(wrapped, wrapRunes(line, inner)...)
	}
	bodyBlock := strings.Join(wrapped, "\n")
	panel := t.Box.Width(width).Render(head + "\n" + bodyBlock)
	return panel
}

// BrandBar is the sticky-style header used by CLI table mode and the TUI.
func (t Theme) BrandBar(subtitle string) string {
	mark := t.Accent.Render("◆")
	name := t.Title.Render("TALON")
	sub := t.Mute.Render(" · " + subtitle)
	return mark + " " + name + sub
}

// StatusDot returns a colored ● for ok/fail/warn states.
func (t Theme) StatusDot(status string) string {
	switch strings.ToLower(status) {
	case "ok", "healthy", "completed", "success", "done", "running", "live":
		return t.Live.Render("●")
	case "degraded", "awaiting_approval", "warn", "warning":
		return t.Warn.Render("●")
	case "fail", "failed", "error", "down":
		return t.Danger.Render("●")
	default:
		return t.Mute.Render("●")
	}
}

// StatusLabel colors a status string.
func (t Theme) StatusLabel(status string) string {
	s := strings.ToLower(status)
	switch s {
	case "ok", "healthy", "completed", "success", "done":
		return t.OK.Render(status)
	case "running", "live":
		return t.Live.Render(status)
	case "degraded", "awaiting_approval":
		return t.Warn.Render(status)
	case "fail", "failed", "error", "down":
		return t.Danger.Render(status)
	default:
		return t.Dim.Render(status)
	}
}

// KeyValueBox renders a titled key/value panel.
func (t Theme) KeyValueBox(title string, rows [][2]string, width int) string {
	if width <= 0 {
		width = 72
	}
	var b strings.Builder
	labelW := 0
	for _, r := range rows {
		if n := utf8.RuneCountInString(r[0]); n > labelW {
			labelW = n
		}
	}
	if labelW < 8 {
		labelW = 8
	}
	if labelW > 22 {
		labelW = 22
	}
	for i, r := range rows {
		if i > 0 {
			b.WriteByte('\n')
		}
		key := t.Label.Render(padRight(r[0], labelW))
		val := t.Value.Render(r[1])
		b.WriteString(key)
		b.WriteString("  ")
		b.WriteString(val)
	}
	return t.BoxTitle(title, b.String(), width)
}

// ServiceRowsBox formats stack probes as a services panel.
func (t Theme) ServiceRowsBox(title string, rows []serviceProbe, width int) string {
	var b strings.Builder
	for i, s := range rows {
		if i > 0 {
			b.WriteByte('\n')
		}
		dot := t.StatusDot(s.Status)
		name := t.Value.Render(fmt.Sprintf("%-16s", s.Name))
		st := t.StatusLabel(s.Status)
		extra := ""
		if s.Latency != "" {
			extra = t.Mute.Render(" " + s.Latency)
		}
		if s.Detail != "" {
			extra += t.Mute.Render(" — " + s.Detail)
		}
		b.WriteString(dot + " " + name + " " + st + extra)
	}
	return t.BoxTitle(title, b.String(), width)
}

func padRight(s string, n int) string {
	r := utf8.RuneCountInString(s)
	if r >= n {
		return s
	}
	return s + strings.Repeat(" ", n-r)
}

func wrapRunes(s string, width int) []string {
	if width <= 0 || utf8.RuneCountInString(s) <= width {
		return []string{s}
	}
	var out []string
	runes := []rune(s)
	for len(runes) > 0 {
		if len(runes) <= width {
			out = append(out, string(runes))
			break
		}
		out = append(out, string(runes[:width]))
		runes = runes[width:]
	}
	return out
}

// termWidth guesses terminal width for boxed layout.
func termWidth() int {
	// lipgloss can read from env; fall back to 88 for readable panels.
	w := lipgloss.Width(strings.Repeat("x", 88))
	_ = w
	if cols := os.Getenv("COLUMNS"); cols != "" {
		var n int
		if _, err := fmt.Sscanf(cols, "%d", &n); err == nil && n >= 40 {
			return n
		}
	}
	return 88
}
