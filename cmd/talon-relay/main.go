// Command talon-relay is the Go port of executor.py: an AMQP consumer
// replacing the Celery task pairing (execute_agent_task / get_pentest_ouput)
// with a plain consumer/publisher against the same RabbitMQ queue names.
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

// codeModelID mirrors code_gen.py's dedicated code_model.
const codeModelID = "us.meta.llama4-maverick-17b-instruct-v1:0"

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
		log.Fatal("talon-relay: AMQP_URL is not set")
	}

	region := getenv("BEDROCK_REGION", "us-east-1")

	model, err := llm.NewBedrock(ctx, getenv("AGENT_MODEL_ID", "qwen.qwen3-vl-235b-a22b"), region, 0.3, 1000)
	if err != nil {
		log.Fatalf("talon-relay: init agent model: %v", err)
	}
	judge, err := llm.NewBedrock(ctx, getenv("JUDGE_MODEL_ID", "openai.gpt-oss-120b-1:0"), region, 0.3, 1000)
	if err != nil {
		log.Fatalf("talon-relay: init judge model: %v", err)
	}
	codeModel, err := llm.NewBedrock(ctx, getenv("CODE_MODEL_ID", codeModelID), region, 0.3, 1000)
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
