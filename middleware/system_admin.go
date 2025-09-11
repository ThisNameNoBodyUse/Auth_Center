package middleware

import (
	"net/http"

	"auth-center/config"
	"auth-center/models"

	"github.com/gin-gonic/gin"
)

// SystemAdminMiddleware 系统级超级管理员中间件
// 只有系统级超级管理员可以访问，用于管理应用、应用密钥等系统级资源
func SystemAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			c.Abort()
			return
		}

		// 获取应用ID
		_, exists = c.Get("app_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "应用未认证"})
			c.Abort()
			return
		}

		// 检查是否为系统级超级管理员（必须是system-admin应用的超级管理员）
		var user models.User
		if err := config.DB.Where("id = ? AND app_id = ? AND is_super_admin = ? AND status = 1",
			userID, "system-admin", true).First(&user).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要系统级超级管理员权限"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("system_admin_user", user)
		c.Next()
	}
}

// AppAdminMiddleware 应用级超级管理员中间件
// 应用级超级管理员可以访问，用于管理自己应用内的角色、权限、用户等
func AppAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			c.Abort()
			return
		}

		// 获取应用ID
		appID, exists := c.Get("app_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "应用未认证"})
			c.Abort()
			return
		}

		// 检查是否为超级管理员（系统级或应用级都可以）
		var user models.User
		if err := config.DB.Where("id = ? AND app_id = ? AND is_super_admin = ? AND status = 1",
			userID, appID, true).First(&user).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要超级管理员权限"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("admin_user", user)
		c.Set("is_system_admin", appID == "system-admin")
		c.Next()
	}
}

// FlexibleAdminMiddleware 灵活的超级管理员中间件
// 系统级超级管理员可以访问任何应用的数据，应用级超级管理员只能访问自己应用的数据
func FlexibleAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			c.Abort()
			return
		}

		// 获取应用ID
		appID, exists := c.Get("app_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "应用未认证"})
			c.Abort()
			return
		}

		// 首先检查是否为系统级超级管理员
		var systemUser models.User
		if err := config.DB.Where("id = ? AND app_id = ? AND is_super_admin = ? AND status = 1",
			userID, "system-admin", true).First(&systemUser).Error; err == nil {
			// 系统级超级管理员，可以访问任何应用的数据
			c.Set("admin_user", systemUser)
			c.Set("is_system_admin", true)
			c.Set("can_access_any_app", true)
			c.Next()
			return
		}

		// 检查是否为应用级超级管理员
		var user models.User
		if err := config.DB.Where("id = ? AND app_id = ? AND is_super_admin = ? AND status = 1",
			userID, appID, true).First(&user).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要超级管理员权限"})
			c.Abort()
			return
		}

		// 应用级超级管理员，只能访问自己应用的数据
		c.Set("admin_user", user)
		c.Set("is_system_admin", false)
		c.Set("can_access_any_app", false)
		c.Next()
	}
}

// GetSystemAdminUser 从上下文获取系统级超级管理员用户
func GetSystemAdminUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("system_admin_user")
	if !exists {
		return nil, false
	}
	return user.(*models.User), true
}

// GetAdminUser 从上下文获取管理员用户
func GetAdminUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("admin_user")
	if !exists {
		return nil, false
	}
	return user.(*models.User), true
}

