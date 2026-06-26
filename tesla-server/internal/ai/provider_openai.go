package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tesla-server/config"
	"time"
)

// openAICompatProvider 通过 OpenAI 兼容的 /chat/completions 接口调用，
// 覆盖智谱 GLM、OpenAI、DeepSeek、Moonshot 等服务。
type openAICompatProvider struct {
	cfg *config.Config
}

func (p *openAICompatProvider) Chat(systemPrompt, userPrompt string) (*ChatResponse, error) {
	cfg := p.cfg

	reqBody := ChatRequest{
		Model: cfg.AI.Model,
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.7,
		MaxTokens:   4096,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := cfg.AI.BaseURL + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.AI.APIKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.Unmarshal(body, &errResp)
		errMsg := string(body)
		if e, ok := errResp["error"].(map[string]interface{}); ok {
			if msg, ok := e["message"].(string); ok {
				errMsg = msg
			}
		}
		return nil, &APIError{StatusCode: resp.StatusCode, Message: errMsg}
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &chatResp, nil
}
