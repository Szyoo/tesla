package ai

import (
	"fmt"
	"tesla-server/config"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("AI API error (status %d): %s", e.StatusCode, e.Message)
}

// Provider 抽象云端 AI 提供方，便于在 OpenAI 兼容服务与 Claude 之间切换。
// 实现需返回填充好 Choices[0].Message.Content 与 Usage 的 ChatResponse，
// 以保证上层 handler 的契约不变。
type Provider interface {
	Chat(systemPrompt, userPrompt string) (*ChatResponse, error)
}

// Chat 按 AI_PROVIDER 选择提供方并转发。保持原有函数签名，handler 调用零改动。
func Chat(systemPrompt, userPrompt string) (*ChatResponse, error) {
	cfg := config.Load()
	switch cfg.AI.Provider {
	case "anthropic":
		return (&anthropicProvider{cfg: cfg}).Chat(systemPrompt, userPrompt)
	default: // openai-compat
		return (&openAICompatProvider{cfg: cfg}).Chat(systemPrompt, userPrompt)
	}
}
