package dto

// DeepSeekChatMessage 聊天消息
type DeepSeekChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepSeekChatCompletionRequest DeepSeek 聊天完成请求 (OpenAI 兼容格式)
type DeepSeekChatCompletionRequest struct {
	Model       string                  `json:"model"`
	Messages    []DeepSeekChatMessage   `json:"messages"`
	Stream      bool                    `json:"stream"`
	Thinking    *DeepSeekThinkingConfig `json:"thinking,omitempty"`
}

// DeepSeekThinkingConfig 思考配置
type DeepSeekThinkingConfig struct {
	Type string `json:"type"` // "enabled" 或 "disabled"
}

// DeepSeekChatCompletionResponse DeepSeek 聊天完成响应
type DeepSeekChatCompletionResponse struct {
	ID      string                    `json:"id"`
	Object  string                    `json:"object"`
	Created int64                     `json:"created"`
	Model   string                    `json:"model"`
	Choices []DeepSeekChatChoice      `json:"choices"`
	Usage   DeepSeekUsageInfo          `json:"usage"`
}

// DeepSeekChatChoice 聊天选择项
type DeepSeekChatChoice struct {
	Index        int                     `json:"index"`
	Message      DeepSeekChatMessage    `json:"message"`
	FinishReason string                 `json:"finish_reason"`
}

// DeepSeekUsageInfo 使用信息
type DeepSeekUsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
