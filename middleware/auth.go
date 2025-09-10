package middleware

import (
	"net/http"
	"strings"

	"auth-center/service"
	"auth-center/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// 解析令牌
		claims, err := utils.ParseAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 检查令牌是否在黑名单中
		blacklistKey := utils.TokenBlacklistPrefix + claims.JTI
		exists, err := utils.Exists(blacklistKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			c.Abort()
			return
		}
		if exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("app_id", claims.AppID)
		c.Set("roles", claims.Roles)
		c.Set("jti", claims.JTI)

		c.Next()
	}
}

// AppAuthMiddleware 应用认证中间件
func AppAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取应用ID和密钥
		appID := c.GetHeader("X-App-Id")
		appSecret := c.GetHeader("X-App-Secret")

		if appID == "" || appSecret == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "App ID and Secret are required"})
			c.Abort()
			return
		}

		// 验证应用凭据
		app, err := service.ValidateAppCredentials(appID, appSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid app credentials"})
			c.Abort()
			return
		}

		// 检查应用状态
		if app.Status != 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "App is disabled"})
			c.Abort()
			return
		}

		// 将应用信息存储到上下文中
		c.Set("app_id", app.AppID)
		c.Set("app_name", app.Name)

		c.Next()
	}
}

// PermissionMiddleware 权限中间件
func PermissionMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID和应用ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		appID, exists := c.Get("app_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "App not authenticated"})
			c.Abort()
			return
		}

		// 检查用户权限
		hasPermission, err := service.CheckUserPermission(userID.(uint), appID.(string), permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIPermissionMiddleware API权限中间件
func APIPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID和应用ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		appID, exists := c.Get("app_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "App not authenticated"})
			c.Abort()
			return
		}

		// 获取请求的API路径和方法
		apiPath := c.FullPath()
		apiMethod := c.Request.Method

		// 检查用户是否有权限访问该API
		hasPermission, err := service.CheckAPIPermission(userID.(uint), appID.(string), apiPath, apiMethod)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API permission check failed"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "API access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetAppID 从上下文获取应用ID
func GetAppID(c *gin.Context) (string, bool) {
	appID, exists := c.Get("app_id")
	if !exists {
		return "", false
	}
	return appID.(string), true
}

// GetRoles 从上下文获取用户角色
func GetRoles(c *gin.Context) ([]uint, bool) {
	roles, exists := c.Get("roles")
	if !exists {
		return nil, false
	}
	return roles.([]uint), true
}
