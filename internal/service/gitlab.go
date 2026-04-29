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
	"strings"
)

// GitLabService GitLab 服务
type GitLabService struct {
	baseURL         string
	privateToken    string
	client          *http.Client
	codeReviewService *CodeReviewService
}

// NewGitLabService 创建新的 GitLab 服务
func NewGitLabService(cfg config.GitLabConfig, codeReviewService *CodeReviewService) *GitLabService {
	return &GitLabService{
		baseURL:         cfg.BaseURL,
		privateToken:    cfg.PrivateToken,
		client:          &http.Client{},
		codeReviewService: codeReviewService,
	}
}

// GetMergeRequestDiff 获取合并请求差异
func (s *GitLabService) GetMergeRequestDiff(projectID int64, mrIID int64) (string, error) {
	slog.Info("获取 GitLab MR 差异", "project_id", projectID, "mr_iid", mrIID)

	url := fmt.Sprintf("%s/api/v4/projects/%d/merge_requests/%d/diffs", s.baseURL, projectID, mrIID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Private-Token", s.privateToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitLab 返回错误状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var diffs []dto.GitLabDiffResponse
	if err := json.Unmarshal(body, &diffs); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	var diffParts []string
	for _, diff := range diffs {
		diffParts = append(diffParts, diff.Diff)
	}

	slog.Info("成功获取 GitLab MR 差异")
	return strings.Join(diffParts, "\n---\n"), nil
}

// AddMergeRequestComment 添加合并请求评论
func (s *GitLabService) AddMergeRequestComment(projectID int64, mrIID int64, comment string) error {
	slog.Info("添加 GitLab MR 评论", "project_id", projectID, "mr_iid", mrIID)

	reqBody := dto.GitLabCommentRequest{Body: comment}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	url := fmt.Sprintf("%s/api/v4/projects/%d/merge_requests/%d/notes", s.baseURL, projectID, mrIID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Private-Token", s.privateToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("GitLab 返回错误状态码: %d", resp.StatusCode)
	}

	slog.Info("成功添加 GitLab MR 评论")
	return nil
}

// StartCodeReview 启动代码审查
func (s *GitLabService) StartCodeReview(eventHeader string, event dto.GitLabWebhookMergeRequestEvent) {
	if eventHeader != "Merge Request Hook" {
		slog.Info("忽略非 MR 事件", "event_header", eventHeader)
		return
	}

	if event.ObjectAttributes.Action != "open" && event.ObjectAttributes.Action != "update" {
		slog.Info("忽略 MR 动作", "action", event.ObjectAttributes.Action)
		return
	}

	slog.Info("收到 MR 审查请求", "url", event.ObjectAttributes.URL)

	projectID := event.Project.ID
	mrIID := event.ObjectAttributes.IID
	mrTitle := event.ObjectAttributes.Title

	// 异步执行审查
	go func() {
		diff, err := s.GetMergeRequestDiff(projectID, mrIID)
		if err != nil {
			slog.Error("获取 MR 差异失败", "error", err)
			return
		}

		if diff == "" {
			slog.Info("没有找到差异内容", "title", mrTitle)
			return
		}

		slog.Info("差异获取成功，开始 AI 分析", "title", mrTitle)

		reviewComment, err := s.codeReviewService.RequestReview(mrTitle, diff)
		if err != nil {
			slog.Error("代码审查失败", "error", err)
			return
		}

		slog.Info("AI 分析完成，正在添加 GitLab 评论")

		if err := s.AddMergeRequestComment(projectID, mrIID, reviewComment); err != nil {
			slog.Error("添加评论失败", "error", err)
			return
		}

		slog.Info("审查完成", "title", mrTitle)
	}()
}
