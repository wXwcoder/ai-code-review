package main

import (
	"embed"
	"git-svn-reviewbot/internal/config"
	"git-svn-reviewbot/internal/controller"
	"git-svn-reviewbot/internal/service"
	"git-svn-reviewbot/internal/storage"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var embeddedFS embed.FS

func main() {
	// 配置日志
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("启动代码审查机器人")

	// 加载配置
	cfg := config.Load()

	// 初始化存储
	fileStorage := storage.NewFileStorage(cfg.SVN.Path)

	// 初始化 LLM 服务
	var llmService service.LLMService
	switch cfg.LLM.Backend {
	case "deepseek":
		slog.Info("使用 DeepSeek API", "model", cfg.DeepSeek.Model)
		llmService = service.NewDeepSeekService(cfg.DeepSeek)
	case "ollama":
		fallthrough
	default:
		slog.Info("使用 Ollama", "model", cfg.Ollama.Model)
		llmService = service.NewOllamaService(cfg.Ollama)
	}

	// 初始化服务
	codeReviewService := service.NewCodeReviewService(llmService, cfg.Review)
	gitLabService := service.NewGitLabService(cfg.GitLab, codeReviewService)
	svnService := service.NewSvnService(cfg.SVN, fileStorage, codeReviewService)

	// 初始化控制器
	gitLabController := controller.NewGitLabController(gitLabService)
	svnController := controller.NewSvnController(svnService, fileStorage, codeReviewService, *cfg)
	healthController := controller.NewHealthController()

	// 设置路由
	r := gin.Default()

	// 健康检查
	r.GET("/health", healthController.Check)

	// GitLab Webhook
	r.POST("/webhook/gitlab", gitLabController.HandleWebhook)

	// SVN Webhook
	r.POST("/webhook/svn", svnController.HandleWebhook)

	// 手动触发 SVN 代码审查 - 通过查询参数传递路径
	r.POST("/api/reviews/svn/trigger", svnController.TriggerManualReview)
	r.GET("/api/reviews/svn/trigger", svnController.TriggerManualReview)

	// SVN 审查结果查询 - 通过查询参数传递路径
	r.GET("/api/reviews/svn", svnController.GetReviewMarkdown)

	// 列出所有审查记录
	r.GET("/api/reviews/list", svnController.ListReviews)

	// Pre-commit 审查 API（同步调用）
	r.POST("/api/reviews/svn/pre-commit", svnController.PreCommitReview)

	// 获取审查报告详情 API
	r.GET("/api/review/detail", svnController.GetReviewDetail)

	// 从嵌入的文件系统中获取dist子目录
	fsys, err := fs.Sub(embeddedFS, "dist")
	if err != nil {
		log.Fatal("无法访问嵌入的dist目录:", err)
	}

	// 设置静态文件服务
	fileServer := http.FileServer(http.FS(fsys))

	// 使用 NoRoute 处理静态文件
	r.NoRoute(func(c *gin.Context) {
		// 先尝试提供静态文件
		path := c.Request.URL.Path
		if path != "/" && path != "" {
			if _, err := fs.Stat(fsys, strings.TrimLeft(path, "/")); err == nil {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		// 如果找不到静态文件，返回 index.html 让前端路由处理
		indexData, err := fs.ReadFile(fsys, "index.html")
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}

		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(indexData)
	})

	// 启动服务器
	addr := ":" + cfg.Server.Port
	slog.Info("服务器启动", "address", addr)
	if err := r.Run(addr); err != nil {
		slog.Error("服务器启动失败", "error", err)
		os.Exit(1)
	}
}
