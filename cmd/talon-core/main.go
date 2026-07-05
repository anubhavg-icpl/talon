// Command talon-core is the Go port of pentest_core/fast_api.py: an HTTP
// front end over the agent orchestrator (pentest_core/final.py's
// build_agent()), mirroring uvicorn.run(app, host="0.0.0.0", port=8000).
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
	// codeModelID mirrors code_gen.py's dedicated code_model, kept distinct
	// from mainModelID since custom exploit generation is a different task
	// profile than orchestration.
	codeModelID        = "us.meta.llama4-maverick-17b-instruct-v1:0"
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
		log.Println("talon-core: warning: MSF_PASSWORD is not set; Metasploit MCP tool calls will fail")
	}
	log.Printf("talon-core: hexstrike server %s (timeout %ds)", hexCfg.ServerURL, hexCfg.Timeout)

	ctx := context.Background()

	tools, err := mcpclient.NewMulti(ctx, []mcpclient.ServerSpec{
		{Name: "hexstrike", Command: mcpBinaryPath("HEXSTRIKE_MCP_PATH", "talon-arsenal")},
		{Name: "metasploit", Command: mcpBinaryPath("METASPLOIT_MCP_PATH", "talon-strike")},
	})
	if err != nil {
		log.Fatalf("talon-core: start mcp servers: %v", err)
	}
	defer tools.Close()

	model, err := llm.NewBedrock(ctx, mainModelID, bedrockRegion, bedrockTemperature, bedrockMaxTokens)
	if err != nil {
		log.Fatalf("talon-core: init main model: %v", err)
	}
	judge, err := llm.NewBedrock(ctx, judgeModelID, bedrockRegion, bedrockTemperature, bedrockMaxTokens)
	if err != nil {
		log.Fatalf("talon-core: init judge model: %v", err)
	}
	codeModel, err := llm.NewBedrock(ctx, codeModelID, bedrockRegion, bedrockTemperature, bedrockMaxTokens)
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
