package controllers

import (
	"net/http"
	"strconv"

	"auth-center/config"
	"auth-center/middleware"
	"auth-center/models"
	"auth-center/service"

	"github.com/gin-gonic/gin"
)

// AppManagementController 应用管理控制器
type AppManagementController struct{}

// ListApps 获取应用列表
// 系统级超级管理员：可以看到所有应用
// 应用级超级管理员：只能看到自己应用的信息
func (c *AppManagementController) ListApps(ctx *gin.Context) {
	// 检查是否为系统级超级管理员
	if middleware.IsSystemAdmin(ctx) {
		// 系统级超级管理员可以查看所有应用
		var apps []models.Application
		if err := config.DB.Find(&apps).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取应用列表失败"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": apps})
		return
	}

	// 应用级超级管理员只能查看自己应用的信息
	appID, _ := middleware.GetAppID(ctx)
	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "应用不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": []models.Application{app}})
}

// GetApp 获取应用详情
func (c *AppManagementController) GetApp(ctx *gin.Context) {
	appID := ctx.Param("app_id")

	// 检查权限
	if !middleware.IsSystemAdmin(ctx) {
		// 应用级超级管理员只能查看自己应用
		currentAppID, _ := middleware.GetAppID(ctx)
		if appID != currentAppID {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "无权查看其他应用信息"})
			return
		}
	}

	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "应用不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": app})
}

// CreateApp 创建应用（仅系统级超级管理员）
func (c *AppManagementController) CreateApp(ctx *gin.Context) {
	var req service.CreateAppRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用服务层创建应用
	appService := &service.AppService{}
	response, err := appService.CreateApp(&req)
	if err != nil {
		if err.Error() == "应用ID已存在" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "应用创建成功",
		"data":    response,
	})
}

// UpdateApp 更新应用（仅系统级超级管理员）
func (c *AppManagementController) UpdateApp(ctx *gin.Context) {
	appID := ctx.Param("app_id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "应用不存在"})
		return
	}

	// 更新应用信息
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status >= 0 {
		updates["status"] = req.Status
	}

	if err := config.DB.Model(&app).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新应用失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": app})
}

// DeleteApp 删除应用（仅系统级超级管理员）
func (c *AppManagementController) DeleteApp(ctx *gin.Context) {
	appID := ctx.Param("app_id")

	// 检查应用是否存在
	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "应用不存在"})
		return
	}

	// 不能删除系统管理应用
	if appID == "system-admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "不能删除系统管理应用"})
		return
	}

	// 删除应用（级联删除相关数据）
	if err := config.DB.Where("app_id = ?", appID).Delete(&models.Application{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除应用失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "应用删除成功"})
}

// RegenerateAppSecret 重新生成应用密钥（仅系统级超级管理员）
func (c *AppManagementController) RegenerateAppSecret(ctx *gin.Context) {
	appID := ctx.Param("app_id")

	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "应用不存在"})
		return
	}

	// 生成新的应用密钥
	newSecret, err := service.GenerateAppSecret()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成应用密钥失败"})
		return
	}

	// 更新应用密钥
	if err := config.DB.Model(&app).Update("app_secret", newSecret).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新应用密钥失败"})
		return
	}

	app.AppSecret = newSecret
	ctx.JSON(http.StatusOK, gin.H{"data": app})
}

// ListAppUsers 获取应用用户列表
// 系统级超级管理员：可以查看所有应用的用户
// 应用级超级管理员：只能查看自己应用的用户
func (c *AppManagementController) ListAppUsers(ctx *gin.Context) {
	appID := ctx.Param("app_id")

	// 检查权限
	if !middleware.IsSystemAdmin(ctx) {
		// 应用级超级管理员只能查看自己应用的用户
		currentAppID, _ := middleware.GetAppID(ctx)
		if appID != currentAppID {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "无权查看其他应用用户"})
			return
		}
	}

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	offset := (page - 1) * size

	var users []models.User
	var total int64

	// 获取用户列表
	if err := config.DB.Where("app_id = ?", appID).Offset(offset).Limit(size).Find(&users).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	// 获取总数
	config.DB.Model(&models.User{}).Where("app_id = ?", appID).Count(&total)

	ctx.JSON(http.StatusOK, gin.H{
		"data": users,
		"pagination": gin.H{
			"page":      page,
			"page_size": size,
			"total":     total,
		},
	})
}
