# Auth Center API 文档

## 概述

认证授权中心提供完整的用户认证、授权和权限管理功能。支持多级权限架构，包括系统级超级管理员和应用级超级管理员。

## 基础信息

- **Base URL**: `http://localhost:8080`
- **API Version**: v1
- **认证方式**: JWT Bearer Token
- **Content-Type**: `application/json`

## 权限说明

### 系统级超级管理员
- 可以管理所有应用
- 可以查看和管理所有应用的用户、角色、权限
- 可以管理应用密钥
- 通过 `app_id` 查询参数指定目标应用

### 应用级超级管理员
- 只能管理自己应用内的用户、角色、权限
- 不能管理应用密钥
- 不能访问其他应用的数据

## 认证流程

1. 用户登录获取访问令牌
2. 在请求头中携带 `Authorization: Bearer <token>`
3. 系统验证令牌有效性
4. 根据用户权限返回相应数据

## API 接口

### 1. 认证管理

#### 1.1 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "superadmin",
  "password": "admin123",
  "app_id": "system-admin",
  "app_secret": "system-admin-secret-key-change-in-production"
}
```

**响应示例:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "user": {
    "id": 1,
    "username": "superadmin",
    "email": "admin@auth-center.com",
    "is_super_admin": true
  }
}
```

#### 1.2 用户注册
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "newuser",
  "password": "password123",
  "email": "user@example.com",
  "phone": "13800138000",
  "app_id": "default-app",
  "app_secret": "default-secret"
}
```

#### 1.3 刷新令牌
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 1.4 用户登出
```http
POST /api/v1/auth/logout
Authorization: Bearer <access_token>
```

#### 1.5 获取当前用户信息
```http
GET /api/v1/auth/user
Authorization: Bearer <access_token>
```

### 2. 应用管理（仅系统级超级管理员）

#### 2.1 获取应用列表
```http
GET /api/v1/apps?page=1&page_size=10&name=&status=
Authorization: Bearer <access_token>
```

**查询参数:**
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）
- `name`: 应用名称（模糊搜索）
- `status`: 应用状态（0-禁用，1-启用）

#### 2.2 创建应用
```http
POST /api/v1/apps
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "新应用",
  "description": "应用描述"
}
```

#### 2.3 获取应用详情
```http
GET /api/v1/apps/{app_id}
Authorization: Bearer <access_token>
```

#### 2.4 更新应用
```http
PUT /api/v1/apps/{app_id}
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "更新后的应用名",
  "description": "更新后的描述",
  "status": 1
}
```

#### 2.5 删除应用
```http
DELETE /api/v1/apps/{app_id}
Authorization: Bearer <access_token>
```

#### 2.6 重新生成应用密钥
```http
POST /api/v1/apps/{app_id}/regenerate-secret
Authorization: Bearer <access_token>
```

#### 2.7 获取应用用户列表
```http
GET /api/v1/apps/{app_id}/users?page=1&page_size=10
Authorization: Bearer <access_token>
```

### 3. 应用资源管理（系统级和应用级超级管理员）

#### 3.1 角色管理

##### 获取角色列表
```http
GET /api/v1/app/roles?app_id=default-app&page=1&page_size=10
Authorization: Bearer <access_token>
```

**查询参数:**
- `app_id`: 应用ID（系统级超级管理员可指定，应用级超级管理员忽略）
- `page`: 页码
- `page_size`: 每页数量

##### 创建角色
```http
POST /api/v1/app/roles?app_id=default-app
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "管理员",
  "code": "admin",
  "description": "管理员角色"
}
```

##### 更新角色
```http
PUT /api/v1/app/roles/{id}?app_id=default-app
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "更新后的角色名",
  "code": "updated_admin",
  "description": "更新后的描述",
  "status": 1
}
```

##### 删除角色
```http
DELETE /api/v1/app/roles/{id}?app_id=default-app
Authorization: Bearer <access_token>
```

##### 获取角色权限
```http
GET /api/v1/app/roles/{id}/permissions?app_id=default-app
Authorization: Bearer <access_token>
```

##### 分配角色权限
```http
POST /api/v1/app/roles/{id}/permissions?app_id=default-app
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "permission_ids": [1, 2, 3, 4, 5]
}
```

#### 3.2 权限管理

##### 获取权限列表
```http
GET /api/v1/app/permissions?app_id=default-app&page=1&page_size=10
Authorization: Bearer <access_token>
```

##### 创建权限
```http
POST /api/v1/app/permissions?app_id=default-app
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "用户管理",
  "code": "user:manage",
  "resource": "user",
  "action": "manage",
  "description": "用户管理权限"
}
```

##### 更新权限
```http
PUT /api/v1/app/permissions/{id}?app_id=default-app
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "更新后的权限名",
  "code": "user:updated",
  "resource": "user",
  "action": "updated",
  "description": "更新后的描述",
  "status": 1
}
```

##### 删除权限
```http
DELETE /api/v1/app/permissions/{id}?app_id=default-app
Authorization: Bearer <access_token>
```

#### 3.3 用户管理

##### 获取用户列表
```http
GET /api/v1/app/users?app_id=default-app&page=1&page_size=10
Authorization: Bearer <access_token>
```

##### 创建用户
```http
POST /api/v1/app/users?app_id=default-app
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "username": "newuser",
  "password": "password123",
  "email": "user@example.com",
  "phone": "13800138000",
  "is_super_admin": false
}
```

##### 更新用户
```http
PUT /api/v1/app/users/{id}?app_id=default-app
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "username": "updateduser",
  "email": "updated@example.com",
  "phone": "13900139000",
  "password": "newpassword123",
  "is_super_admin": false,
  "status": 1
}
```

##### 删除用户
```http
DELETE /api/v1/app/users/{id}?app_id=default-app
Authorization: Bearer <access_token>
```

##### 获取用户角色
```http
GET /api/v1/app/users/{id}/roles?app_id=default-app
Authorization: Bearer <access_token>
```

##### 分配用户角色
```http
POST /api/v1/app/users/{id}/roles?app_id=default-app
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "role_ids": [1, 2]
}
```

## 错误码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 认证失败 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突（如用户名已存在） |
| 500 | 服务器内部错误 |

## 响应格式

### 成功响应
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 100,
    "total_page": 10
  }
}
```

