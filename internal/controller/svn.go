package controller

import (
	"fmt"
	"git-svn-reviewbot/internal/config"
	"git-svn-reviewbot/internal/dto"
	"git-svn-reviewbot/internal/service"
	"git-svn-reviewbot/internal/storage"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SvnController SVN 控制器
type SvnController struct {
	svnService        *service.SvnService
	storage           *storage.FileStorage
	codeReviewService *service.CodeReviewService
	keyword           string
	baseURL           string
}

// NewSvnController 创建新的 SVN 控制器
func NewSvnController(svnService *service.SvnService, storage *storage.FileStorage, codeReviewService *service.CodeReviewService, cfg config.Config) *SvnController {
	return &SvnController{
		svnService:        svnService,
		storage:           storage,
		codeReviewService: codeReviewService,
		keyword:           cfg.SVN.Keyword,
		baseURL:           cfg.Server.BaseURL,
	}
}

// HandleWebhook 处理 SVN Webhook
func (c *SvnController) HandleWebhook(ctx *gin.Context) {
	var request dto.SvnWebhookRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体"})
		return
	}

	// 检查是否包含关键字
	if request.Message == "" || !containsKeyword(request.Message, c.keyword) {
		ctx.JSON(http.StatusOK, dto.SvnWebhookResponse{
			ReviewID:  request.Revision,
			ReviewURL: "",
			Message:   "Skip",
		})
		return
	}

	// 启动处理
	c.svnService.ProcessSvnReviewAsync(request)

	slog.Info("SVN 钩子已接收", "repo", request.Repo, "revision", request.Revision)

	// 使用完整的 repo 路径构建 URL（现在都通过查询参数传递）
	var reviewURL string
	if request.Repo == "" {
		reviewURL = fmt.Sprintf("%s/?version=%d", c.baseURL, request.Revision)
	} else {
		reviewURL = fmt.Sprintf("%s/?repo=%s&version=%d", c.baseURL, request.Repo, request.Revision)
	}

	ctx.JSON(http.StatusOK, dto.SvnWebhookResponse{
		ReviewID:  request.Revision,
		ReviewURL: reviewURL,
		Message:   "Review request accepted.",
	})
}

// ListReviews 列出所有审查记录
func (c *SvnController) ListReviews(ctx *gin.Context) {
	reviews, err := c.storage.ListReviews()
	if err != nil {
		slog.Error("获取审查列表失败", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取审查列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

// GetReviewMarkdown 获取审查结果 Markdown
func (c *SvnController) GetReviewMarkdown(ctx *gin.Context) {
	// repo 参数可选：
	// - 如果为空，认为 SVN_BASE_URL 已经是完整路径
	// - 如果不为空，拼接到 SVN_BASE_URL 后面
	repoPath := ctx.Query("repo")

	revisionStr := ctx.Query("revision")
	if revisionStr == "" {
		ctx.String(http.StatusBadRequest, "必须提供 revision 参数")
		return
	}

	revision, err := strconv.ParseInt(revisionStr, 10, 64)
	if err != nil {
		ctx.String(http.StatusBadRequest, "无效的版本号")
		return
	}

	// 如果 repoPath 为空，使用 baseURL 的文件名部分作为存储标识
	storageKey := repoPath
	if storageKey == "" {
		storageKey = "default"
	}

	finalFile := c.storage.GetReviewFile(storageKey, revision, ".md")
	processingFile := c.storage.GetReviewFile(storageKey, revision, ".processing")

	if c.storage.FileExists(finalFile) {
		content, err := c.storage.ReadFile(finalFile)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "读取文件失败")
			return
		}
		ctx.Header("Content-Type", "text/markdown")
		ctx.String(http.StatusOK, content)
		return
	}

	if c.storage.FileExists(processingFile) {
		ctx.Header("Content-Type", "text/markdown")
		ctx.String(http.StatusOK, "# 处理中...\n\nAI 正在分析代码，请稍后刷新。")
		return
	}

	ctx.String(http.StatusNotFound, "未找到")
}

// containsKeyword 检查消息是否包含关键字
func containsKeyword(message, keyword string) bool {
	return len(message) >= len(keyword) && message[0:len(keyword)] == keyword
}

// TriggerManualReview 手动触发 SVN 审查
func (c *SvnController) TriggerManualReview(ctx *gin.Context) {
	// repo 参数可选：
	// - 如果为空，认为 SVN_BASE_URL 已经是完整路径
	// - 如果不为空，拼接到 SVN_BASE_URL 后面
	repoPath := ctx.Query("repo")

	revisionStr := ctx.Query("revision")
	if revisionStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "必须提供 revision 参数"})
		return
	}

	// 获取查询参数
	author := ctx.DefaultQuery("author", "manual")
	message := ctx.DefaultQuery("message", "[review] Manual review trigger")

	revision, err := strconv.ParseInt(revisionStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的版本号"})
		return
	}

	// 构建请求
	request := dto.SvnWebhookRequest{
		Repo:     repoPath,
		Revision: revision,
		Author:   author,
		Message:  message,
	}

	// 检查是否有关键字 (手动触发可以不强制要求，但我们还是检查一下)
	if message == "" || !containsKeyword(message, c.keyword) {
		// 如果没有关键字，自动添加
		request.Message = fmt.Sprintf("%s %s", c.keyword, message)
	}

	// 启动处理
	c.svnService.ProcessSvnReviewAsync(request)

	slog.Info("手动触发 SVN 审查", "repo", repoPath, "revision", revision)

	// 构建查询结果的 URL - 指向首页，带 URL 参数
	var reviewURL string
	if repoPath == "" {
		reviewURL = fmt.Sprintf("%s/?version=%d", c.baseURL, revision)
	} else {
		reviewURL = fmt.Sprintf("%s/?repo=%s&version=%d", c.baseURL, repoPath, revision)
	}

	ctx.JSON(http.StatusOK, dto.SvnWebhookResponse{
		ReviewID:  revision,
		ReviewURL: reviewURL,
		Message:   "Manual review triggered successfully.",
	})
}

