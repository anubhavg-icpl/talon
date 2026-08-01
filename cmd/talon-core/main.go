// Command talon-core is the HTTP control plane: a front end over the agent
// orchestrator, serving on :8000.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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

	specs := []mcpclient.ServerSpec{
		{Name: "hexstrike", Command: mcpBinaryPath("HEXSTRIKE_MCP_PATH", "talon-arsenal")},
		{Name: "metasploit", Command: mcpBinaryPath("METASPLOIT_MCP_PATH", "talon-strike")},
	}
	// Lightpanda browser-automation MCP (optional): headless browser tools
	// (goto/markdown/links/evaluate/click/…) for web-facing recon. Self-gating:
	// added only when the `lightpanda` binary is present, and Optional so a
	// missing/broken binary is skipped with a warning instead of crashing core.
	if cmd := lightpandaCommand(); cmd != "" {
		log.Printf("talon-core: lightpanda mcp enabled (%s)", cmd)
		specs = append(specs, mcpclient.ServerSpec{Name: "lightpanda", Command: cmd, Args: []string{"mcp"}, Optional: true})
	}
	tools, err := mcpclient.NewMulti(ctx, specs)
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
	dataDir := os.Getenv("TALON_DATA_DIR")
	if dataDir == "" {
		dataDir = "talon-data"
	}

	// Persistence: Postgres when TALON_DATABASE_URL works, else JSON file.
	var db *control.DB
	if url := os.Getenv("TALON_DATABASE_URL"); url != "" {
		pgCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		d, err := control.ConnectDB(pgCtx, url)
		cancel()
		if err != nil {
			log.Printf("talon-core: postgres unavailable (%v) — falling back to JSON persistence", err)
		} else {
			db = d
			defer db.Close()
			log.Println("talon-core: postgres connected")
		}
	}
	if db != nil {
		if err := store.EnablePostgres(ctx, db, dataDir); err != nil {
			log.Printf("talon-core: pg load failed (%v) — falling back to JSON persistence", err)
			db.Close()
			db = nil
		}
	}
	if db == nil {
		if err := store.EnablePersistence(dataDir); err != nil {
			log.Printf("talon-core: persistence disabled: %v", err)
		}
	}

	plat := control.NewPlatform(dataDir)
	opts := []control.ServerOption{
		control.WithAnalyzer(model),
		control.WithTools(tools),
		control.WithSettings(control.NewSettings(db)),
		control.WithPlatform(plat),
	}
	if db != nil {
		opts = append(opts, control.WithDB(db))
	}

	// Cache: Redis when REDIS_URL works, else endpoints run uncached.
	var cache control.Cache
	if url := os.Getenv("REDIS_URL"); url != "" {
		cacheCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		c, err := control.ConnectCache(cacheCtx, url)
		cancel()
		if err != nil {
			log.Printf("talon-core: redis unavailable (%v) — running uncached", err)
		} else {
			cache = c
			opts = append(opts, control.WithCache(cache))
			log.Println("talon-core: redis connected")
		}
	}
	authDisabled := os.Getenv("TALON_AUTH_DISABLED") == "true"
	switch {
	case authDisabled:
		log.Println("talon-core: warning: auth DISABLED (TALON_AUTH_DISABLED=true)")
	case db == nil:
		log.Println("talon-core: warning: auth disabled — no postgres (set TALON_DATABASE_URL)")
	case os.Getenv("TALON_ADMIN_PASSWORD") == "":
		log.Println("talon-core: warning: auth disabled — TALON_ADMIN_PASSWORD not set")
	default:
		adminUser := os.Getenv("TALON_ADMIN_USERNAME")
		if adminUser == "" {
			adminUser = "admin"
		}
		a, err := control.NewAuth(ctx, db, adminUser, os.Getenv("TALON_ADMIN_PASSWORD"))
		if err != nil {
			log.Fatalf("talon-core: init auth: %v", err)
		}
		if cache != nil {
			a.SetCache(cache)
		}
		opts = append(opts, control.WithAuth(a))
	}
	srv := control.NewServer(orch, store, opts...)

	log.Println("talon-core: listening on :8000")
	if err := http.ListenAndServe(":8000", srv.Handler()); err != nil {
		log.Fatalf("talon-core: %v", err)
	}
}

// lightpandaCommand resolves the lightpanda binary for the browser MCP server:
// LIGHTPANDA_MCP_PATH override, else PATH lookup, else "" (feature skipped).
func lightpandaCommand() string {
	if v := os.Getenv("LIGHTPANDA_MCP_PATH"); v != "" {
		return v
	}
	if p, err := exec.LookPath("lightpanda"); err == nil {
		return p
	}
	return ""
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