### 错误响应
```json
{
  "error": "错误信息"
}
```

## 使用示例

### 1. 系统级超级管理员登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "superadmin",
    "password": "admin123",
    "app_id": "system-admin",
    "app_secret": "system-admin-secret-key-change-in-production"
  }'
```

### 2. 获取应用列表
```bash
curl -X GET "http://localhost:8080/api/v1/apps?page=1&page_size=10" \
  -H "Authorization: Bearer <access_token>"
```

### 3. 创建角色（在指定应用中）
```bash
curl -X POST "http://localhost:8080/api/v1/app/roles?app_id=default-app" \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "编辑者",
    "code": "editor",
    "description": "内容编辑者角色"
  }'
```

## 注意事项

1. 所有需要认证的接口都需要在请求头中携带 `Authorization: Bearer <token>`
2. 系统级超级管理员可以通过 `app_id` 查询参数指定目标应用
3. 应用级超级管理员只能操作自己应用内的资源
4. 密码使用 bcrypt 加密存储
5. 令牌有过期时间，需要定期刷新
6. 登出后令牌会被加入黑名单，无法再次使用

## 导入 Apifox

可以将 `docs/api-documentation.json` 文件直接导入到 Apifox 中：

1. 打开 Apifox
2. 选择 "导入" -> "OpenAPI"
3. 选择 `api-documentation.json` 文件
4. 确认导入

导入后即可在 Apifox 中查看和测试所有 API 接口。
