package controller

import (
	"git-svn-reviewbot/internal/dto"
	"git-svn-reviewbot/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GitLabController GitLab 控制器
type GitLabController struct {
	gitLabService *service.GitLabService
}

// NewGitLabController 创建新的 GitLab 控制器
func NewGitLabController(gitLabService *service.GitLabService) *GitLabController {
	return &GitLabController{
		gitLabService: gitLabService,
	}
}

// HandleWebhook 处理 GitLab Webhook
func (c *GitLabController) HandleWebhook(ctx *gin.Context) {
	eventHeader := ctx.GetHeader("X-Gitlab-Event")

	var event dto.GitLabWebhookMergeRequestEvent
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体"})
		return
	}

	c.gitLabService.StartCodeReview(eventHeader, event)
	ctx.String(http.StatusOK, "Processing start")
}
