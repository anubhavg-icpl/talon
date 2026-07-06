package llm

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

// Bedrock implements ChatModel via AWS Bedrock's Converse API.
type Bedrock struct {
	client      *bedrockruntime.Client
	modelID     string
	temperature float32
	maxTokens   int32
}

func NewBedrock(ctx context.Context, modelID, region string, temperature float32, maxTokens int32) (*Bedrock, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("llm: load aws config: %w", err)
	}
	return &Bedrock{
		client:      bedrockruntime.NewFromConfig(cfg),
		modelID:     modelID,
		temperature: temperature,
		maxTokens:   maxTokens,
	}, nil
}

func (b *Bedrock) Converse(ctx context.Context, systemPrompt string, messages []Message, tools []ToolSpec) (Message, error) {
	sdkMessages := make([]types.Message, 0, len(messages))
	for _, m := range messages {
		sdkMessages = append(sdkMessages, toSDKMessage(m))
	}

	input := &bedrockruntime.ConverseInput{
		ModelId:  aws.String(b.modelID),
		Messages: sdkMessages,
		InferenceConfig: &types.InferenceConfiguration{
			Temperature: aws.Float32(b.temperature),
			MaxTokens:   aws.Int32(b.maxTokens),
		},
	}
	if systemPrompt != "" {
		input.System = []types.SystemContentBlock{&types.SystemContentBlockMemberText{Value: systemPrompt}}
	}
	if len(tools) > 0 {
		sdkTools := make([]types.Tool, 0, len(tools))
		for _, t := range tools {
			schema := t.InputSchema
			if schema == nil {
				schema = map[string]any{"type": "object", "properties": map[string]any{}}
			}
			sdkTools = append(sdkTools, &types.ToolMemberToolSpec{Value: types.ToolSpecification{
				Name:        aws.String(t.Name),
				Description: aws.String(t.Description),
				InputSchema: &types.ToolInputSchemaMemberJson{Value: document.NewLazyDocument(schema)},
			}})
		}
		input.ToolConfig = &types.ToolConfiguration{Tools: sdkTools}
	}

	out, err := b.client.Converse(ctx, input)
	if err != nil {
		return Message{}, fmt.Errorf("llm: converse: %w", err)
	}

	msgOut, ok := out.Output.(*types.ConverseOutputMemberMessage)
	if !ok {
		return Message{}, fmt.Errorf("llm: unexpected converse output type %T", out.Output)
	}

	result := Message{Role: RoleAssistant}
	for _, block := range msgOut.Value.Content {
		switch c := block.(type) {
		case *types.ContentBlockMemberText:
			result.Text += c.Value
		case *types.ContentBlockMemberToolUse:
			var args map[string]any
			if c.Value.Input != nil {
				if err := c.Value.Input.UnmarshalSmithyDocument(&args); err != nil {
					return Message{}, fmt.Errorf("llm: decode tool_use input: %w", err)
				}
			}
			result.ToolCalls = append(result.ToolCalls, ToolCall{
				ID:   aws.ToString(c.Value.ToolUseId),
				Name: aws.ToString(c.Value.Name),
				Args: args,
			})
		}
	}
	return result, nil
}

func toSDKMessage(m Message) types.Message {
	switch m.Role {
	case RoleAssistant:
		var blocks []types.ContentBlock
		if m.Text != "" {
			blocks = append(blocks, &types.ContentBlockMemberText{Value: m.Text})
		}
		for _, tc := range m.ToolCalls {
			blocks = append(blocks, &types.ContentBlockMemberToolUse{Value: types.ToolUseBlock{
				ToolUseId: aws.String(tc.ID),
				Name:      aws.String(tc.Name),
				Input:     document.NewLazyDocument(tc.Args),
			}})
		}
		return types.Message{Role: types.ConversationRoleAssistant, Content: blocks}
	case RoleTool:
		var blocks []types.ContentBlock
		for _, tr := range m.ToolResults {
			status := types.ToolResultStatusSuccess
			if tr.IsError {
				status = types.ToolResultStatusError
			}
			blocks = append(blocks, &types.ContentBlockMemberToolResult{Value: types.ToolResultBlock{
				ToolUseId: aws.String(tr.ToolCallID),
				Content:   []types.ToolResultContentBlock{&types.ToolResultContentBlockMemberText{Value: tr.Content}},
				Status:    status,
			}})
		}
		// Converse API expects tool results as a "user" turn following the
		// assistant's tool_use turn.
		return types.Message{Role: types.ConversationRoleUser, Content: blocks}
	default: // RoleUser / RoleSystem (system prompt is sent via input.System, not a message)
		return types.Message{Role: types.ConversationRoleUser, Content: []types.ContentBlock{&types.ContentBlockMemberText{Value: m.Text}}}
	}
}
