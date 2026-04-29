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

// OllamaService Ollama 服务
type OllamaService struct {
	baseURL string
	model   string
	client  *http.Client
}

// NewOllamaService 创建新的 Ollama 服务
func NewOllamaService(cfg config.OllamaConfig) *OllamaService {
	return &OllamaService{
		baseURL: cfg.BaseURL,
		model:   cfg.Model,
		client:  &http.Client{},
	}
}

// GenerateResponse 生成响应
func (s *OllamaService) GenerateResponse(prompt string) (string, error) {
	slog.Info("调用 Ollama 生成响应")

	reqBody := dto.OllamaGenerateRequest{
		Model:  s.model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", s.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama 返回错误状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var ollamaResp dto.OllamaGenerateResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	slog.Info("Ollama 响应生成成功")
	return ollamaResp.Response, nil
}
