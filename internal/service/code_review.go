package service

import (
	"encoding/json"
	"fmt"
	"git-svn-reviewbot/internal/config"
	"log/slog"
	"regexp"
)

// LLMService LLM 服务接口
type LLMService interface {
	GenerateResponse(prompt string) (string, error)
}

// ReviewResult 审查结果
type ReviewResult struct {
	Score  int    // 评分
	Report string // 报告内容
}

// CodeReviewService 代码审查服务
type CodeReviewService struct {
	llmService LLMService
	prompt     string
}

// NewCodeReviewService 创建新的代码审查服务
func NewCodeReviewService(llmService LLMService, cfg config.ReviewConfig) *CodeReviewService {
	return &CodeReviewService{
		llmService: llmService,
		prompt:     cfg.Prompt,
	}
}

// RequestReview 发起审查请求（仅返回报告，不解析评分）
func (s *CodeReviewService) RequestReview(title string, diffContent string) (string, error) {
	slog.Info("发起代码审查", "title", title)

	if diffContent == "" {
		return "审查失败：没有差异内容", nil
	}

	reviewPrompt := fmt.Sprintf("%s\n\n---\n标题：%s\n---\n代码变更内容（Diff）：\n%s\n---",
		s.prompt, title, diffContent)

	response, err := s.llmService.GenerateResponse(reviewPrompt)
	if err != nil {
		return "", fmt.Errorf("调用 LLM 失败: %w", err)
	}

	slog.Info("代码审查完成")
	return response, nil
}

// RequestReviewWithScore 发起审查请求并解析评分
func (s *CodeReviewService) RequestReviewWithScore(title string, diffContent string) (*ReviewResult, error) {
	slog.Info("发起带评分的代码审查", "title", title)

	if diffContent == "" {
		return &ReviewResult{Score: 10, Report: "没有变更内容，默认通过"}, nil
	}

	reviewPrompt := fmt.Sprintf("%s\n\n---\n标题：%s\n---\n代码变更内容（Diff）：\n%s\n---",
		s.prompt, title, diffContent)

	response, err := s.llmService.GenerateResponse(reviewPrompt)
	if err != nil {
		return nil, fmt.Errorf("调用 LLM 失败: %w", err)
	}

	// 解析评分
	score := s.extractScore(response)
	
	// 清理 JSON 标记，保留纯粹的报告内容
	cleanReport := s.cleanReport(response)
	
	slog.Info("带评分的代码审查完成", "score", score)
	return &ReviewResult{
		Score:  score,
		Report: cleanReport,
	}, nil
}

// extractScore 从响应中提取评分
func (s *CodeReviewService) extractScore(response string) int {
	// 尝试解析 JSON 评分 - 支持多种标记
	patterns := []string{
		`(?s)'''json.*?({.*?"score".*?})'''`,
		`(?s)'''.*?({.*?"score".*?})'''`,
		`(?s)'''.*?({.*?score.*?})'''`,
	}
	
	for _, pattern := range patterns {
		jsonRegex := regexp.MustCompile(pattern)
		matches := jsonRegex.FindStringSubmatch(response)
		
		if len(matches) >= 2 {
			var result struct {
				Score int `json:"score"`
			}
			err := json.Unmarshal([]byte(matches[1]), &result)
			if err == nil && result.Score >= 0 && result.Score <= 10 {
				return result.Score
			}
		}
	}
	
	// 如果没有找到 JSON，尝试简单的正则匹配
	scoreRegex := regexp.MustCompile(`(?:评分|Score|score).*?(?:[:：]|是|为)?\s*(\d+)`)
	scoreMatch := scoreRegex.FindStringSubmatch(response)
	
	if len(scoreMatch) >= 2 {
		var score int
		_, err := fmt.Sscanf(scoreMatch[1], "%d", &score)
		if err == nil && score >= 0 && score <= 10 {
			return score
		}
	}
	
	// 默认评分
	return 7
}

// cleanReport 清理报告内容，移除 JSON 标记
func (s *CodeReviewService) cleanReport(response string) string {
	// 移除开头的 JSON 代码块
	patterns := []string{
		`(?s)'''json.*?{.*?"score".*?}'''`,
		`(?s)'''.*?{.*?"score".*?}'''`,
		`(?s)'''json.*?{.*?score.*?}'''`,
		`(?s)'''.*?{.*?score.*?}'''`,
	}
	
	cleaned := response
	for _, pattern := range patterns {
		jsonRegex := regexp.MustCompile(pattern)
		cleaned = jsonRegex.ReplaceAllString(cleaned, "")
	}
	
	return cleaned
}
