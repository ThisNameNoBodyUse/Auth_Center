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
		// 应用认证路由（外部应用使用）
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

		// 系统管理路由（系统内部使用）
		system := v1.Group("/system")
		{
			systemAdminController := &controllers.SystemAdminController{}
			system.POST("/login", systemAdminController.SystemLogin)
			system.POST("/register", systemAdminController.SystemRegister)
			system.POST("/refresh", systemAdminController.SystemRefreshToken)
			system.POST("/logout", systemAdminController.SystemLogout)
			system.GET("/admin/info", middleware.SystemAdminAuthMiddleware(), systemAdminController.GetSystemAdminInfo)
		}

		// 系统级应用管理路由（仅系统级超级管理员）
		appManagementController := &controllers.AppManagementController{}
		apps := v1.Group("/apps")
		apps.Use(middleware.SystemAdminAuthMiddleware(), middleware.SystemAdminOnlyMiddleware())
		{
			apps.POST("", appManagementController.CreateApp)
			apps.GET("", appManagementController.ListApps)
			apps.GET("/:app_id", appManagementController.GetApp)
			apps.PUT("/:app_id", appManagementController.UpdateApp)
			apps.DELETE("/:app_id", appManagementController.DeleteApp)
			apps.POST("/:app_id/regenerate-secret", appManagementController.RegenerateAppSecret)
			apps.GET("/:app_id/users", appManagementController.ListAppUsers)
		}

		// 系统管理员管理路由（仅系统级超级管理员）
		systemAdminManagementController := &controllers.SystemAdminManagementController{}
		systemAdmins := v1.Group("/system-admins")
		systemAdmins.Use(middleware.SystemAdminAuthMiddleware(), middleware.SystemAdminOnlyMiddleware())
		{
			systemAdmins.GET("", systemAdminManagementController.ListSystemAdmins)
			systemAdmins.POST("", systemAdminManagementController.CreateSystemAdmin)
			systemAdmins.PUT("/:id", systemAdminManagementController.UpdateSystemAdmin)
			systemAdmins.DELETE("/:id", systemAdminManagementController.DeleteSystemAdmin)
			systemAdmins.POST("/:id/reset-password", systemAdminManagementController.ResetSystemAdminPassword)
		}

		// 应用内资源管理路由（系统级和应用级超级管理员）
		appResourceController := &controllers.AppResourceController{}
		appResources := v1.Group("/app")
		appResources.Use(middleware.SystemAdminAuthMiddleware(), middleware.FlexibleSystemAdminMiddleware())
		{
			// 角色管理
			roles := appResources.Group("/roles")
			{
				roles.GET("", appResourceController.ListRoles)
				roles.POST("", appResourceController.CreateRole)
				roles.PUT("/:id", appResourceController.UpdateRole)
				roles.DELETE("/:id", appResourceController.DeleteRole)
				roles.POST("/:id/permissions", appResourceController.AssignRolePermissions)
				roles.GET("/:id/permissions", appResourceController.GetRolePermissions)
			}

			// 权限管理
			permissions := appResources.Group("/permissions")
			{
				permissions.GET("", appResourceController.ListPermissions)
				permissions.POST("", appResourceController.CreatePermission)
				permissions.PUT("/:id", appResourceController.UpdatePermission)
				permissions.DELETE("/:id", appResourceController.DeletePermission)
			}

			// 用户管理
			users := appResources.Group("/users")
			{
				users.GET("", appResourceController.ListUsers)
				users.POST("", appResourceController.CreateUser)
				users.PUT("/:id", appResourceController.UpdateUser)
				users.DELETE("/:id", appResourceController.DeleteUser)
				users.POST("/:id/roles", appResourceController.AssignUserRoles)
				users.GET("/:id/roles", appResourceController.GetUserRoles)
			}
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
