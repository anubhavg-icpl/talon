// Command api is the Go port of pentest_core/fast_api.py: an HTTP front end
// over the agent orchestrator (pentest_core/final.py's build_agent()),
// mirroring uvicorn.run(app, host="0.0.0.0", port=8000).
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/anubhavg-icpl/pentester2/go/internal/agent"
	"github.com/anubhavg-icpl/pentester2/go/internal/api"
	"github.com/anubhavg-icpl/pentester2/go/internal/config"
	"github.com/anubhavg-icpl/pentester2/go/internal/llm"
	"github.com/anubhavg-icpl/pentester2/go/internal/mcpclient"
)

const (
	mainModelID        = "qwen.qwen3-vl-235b-a22b"
	judgeModelID       = "openai.gpt-oss-120b-1:0"
	bedrockRegion      = "us-east-1"
	bedrockTemperature = 0.3
	bedrockMaxTokens   = 1000
)

func main() {
	// Load process-wide config up front so missing/invalid env (e.g. no
	// MSF_PASSWORD) fails fast instead of surfacing later as a confusing MCP
	// tool-call error, mirroring the config reads scattered through final.py.
	hexCfg := config.LoadHexstrikeConfig()
	msfCfg := config.LoadMSFConfig()
	if msfCfg.Password == "" {
		log.Println("api: warning: MSF_PASSWORD is not set; Metasploit MCP tool calls will fail")
	}
	log.Printf("api: hexstrike server %s (timeout %ds)", hexCfg.ServerURL, hexCfg.Timeout)

	ctx := context.Background()

	tools, err := mcpclient.NewMulti(ctx, []mcpclient.ServerSpec{
		{Name: "hexstrike", Command: mcpBinaryPath("HEXSTRIKE_MCP_PATH", "hexstrike-mcp")},
		{Name: "metasploit", Command: mcpBinaryPath("METASPLOIT_MCP_PATH", "metasploit-mcp")},
	})
	if err != nil {
		log.Fatalf("api: start mcp servers: %v", err)
	}
	defer tools.Close()

	model, err := llm.NewBedrock(ctx, mainModelID, bedrockRegion, bedrockTemperature, bedrockMaxTokens)
	if err != nil {
		log.Fatalf("api: init main model: %v", err)
	}
	judge, err := llm.NewBedrock(ctx, judgeModelID, bedrockRegion, bedrockTemperature, bedrockMaxTokens)
	if err != nil {
		log.Fatalf("api: init judge model: %v", err)
	}

	// ponytail: nil CodegenTool, upgrade when internal/codegen (the
	// code_gen.py docker-sandbox port) lands -- agent.Orchestrator only
	// invokes it from the codegen subagent path, which no run reaches yet.
	orch := agent.New(model, judge, tools, nil)

	store := api.NewStore()
	srv := api.NewServer(orch, store)

	log.Println("api: listening on :8000")
	if err := http.ListenAndServe(":8000", srv.Mux()); err != nil {
		log.Fatalf("api: %v", err)
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
