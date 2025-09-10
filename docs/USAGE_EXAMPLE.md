# 认证授权中心使用示例（从零到一全流程）

> 目标：用一个“用户中心”应用为例，按真实使用顺序演示如何从创建应用、注册登录、角色权限与 API 管理、到权限校验、刷新令牌、登出等完整流程。所有请求均给出 cURL 示例与期望响应，便于直接照搬执行。

## 前置准备
- 服务地址：`http://localhost:8080`，基础前缀：`/api/v1`
- 内容类型：`application/json`
- 认证方式：`Authorization: Bearer <access_token>`（需要时）

### 通用请求头与参数（重要）
- 用户身份认证：
  - `Authorization: Bearer <access_token>`（访问受保护资源、权限校验等需要）
- 应用级鉴权（S2S，可选，启用后才必填）：
  - `X-App-Id: <app_id>`
  - `X-App-Secret: <app_secret>`
  - 说明：当某些路由组启用了应用级中间件（`AppAuthMiddleware`）时，必须携带上述头部；当前仓库已对 `/api/v1/apps`、`/api/v1/admin` 启用该中间件。

### 登录方式（providers）
- 每个应用可以配置登录方式（`providers` 表）：
  - `login_method=0`：账号密码登录（默认）
  - `login_method=1`：手机验证码登录
- 账号密码登录时，需要提交 `username/password`；短信登录时，需要提交 `phone/code`。

## 步骤 1：创建应用（获取 app_id/app_secret）
用于区分多应用/多租户，接入方需要先创建应用并获取 `app_id` 与 `app_secret`。

请求：
```bash
curl -X POST http://localhost:8080/api/v1/apps \
  -H "Content-Type: application/json" \
  -H "X-App-Id: app_12345678" \
  -H "X-App-Secret: secret_c4e31b98-..." \
  -d '{
    "name": "用户中心",
    "description": "公司内部用户中心"
  }'
```

期望响应（示例）：
```json
{
  "id": 1,
  "name": "用户中心",
  "app_id": "app_12345678",
  "app_secret": "secret_c4e31b98-...",
  "description": "公司内部用户中心",
  "status": 1,
  "created_at": "2024-01-01T00:00:00Z"
}
```
记录：`app_id=app_12345678`，`app_secret=secret_c4e31b98-...`

> 说明：生产环境建议对 `/apps` 系列接口增加管理员鉴权或使用 `AppAuthMiddleware`（本仓库已启用）。

## 步骤 2：初始化角色/权限/API（一次性）
你可以用 SQL（见 `docker/init.sql`）或后台管理接口（如后续扩展）来初始化。这里演示插入的目标状态：
- 角色：`admin`（管理员），`user`（普通用户）
- 权限：`user:read`、`user:create`、`user:update`、`user:delete` 等
- API 绑定（示例）：
  - GET `/api/v1/users` -> `user:read`
  - POST `/api/v1/users` -> `user:create`
  - PUT `/api/v1/users/:id` -> `user:update`
  - DELETE `/api/v1/users/:id` -> `user:delete`

> 路径规范务必与运行时校验一致，建议参数位使用 `:id` 统一格式。

## 步骤 3：注册首个用户（应用内）
请求：
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "app_id": "app_12345678",
    "username": "alice",
    "email": "alice@example.com",
    "password": "P@ssw0rd"
  }'
```

期望响应：
```json
{ "message": "注册成功" }
```

> 说明：默认策略是知道 `app_id` 即可注册，生产中可改为仅管理员创建或加验证码。

## 步骤 4：分配角色给用户（user_roles）
将用户与角色建立绑定关系（通常通过管理后台或 SQL 完成）。示例 SQL（请根据实际自增ID调整）：
```sql
-- 假设 alice 的 user_id=1，admin 角色的 role_id=1
INSERT INTO user_roles (user_id, role_id, app_id) VALUES (1, 1, 'app_12345678');
```

> 变更角色/权限后，请失效相关 Redis 缓存（用户权限、角色权限、API 权限集合），确保新权限即时生效。

## 步骤 5A：账号密码登录（login_method=0）
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "app_id": "app_12345678",
    "username": "alice",
    "password": "P@ssw0rd"
  }'
```

## 步骤 5B：短信验证码登录（login_method=1）
- 先通过你自己的短信服务生成并下发验证码（中心示例从 Redis 读取 key: `otp:<app_id>:<phone>`）。
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "app_id": "app_12345678",
    "phone": "13800001111",
    "code": "123456"
  }'
```

成功登录示例响应（两种方式相同）：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600,
  "token_type": "Bearer",
  "user": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "avatar": "",
    "roles": [
      {"id": 1, "name": "管理员", "code": "admin"}
    ]
  }
}
```

将 `access_token` 保存到客户端（短期），`refresh_token` 安全保存（长期）。

## 步骤 6：访问受保护资源（携带令牌）
请求（示例：获取当前用户信息）：
```bash
curl -X GET http://localhost:8080/api/v1/auth/user \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

期望响应（示例）：
```json
{
  "id": 1,
  "username": "alice",
  "email": "alice@example.com",
  "avatar": "",
  "roles": [
    {"id": 1, "name": "管理员", "code": "admin"}
  ]
}
```

## 步骤 7：按权限点校验（permission code）
请求：
```bash
curl -X GET "http://localhost:8080/api/v1/permissions/check?permission=user:read" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

期望响应：
```json
{ "has_permission": true }
```

> 用于在前后端按权限点（如 `user:read`）精细控制按钮/菜单/功能。

## 步骤 8：按 API 校验（path + method）
请求：
```bash
curl -X GET "http://localhost:8080/api/v1/permissions/check-api?path=/api/v1/users&method=GET" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

期望响应：
```json
{ "has_permission": true }
```

> 适合在网关层或后端统一拦截，对应 `apis` 表维护的受控接口。

## 步骤 9：刷新令牌（access 过期后）
请求：
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<REFRESH_TOKEN>"
  }'
```

期望响应（返回新的一对令牌）：
```json
{
  "access_token": "...",
  "refresh_token": "...",
  "expires_in": 3600,
  "token_type": "Bearer",
  "user": { ... }
}
```

## 步骤 10：登出（使当前 access 失效）
请求：
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "token": "<ACCESS_TOKEN>"
  }'
```

期望响应：
```json
{ "message": "登出成功" }
```

> 该操作会将 access token 的 `jti` 加入黑名单；如需彻底退出，建议对 refresh token 也实施吊销策略。

## 附录：常见注意事项
- 多租户隔离：用户/角色/权限/API 索引建议以 `(app_id, ...)` 作为组合唯一；查询需带 `app_id` 过滤。
- 路径一致性：上报与校验时，对路径做一致的规范化（如参数统一为 `:id`）。
- 缓存一致性：当 `user_roles / role_permissions / apis` 发生变动时，请失效对应 Redis 集合缓存。
- 安全：生产限制 CORS 来源，使用强 JWT 密钥与合理 TTL，`/apps` 接口加管理员保护，`/register` 可按需改为邀请制或验证码。
