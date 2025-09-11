package controllers

import (
	"net/http"
	"strconv"

	"auth-center/config"
	"auth-center/middleware"
	"auth-center/models"
	"auth-center/utils"

	"github.com/gin-gonic/gin"
)

// AppResourceController 应用内资源管理控制器
type AppResourceController struct{}

// getTargetAppID 获取目标应用ID
// 系统级超级管理员可以通过查询参数指定应用ID，应用级超级管理员只能使用当前应用ID
func (c *AppResourceController) getTargetAppID(ctx *gin.Context) string {
	return middleware.GetTargetAppID(ctx)
}

// GetSelfApp 获取当前应用级管理员所属应用的信息（仅应用级管理员可访问）
func (c *AppResourceController) GetSelfApp(ctx *gin.Context) {
	adminType, exists := ctx.Get("admin_type")
	if !exists || adminType != "app" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "仅应用级管理员可访问"})
		return
	}

	appID := c.getTargetAppID(ctx)

	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "应用不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":          app.ID,
			"name":        app.Name,
			"app_id":      app.AppID,
			"description": app.Description,
			"status":      app.Status,
		},
	})
}

// ListRoles 获取应用角色列表
func (c *AppResourceController) ListRoles(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	offset := (page - 1) * size

	var roles []models.Role
	var total int64

	// 获取角色列表
	if err := config.DB.Where("app_id = ?", appID).Offset(offset).Limit(size).Find(&roles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取角色列表失败"})
		return
	}

	// 获取总数
	config.DB.Model(&models.Role{}).Where("app_id = ?", appID).Count(&total)

	ctx.JSON(http.StatusOK, gin.H{
		"data": roles,
		"pagination": gin.H{
			"page":      page,
			"page_size": size,
			"total":     total,
		},
	})
}

// CreateRole 创建角色
func (c *AppResourceController) CreateRole(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)

	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查角色编码是否已存在
	var existingRole models.Role
	if err := config.DB.Where("app_id = ? AND code = ?", appID, req.Code).First(&existingRole).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "角色编码已存在"})
		return
	}

	// 创建角色
	role := models.Role{
		AppID:       appID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := config.DB.Create(&role).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建角色失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": role})
}

// UpdateRole 更新角色
func (c *AppResourceController) UpdateRole(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)
	roleID := ctx.Param("id")

	var req struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var role models.Role
	if err := config.DB.Where("id = ? AND app_id = ?", roleID, appID).First(&role).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	// 更新角色信息
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status >= 0 {
		updates["status"] = req.Status
	}

	if err := config.DB.Model(&role).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新角色失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": role})
}

// DeleteRole 删除角色
func (c *AppResourceController) DeleteRole(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)
	roleID := ctx.Param("id")

	// 检查角色是否存在
	var role models.Role
	if err := config.DB.Where("id = ? AND app_id = ?", roleID, appID).First(&role).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	// 删除角色（级联删除相关数据）
	if err := config.DB.Where("id = ? AND app_id = ?", roleID, appID).Delete(&models.Role{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除角色失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "角色删除成功"})
}

// ListPermissions 获取应用权限列表
func (c *AppResourceController) ListPermissions(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	offset := (page - 1) * size

	var permissions []models.Permission
	var total int64

	// 获取权限列表
	if err := config.DB.Where("app_id = ?", appID).Offset(offset).Limit(size).Find(&permissions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取权限列表失败"})
		return
	}

	// 获取总数
	config.DB.Model(&models.Permission{}).Where("app_id = ?", appID).Count(&total)

	ctx.JSON(http.StatusOK, gin.H{
		"data": permissions,
		"pagination": gin.H{
			"page":      page,
			"page_size": size,
			"total":     total,
		},
	})
}

// CreatePermission 创建权限
func (c *AppResourceController) CreatePermission(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)

	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Resource    string `json:"resource" binding:"required"`
		Action      string `json:"action" binding:"required"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查权限编码是否已存在
	var existingPermission models.Permission
	if err := config.DB.Where("app_id = ? AND code = ?", appID, req.Code).First(&existingPermission).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "权限编码已存在"})
		return
	}

	// 创建权限
	permission := models.Permission{
		AppID:       appID,
		Name:        req.Name,
		Code:        req.Code,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := config.DB.Create(&permission).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建权限失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": permission})
}

// UpdatePermission 更新权限
func (c *AppResourceController) UpdatePermission(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)
	permissionID := ctx.Param("id")

	var req struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Resource    string `json:"resource"`
		Action      string `json:"action"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var permission models.Permission
	if err := config.DB.Where("id = ? AND app_id = ?", permissionID, appID).First(&permission).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "权限不存在"})
		return
	}

	// 更新权限信息
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.Resource != "" {
		updates["resource"] = req.Resource
	}
	if req.Action != "" {
		updates["action"] = req.Action
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status >= 0 {
		updates["status"] = req.Status
	}

	if err := config.DB.Model(&permission).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新权限失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": permission})
}

