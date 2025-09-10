# 认证授权中心 API 文档

## 概述

认证授权中心提供统一的认证授权服务，支持多应用、多租户的权限管理。

## 基础信息

- **基础URL**: `http://localhost:8080/api/v1`
- **认证方式**: Bearer Token
- **内容类型**: `application/json`

## 认证流程

1. 应用注册获取 `app_id` 和 `app_secret`
2. 用户注册/登录获取访问令牌
3. 使用访问令牌访问受保护的资源
4. 令牌过期时使用刷新令牌获取新令牌

## API 接口

### 1. 认证相关

#### 1.1 用户登录

**POST** `/auth/login`

**请求体:**
```json
{
  "app_id": "your-app-id",
  "username": "username",
  "password": "password"
}
```

**响应:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "token_type": "Bearer",
  "user": {
    "id": 1,
    "username": "username",
    "email": "user@example.com",
    "avatar": "avatar_url",
    "roles": [
      {
        "id": 1,
        "name": "管理员",
        "code": "admin"
      }
    ]
  }
}
```

#### 1.2 用户注册

**POST** `/auth/register`

**请求体:**
```json
{
  "app_id": "your-app-id",
  "username": "username",
  "email": "user@example.com",
  "password": "password"
}
```

**响应:**
```json
{
  "message": "注册成功"
}
```

#### 1.3 刷新令牌

**POST** `/auth/refresh`

**请求体:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "token_type": "Bearer",
  "user": {
    "id": 1,
    "username": "username",
    "email": "user@example.com",
    "avatar": "avatar_url",
    "roles": [...]
  }
}
```

#### 1.4 用户登出

**POST** `/auth/logout`

**请求头:**
```
Authorization: Bearer <access_token>
```

**请求体:**
```json
{
  "token": "access_token"
}
```

**响应:**
```json
{
  "message": "登出成功"
}
```

#### 1.5 获取用户信息

**GET** `/auth/user`

**请求头:**
```
Authorization: Bearer <access_token>
```

**响应:**
```json
{
  "id": 1,
  "username": "username",
  "email": "user@example.com",
  "avatar": "avatar_url",
  "roles": [
    {
      "id": 1,
      "name": "管理员",
      "code": "admin"
    }
  ]
}
```

### 2. 应用管理

#### 2.1 创建应用

**POST** `/apps`

**请求体:**
```json
{
  "name": "应用名称",
  "description": "应用描述"
}
```

**响应:**
```json
{
  "id": 1,
  "name": "应用名称",
  "app_id": "app_12345678",
  "app_secret": "secret_12345678",
  "description": "应用描述",
  "status": 1,
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### 2.2 获取应用信息

**GET** `/apps/{app_id}`

**响应:**
```json
{
  "id": 1,
  "name": "应用名称",
  "app_id": "app_12345678",
  "description": "应用描述",
  "status": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### 2.3 更新应用

**PUT** `/apps/{app_id}`

**请求体:**
```json
{
  "name": "新应用名称",
  "description": "新应用描述",
  "status": 1
}
```

**响应:**
```json
{
  "message": "更新成功"
}
```

#### 2.4 删除应用

**DELETE** `/apps/{app_id}`

**响应:**
```json
{
  "message": "删除成功"
}
```

#### 2.5 获取应用列表

**GET** `/apps?page=1&page_size=10`

**响应:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "应用名称",
      "app_id": "app_12345678",
      "description": "应用描述",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 1,
    "total_page": 1
  }
}
```

#### 2.6 重新生成应用密钥

**POST** `/apps/{app_id}/regenerate-secret`

**响应:**
```json
{
  "app_secret": "new_secret_12345678"
}
```

### 3. 权限管理

#### 3.1 检查权限

**GET** `/permissions/check?permission=user:read`

**请求头:**
```
Authorization: Bearer <access_token>
```

**响应:**
```json
{
  "has_permission": true
}
```

#### 3.2 检查API权限

**GET** `/permissions/check-api?path=/api/v1/users&method=GET`

**请求头:**
```
Authorization: Bearer <access_token>
```

**响应:**
```json
{
  "has_permission": true
}
```

#### 3.3 获取用户权限列表

**GET** `/permissions/user`

**请求头:**
```
Authorization: Bearer <access_token>
```

**响应:**
```json
{
  "permissions": ["user:read", "user:create", "user:update"]
}
```

#### 3.4 获取用户角色列表

**GET** `/permissions/roles`

**请求头:**
```
Authorization: Bearer <access_token>
```

**响应:**
```json
{
  "roles": [
    {
      "id": 1,
      "name": "管理员",
      "code": "admin"
    }
  ]
}
```

## 错误码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或令牌无效 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突（如用户名已存在） |
| 500 | 服务器内部错误 |

## 错误响应格式

```json
{
  "error": "错误描述"
}
```

## 使用示例

### Go 客户端示例

```go
package main

import (
    "fmt"
    "log"
    "auth-sdk"
)

func main() {
    // 创建认证客户端
    client := auth.NewAuthClient("http://localhost:8080", "your-app-id", "your-app-secret")

    // 用户登录
    loginResp, err := client.Login(&auth.LoginRequest{
        Username: "username",
        Password: "password",
    })
    if err != nil {
        log.Fatal(err)
    }

    // 使用访问令牌
    userInfo, err := client.GetUserInfo(loginResp.AccessToken)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("用户信息: %+v\n", userInfo)
}
```

### cURL 示例

```bash
# 用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "app_id": "your-app-id",
    "username": "username",
    "password": "password"
  }'

# 获取用户信息
curl -X GET http://localhost:8080/api/v1/auth/user \
  -H "Authorization: Bearer <access_token>"

# 检查权限
curl -X GET "http://localhost:8080/api/v1/permissions/check?permission=user:read" \
  -H "Authorization: Bearer <access_token>"
```

## 部署说明

### Docker 部署

```bash
# 构建镜像
docker build -t auth-center .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=mysql \
  -e DB_PASSWORD=password \
  -e REDIS_HOST=redis \
  auth-center
```

### Docker Compose 部署

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f auth-center
```

## 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| SERVER_PORT | 服务端口 | 8080 |
| SERVER_MODE | 运行模式 | debug |
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 3306 |
| DB_USER | 数据库用户名 | root |
| DB_PASSWORD | 数据库密码 | - |
| DB_DATABASE | 数据库名 | auth_center |
| REDIS_HOST | Redis主机 | localhost |
| REDIS_PORT | Redis端口 | 6379 |
| REDIS_PASSWORD | Redis密码 | - |
| REDIS_DB | Redis数据库 | 0 |
| JWT_SECRET_KEY | JWT密钥 | - |
| JWT_TTL | 访问令牌有效期(秒) | 3600 |
| JWT_REFRESH_SECRET_KEY | 刷新令牌密钥 | - |
| JWT_REFRESH_TTL | 刷新令牌有效期(秒) | 7200 |

## 安全建议

1. **生产环境配置**：
   - 修改默认的JWT密钥
   - 使用强密码
   - 启用HTTPS
   - 配置防火墙规则

2. **令牌管理**：
   - 设置合理的令牌过期时间
   - 实现令牌刷新机制
   - 支持令牌撤销

3. **权限控制**：
   - 遵循最小权限原则
   - 定期审查权限配置
   - 实现权限审计日志

4. **监控告警**：
   - 监控认证失败次数
   - 监控异常访问模式
   - 设置安全告警
