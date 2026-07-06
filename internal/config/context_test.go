package config

import (
	"os"
	"testing"
)

// TestResolveModel covers the .env -> provider/model wiring that talon-core
// and talon-relay both route through. In particular: per-role env vars, the
// judge-falls-back-to-main default, and the LLM_CODE_PROVIDER override (the
// best-of-both path where codegen runs on a local uncensored model while the
// orchestrator runs hosted).
func TestResolveModel(t *testing.T) {
	cases := []struct {
		name                                           string
		provider, codeProvider                         string
		mainM, judgeM, codeM, ollamaMain, ollamaCode   string
		role                                           string
		wantProvider, wantModel                        string
	}{
		// OpenAI / hosted GLM defaults
		{"openai main default", "openai", "", "", "", "", "", "", RoleMain, "openai", "glm-5.2"},
		{"openai code default", "openai", "", "", "", "", "", "", RoleCode, "openai", "glm-5.2"},
		{"openai main override", "openai", "", "glm-4.5", "", "", "", "", RoleMain, "openai", "glm-4.5"},

		// Ollama defaults to the single bundled "talon" model for every role
		{"ollama main default", "ollama", "", "", "", "", "", "", RoleMain, "ollama", "talon"},
		{"ollama code default", "ollama", "", "", "", "", "", "", RoleCode, "ollama", "talon"},

		// LLM_CODE_PROVIDER override: orchestrator on openai, codegen on ollama
		{"code override provider+model", "openai", "ollama", "", "", "", "", "", RoleCode, "ollama", "talon"},
		{"code override keeps main on host", "openai", "ollama", "", "", "", "", "", RoleMain, "openai", "glm-5.2"},

		// Bedrock defaults
		{"bedrock main default", "bedrock", "", "", "", "", "", "", RoleMain, "bedrock", "qwen.qwen3-vl-235b-a22b"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			os.Unsetenv("LLM_PROVIDER")
			os.Unsetenv("LLM_CODE_PROVIDER")
			os.Unsetenv("OPENAI_MAIN_MODEL")
			os.Unsetenv("OPENAI_CODE_MODEL")
			os.Unsetenv("OPENAI_JUDGE_MODEL")
			os.Unsetenv("OLLAMA_MAIN_MODEL")
			os.Unsetenv("OLLAMA_CODE_MODEL")
			os.Unsetenv("OLLAMA_JUDGE_MODEL")

			if c.provider != "" {
				os.Setenv("LLM_PROVIDER", c.provider)
			}
			if c.codeProvider != "" {
				os.Setenv("LLM_CODE_PROVIDER", c.codeProvider)
			}
			if c.mainM != "" {
				os.Setenv("OPENAI_MAIN_MODEL", c.mainM)
			}
			if c.ollamaMain != "" {
				os.Setenv("OLLAMA_MAIN_MODEL", c.ollamaMain)
			}
			if c.ollamaCode != "" {
				os.Setenv("OLLAMA_CODE_MODEL", c.ollamaCode)
			}

			cfg := LoadLLMConfig()
			gotProvider, gotModel := ResolveModel(cfg, c.role)
			if gotProvider != c.wantProvider || gotModel != c.wantModel {
				t.Fatalf("ResolveModel(%q) = (%q,%q), want (%q,%q)", c.role, gotProvider, gotModel, c.wantProvider, c.wantModel)
			}
		})
	}
}
