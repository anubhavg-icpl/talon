package llm

import (
	"context"
	"fmt"

	"github.com/anubhavg-icpl/talon/internal/config"
)

// NewModel is the single construction point for every ChatModel in the
// platform. Both cmd binaries (talon-core, talon-relay) route through it so
// the provider switch lives in one place -- previously each binary had its
// own copy, and talon-core's silently ignored OLLAMA_MAIN_MODEL.
//
// provider is the *effective* provider for this call (honoring a per-role
// override from config.ProviderFor); modelID is the already-resolved model
// string from config.ModelIDFor.
func NewModel(ctx context.Context, cfg config.LLMConfig, provider, modelID string) (ChatModel, error) {
	switch provider {
	case "ollama":
		return NewOllama(cfg.OllamaURL, modelID), nil
	case "openai":
		if cfg.OpenAIAPIKey == "" {
			return nil, fmt.Errorf("llm: OPENAI_API_KEY is not set (required when LLM_PROVIDER=openai; no hardcoded fallback)")
		}
		return NewOpenAI(cfg.OpenAIBaseURL, cfg.OpenAIAPIKey, modelID), nil
	case "bedrock", "":
		return NewBedrock(ctx, modelID, cfg.BedrockRegion, cfg.Temperature, cfg.MaxTokens)
	default:
		return nil, fmt.Errorf("llm: unknown provider %q (want bedrock|ollama|openai)", provider)
	}
}
