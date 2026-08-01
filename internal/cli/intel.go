package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

func newSkillsCmd(opts *RootOptions) *cobra.Command {
	var brief bool
	var stage, category, query string
	var limit, offset int
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Browse CyberStrike + builtin methodology skills",
		Example: `  talon skills --brief
  talon skills --category WEB --limit 20
  talon skills --q ssrf -o json
  talon skills --stage exploit --category attack`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Build query path
			q := []string{}
			if brief {
				q = append(q, "brief=1")
			} else {
				q = append(q, "brief=1") // list is always brief; use GET by id for body
			}
			if stage != "" {
				q = append(q, "stage="+stage)
			}
			if category != "" {
				q = append(q, "category="+category)
			}
			if query != "" {
				q = append(q, "q="+query)
			}
			if limit > 0 {
				q = append(q, fmt.Sprintf("limit=%d", limit))
			} else {
				q = append(q, "limit=50")
			}
			if offset > 0 {
				q = append(q, fmt.Sprintf("offset=%d", offset))
			}
			path := "/skills?" + strings.Join(q, "&")
			var out map[string]any
			if err := opts.Client.GetJSON(cmd.Context(), path, &out); err != nil {
				return err
			}
			return opts.Printer.PrintValue(out, func(w io.Writer) error {
				total, _ := out["total"].(float64)
				count, _ := out["count"].(float64)
				fmt.Fprintf(w, "showing=%d total=%d\n", int(count), int(total))
				if stats, ok := out["stats"].(map[string]any); ok {
					fmt.Fprintf(w, "catalog: total=%v disk=%v builtin=%v\n",
						stats["total"], stats["src_disk"], stats["src_builtin"])
				}
				skills, _ := out["skills"].([]any)
				for _, s := range skills {
					row, _ := s.(map[string]any)
					if row == nil {
						continue
					}
					fmt.Fprintf(w, "  %-12v  %-18v  %v\n", row["stage"], row["category"], row["name"])
					fmt.Fprintf(w, "    id=%v\n", row["id"])
				}
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&brief, "brief", true, "metadata only (default)")
	cmd.Flags().StringVar(&stage, "stage", "", "filter stage")
	cmd.Flags().StringVar(&category, "category", "", "filter category (WEB, mitre_attack, …)")
	cmd.Flags().StringVar(&query, "q", "", "search query")
	cmd.Flags().IntVar(&limit, "limit", 50, "page size")
	cmd.Flags().IntVar(&offset, "offset", 0, "page offset")
	return cmd
}

func newAgentsCmd(opts *RootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "agents",
		Short: "List specialist agent modes",
		Example: `  talon agents
  talon agents -o json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := opts.Client.Agents(cmd.Context())
			if err != nil {
				return err
			}
			return opts.Printer.PrintValue(out, func(w io.Writer) error {
				agents, _ := out["agents"].([]any)
				for _, a := range agents {
					row, _ := a.(map[string]any)
					if row == nil {
						continue
					}
					fmt.Fprintf(w, "  %-10v  %-14v  %v\n", row["id"], row["codename"], row["name"])
					if d, ok := row["description"].(string); ok {
						fmt.Fprintf(w, "    %s\n", d)
					}
				}
				return nil
			})
		},
	}
}

func newFindingsCmd(opts *RootOptions) *cobra.Command {
	var severity string
	var limit int
	cmd := &cobra.Command{
		Use:   "findings",
		Short: "List structured findings across all runs",
		Example: `  talon findings
  talon findings --severity critical
  talon findings --limit 20 -o json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := opts.Client.GlobalFindings(cmd.Context(), severity, limit)
			if err != nil {
				return err
			}
			return opts.Printer.PrintValue(out, func(w io.Writer) error {
				count, _ := out["count"].(float64)
				fmt.Fprintf(w, "findings=%d\n", int(count))
				items, _ := out["findings"].([]any)
				for _, it := range items {
					row, _ := it.(map[string]any)
					if row == nil {
						continue
					}
					runID, _ := row["run_id"].(string)
					target, _ := row["target"].(string)
					f, _ := row["finding"].(map[string]any)
					if f == nil {
						continue
					}
					fmt.Fprintf(w, "  [%v] %v  target=%s  run=%s  id=%v\n",
						f["severity"], f["title"], target, shortRun(runID), f["id"])
				}
				if int(count) == 0 {
					fmt.Fprintln(w, "(no findings yet)")
				}
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&severity, "severity", "", "filter: critical|high|medium|low|info")
	cmd.Flags().IntVar(&limit, "limit", 50, "max findings to return")
	return cmd
}

func shortRun(id string) string {
	if len(id) > 8 {
		return id[:8]
	}
	return id
}
