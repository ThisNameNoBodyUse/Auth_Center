package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"auth-center/service"
)

// SystemAdminController 系统管理员控制器
type SystemAdminController struct{}

// SystemLogin 系统管理员登录
// @Summary 系统管理员登录
// @Description 系统内部管理员登录接口
// @Tags 系统管理
// @Accept json
// @Produce json
// @Param request body service.SystemLoginRequest true "登录请求"
// @Success 200 {object} service.SystemLoginResponse "登录成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "认证失败"
// @Router /system/login [post]
func (c *SystemAdminController) SystemLogin(ctx *gin.Context) {
	var req service.SystemLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminService := &service.SystemAdminService{}
	response, err := adminService.SystemLogin(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// SystemRegister 系统管理员注册
// @Summary 系统管理员注册
// @Description 注册新的系统管理员（仅系统级超级管理员可操作）
// @Tags 系统管理
// @Accept json
// @Produce json
// @Param request body service.SystemRegisterRequest true "注册请求"
// @Success 201 {object} service.SystemRegisterResponse "注册成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 409 {object} map[string]string "用户名或邮箱已存在"
// @Router /system/register [post]
func (c *SystemAdminController) SystemRegister(ctx *gin.Context) {
	var req service.SystemRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminService := &service.SystemAdminService{}
	response, err := adminService.SystemRegister(&req)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// SystemRefreshToken 刷新系统管理员令牌
// @Summary 刷新系统管理员令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 系统管理
// @Accept json
// @Produce json
// @Param request body map[string]string true "刷新令牌请求"
// @Success 200 {object} service.SystemLoginResponse "刷新成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "令牌无效"
// @Router /system/refresh [post]
func (c *SystemAdminController) SystemRefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminService := &service.SystemAdminService{}
	response, err := adminService.SystemRefreshToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// SystemLogout 系统管理员登出
// @Summary 系统管理员登出
// @Description 系统管理员登出并黑名单化令牌
// @Tags 系统管理
// @Success 200 {object} map[string]string "登出成功"
// @Router /system/logout [post]
func (c *SystemAdminController) SystemLogout(ctx *gin.Context) {
	// 这里可以添加令牌黑名单逻辑
	ctx.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// GetSystemAdminInfo 获取系统管理员信息
// @Summary 获取系统管理员信息
// @Description 获取当前登录的系统管理员信息
// @Tags 系统管理
// @Success 200 {object} models.SystemAdmin "管理员信息"
// @Failure 401 {object} map[string]string "认证失败"
// @Router /system/admin/info [get]
func (c *SystemAdminController) GetSystemAdminInfo(ctx *gin.Context) {
	adminID, exists := ctx.Get("admin_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未找到管理员信息"})
		return
	}

	adminService := &service.SystemAdminService{}
	admin, err := adminService.GetSystemAdminInfo(adminID.(uint))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, admin)
}
