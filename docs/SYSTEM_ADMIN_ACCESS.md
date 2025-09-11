# 系统级超级管理员跨应用访问说明

## 问题背景

您提出了一个很好的问题：应用级接口使用了 `AppAdminMiddleware`，系统级超级管理员如何访问应用级接口？

## 解决方案

我们创建了一个新的中间件 `FlexibleAdminMiddleware`，它允许：

1. **系统级超级管理员**：可以访问任何应用的数据
2. **应用级超级管理员**：只能访问自己应用的数据

## 技术实现

### 1. 新的中间件：FlexibleAdminMiddleware

```go
func FlexibleAdminMiddleware() gin.HandlerFunc {
    // 首先检查是否为系统级超级管理员
    if 是系统级超级管理员 {
        // 可以访问任何应用的数据
        c.Set("can_access_any_app", true)
    } else if 是应用级超级管理员 {
        // 只能访问自己应用的数据
        c.Set("can_access_any_app", false)
    }
}
```

### 2. 控制器中的跨应用访问支持

```go
func (c *AppResourceController) getTargetAppID(ctx *gin.Context) string {
    if middleware.IsSystemAdmin(ctx) {
        // 系统级超级管理员可以通过查询参数指定应用ID
        if targetAppID := ctx.Query("app_id"); targetAppID != "" {
            return targetAppID
        }
    }
    // 应用级超级管理员只能使用当前应用ID
    appID, _ := middleware.GetAppID(ctx)
    return appID
}
```

## 使用方式

### 系统级超级管理员访问其他应用的数据

```bash
# 查看 my-app 应用的角色列表
curl -X GET "http://localhost:8080/api/v1/app/roles?app_id=my-app" \
  -H "Authorization: Bearer <system-admin-token>" \
  -H "X-App-Id: system-admin" \
  -H "X-App-Secret: system-admin-secret"

# 查看 other-app 应用的用户列表
curl -X GET "http://localhost:8080/api/v1/app/users?app_id=other-app" \
  -H "Authorization: Bearer <system-admin-token>" \
  -H "X-App-Id: system-admin" \
  -H "X-App-Secret: system-admin-secret"

# 在 other-app 应用中创建角色
curl -X POST "http://localhost:8080/api/v1/app/roles?app_id=other-app" \
  -H "Authorization: Bearer <system-admin-token>" \
  -H "X-App-Id: system-admin" \
  -H "X-App-Secret: system-admin-secret" \
  -d '{"name": "新角色", "code": "new_role", "description": "系统管理员创建的角色"}'
```

### 应用级超级管理员访问自己应用的数据

```bash
# 查看自己应用的角色列表（不需要 app_id 参数）
curl -X GET "http://localhost:8080/api/v1/app/roles" \
  -H "Authorization: Bearer <app-admin-token>" \
  -H "X-App-Id: my-app" \
  -H "X-App-Secret: my-app-secret"

# 创建角色（自动在当前应用中创建）
curl -X POST "http://localhost:8080/api/v1/app/roles" \
  -H "Authorization: Bearer <app-admin-token>" \
  -H "X-App-Id: my-app" \
  -H "X-App-Secret: my-app-secret" \
  -d '{"name": "编辑者", "code": "editor", "description": "内容编辑角色"}'
```

## 权限对比

| 操作 | 系统级超级管理员 | 应用级超级管理员 |
|------|------------------|------------------|
| 访问自己应用数据 | ✅ | ✅ |
| 访问其他应用数据 | ✅ (通过 ?app_id=xxx) | ❌ |
| 创建应用 | ✅ | ❌ |
| 管理应用密钥 | ✅ | ❌ |
| 管理应用内资源 | ✅ (任何应用) | ✅ (仅自己应用) |

## 路由配置

```go
// 应用内资源管理路由（系统级和应用级超级管理员）
appResources := v1.Group("/app")
appResources.Use(middleware.AuthMiddleware(), middleware.AppAuthMiddleware(), middleware.FlexibleAdminMiddleware())
```

## 安全特性

1. **权限验证**：系统级超级管理员必须通过 `system-admin` 应用登录
2. **应用隔离**：应用级超级管理员无法访问其他应用的数据
3. **参数验证**：系统级超级管理员必须明确指定 `app_id` 参数
4. **自动回退**：如果没有指定 `app_id`，系统级超级管理员会使用当前应用ID

## 总结

现在系统级超级管理员可以：

1. **通过系统级接口**管理应用（`/api/v1/apps/*`）
2. **通过应用级接口**管理任何应用内的资源（`/api/v1/app/*?app_id=xxx`）
3. **跨应用访问**任何应用的角色、权限、用户数据

而应用级超级管理员只能：

1. **通过应用级接口**管理自己应用内的资源（`/api/v1/app/*`）
2. **无法访问**其他应用的数据
3. **无法管理**应用本身

这样就实现了您要求的多层级权限架构！
