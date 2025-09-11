package middleware

import (
	"net/http"
	"strings"

	"auth-center/config"
	"auth-center/models"
	"auth-center/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// SystemAdminAuthMiddleware 系统管理员认证中间件
func SystemAdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证令牌"})
			c.Abort()
			return
		}

		// 检查Bearer格式
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证格式"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// 验证令牌
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌声明"})
			c.Abort()
			return
		}

		// 检查令牌类型
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "system_admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌类型"})
			c.Abort()
			return
		}

		// 检查令牌是否在黑名单中
		blacklistKey := "blacklist:" + tokenString
		exists, err := utils.Exists(blacklistKey)
		if err == nil && exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌已失效"})
			c.Abort()
			return
		}

		// 获取管理员信息
		adminID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户ID"})
			c.Abort()
			return
		}

		// 验证管理员是否存在且激活
		var admin models.SystemAdmin
		if err := config.DB.Where("id = ? AND is_active = ?", uint(adminID), true).First(&admin).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "管理员不存在或已禁用"})
			c.Abort()
			return
		}

		// 将管理员信息存储到上下文
		c.Set("admin_id", admin.ID)
		c.Set("admin_username", admin.Username)
		c.Set("admin_type", admin.AdminType)
		c.Set("admin_app_id", admin.AppID)
		c.Set("is_system_admin", admin.AdminType == "system")
		c.Set("is_app_admin", admin.AdminType == "app")

		c.Next()
	}
}

// SystemAdminOnlyMiddleware 仅系统级管理员中间件
func SystemAdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminType, exists := c.Get("admin_type")
		if !exists || adminType != "system" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要系统级管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AppAdminOnlyMiddleware 仅应用级管理员中间件
func AppAdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminType, exists := c.Get("admin_type")
		if !exists || adminType != "app" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要应用级管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// FlexibleSystemAdminMiddleware 灵活的系统管理员中间件
// 系统级管理员可以访问所有资源，应用级管理员只能访问自己应用的资源
func FlexibleSystemAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminType, exists := c.Get("admin_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到管理员类型"})
			c.Abort()
			return
		}

		// 系统级管理员可以访问所有资源
		if adminType == "system" {
			c.Set("can_access_any_app", true)
			c.Next()
			return
		}

		// 应用级管理员只能访问自己应用的资源
		if adminType == "app" {
			adminAppID, exists := c.Get("admin_app_id")
			if !exists {
				c.JSON(http.StatusForbidden, gin.H{"error": "未找到管理员应用ID"})
				c.Abort()
				return
			}

			// 检查请求的应用ID是否匹配
			requestAppID := c.Query("app_id")
			if requestAppID != "" && requestAppID != adminAppID {
				c.JSON(http.StatusForbidden, gin.H{"error": "无权访问其他应用的数据"})
				c.Abort()
				return
			}

			c.Set("can_access_any_app", false)
			c.Set("target_app_id", adminAppID)
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "无效的管理员类型"})
		c.Abort()
	}
}

// GetSystemAdminID 获取系统管理员ID
func GetSystemAdminID(c *gin.Context) (uint, bool) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		return 0, false
	}
	id, ok := adminID.(uint)
	return id, ok
}

// GetSystemAdminType 获取系统管理员类型
func GetSystemAdminType(c *gin.Context) (string, bool) {
	adminType, exists := c.Get("admin_type")
	if !exists {
		return "", false
	}
	t, ok := adminType.(string)
	return t, ok
}

// IsSystemAdmin 检查是否为系统级管理员
func IsSystemAdmin(c *gin.Context) bool {
	isSystemAdmin, exists := c.Get("is_system_admin")
	return exists && isSystemAdmin.(bool)
}

// IsAppAdmin 检查是否为应用级管理员
func IsAppAdmin(c *gin.Context) bool {
	isAppAdmin, exists := c.Get("is_app_admin")
	return exists && isAppAdmin.(bool)
}

// CanAccessAnyApp 检查是否可以访问任何应用
func CanAccessAnyApp(c *gin.Context) bool {
	canAccess, exists := c.Get("can_access_any_app")
	return exists && canAccess.(bool)
}

// GetTargetAppID 获取目标应用ID
func GetTargetAppID(c *gin.Context) string {
	if CanAccessAnyApp(c) {
		// 系统级管理员可以通过查询参数指定应用ID
		if targetAppID := c.Query("app_id"); targetAppID != "" {
			return targetAppID
		}
	}

	// 应用级管理员或未指定应用ID时，使用管理员的应用ID
	targetAppID, exists := c.Get("target_app_id")
	if exists {
		return targetAppID.(string)
	}

	adminAppID, exists := c.Get("admin_app_id")
	if exists {
		return adminAppID.(string)
	}

	return ""
}
