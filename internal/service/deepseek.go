package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git-svn-reviewbot/internal/config"
	"git-svn-reviewbot/internal/dto"
	"io"
	"log/slog"
	"net/http"
)

// DeepSeekService DeepSeek API 服务
type DeepSeekService struct {
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

// NewDeepSeekService 创建新的 DeepSeek 服务
func NewDeepSeekService(cfg config.DeepSeekConfig) *DeepSeekService {
	return &DeepSeekService{
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
		model:   cfg.Model,
		client:  &http.Client{},
	}
}

// GenerateResponse 生成响应
func (s *DeepSeekService) GenerateResponse(prompt string) (string, error) {
	slog.Info("调用 DeepSeek API", "model", s.model)

	// 构建请求
	messages := []dto.DeepSeekChatMessage{
		{
			Role:    "system",
			Content: "You are a helpful code review assistant.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	reqBody := dto.DeepSeekChatCompletionRequest{
		Model:    s.model,
		Messages: messages,
		Stream:   false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", s.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("DeepSeek API 返回错误状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var deepseekResp dto.DeepSeekChatCompletionResponse
	if err := json.Unmarshal(body, &deepseekResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(deepseekResp.Choices) == 0 {
		return "", fmt.Errorf("DeepSeek API 返回空响应")
	}

	slog.Info("DeepSeek API 响应生成成功", 
		"tokens_prompt", deepseekResp.Usage.PromptTokens,
		"tokens_completion", deepseekResp.Usage.CompletionTokens,
		"tokens_total", deepseekResp.Usage.TotalTokens)

	return deepseekResp.Choices[0].Message.Content, nil
}
