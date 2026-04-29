package service

import (
	"bytes"
	"fmt"
	"git-svn-reviewbot/internal/config"
	"git-svn-reviewbot/internal/dto"
	"git-svn-reviewbot/internal/storage"
	"log/slog"
	"os/exec"
)

// SvnService SVN 服务
type SvnService struct {
	baseURL           string
	storage           *storage.FileStorage
	codeReviewService *CodeReviewService
}

// NewSvnService 创建新的 SVN 服务
func NewSvnService(cfg config.SVNConfig, storage *storage.FileStorage, codeReviewService *CodeReviewService) *SvnService {
	return &SvnService{
		baseURL:           cfg.BaseURL,
		storage:           storage,
		codeReviewService: codeReviewService,
	}
}

// GetBaseURL 获取 SVN 基础 URL
func (s *SvnService) GetBaseURL() string {
	return s.baseURL
}

// ProcessSvnReviewAsync 异步处理 SVN 审查
func (s *SvnService) ProcessSvnReviewAsync(request dto.SvnWebhookRequest) {
	repoPath := request.Repo
	title := fmt.Sprintf("SVN [%s] r%d by %s", repoPath, request.Revision, request.Author)

	slog.Info("开始 SVN 审查", "title", title)

	processingFile := s.storage.GetReviewFile(repoPath, request.Revision, ".processing")
	finalFile := s.storage.GetReviewFile(repoPath, request.Revision, ".md")

	// 创建处理中文件
	processingContent := fmt.Sprintf("# AI 审查生成中...\n\n请稍候。（仓库：%s，版本：%d）", repoPath, request.Revision)
	if err := s.storage.WriteFile(processingFile, processingContent); err != nil {
		slog.Error("创建处理中文件失败", "error", err)
	}

	// 异步执行审查
	go func() {
		defer func() {
			// 清理处理中文件
			if err := s.storage.DeleteFile(processingFile); err != nil {
				slog.Warn("删除处理中文件失败", "error", err)
			}
		}()

		// 构建完整 SVN URL
		// 如果 repoPath 为空，直接使用 baseURL
		var fullSvnURL string
		if repoPath == "" {
			fullSvnURL = s.baseURL
		} else {
			fullSvnURL = fmt.Sprintf("%s/%s", s.baseURL, repoPath)
		}
		diff, err := s.FetchSvnDiff(fullSvnURL, request.Revision)
		if err != nil {
			slog.Error("获取 SVN 差异失败", "error", err)
			errorContent := fmt.Sprintf("# 错误\n\n审查生成过程中出错：\n\n`%s`", err.Error())
			s.storage.WriteFile(finalFile, errorContent)
			return
		}

		if diff == "" {
			slog.Info("没有差异内容")
			s.storage.WriteFile(finalFile, "# 审查失败\n\n没有可审查的文本变更。")
			return
		}

		reviewComment, err := s.codeReviewService.RequestReview(title, diff)
		if err != nil {
			slog.Error("代码审查失败", "error", err)
			errorContent := fmt.Sprintf("# 错误\n\nAI 响应失败：\n\n`%s`", err.Error())
			s.storage.WriteFile(finalFile, errorContent)
			return
		}

		if err := s.storage.WriteFile(finalFile, reviewComment); err != nil {
			slog.Error("保存审查结果失败", "error", err)
			return
		}

		slog.Info("SVN 审查完成", "file", finalFile)
	}()
}

// FetchSvnDiff 获取 SVN 差异
func (s *SvnService) FetchSvnDiff(url string, revision int64) (string, error) {
	slog.Info("获取 SVN 差异", "url", url, "revision", revision)

	prevRevision := revision - 1
	cmd := exec.Command("svn", "diff", "-r", fmt.Sprintf("%d:%d", prevRevision, revision), url)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("执行 svn 命令失败: %w, stderr: %s", err, stderr.String())
	}

	slog.Info("SVN 差异获取成功")
	return stdout.String(), nil
}
