// Package config holds run-scoped attacker context and process-wide,
// env-var-driven configuration. No setting here has a hardcoded credential
// fallback -- anything required fails fast if unset.
package config

import (
	"os"
	"strconv"
)

// Context is the attacker context for a single run (LHOST/LPORT for
// reverse shells and listeners).
type Context struct {
	LHOST string
	LPORT int
}

func DefaultContext() Context {
	return Context{
		LHOST: getenv("LHOST", "127.0.0.1"),
		LPORT: getenvInt("LPORT", 4444),
	}
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
	// PayloadSaveDir is where generated Metasploit payloads are written.
	PayloadSaveDir string
}

func LoadMSFConfig() MSFConfig {
	home, _ := os.UserHomeDir()
	return MSFConfig{
		// ponytail: no hardcoded password fallback -- required env var, fail
		// fast instead of silently using a guessable default in production.
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

// LLMConfig selects and configures the ChatModel backend. Provider picks the
// implementation: "bedrock" (default), "ollama" (local GPU, zero cloud), or
// "openai" (any OpenAI-compatible /chat/completions endpoint -- OpenAI, Azure,
// z.ai/GLM, vLLM, LiteLLM). Only the fields relevant to the active provider
// are read; the others may stay blank. No credential has a hardcoded fallback
// -- OPENAI_API_KEY / AWS creds are required when their provider is selected.
type LLMConfig struct {
	Provider string // bedrock | ollama | openai

	// Ollama (LLM_PROVIDER=ollama)
	OllamaURL string

	// OpenAI-compatible (LLM_PROVIDER=openai)
	OpenAIBaseURL string
	OpenAIAPIKey  string

	// Bedrock (LLM_PROVIDER=bedrock, the default)
	BedrockRegion string

	// Shared sampling params (Bedrock path uses them today; Ollama/OpenAI
	// rely on the model's own defaults since those servers set them at the
	// model/Modelfile level).
	Temperature float32
	MaxTokens   int32
}

func LoadLLMConfig() LLMConfig {
	return LLMConfig{
		Provider:      getenv("LLM_PROVIDER", "bedrock"),
		OllamaURL:     getenv("OLLAMA_URL", "http://localhost:11434"),
		OpenAIBaseURL: getenv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		OpenAIAPIKey:  os.Getenv("OPENAI_API_KEY"),
		BedrockRegion: getenv("AWS_REGION", "us-east-1"),
		Temperature:   getenvFloat("LLM_TEMPERATURE", 0.3),
		MaxTokens:     int32(getenvInt("LLM_MAX_TOKENS", 1000)),
	}
}

// Role names for per-role model resolution. The codegen ("code") role can run
// under a different provider than the rest via LLM_CODE_PROVIDER -- the
// production-default best-of-both: a hosted model for orchestration/judging
// and an uncensored local model (e.g. talon-cyber) for exploit generation.
const (
	RoleMain  = "main"
	RoleJudge = "judge"
	RoleCode  = "code"
)

// ProviderFor returns the effective provider for a role, honoring the
// LLM_CODE_PROVIDER override (so codegen can run on ollama+talon-cyber while
// the orchestrator and judge run on a hosted provider).
func ProviderFor(cfg LLMConfig, role string) string {
	if role == RoleCode {
		if v := os.Getenv("LLM_CODE_PROVIDER"); v != "" {
			return v
		}
	}
	return cfg.Provider
}

// ModelIDFor resolves the model string for a role under its effective provider.
// Env vars, by provider and role:
//
//	bedrock: AGENT_MODEL_ID (main), JUDGE_MODEL_ID (judge), CODE_MODEL_ID (code)
//	ollama:  OLLAMA_MAIN_MODEL, OLLAMA_JUDGE_MODEL, OLLAMA_CODE_MODEL
//	openai:  OPENAI_MAIN_MODEL, OPENAI_JUDGE_MODEL, OPENAI_CODE_MODEL
//
// Judge defaults to the main model when its own env is unset, since judging is
// a lighter task; codegen always needs its own (specialized) model.
func ModelIDFor(cfg LLMConfig, role string) string {
	provider := ProviderFor(cfg, role)
	switch provider {
	case "ollama":
		// Single local model ("talon", built from models/Modelfile) covers all
		// three roles by default; per-role env still lets an operator split if
		// they ever want to.
		if role == RoleMain {
			return getenv("OLLAMA_MAIN_MODEL", "talon")
		}
		if role == RoleCode {
			return getenv("OLLAMA_CODE_MODEL", "talon")
		}
		// judge
		return getenv("OLLAMA_JUDGE_MODEL", getenv("OLLAMA_MAIN_MODEL", "talon"))
	case "openai":
		if role == RoleMain {
			return getenv("OPENAI_MAIN_MODEL", "glm-5.2")
		}
		if role == RoleCode {
			return getenv("OPENAI_CODE_MODEL", "glm-5.2")
		}
		return getenv("OPENAI_JUDGE_MODEL", getenv("OPENAI_MAIN_MODEL", "glm-5.2"))
	default: // bedrock
		if role == RoleMain {
			return getenv("AGENT_MODEL_ID", "qwen.qwen3-vl-235b-a22b")
		}
		if role == RoleCode {
			return getenv("CODE_MODEL_ID", "us.meta.llama4-maverick-17b-instruct-v1:0")
		}
		return getenv("JUDGE_MODEL_ID", "openai.gpt-oss-120b-1:0")
	}
}

// ResolveModel is the convenience pair (provider, modelID) for one role, the
// shape cmd binaries hand to llm.NewModel.
func ResolveModel(cfg LLMConfig, role string) (provider, modelID string) {
	return ProviderFor(cfg, role), ModelIDFor(cfg, role)
}

func getenvFloat(key string, fallback float64) float32 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 32); err == nil {
			return float32(f)
		}
	}
	return float32(fallback)
}
