package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthController 健康检查控制器
type HealthController struct{}

// NewHealthController 创建新的健康检查控制器
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Check 健康检查
func (c *HealthController) Check(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
