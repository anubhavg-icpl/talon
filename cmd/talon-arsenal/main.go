// Command talon-arsenal is an MCP stdio server exposing the Talon Arsenal
// Engine's HTTP API as 151 tools.
package main

import (
	"log"
	"time"

	"github.com/anubhavg-icpl/talon/internal/arsenal"
	"github.com/anubhavg-icpl/talon/internal/config"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg := config.LoadHexstrikeConfig()
	client := arsenal.NewClient(cfg.ServerURL, time.Duration(cfg.Timeout)*time.Second)

	srv := server.NewMCPServer("talon-arsenal", "1.0.0")
	arsenal.Register(srv, client)

	if err := server.ServeStdio(srv); err != nil {
		log.Fatalf("talon-arsenal: %v", err)
	}
}
