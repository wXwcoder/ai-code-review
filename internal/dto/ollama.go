package dto

// OllamaGenerateRequest Ollama 生成请求
type OllamaGenerateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// OllamaGenerateResponse Ollama 生成响应
type OllamaGenerateResponse struct {
	Response string `json:"response"`
}
