# 认证授权中心超级管理员使用指南

## 概述

认证授权中心采用**预创建应用模式**，只有内置的超级管理员才能创建和管理应用。具体的应用不需要自己注册，只需要使用下发的凭据进行用户认证等操作。

## 架构设计

```
认证授权中心
├── 超级管理员 (system-admin)
│   ├── 创建应用 (app_id, app_secret)
│   ├── 管理应用信息
│   └── 下发应用凭据
└── 具体应用
    ├── 使用下发的凭据
    ├── 用户认证
    ├── 权限验证
    └── 用户管理
```

## 初始化超级管理员

### 1. 执行初始化脚本

```bash
# 执行超级管理员初始化脚本
mysql -u root -p auth_center < scripts/init_super_admin.sql
```

### 2. 默认凭据

- **应用ID**: `system-admin`
- **应用密钥**: `system-admin-secret-key-change-in-production`
- **用户名**: `superadmin`
- **密码**: `admin123`

⚠️ **重要**: 生产环境中必须修改默认密码和应用密钥！

## 使用流程

### 1. 超级管理员登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-App-Id: system-admin" \
  -H "X-App-Secret: system-admin-secret-key-change-in-production" \
  -d '{
    "app_id": "system-admin",
    "username": "superadmin",
    "password": "admin123"
  }'
```

### 2. 创建应用

```bash
curl -X POST http://localhost:8080/api/v1/apps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "X-App-Id: system-admin" \
  -H "X-App-Secret: system-admin-secret-key-change-in-production" \
  -d '{
    "name": "用户管理系统",
    "description": "企业内部用户管理系统"
  }'
```

**响应示例**:
```json
{
  "message": "应用创建成功",
  "app": {
    "id": 2,
    "name": "用户管理系统",
    "app_id": "user-mgmt-001",
    "app_secret": "sk_1234567890abcdef",
    "description": "企业内部用户管理系统",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 3. 下发应用凭据

将创建的应用凭据下发给具体的应用：

```json
{
  "app_id": "user-mgmt-001",
  "app_secret": "sk_1234567890abcdef",
  "name": "用户管理系统",
  "description": "企业内部用户管理系统"
}
```

### 4. 具体应用使用凭据

具体应用使用下发的凭据进行用户认证：

```bash
# 用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-App-Id: user-mgmt-001" \
  -H "X-App-Secret: sk_1234567890abcdef" \
  -d '{
    "app_id": "user-mgmt-001",
    "username": "user1",
    "password": "password123"
  }'

# 用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "X-App-Id: user-mgmt-001" \
  -H "X-App-Secret: sk_1234567890abcdef" \
  -d '{
    "app_id": "user-mgmt-001",
    "username": "user1",
    "email": "user1@example.com",
    "password": "password123"
  }'
```

## 管理接口

### 应用管理

| 接口 | 方法 | 描述 | 权限 |
|------|------|------|------|
| `/api/v1/apps` | POST | 创建应用 | 超级管理员 |
| `/api/v1/apps` | GET | 获取应用列表 | 超级管理员 |
| `/api/v1/apps/:app_id` | GET | 获取应用详情 | 超级管理员 |
| `/api/v1/apps/:app_id` | PUT | 更新应用 | 超级管理员 |
| `/api/v1/apps/:app_id` | DELETE | 删除应用 | 超级管理员 |
| `/api/v1/apps/:app_id/regenerate-secret` | POST | 重新生成密钥 | 超级管理员 |

### 用户认证（具体应用使用）

| 接口 | 方法 | 描述 | 权限 |
|------|------|------|------|
| `/api/v1/auth/login` | POST | 用户登录 | 应用认证 |
| `/api/v1/auth/register` | POST | 用户注册 | 应用认证 |
| `/api/v1/auth/refresh` | POST | 刷新令牌 | 应用认证 |
| `/api/v1/auth/logout` | POST | 用户登出 | 应用认证 |
| `/api/v1/auth/user` | GET | 获取用户信息 | 用户认证 |

### 权限验证（具体应用使用）

| 接口 | 方法 | 描述 | 权限 |
|------|------|------|------|
| `/api/v1/permissions/check` | GET | 检查权限 | 用户认证 |
| `/api/v1/permissions/check-api` | GET | 检查API权限 | 用户认证 |
| `/api/v1/permissions/user` | GET | 获取用户权限 | 用户认证 |
| `/api/v1/permissions/roles` | GET | 获取用户角色 | 用户认证 |

## 安全特性

### 1. 权限隔离
- 超级管理员只能管理应用，不能直接操作用户数据
- 具体应用只能操作用户认证和权限验证
- 应用间数据完全隔离

### 2. 认证层级
```
超级管理员认证
├── 用户认证 (JWT Token)
├── 应用认证 (X-App-Id + X-App-Secret)
└── 超级管理员验证 (is_super_admin = true)

具体应用认证
├── 用户认证 (JWT Token)
└── 应用认证 (X-App-Id + X-App-Secret)
```

### 3. 密钥管理
- 应用密钥由超级管理员生成和管理
- 支持密钥重新生成
- 密钥泄露时可立即重新生成

## 部署建议

### 1. 环境隔离
- 超级管理员接口限制在内网访问
- 具体应用接口可对外开放
- 使用防火墙限制管理接口访问

### 2. 监控审计
- 记录所有超级管理员操作日志
- 监控应用创建和密钥生成
- 定期审计应用使用情况

### 3. 备份恢复
- 定期备份数据库
- 保存应用密钥备份
- 建立灾难恢复流程

## 故障排除

### 1. 无法创建应用
- 检查是否为超级管理员登录
- 检查应用ID是否已存在
- 检查数据库连接是否正常

### 2. 应用无法认证
- 检查应用ID和密钥是否正确
- 检查应用状态是否启用
- 检查网络连接是否正常

### 3. 用户认证失败
- 检查用户是否存在
- 检查密码是否正确
- 检查用户状态是否正常
