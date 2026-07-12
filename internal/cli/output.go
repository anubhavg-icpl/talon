package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// OutputFormat is how command results are printed.
type OutputFormat string

const (
	OutputTable OutputFormat = "table"
	OutputJSON  OutputFormat = "json"
	OutputRaw   OutputFormat = "raw"
)

// ParseOutputFormat accepts table|json|raw (case-insensitive).
func ParseOutputFormat(s string) (OutputFormat, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "table", "text":
		return OutputTable, nil
	case "json":
		return OutputJSON, nil
	case "raw":
		return OutputRaw, nil
	default:
		return "", fmt.Errorf("unknown output format %q (want table|json|raw)", s)
	}
}

// Printer writes command results.
type Printer struct {
	Format OutputFormat
	Out    io.Writer
	Err    io.Writer
}

// NewPrinter builds a Printer writing to stdout/stderr.
func NewPrinter(format OutputFormat) *Printer {
	return &Printer{Format: format, Out: os.Stdout, Err: os.Stderr}
}

// JSON emits v as pretty JSON (or compact if format is raw).
func (p *Printer) JSON(v any) error {
	enc := json.NewEncoder(p.Out)
	if p.Format != OutputRaw {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(v)
}

// PrintValue prints v as JSON when format is json/raw, otherwise uses tableFn.
func (p *Printer) PrintValue(v any, tableFn func(io.Writer) error) error {
	switch p.Format {
	case OutputJSON, OutputRaw:
		return p.JSON(v)
	default:
		if tableFn == nil {
			return p.JSON(v)
		}
		return tableFn(p.Out)
	}
}

// KeyValueTable writes key/value pairs as a simple two-column table.
func KeyValueTable(w io.Writer, rows [][2]string) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	for _, r := range rows {
		if _, err := fmt.Fprintf(tw, "%s\t%s\n", r[0], r[1]); err != nil {
			return err
		}
	}
	return tw.Flush()
}

// ToolsTable prints a compact tool log.
func ToolsTable(w io.Writer, tools []ToolCallRecord) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "INDEX\tTOOL\tOUTPUT_PREVIEW"); err != nil {
		return err
	}
	for _, t := range tools {
		preview := strings.ReplaceAll(t.Output, "\n", " ")
		if len(preview) > 80 {
			preview = preview[:77] + "..."
		}
		if _, err := fmt.Fprintf(tw, "%d\t%s\t%s\n", t.Index, t.ToolName, preview); err != nil {
			return err
		}
	}
	return tw.Flush()
}

// Exit codes for scripts and automation.
const (
	ExitOK               = 0
	ExitError            = 1
	ExitAwaitingApproval = 2
	ExitNotFound         = 3
	ExitUsage            = 64
)

// IsTerminalStatus reports whether a run status is done (success or failure).
func IsTerminalStatus(status string) bool {
	switch strings.ToLower(status) {
	case "completed", "error", "failed", "done", "success":
		return true
	default:
		return false
	}
}

// ExitCodeForStatus maps a run status to a process exit code.
func ExitCodeForStatus(status string) int {
	switch strings.ToLower(status) {
	case "completed", "done", "success":
		return ExitOK
	case "awaiting_approval":
		return ExitAwaitingApproval
	case "not_found":
		return ExitNotFound
	case "error", "failed":
		return ExitError
	default:
		return ExitOK // still running is not a hard failure for `status`
	}
}
