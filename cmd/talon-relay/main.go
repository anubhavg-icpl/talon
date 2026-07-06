// Command talon-relay is the AMQP queue worker: a plain consumer/publisher
// that runs the agent orchestrator against jobs pulled off RabbitMQ.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anubhavg-icpl/pentester2/internal/config"
	"github.com/anubhavg-icpl/pentester2/internal/core"
	"github.com/anubhavg-icpl/pentester2/internal/forge"
	"github.com/anubhavg-icpl/pentester2/internal/llm"
	"github.com/anubhavg-icpl/pentester2/internal/mcpclient"
	"github.com/anubhavg-icpl/pentester2/internal/relay"
)

// codeModelID is the dedicated model used for custom exploit generation,
// kept distinct from the main orchestrator model.
const codeModelID = "us.meta.llama4-maverick-17b-instruct-v1:0"

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// newModel builds a ChatModel per llmCfg.Provider: Bedrock (default) or a
// local Ollama server, so the whole platform can run with zero AWS
// dependency for inference if LLM_PROVIDER=ollama.
func newModel(ctx context.Context, llmCfg config.LLMConfig, region, bedrockModelID, ollamaModel string) (llm.ChatModel, error) {
	if llmCfg.Provider == "ollama" {
		return llm.NewOllama(llmCfg.OllamaURL, ollamaModel), nil
	}
	return llm.NewBedrock(ctx, bedrockModelID, region, 0.3, 1000)
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

	region := getenv("BEDROCK_REGION", "us-east-1")

	model, err := newModel(ctx, llmCfg, region, getenv("AGENT_MODEL_ID", "qwen.qwen3-vl-235b-a22b"), getenv("OLLAMA_MAIN_MODEL", "qwen2.5:14b"))
	if err != nil {
		log.Fatalf("talon-relay: init agent model: %v", err)
	}
	judge, err := newModel(ctx, llmCfg, region, getenv("JUDGE_MODEL_ID", "openai.gpt-oss-120b-1:0"), getenv("OLLAMA_MAIN_MODEL", "qwen2.5:14b"))
	if err != nil {
		log.Fatalf("talon-relay: init judge model: %v", err)
	}
	codeModel, err := newModel(ctx, llmCfg, region, getenv("CODE_MODEL_ID", codeModelID), getenv("OLLAMA_CODE_MODEL", "qwen2.5-coder:14b"))
	if err != nil {
		log.Fatalf("talon-relay: init code model: %v", err)
	}

	tools, err := mcpclient.NewMulti(ctx, []mcpclient.ServerSpec{
		{Name: "hexstrike", Command: getenv("HEXSTRIKE_MCP_BIN", "talon-arsenal")},
		{Name: "metasploit", Command: getenv("METASPLOIT_MCP_BIN", "talon-strike"), Args: []string{"--transport", "stdio"}},
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
