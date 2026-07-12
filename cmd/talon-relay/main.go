// Command talon-relay is the AMQP queue worker: a plain consumer/publisher
// that runs the agent orchestrator against jobs pulled off RabbitMQ.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/anubhavg-icpl/talon/internal/config"
	"github.com/anubhavg-icpl/talon/internal/core"
	"github.com/anubhavg-icpl/talon/internal/forge"
	"github.com/anubhavg-icpl/talon/internal/llm"
	"github.com/anubhavg-icpl/talon/internal/mcpclient"
	"github.com/anubhavg-icpl/talon/internal/relay"
)

// newModel routes through the shared llm.NewModel factory so the provider
// switch (bedrock|ollama|openai) and per-role model resolution live in one
// place, identical to talon-core.
func newModel(ctx context.Context, llmCfg config.LLMConfig, role string) (llm.ChatModel, error) {
	provider, modelID := config.ResolveModel(llmCfg, role)
	return llm.NewModel(ctx, llmCfg, provider, modelID)
}

// mcpBinaryPath matches talon-core: honor the Dockerfile env names
// (HEXSTRIKE_MCP_PATH / METASPLOIT_MCP_PATH), then a sibling of this
// executable. Also accept the legacy *_MCP_BIN names so older compose
// overrides keep working.
func mcpBinaryPath(envVar, legacyEnvVar, fallbackName string) string {
	if v := os.Getenv(envVar); v != "" {
		return v
	}
	if v := os.Getenv(legacyEnvVar); v != "" {
		return v
	}
	exe, err := os.Executable()
	if err != nil {
		return fallbackName
	}
	return filepath.Join(filepath.Dir(exe), fallbackName)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	amqpCfg := config.LoadAMQPConfig()
	if amqpCfg.URL == "" {
		log.Fatal("talon-relay: AMQP_URL is not set")
	}
	llmCfg := config.LoadLLMConfig()
	log.Printf("talon-relay: llm provider %s", llmCfg.Provider)

	model, err := newModel(ctx, llmCfg, config.RoleMain)
	if err != nil {
		log.Fatalf("talon-relay: init agent model: %v", err)
	}
	judge, err := newModel(ctx, llmCfg, config.RoleJudge)
	if err != nil {
		log.Fatalf("talon-relay: init judge model: %v", err)
	}
	codeModel, err := newModel(ctx, llmCfg, config.RoleCode)
	if err != nil {
		log.Fatalf("talon-relay: init code model: %v", err)
	}

	tools, err := mcpclient.NewMulti(ctx, []mcpclient.ServerSpec{
		{Name: "hexstrike", Command: mcpBinaryPath("HEXSTRIKE_MCP_PATH", "HEXSTRIKE_MCP_BIN", "talon-arsenal")},
		{Name: "metasploit", Command: mcpBinaryPath("METASPLOIT_MCP_PATH", "METASPLOIT_MCP_BIN", "talon-strike"), Args: []string{"--transport", "stdio"}},
	})
	if err != nil {
		log.Fatalf("talon-relay: init mcp servers: %v", err)
	}
	defer tools.Close()

	orchestrator := core.New(model, judge, tools, forge.NewCustomExploitTool(codeModel))

	w, err := relay.NewWorker(amqpCfg.URL)
	if err != nil {
		log.Fatalf("talon-relay: %v", err)
	}
	defer w.Close()

	log.Println("talon-relay: consuming execute_agent_task")
	if err := w.Consume(ctx, orchestrator); err != nil && ctx.Err() == nil {
		log.Fatalf("talon-relay: consume: %v", err)
	}
}
