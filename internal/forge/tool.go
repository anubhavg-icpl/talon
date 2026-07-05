package forge

import (
	"context"

	"github.com/anubhavg-icpl/pentester2/internal/llm"
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

// Call runs the generate/execute/feedback retry loop and returns the final
// execution output, mirroring custom_exploit()'s `return result` in
// final.py (there result was the whole final_state dict; here we return the
// last observed output text, which is what the caller LLM actually needs).
func (t *CustomExploitTool) Call(ctx context.Context, query string) (string, error) {
	return Coder(ctx, t.Model, query)
}
