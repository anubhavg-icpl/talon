// Command talon-strike is an MCP server exposing Metasploit RPC
// functionality (msfrpcd) as 12 tools.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/anubhavg-icpl/pentester2/internal/config"
	"github.com/anubhavg-icpl/pentester2/internal/strike"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// ponytail: only stdio transport is implemented -- the only caller (the
	// agent orchestrator) always spawns MCP servers over stdio. Upgrade to
	// add an HTTP/SSE transport if a network-facing deployment needs one.
	transport := flag.String("transport", "stdio", "MCP transport mode (only stdio is supported)")
	flag.Parse()
	if *transport != "stdio" {
		log.Fatalf("talon-strike: unsupported transport %q (only stdio is supported)", *transport)
	}

	cfg := config.LoadMSFConfig()
	client, err := strike.NewClient(context.Background(), cfg)
	if err != nil {
		log.Fatalf("talon-strike: failed to connect to Metasploit RPC: %v", err)
	}

	srv := server.NewMCPServer("talon-strike", "1.6.0")
	strike.Register(srv, client)

	if err := server.ServeStdio(srv); err != nil {
		log.Fatalf("talon-strike: %v", err)
	}
}
