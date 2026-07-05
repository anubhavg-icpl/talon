// Command hexstrike-mcp is the Go port of pentest_core/hexstrike_mcp.py: an
// MCP stdio server exposing the HexStrike AI HTTP API as 151 tools.
package main

import (
	"log"
	"time"

	"github.com/anubhavg-icpl/pentester2/go/internal/config"
	"github.com/anubhavg-icpl/pentester2/go/internal/hexstrike"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg := config.LoadHexstrikeConfig()
	client := hexstrike.NewClient(cfg.ServerURL, time.Duration(cfg.Timeout)*time.Second)

	srv := server.NewMCPServer("hexstrike-ai-mcp", "1.0.0")
	hexstrike.Register(srv, client)

	if err := server.ServeStdio(srv); err != nil {
		log.Fatalf("hexstrike-mcp: %v", err)
	}
}
