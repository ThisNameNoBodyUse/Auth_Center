package middleware

import (
	"net/http"

	"auth-center/config"
	"auth-center/models"

	"github.com/gin-gonic/gin"
)

// SuperAdminMiddleware 超级管理员鉴权中间件
func SuperAdminMiddleware() gin.HandlerFunc {
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

		// 查询用户是否为超级管理员
		var user models.User
		if err := config.DB.Where("id = ? AND app_id = ? AND is_super_admin = ? AND status = 1",
			userID, appID, true).First(&user).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要超级管理员权限"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("super_admin_user", user)
		c.Next()
	}
}
