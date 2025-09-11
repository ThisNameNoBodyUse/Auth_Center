package controllers

import (
	"auth-center/config"
	"auth-center/models"
	"auth-center/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SystemAdminManagementController struct{}

// ListSystemAdmins 获取系统管理员列表
func (c *SystemAdminManagementController) ListSystemAdmins(ctx *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	adminType := ctx.Query("admin_type")
	appID := ctx.Query("app_id")
	username := ctx.Query("username")
	email := ctx.Query("email")

	// 构建查询条件
	query := config.DB.Model(&models.SystemAdmin{})
	
	if adminType != "" {
		query = query.Where("admin_type = ?", adminType)
	}
	if appID != "" {
		query = query.Where("app_id = ?", appID)
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var admins []models.SystemAdmin
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&admins).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询系统管理员失败"})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"data": admins,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// CreateSystemAdmin 创建系统管理员
func (c *SystemAdminManagementController) CreateSystemAdmin(ctx *gin.Context) {
	var req struct {
		Username  string `json:"username" binding:"required"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Password  string `json:"password" binding:"required"`
		AdminType string `json:"admin_type" binding:"required,oneof=system app"`
		AppID     string `json:"app_id"`
		IsActive  *bool  `json:"is_active"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证应用级管理员必须指定app_id
	if req.AdminType == "app" && req.AppID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "应用级管理员必须指定应用ID"})
		return
	}

	// 验证系统级管理员不能指定app_id
	if req.AdminType == "system" && req.AppID != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "系统级管理员不能指定应用ID"})
		return
	}

	// 检查用户名是否已存在
	var existingAdmin models.SystemAdmin
	if err := config.DB.Where("username = ?", req.Username).First(&existingAdmin).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		if err := config.DB.Where("email = ?", req.Email).First(&existingAdmin).Error; err == nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": "邮箱已存在"})
			return
		}
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 设置默认值
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// 创建管理员
	admin := models.SystemAdmin{
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  hashedPassword,
		AdminType: req.AdminType,
		AppID:     req.AppID,
		IsActive:  isActive,
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建系统管理员失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "系统管理员创建成功",
		"data":    admin,
	})
}

// UpdateSystemAdmin 更新系统管理员
func (c *SystemAdminManagementController) UpdateSystemAdmin(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的管理员ID"})
		return
	}

	var req struct {
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
		Password  *string `json:"password"`
		IsActive  *bool   `json:"is_active"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找管理员
	var admin models.SystemAdmin
	if err := config.DB.First(&admin, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "管理员不存在"})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	
	if req.Email != nil {
		// 检查邮箱是否已被其他管理员使用
		var existingAdmin models.SystemAdmin
		if err := config.DB.Where("email = ? AND id != ?", *req.Email, id).First(&existingAdmin).Error; err == nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": "邮箱已被其他管理员使用"})
			return
		}
		updates["email"] = *req.Email
	}
	
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	
	if req.Password != nil {
		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}
		updates["password"] = hashedPassword
	}
	
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// 执行更新
	if err := config.DB.Model(&admin).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新系统管理员失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "系统管理员更新成功",
		"data":    admin,
	})
}

// DeleteSystemAdmin 删除系统管理员
func (c *SystemAdminManagementController) DeleteSystemAdmin(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的管理员ID"})
		return
	}

	// 查找管理员
	var admin models.SystemAdmin
	if err := config.DB.First(&admin, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "管理员不存在"})
		return
	}

	// 不能删除自己
	userID, exists := ctx.Get("admin_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取当前用户信息"})
		return
	}
	if admin.ID == userID.(uint) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不能删除自己的账号"})
		return
	}

	// 软删除
	if err := config.DB.Delete(&admin).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除系统管理员失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "系统管理员删除成功"})
}

// ResetSystemAdminPassword 重置系统管理员密码
func (c *SystemAdminManagementController) ResetSystemAdminPassword(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的管理员ID"})
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找管理员
	var admin models.SystemAdmin
	if err := config.DB.First(&admin, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "管理员不存在"})
		return
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 更新密码
	if err := config.DB.Model(&admin).Update("password", hashedPassword).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "重置密码失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}
