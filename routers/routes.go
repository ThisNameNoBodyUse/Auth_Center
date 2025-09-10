package routers

import (
	"auth-center/controllers"
	"auth-center/middleware"

	"github.com/gin-gonic/gin"
)

// InitRoutes 初始化路由
func InitRoutes(r *gin.Engine) {
	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := v1.Group("/auth")
		{
			authController := &controllers.AuthController{}
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
			auth.POST("/refresh", authController.RefreshToken)
			auth.POST("/logout", authController.Logout)

			// 需要认证的路由
			auth.Use(middleware.AuthMiddleware())
			auth.GET("/user", authController.GetUserInfo)
		}

		// 应用管理路由（需要超级管理员权限）
		appController := &controllers.AppController{}
		apps := v1.Group("/apps")
		apps.Use(middleware.AuthMiddleware(), middleware.AppAuthMiddleware(), middleware.SuperAdminMiddleware())
		{
			apps.POST("", appController.CreateApp)
			apps.GET("", appController.ListApps)
			apps.GET("/:app_id", appController.GetApp)
			apps.PUT("/:app_id", appController.UpdateApp)
			apps.DELETE("/:app_id", appController.DeleteApp)
			apps.POST("/:app_id/regenerate-secret", appController.RegenerateAppSecret)
		}

		// 权限管理路由
		permissions := v1.Group("/permissions")
		permissions.Use(middleware.AuthMiddleware())
		{
			permissionController := &controllers.PermissionController{}
			permissions.GET("/check", permissionController.CheckPermission)
			permissions.GET("/check-api", permissionController.CheckAPIPermission)
			permissions.GET("/user", permissionController.GetUserPermissions)
			permissions.GET("/roles", permissionController.GetUserRoles)
		}

		// 管理后台路由（需要应用认证）
		admin := v1.Group("/admin")
		admin.Use(middleware.AppAuthMiddleware())
		{
			// 这里可以添加管理后台相关的路由
			// 例如：用户管理、角色管理、权限管理等
		}
	}
}
