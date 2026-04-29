package config

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

// Config 保存应用配置
type Config struct {
	Server   ServerConfig
	GitLab   GitLabConfig
	SVN      SVNConfig
	Ollama   OllamaConfig
	DeepSeek DeepSeekConfig
	LLM      LLMConfig
	Review   ReviewConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	BaseURL string
	Port    string
}

// GitLabConfig GitLab 配置
type GitLabConfig struct {
	BaseURL      string
	PrivateToken string
}

// SVNConfig SVN 配置
type SVNConfig struct {
	BaseURL     string
	Keyword     string
	SkipKeyword string
	Path        string
}

// OllamaConfig Ollama 配置
type OllamaConfig struct {
	BaseURL string
	Model   string
}

// DeepSeekConfig DeepSeek API 配置
type DeepSeekConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

// LLMConfig LLM 后端选择配置
type LLMConfig struct {
	Backend string // "ollama" 或 "deepseek"
}

// ReviewConfig 审查配置
type ReviewConfig struct {
	Prompt    string
	PassScore int
}

// Load 加载配置
func Load() *Config {
	viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("SERVER_BASE_URL", "http://localhost:8080")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("GITLAB_BASE_URL", "http://localhost:8929")
	viper.SetDefault("SVN_BASE_URL", "svn://localhost:3690")
	viper.SetDefault("SVN_KEYWORD", "[review]")
	viper.SetDefault("SKIP_KEYWORD", "[skip-review]")
	viper.SetDefault("SVN_PATH", "./reviews")
	viper.SetDefault("OLLAMA_BASE_URL", "http://localhost:11434")
	viper.SetDefault("OLLAMA_MODEL", "llama3")
	viper.SetDefault("DEEPSEEK_BASE_URL", "https://api.deepseek.com")
	viper.SetDefault("DEEPSEEK_MODEL", "deepseek-v4-flash")
	viper.SetDefault("LLM_BACKEND", "ollama") // 默认使用 Ollama
	viper.SetDefault("REVIEW_PROMPT", getDefaultPrompt())
	viper.SetDefault("REVIEW_PASS_SCORE", 6)

	// 确保 reviews 目录存在
	storagePath := viper.GetString("SVN_PATH")
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		slog.Warn("无法创建 reviews 目录", "error", err)
	}

	return &Config{
		Server: ServerConfig{
			BaseURL: viper.GetString("SERVER_BASE_URL"),
			Port:    viper.GetString("SERVER_PORT"),
		},
		GitLab: GitLabConfig{
			BaseURL:      viper.GetString("GITLAB_BASE_URL"),
			PrivateToken: viper.GetString("GITLAB_PRIVATE_TOKEN"),
		},
		SVN: SVNConfig{
			BaseURL:     viper.GetString("SVN_BASE_URL"),
			Keyword:     viper.GetString("SVN_KEYWORD"),
			SkipKeyword: viper.GetString("SKIP_KEYWORD"),
			Path:        viper.GetString("SVN_PATH"),
		},
		Ollama: OllamaConfig{
			BaseURL: viper.GetString("OLLAMA_BASE_URL"),
			Model:   viper.GetString("OLLAMA_MODEL"),
		},
		DeepSeek: DeepSeekConfig{
			BaseURL: viper.GetString("DEEPSEEK_BASE_URL"),
			APIKey:  viper.GetString("DEEPSEEK_API_KEY"),
			Model:   viper.GetString("DEEPSEEK_MODEL"),
		},
		LLM: LLMConfig{
			Backend: viper.GetString("LLM_BACKEND"),
		},
		Review: ReviewConfig{
			Prompt:    viper.GetString("REVIEW_PROMPT"),
			PassScore: viper.GetInt("REVIEW_PASS_SCORE"),
		},
	}
}

// getDefaultPrompt 获取默认审查提示词
func getDefaultPrompt() string {
	return `请扮演一位资深的代码审查专家，对以下代码变更进行全面审查。

审查要点：
1. 代码质量和可读性
2. 潜在的 bug 和逻辑问题
3. 性能问题
4. 安全性问题
5. 最佳实践遵循情况
6. 建议和改进意见

请在报告的开头，使用 JSON 格式返回评分，示例：
'''json
{
  "score": 8
}
'''
评分规则：
- 10 分：完美，没有任何问题
- 8-9 分：优秀，只有很小的改进空间
- 6-7 分：合格，仅有轻微问题，不影响功能
- 0-5 分：存在代码错误、逻辑问题或安全漏洞，必须修改
【重要】只要发现代码错误（包括但不限于：语法错误、逻辑错误、潜在 bug、安全漏洞），评分不得超过5分！

评分后再给出详细的 Markdown 格式审查报告。`
}
