package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"auth-center/service"
)

// AppController 应用控制器
type AppController struct{}

// CreateApp 创建应用
// @Summary 创建应用
// @Description 创建新的应用并获取应用凭据
// @Tags 应用管理
// @Accept json
// @Produce json
// @Param request body service.CreateAppRequest true "创建应用请求"
// @Success 201 {object} service.CreateAppResponse "创建成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 409 {object} map[string]string "应用名称已存在"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /apps [post]
func (c *AppController) CreateApp(ctx *gin.Context) {
	var req service.CreateAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appService := &service.AppService{}
	response, err := appService.CreateApp(&req)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// GetApp 获取应用信息
// @Summary 获取应用信息
// @Description 根据应用ID获取应用详细信息
// @Tags 应用管理
// @Produce json
// @Param app_id path string true "应用ID"
// @Success 200 {object} service.AppListResponse "应用信息"
// @Failure 404 {object} map[string]string "应用不存在"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /apps/{app_id} [get]
func (c *AppController) GetApp(ctx *gin.Context) {
	appID := ctx.Param("app_id")
	if appID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "应用ID不能为空"})
		return
	}

	appService := &service.AppService{}
	response, err := appService.GetApp(appID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateApp 更新应用
// @Summary 更新应用
// @Description 更新应用信息
// @Tags 应用管理
// @Accept json
// @Produce json
// @Param app_id path string true "应用ID"
// @Param request body service.UpdateAppRequest true "更新应用请求"
// @Success 200 {object} map[string]string "更新成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 404 {object} map[string]string "应用不存在"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /apps/{app_id} [put]
func (c *AppController) UpdateApp(ctx *gin.Context) {
	appID := ctx.Param("app_id")
	if appID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "应用ID不能为空"})
		return
	}

	var req service.UpdateAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appService := &service.AppService{}
	if err := appService.UpdateApp(appID, &req); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteApp 删除应用
// @Summary 删除应用
// @Description 删除应用（软删除）
// @Tags 应用管理
// @Param app_id path string true "应用ID"
// @Success 200 {object} map[string]string "删除成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 404 {object} map[string]string "应用不存在"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /apps/{app_id} [delete]
func (c *AppController) DeleteApp(ctx *gin.Context) {
	appID := ctx.Param("app_id")
	if appID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "应用ID不能为空"})
		return
	}

	appService := &service.AppService{}
	if err := appService.DeleteApp(appID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ListApps 获取应用列表
// @Summary 获取应用列表
// @Description 分页获取应用列表
// @Tags 应用管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{} "应用列表"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /apps [get]
func (c *AppController) ListApps(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	appService := &service.AppService{}
	apps, total, err := appService.ListApps(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": apps,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// RegenerateAppSecret 重新生成应用密钥
// @Summary 重新生成应用密钥
// @Description 重新生成应用密钥
// @Tags 应用管理
// @Produce json
// @Param app_id path string true "应用ID"
// @Success 200 {object} map[string]string "新密钥"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 404 {object} map[string]string "应用不存在"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /apps/{app_id}/regenerate-secret [post]
func (c *AppController) RegenerateAppSecret(ctx *gin.Context) {
	appID := ctx.Param("app_id")
	if appID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "应用ID不能为空"})
		return
	}

	appService := &service.AppService{}
	newSecret, err := appService.RegenerateAppSecret(appID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"app_secret": newSecret})
}
