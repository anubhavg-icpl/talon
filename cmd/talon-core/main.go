// Command talon-core is the HTTP control plane: a front end over the agent
// orchestrator, serving on :8000.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/anubhavg-icpl/talon/internal/config"
	"github.com/anubhavg-icpl/talon/internal/control"
	"github.com/anubhavg-icpl/talon/internal/core"
	"github.com/anubhavg-icpl/talon/internal/forge"
	"github.com/anubhavg-icpl/talon/internal/llm"
	"github.com/anubhavg-icpl/talon/internal/mcpclient"
)

// Models are built via the shared llm.NewModel factory (single provider
// switch for bedrock|ollama|openai), with per-role model IDs resolved from
// env by config.ResolveModel -- so OLLAMA_MAIN_MODEL / OPENAI_CODE_MODEL /
// AGENT_MODEL_ID etc. are honored uniformly here and in talon-relay.
func newModel(ctx context.Context, llmCfg config.LLMConfig, role string) (llm.ChatModel, error) {
	provider, modelID := config.ResolveModel(llmCfg, role)
	return llm.NewModel(ctx, llmCfg, provider, modelID)
}

func main() {
	// Load process-wide config up front so missing/invalid env (e.g. no
	// MSF_PASSWORD) fails fast instead of surfacing later as a confusing
	// MCP tool-call error.
	hexCfg := config.LoadHexstrikeConfig()
	msfCfg := config.LoadMSFConfig()
	llmCfg := config.LoadLLMConfig()
	if msfCfg.Password == "" {
		log.Println("talon-core: warning: MSF_PASSWORD is not set; Metasploit MCP tool calls will fail")
	}
	log.Printf("talon-core: hexstrike server %s (timeout %ds)", hexCfg.ServerURL, hexCfg.Timeout)
	log.Printf("talon-core: llm provider %s", llmCfg.Provider)

	ctx := context.Background()

	tools, err := mcpclient.NewMulti(ctx, []mcpclient.ServerSpec{
		{Name: "hexstrike", Command: mcpBinaryPath("HEXSTRIKE_MCP_PATH", "talon-arsenal")},
		{Name: "metasploit", Command: mcpBinaryPath("METASPLOIT_MCP_PATH", "talon-strike")},
	})
	if err != nil {
		log.Fatalf("talon-core: start mcp servers: %v", err)
	}
	defer tools.Close()

	model, err := newModel(ctx, llmCfg, config.RoleMain)
	if err != nil {
		log.Fatalf("talon-core: init main model: %v", err)
	}
	judge, err := newModel(ctx, llmCfg, config.RoleJudge)
	if err != nil {
		log.Fatalf("talon-core: init judge model: %v", err)
	}
	codeModel, err := newModel(ctx, llmCfg, config.RoleCode)
	if err != nil {
		log.Fatalf("talon-core: init code model: %v", err)
	}

	orch := core.New(model, judge, tools, forge.NewCustomExploitTool(codeModel))

	store := control.NewStore()
	srv := control.NewServer(orch, store)

	log.Println("talon-core: listening on :8000")
	if err := http.ListenAndServe(":8000", srv.Mux()); err != nil {
		log.Fatalf("talon-core: %v", err)
	}
}

// mcpBinaryPath resolves an MCP server binary path from an env var override,
// falling back to a sibling of this executable with the given name.
func mcpBinaryPath(envVar, fallbackName string) string {
	if v := os.Getenv(envVar); v != "" {
		return v
	}
	exe, err := os.Executable()
	if err != nil {
		return fallbackName
	}
	return filepath.Join(filepath.Dir(exe), fallbackName)
}
