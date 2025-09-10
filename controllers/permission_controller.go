package controllers

import (
	"net/http"

	"auth-center/service"

	"github.com/gin-gonic/gin"
)

// PermissionController 权限控制器
type PermissionController struct{}

// CheckPermission 检查权限
// @Summary 检查权限
// @Description 检查用户是否具有指定权限
// @Tags 权限管理
// @Produce json
// @Security BearerAuth
// @Param permission query string true "权限代码"
// @Success 200 {object} map[string]bool "权限检查结果"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "未认证"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /permissions/check [get]
func (c *PermissionController) CheckPermission(ctx *gin.Context) {
	permission := ctx.Query("permission")
	if permission == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "权限代码不能为空"})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	appID, exists := ctx.Get("app_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "应用未认证"})
		return
	}

	permissionService := &service.PermissionService{}
	hasPermission, err := permissionService.CheckUserPermission(userID.(uint), appID.(string), permission)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"has_permission": hasPermission})
}

// CheckAPIPermission 检查API权限
// @Summary 检查API权限
// @Description 检查用户是否有权限访问指定API
// @Tags 权限管理
// @Produce json
// @Security BearerAuth
// @Param path query string true "API路径"
// @Param method query string true "HTTP方法"
// @Success 200 {object} map[string]bool "权限检查结果"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "未认证"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /permissions/check-api [get]
func (c *PermissionController) CheckAPIPermission(ctx *gin.Context) {
	path := ctx.Query("path")
	method := ctx.Query("method")

	if path == "" || method == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "API路径和方法不能为空"})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	appID, exists := ctx.Get("app_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "应用未认证"})
		return
	}

	permissionService := &service.PermissionService{}
	hasPermission, err := permissionService.CheckAPIPermission(userID.(uint), appID.(string), path, method)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"has_permission": hasPermission})
}

// GetUserPermissions 获取用户权限列表
// @Summary 获取用户权限列表
// @Description 获取当前用户的所有权限
// @Tags 权限管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "用户权限列表"
// @Failure 401 {object} map[string]string "未认证"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /permissions/user [get]
func (c *PermissionController) GetUserPermissions(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	appID, exists := ctx.Get("app_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "应用未认证"})
		return
	}

	permissionService := &service.PermissionService{}
	permissions, err := permissionService.GetUserPermissionsFromDB(userID.(uint), appID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// GetUserRoles 获取用户角色列表
// @Summary 获取用户角色列表
// @Description 获取当前用户的所有角色
// @Tags 权限管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "用户角色列表"
// @Failure 401 {object} map[string]string "未认证"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /permissions/roles [get]
func (c *PermissionController) GetUserRoles(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	_, exists = ctx.Get("app_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "应用未认证"})
		return
	}

	authService := &service.AuthService{}
	roles, err := authService.GetUserRoles(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roleInfos, err := authService.GetRoleInfos(roles)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"roles": roleInfos})
}
