# 系统架构更新说明

## 概述

本次更新将系统内部登录和外部应用登录完全分离，创建了独立的系统管理员表和管理体系，实现了更清晰的权限架构。

## 主要变更

### 1. 数据库结构变更

#### 新增系统管理员表 (`system_admins`)
```sql
CREATE TABLE system_admins (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(128) NOT NULL UNIQUE,
  email VARCHAR(255) NULL UNIQUE,
  phone VARCHAR(32) NULL,
  password VARCHAR(255) NOT NULL,
  admin_type VARCHAR(20) NOT NULL COMMENT 'system: 系统级管理员, app: 应用级管理员',
  app_id VARCHAR(64) NULL COMMENT '应用级管理员关联的应用ID',
  is_active TINYINT NOT NULL DEFAULT 1,
  last_login_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL DEFAULT NULL
);
```

### 2. 认证体系分离

#### 外部应用认证 (`/api/v1/auth/*`)
- **用途**: 供其他应用调用，进行用户认证
- **认证方式**: 需要提供 `app_id` 和 `app_secret`
- **用户表**: 使用 `users` 表
- **权限**: 基于应用内的角色和权限

#### 系统内部认证 (`/api/v1/system/*`)
- **用途**: 系统内部管理界面使用
- **认证方式**: 仅需要用户名和密码
- **用户表**: 使用 `system_admins` 表
- **权限**: 基于管理员类型（系统级/应用级）

### 3. 权限架构

#### 系统级管理员 (`admin_type = 'system'`)
- 可以管理所有应用
- 可以查看和管理所有应用的用户、角色、权限
- 可以管理应用密钥
- 通过 `app_id` 查询参数指定目标应用

#### 应用级管理员 (`admin_type = 'app'`)
- 只能管理自己应用内的用户、角色、权限
- 不能管理应用密钥
- 不能访问其他应用的数据
- 通过 `app_id` 字段关联特定应用

### 4. 新增文件

#### 服务层
- `service/system_admin_service.go` - 系统管理员服务

#### 控制器层
- `controllers/system_admin_controller.go` - 系统管理员控制器

#### 中间件层
- `middleware/system_admin_auth.go` - 系统管理员认证中间件

### 5. 路由更新

#### 应用认证路由 (`/api/v1/auth/*`)
```go
auth := v1.Group("/auth")
{
    auth.POST("/login", authController.Login)           // 需要 app_id + app_secret
    auth.POST("/register", authController.Register)     // 需要 app_id + app_secret
    auth.POST("/refresh", authController.RefreshToken)
    auth.POST("/logout", authController.Logout)
    auth.GET("/user", middleware.AuthMiddleware(), authController.GetUserInfo)
}
```

#### 系统管理路由 (`/api/v1/system/*`)
```go
system := v1.Group("/system")
{
    system.POST("/login", systemAdminController.SystemLogin)           // 仅需用户名密码
    system.POST("/register", systemAdminController.SystemRegister)     // 注册系统管理员
    system.POST("/refresh", systemAdminController.SystemRefreshToken)
    system.POST("/logout", systemAdminController.SystemLogout)
    system.GET("/admin/info", middleware.SystemAdminAuthMiddleware(), systemAdminController.GetSystemAdminInfo)
}
```

#### 应用管理路由 (`/api/v1/apps/*`)
```go
apps := v1.Group("/apps")
apps.Use(middleware.SystemAdminAuthMiddleware(), middleware.SystemAdminOnlyMiddleware())
{
    // 仅系统级管理员可访问
}
```

#### 应用资源管理路由 (`/api/v1/app/*`)
```go
appResources := v1.Group("/app")
appResources.Use(middleware.SystemAdminAuthMiddleware(), middleware.FlexibleSystemAdminMiddleware())
{
    // 系统级和应用级管理员都可访问，但权限不同
}
```

### 6. 中间件更新

#### 系统管理员认证中间件
- `SystemAdminAuthMiddleware()` - 验证系统管理员令牌
- `SystemAdminOnlyMiddleware()` - 仅系统级管理员
- `AppAdminOnlyMiddleware()` - 仅应用级管理员
- `FlexibleSystemAdminMiddleware()` - 灵活的系统管理员中间件

### 7. 初始化数据

#### 系统管理员数据
```sql
INSERT INTO system_admins (username, email, phone, password, admin_type, app_id, is_active) VALUES
('superadmin', 'admin@auth-center.com', '', '$2b$10$UtbwZjygOigggJA.7So9v.cu0S1B.ibbBUNxdtA8GmwFVi86cZSye', 'system', NULL, 1),
('appadmin', 'appadmin@auth-center.com', '13800138000', '$2b$10$UtbwZjygOigggJA.7So9v.cu0S1B.ibbBUNxdtA8GmwFVi86cZSye', 'app', 'default-app', 1);
```

## 使用示例

### 1. 外部应用登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user1",
    "password": "password123",
    "app_id": "default-app",
    "app_secret": "default-secret"
  }'
```

### 2. 系统管理员登录
```bash
curl -X POST http://localhost:8080/api/v1/system/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "superadmin",
    "password": "admin123"
  }'
```

### 3. 系统级管理员管理应用
```bash
curl -X GET "http://localhost:8080/api/v1/apps" \
  -H "Authorization: Bearer <system_admin_token>"
```

### 4. 应用级管理员管理自己应用的用户
```bash
curl -X GET "http://localhost:8080/api/v1/app/users?app_id=default-app" \
  -H "Authorization: Bearer <app_admin_token>"
```

## 安全考虑

1. **应用密钥验证**: 外部应用必须提供正确的 `app_id` 和 `app_secret` 才能进行用户认证
2. **令牌分离**: 系统管理员和应用用户使用不同的令牌体系
3. **权限隔离**: 应用级管理员只能访问自己应用的数据
4. **密码加密**: 所有密码都使用 bcrypt 加密存储

## 迁移指南

1. **数据库迁移**: 运行更新后的 `docker/init.sql` 脚本
2. **代码更新**: 所有相关代码已更新，无需额外修改
3. **前端适配**: 前端需要根据登录类型调用不同的接口
4. **API文档**: 已更新 API 文档，包含新的系统管理接口

## 总结

通过这次更新，系统实现了：
- 清晰的职责分离：外部应用认证 vs 系统内部管理
- 灵活的权限管理：系统级和应用级管理员
- 更好的安全性：应用密钥验证和权限隔离
- 易于维护：独立的认证体系和中间件

这种架构更适合多租户的认证授权中心，既保证了外部应用的安全性，又提供了灵活的内部管理能力。
