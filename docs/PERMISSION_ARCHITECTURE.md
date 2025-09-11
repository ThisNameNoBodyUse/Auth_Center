# 权限架构设计

## 概述

本系统采用多层级权限管理架构，支持系统级超级管理员和应用级超级管理员两种不同的权限级别。

## 权限层级

### 1. 系统级超级管理员 (System Admin)

**权限范围：**
- 管理所有应用（创建、更新、删除应用）
- 管理应用密钥（生成、重新生成）
- 查看所有应用的用户列表
- 管理应用的基本信息

**限制：**
- 不能直接管理其他应用内的角色、权限、用户
- 不能查看其他应用的详细数据

**身份验证：**
- 必须是 `system-admin` 应用的超级管理员
- 通过 `is_super_admin = true` 标识

**API 路径：**
```
/api/v1/apps/* - 应用管理相关接口
```

### 2. 应用级超级管理员 (App Admin)

**权限范围：**
- 管理自己应用内的角色（创建、更新、删除角色）
- 管理自己应用内的权限（创建、更新、删除权限）
- 管理自己应用内的用户
- 为角色分配权限

**限制：**
- 只能管理自己应用内的数据
- 不能查看其他应用的数据
- 不能管理应用密钥
- 不能创建或删除应用

**身份验证：**
- 可以是任何应用的超级管理员
- 通过 `is_super_admin = true` 标识
- 权限范围限制在登录时使用的 `app_id`

**API 路径：**
```
/api/v1/app/* - 应用内资源管理接口
```

## 中间件设计

### SystemAdminMiddleware
- 验证用户是否为 `system-admin` 应用的超级管理员
- 用于系统级管理接口

### AppAdminMiddleware
- 验证用户是否为任何应用的超级管理员
- 用于应用内资源管理接口
- 自动限制权限范围到当前应用

## 路由设计

### 系统级管理路由
```go
// 需要系统级超级管理员权限
apps := v1.Group("/apps")
apps.Use(middleware.AuthMiddleware(), middleware.AppAuthMiddleware(), middleware.SystemAdminMiddleware())
```

### 应用内管理路由
```go
// 需要应用级超级管理员权限
appResources := v1.Group("/app")
appResources.Use(middleware.AuthMiddleware(), middleware.AppAuthMiddleware(), middleware.AppAdminMiddleware())
```

## 权限检查逻辑

### 系统级权限检查
```go
// 检查是否为系统级超级管理员
if middleware.IsSystemAdmin(ctx) {
    // 可以访问所有应用数据
} else {
    // 只能访问自己应用的数据
}
```

### 应用级权限检查
```go
// 获取当前应用ID
appID, _ := middleware.GetAppID(ctx)

// 检查权限范围
if !middleware.IsSystemAdmin(ctx) {
    // 应用级超级管理员只能查看自己应用的数据
    if targetAppID != appID {
        return "无权访问其他应用数据"
    }
}
```

## 数据库设计

### 用户表权限字段
```sql
is_super_admin TINYINT NOT NULL DEFAULT 0 COMMENT '是否为超级管理员 1是 0否'
```

### 权限隔离
- 所有应用内资源都通过 `app_id` 字段进行隔离
- 系统级超级管理员可以跨应用访问
- 应用级超级管理员只能访问自己应用的数据

## 安全特性

### 1. 权限隔离
- 应用间数据完全隔离
- 应用级超级管理员无法访问其他应用数据
- 系统级超级管理员有明确的管理范围限制

### 2. 权限验证
- 多层中间件验证
- 用户认证 + 应用认证 + 权限级别验证
- 自动权限范围限制

### 3. 操作审计
- 所有管理操作都有明确的权限要求
- 可以追踪操作来源（系统级 vs 应用级）
- 支持操作日志记录

## 使用示例

### 系统级超级管理员登录
```bash
# 使用 system-admin 应用登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "X-App-Id: system-admin" \
  -H "X-App-Secret: system-admin-secret-key-change-in-production" \
  -d '{"username": "superadmin", "password": "admin123"}'
```

### 应用级超级管理员登录
```bash
# 使用其他应用登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "X-App-Id: your-app-id" \
  -H "X-App-Secret: your-app-secret" \
  -d '{"username": "appadmin", "password": "password123"}'
```

### 权限验证
```go
// 在控制器中检查权限级别
if middleware.IsSystemAdmin(ctx) {
    // 系统级操作
} else {
    // 应用级操作，自动限制到当前应用
}
```

## 部署建议

### 1. 环境隔离
- 系统管理接口建议限制在内网访问
- 应用管理接口可以对外开放
- 使用防火墙限制管理接口访问

### 2. 监控审计
- 记录所有管理操作日志
- 区分系统级和应用级操作
- 监控权限使用情况

### 3. 密钥管理
- 应用密钥由系统级超级管理员管理
- 定期轮换应用密钥
- 建立密钥泄露应急响应机制
