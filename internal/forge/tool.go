package forge

import (
	"context"

	"github.com/anubhavg-icpl/talon/internal/llm"
)

// CustomExploitTool implements core.CodegenTool (structurally -- this
// package intentionally does not import internal/core to avoid an import
// cycle; callers in cmd/talon-core and cmd/talon-relay wire the concrete type in).
type CustomExploitTool struct {
	Model llm.ChatModel
}

func NewCustomExploitTool(model llm.ChatModel) *CustomExploitTool {
	return &CustomExploitTool{Model: model}
}

func (t *CustomExploitTool) Name() string { return "custom_exploit" }

func (t *CustomExploitTool) Description() string {
	return "Use this tool to generate custom Python code for exploiting vulnerabilities or scanning ports."
}

// Call runs the generate/execute/feedback retry loop and returns the last
// observed execution output -- the text the calling model actually needs.
func (t *CustomExploitTool) Call(ctx context.Context, query string) (string, error) {
	return Coder(ctx, t.Model, query)
}