// DeletePermission 删除权限
func (c *AppResourceController) DeletePermission(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)
	permissionID := ctx.Param("id")

	// 检查权限是否存在
	var permission models.Permission
	if err := config.DB.Where("id = ? AND app_id = ?", permissionID, appID).First(&permission).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "权限不存在"})
		return
	}

	// 删除权限（级联删除相关数据）
	if err := config.DB.Where("id = ? AND app_id = ?", permissionID, appID).Delete(&models.Permission{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除权限失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "权限删除成功"})
}

// AssignRolePermissions 为角色分配权限
func (c *AppResourceController) AssignRolePermissions(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)
	roleID := ctx.Param("id")

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查角色是否存在
	var role models.Role
	if err := config.DB.Where("id = ? AND app_id = ?", roleID, appID).First(&role).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	// 删除现有权限分配
	config.DB.Where("role_id = ? AND app_id = ?", roleID, appID).Delete(&models.RolePermission{})

	// 添加新的权限分配
	for _, permissionID := range req.PermissionIDs {
		rolePermission := models.RolePermission{
			RoleID:       role.ID,
			PermissionID: permissionID,
			AppID:        appID,
		}
		config.DB.Create(&rolePermission)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "权限分配成功"})
}

// GetRolePermissions 获取角色权限
func (c *AppResourceController) GetRolePermissions(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)
	roleID := ctx.Param("id")

	var rolePermissions []models.RolePermission
	if err := config.DB.Where("role_id = ? AND app_id = ?", roleID, appID).Find(&rolePermissions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取角色权限失败"})
		return
	}

	var permissionIDs []uint
	for _, rp := range rolePermissions {
		permissionIDs = append(permissionIDs, rp.PermissionID)
	}

	ctx.JSON(http.StatusOK, gin.H{"permission_ids": permissionIDs})
}

// ListUsers 获取应用用户列表
func (c *AppResourceController) ListUsers(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)

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

// CreateUser 创建用户
func (c *AppResourceController) CreateUser(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)

	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password" binding:"required,min=6"`
		Status   int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.DB.Where("app_id = ? AND username = ?", appID, req.Username).First(&existingUser).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在（如果提供了邮箱）
	if req.Email != "" {
		if err := config.DB.Where("app_id = ? AND email = ?", appID, req.Email).First(&existingUser).Error; err == nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": "邮箱已存在"})
			return
		}
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建用户
	user := models.User{
		AppID:    appID,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Status:   req.Status,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": user})
}

// UpdateUser 更新用户
func (c *AppResourceController) UpdateUser(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)
	userID := ctx.Param("id")

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Status   int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("id = ? AND app_id = ?", userID, appID).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新用户信息
	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Password != "" {
		// 哈希新密码
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}
		updates["password"] = hashedPassword
	}
	if req.Status >= 0 {
		updates["status"] = req.Status
	}

	if err := config.DB.Model(&user).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteUser 删除用户
func (c *AppResourceController) DeleteUser(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)
	userID := ctx.Param("id")

	// 检查用户是否存在
	var user models.User
	if err := config.DB.Where("id = ? AND app_id = ?", userID, appID).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 不能删除超级管理员
	if user.IsSuperAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "不能删除超级管理员"})
		return
	}

	// 删除用户（级联删除相关数据）
	if err := config.DB.Where("id = ? AND app_id = ?", userID, appID).Delete(&models.User{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// AssignUserRoles 为用户分配角色
func (c *AppResourceController) AssignUserRoles(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)
	userID := ctx.Param("id")

	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否存在
	var user models.User
	if err := config.DB.Where("id = ? AND app_id = ?", userID, appID).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 删除现有角色分配
	config.DB.Where("user_id = ? AND app_id = ?", userID, appID).Delete(&models.UserRole{})

	// 添加新的角色分配
	for _, roleID := range req.RoleIDs {
		userRole := models.UserRole{
			UserID: user.ID,
			RoleID: roleID,
			AppID:  appID,
		}
		config.DB.Create(&userRole)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "角色分配成功"})
}

// GetUserRoles 获取用户角色
func (c *AppResourceController) GetUserRoles(ctx *gin.Context) {
	appID := c.getTargetAppID(ctx)
	userID := ctx.Param("id")

	var userRoles []models.UserRole
	if err := config.DB.Where("user_id = ? AND app_id = ?", userID, appID).Find(&userRoles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户角色失败"})
		return
	}

	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	ctx.JSON(http.StatusOK, gin.H{"role_ids": roleIDs})
}
