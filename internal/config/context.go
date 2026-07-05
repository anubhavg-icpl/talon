// Package config holds run-scoped attacker context and process-wide env config,
// mirroring pentest_core.final.Context and the various os.getenv() reads scattered
// through the Python codebase.
package config

import (
	"os"
	"strconv"
)

// Context is the attacker context for a single run (LHOST/LPORT for reverse
// shells and listeners), equivalent to pentest_core.final.Context.
type Context struct {
	LHOST string
	LPORT int
}

func DefaultContext() Context {
	return Context{LHOST: "192.168.122.176", LPORT: 4444}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

// HexstrikeConfig is connection config for the HexStrike AI HTTP API server
// (a separate prebuilt service, not part of this rewrite).
type HexstrikeConfig struct {
	ServerURL string
	Timeout   int // seconds
}

func LoadHexstrikeConfig() HexstrikeConfig {
	return HexstrikeConfig{
		ServerURL: getenv("HEXSTRIKE_SERVER_URL", "http://localhost:8888"),
		Timeout:   getenvInt("HEXSTRIKE_TIMEOUT", 300),
	}
}

// MSFConfig is connection config for the Metasploit RPC daemon (msfrpcd).
type MSFConfig struct {
	Password string
	Server   string
	Port     string
	SSL      bool
	// PayloadSaveDir mirrors PAYLOAD_SAVE_DIR from MetasploitMCP.py.
	PayloadSaveDir string
}

func LoadMSFConfig() MSFConfig {
	home, _ := os.UserHomeDir()
	return MSFConfig{
		// ponytail: no hardcoded password fallback (Python had "network_msf" baked
		// in two places) -- required env var, fail fast instead of silently using
		// a guessable default in production.
		Password:       os.Getenv("MSF_PASSWORD"),
		Server:         getenv("MSF_SERVER", "msf_rpc"),
		Port:           getenv("MSF_PORT", "5554"),
		SSL:            getenv("MSF_SSL", "true") == "true",
		PayloadSaveDir: getenv("PAYLOAD_SAVE_DIR", home+"/payloads"),
	}
}

// AMQPConfig is the broker connection for the queue worker.
type AMQPConfig struct {
	URL string
}

func LoadAMQPConfig() AMQPConfig {
	// ponytail: no hardcoded guest:guest@localhost fallback -- required env var.
	return AMQPConfig{URL: os.Getenv("AMQP_URL")}
}