// PreCommitReview pre-commit 审查 API（同步调用）
func (c *SvnController) PreCommitReview(ctx *gin.Context) {
	var req dto.PreCommitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	slog.Info("收到 pre-commit 审查请求", "repo", req.Repo, "author", req.Author, "message", req.Message)

	if req.Message == "" || !containsKeyword(req.Message, c.keyword) {
		// 如果没有关键字，自动添加
		//request.Message = fmt.Sprintf("%s %s", c.keyword, message)
		//没有关键字，默认不处理
		ctx.JSON(http.StatusOK, dto.PreCommitResponse{
			Allowed: true,  // 审查出错时允许提交
			Score:   10,    // 默认满分
			Report:  "无需代码审查",
			Message: "无需代码审查，允许提交",
		})
		return
	}

	// 构建标题
	title := fmt.Sprintf("Pre-commit 审查 - %s - %s", req.Repo)
	if req.Author != "" {
		title += fmt.Sprintf(" - %s", req.Author)
	}

	// 调用代码审查（带评分）
	reviewResult, err := c.codeReviewService.RequestReviewWithScore(title, req.Diff)
	if err != nil {
		slog.Error("代码审查失败", "error", err)
		ctx.JSON(http.StatusOK, dto.PreCommitResponse{
			Allowed: true,  // 审查出错时允许提交
			Score:   10,    // 默认满分
			Report:  fmt.Sprintf("代码审查服务暂时不可用：%v", err),
			Message: "审查服务异常，允许提交",
		})
		return
	}

	// 判定是否允许提交：评分 <= 6 分不允许
	allowed := reviewResult.Score > 6

	// 保存审查报告（无论是否有 repo，都保存）
	var reviewURL string
	// 生成一个临时的版本号
	tempRevision := time.Now().Unix()
	var storageKey string
	if req.Repo != "" {
		storageKey = fmt.Sprintf("pre-commit-%s", req.Repo)
	} else {
		storageKey = "pre-commit"
	}
	filePath := c.storage.GetReviewFile(storageKey, tempRevision, ".md")
	if err := c.storage.WriteFile(filePath, reviewResult.Report); err == nil {
		// 构建查看 URL - 指向首页
		reviewURL = fmt.Sprintf("%s/?repo=%s&version=%d", c.baseURL, storageKey, tempRevision)
	}

	// 构建响应
	var message string
	if allowed {
		message = fmt.Sprintf("代码审查通过！评分：%d/10 分", reviewResult.Score)
	} else {
		message = fmt.Sprintf("代码审查未通过！评分：%d/10 分（要求大于6分）", reviewResult.Score)
	}

	ctx.JSON(http.StatusOK, dto.PreCommitResponse{
		Allowed: allowed,
		Score:   reviewResult.Score,
		Report:  reviewResult.Report,
		URL:     reviewURL,
		Message: message,
	})
}

// GetReviewDetail 获取审查报告详情
func (c *SvnController) GetReviewDetail(ctx *gin.Context) {
	// 从 URL 查询参数获取 repo 和 revision
	repoPath := ctx.Query("repo")
	revisionStr := ctx.Query("version")
	
	// 如果 repo 可选，如果没提供用默认
	if repoPath == "" {
		repoPath = "default"
	}
	
	// 解析版本号
	revision, err := strconv.ParseInt(revisionStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的版本号"})
		return
	}
	
	// 获取审查报告
	review, err := c.storage.GetReview(repoPath, revision)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "未找到审查报告"})
			return
		}
		slog.Error("获取审查报告失败", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取审查报告失败"})
		return
	}
	
	// 构建响应
	ctx.JSON(http.StatusOK, dto.ReviewDetailResponse{
		RepoPath:  review.RepoPath,
		Revision:  review.Revision,
		Report:    review.Content,
		Timestamp: review.Timestamp.Format("2006-01-02 15:04:05"),
		FileName:  review.FileName,
	})
}