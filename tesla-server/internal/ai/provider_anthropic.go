package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"tesla-server/config"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// anthropicProvider 通过官方 Anthropic SDK 调用 Claude。
// 不走 OpenAI 兼容 shim —— 官方 SDK 正确处理 Claude 的请求/响应差异。
type anthropicProvider struct {
	cfg *config.Config
}

func (p *anthropicProvider) Chat(systemPrompt, userPrompt string) (*ChatResponse, error) {
	cfg := p.cfg

	client := anthropic.NewClient(
		option.WithAPIKey(cfg.AI.AnthropicAPIKey),
		option.WithRequestTimeout(120*time.Second),
	)

	resp, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(cfg.AI.AnthropicModel),
		MaxTokens: 4096,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("anthropic request failed: %w", err)
	}

	// 拼接所有 text 块为最终内容（忽略 thinking 等其他块）。
	var sb strings.Builder
	for _, block := range resp.Content {
		if t, ok := block.AsAny().(anthropic.TextBlock); ok {
			sb.WriteString(t.Text)
		}
	}

	// 映射进通用 ChatResponse，保证上层 handler 契约不变。
	out := &ChatResponse{
		ID:    resp.ID,
		Model: string(resp.Model),
	}
	out.Choices = append(out.Choices, struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	}{})
	out.Choices[0].Message.Role = "assistant"
	out.Choices[0].Message.Content = sb.String()
	out.Choices[0].FinishReason = string(resp.StopReason)
	out.Usage.PromptTokens = int(resp.Usage.InputTokens)
	out.Usage.CompletionTokens = int(resp.Usage.OutputTokens)
	out.Usage.TotalTokens = int(resp.Usage.InputTokens + resp.Usage.OutputTokens)

	return out, nil
}
