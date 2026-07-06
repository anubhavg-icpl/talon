// Command talon-core is the HTTP control plane: a front end over the agent
// orchestrator, serving on :8000.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/anubhavg-icpl/pentester2/internal/config"
	"github.com/anubhavg-icpl/pentester2/internal/control"
	"github.com/anubhavg-icpl/pentester2/internal/core"
	"github.com/anubhavg-icpl/pentester2/internal/forge"
	"github.com/anubhavg-icpl/pentester2/internal/llm"
	"github.com/anubhavg-icpl/pentester2/internal/mcpclient"
)

const (
	mainModelID  = "qwen.qwen3-vl-235b-a22b"
	judgeModelID = "openai.gpt-oss-120b-1:0"
	// codeModelID is the dedicated model for custom exploit generation,
	// kept distinct from mainModelID since it's a different task profile
	// than orchestration.
	codeModelID = "us.meta.llama4-maverick-17b-instruct-v1:0"

	ollamaMainModel = "qwen2.5:14b"
	ollamaCodeModel = "qwen2.5-coder:14b"

	bedrockRegion      = "us-east-1"
	bedrockTemperature = 0.3
	bedrockMaxTokens   = 1000
)

// newModel builds a ChatModel per llmCfg.Provider: Bedrock (default) or a
// local Ollama server, so the whole platform can run with zero AWS
// dependency for inference if LLM_PROVIDER=ollama.
func newModel(ctx context.Context, llmCfg config.LLMConfig, bedrockModelID, ollamaModel string) (llm.ChatModel, error) {
	if llmCfg.Provider == "ollama" {
		return llm.NewOllama(llmCfg.OllamaURL, ollamaModel), nil
	}
	return llm.NewBedrock(ctx, bedrockModelID, bedrockRegion, bedrockTemperature, bedrockMaxTokens)
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

	model, err := newModel(ctx, llmCfg, mainModelID, ollamaMainModel)
	if err != nil {
		log.Fatalf("talon-core: init main model: %v", err)
	}
	judge, err := newModel(ctx, llmCfg, judgeModelID, ollamaMainModel)
	if err != nil {
		log.Fatalf("talon-core: init judge model: %v", err)
	}
	codeModel, err := newModel(ctx, llmCfg, codeModelID, ollamaCodeModel)
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
