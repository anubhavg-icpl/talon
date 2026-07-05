// Command worker is the Go port of executor.py: an AMQP consumer replacing
// the Celery task pairing (execute_agent_task / get_pentest_ouput) with a
// plain consumer/publisher against the same RabbitMQ queue names.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anubhavg-icpl/pentester2/go/internal/agent"
	"github.com/anubhavg-icpl/pentester2/go/internal/config"
	"github.com/anubhavg-icpl/pentester2/go/internal/llm"
	"github.com/anubhavg-icpl/pentester2/go/internal/mcpclient"
	"github.com/anubhavg-icpl/pentester2/go/internal/queue"
)

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	amqpCfg := config.LoadAMQPConfig()
	if amqpCfg.URL == "" {
		log.Fatal("worker: AMQP_URL is not set")
	}

	region := getenv("BEDROCK_REGION", "us-east-1")

	model, err := llm.NewBedrock(ctx, getenv("AGENT_MODEL_ID", "qwen.qwen3-vl-235b-a22b"), region, 0.3, 1000)
	if err != nil {
		log.Fatalf("worker: init agent model: %v", err)
	}
	judge, err := llm.NewBedrock(ctx, getenv("JUDGE_MODEL_ID", "openai.gpt-oss-120b-1:0"), region, 0.3, 1000)
	if err != nil {
		log.Fatalf("worker: init judge model: %v", err)
	}

	tools, err := mcpclient.NewMulti(ctx, []mcpclient.ServerSpec{
		{Name: "hexstrike", Command: getenv("HEXSTRIKE_MCP_BIN", "hexstrike-mcp")},
		{Name: "metasploit", Command: getenv("METASPLOIT_MCP_BIN", "metasploit-mcp"), Args: []string{"--transport", "stdio"}},
	})
	if err != nil {
		log.Fatalf("worker: init mcp servers: %v", err)
	}
	defer tools.Close()

	// ponytail: codegen (custom_exploit) not wired in yet -- that subagent
	// tool is code_gen.py's own LLM-driven retry workflow, ported separately
	// in internal/codegen. Upgrade to a real agent.CodegenTool adapter once
	// that package exposes one.
	orchestrator := agent.New(model, judge, tools, nil)

	w, err := queue.NewWorker(amqpCfg.URL)
	if err != nil {
		log.Fatalf("worker: %v", err)
	}
	defer w.Close()

	log.Println("worker: consuming execute_agent_task")
	if err := w.Consume(ctx, orchestrator); err != nil && ctx.Err() == nil {
		log.Fatalf("worker: consume: %v", err)
	}
}
