# 认证授权中心

一个基于Go语言开发的微服务认证授权中心，提供统一的用户认证、权限管理和令牌服务。

## 特性

- 🔐 **统一认证**：支持多应用、多租户的用户认证
- 🛡️ **权限管理**：基于RBAC的细粒度权限控制
- 🎫 **令牌管理**：JWT访问令牌和刷新令牌机制
- 🚀 **高性能**：Redis缓存提升性能
- 📦 **微服务架构**：独立部署，易于扩展
- 🔌 **SDK支持**：提供多语言客户端SDK
- 📚 **完整文档**：详细的API文档和使用示例

## 架构设计

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   客户端应用1    │    │   客户端应用2    │    │   客户端应用N    │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │      认证授权中心          │
                    │  ┌─────────────────────┐  │
                    │  │   API Gateway       │  │
                    │  └─────────────────────┘  │
                    │  ┌─────────────────────┐  │
                    │  │   认证服务          │  │
                    │  └─────────────────────┘  │
                    │  ┌─────────────────────┐  │
                    │  │   权限服务          │  │
                    │  └─────────────────────┘  │
                    │  ┌─────────────────────┐  │
                    │  │   应用管理服务      │  │
                    │  └─────────────────────┘  │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │        数据层            │
                    │  ┌─────────┐ ┌─────────┐ │
                    │  │  MySQL  │ │  Redis  │ │
                    │  └─────────┘ └─────────┘ │
                    └───────────────────────────┘
```

## 快速开始

### 环境要求

- Go 1.24.2+
- MySQL 8.0+
- Redis 6.0+

### 本地开发

1. **克隆项目**
```bash
git clone <repository-url>
cd auth_center
```

2. **安装依赖**
```bash
go mod download
```

3. **配置数据库**
```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE auth_center CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 导入初始化脚本
mysql -u root -p auth_center < docker/init.sql
```

4. **配置Redis**
```bash
# 启动Redis
redis-server
```

5. **配置应用**
```bash
# 复制配置文件
cp config/app.ini.example config/app.ini

# 编辑配置文件
vim config/app.ini
```

6. **运行应用**
```bash
go run main.go
```

### Docker部署

1. **使用Docker Compose**
```bash
# 启动所有服务
docker-compose -f docker/docker-compose.yml up -d

# 查看日志
docker-compose -f docker/docker-compose.yml logs -f auth-center
```

2. **单独构建镜像**
```bash
# 构建镜像
docker build -f docker/Dockerfile -t auth-center .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=localhost \
  -e DB_PASSWORD=password \
  -e REDIS_HOST=localhost \
  auth-center
```

## 使用指南

### 1. 创建应用

首先需要创建一个应用来获取`app_id`和`app_secret`：

```bash
curl -X POST http://localhost:8080/api/v1/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "我的应用",
    "description": "应用描述"
  }'
```

### 2. 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "app_id": "your-app-id",
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "app_id": "your-app-id",
    "username": "testuser",
    "password": "password123"
  }'
```

### 4. 使用SDK

#### Go SDK

```go
package main

import (
    "fmt"
    "log"
    "auth-sdk"
)

func main() {
    // 创建客户端
    client := auth.NewAuthClient("http://localhost:8080", "your-app-id", "your-app-secret")
    
    // 用户登录
    resp, err := client.Login(&auth.LoginRequest{
        Username: "testuser",
        Password: "password123",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // 使用令牌
    userInfo, err := client.GetUserInfo(resp.AccessToken)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("用户信息: %+v\n", userInfo)
}
```

## API文档

详细的API文档请参考 [API.md](docs/API.md)

## 配置说明

### 配置文件

配置文件位于 `config/app.ini`：

```ini
[server]
port = 8080
mode = debug

[database]
host = localhost
port = 3306
user = root
password = 
database = auth_center
charset = utf8mb4

[redis]
host = localhost
port = 6379
password = 
db = 0

[jwt]
secret_key = your-secret-key
ttl = 3600
refresh_secret_key = your-refresh-secret-key
refresh_ttl = 7200
```

### 环境变量

支持通过环境变量覆盖配置：

```bash
export DB_HOST=mysql
export DB_PASSWORD=password
export REDIS_HOST=redis
export JWT_SECRET_KEY=your-secret-key
```

## 数据库设计

### 核心表结构

- **applications**: 应用表
- **users**: 用户表
- **roles**: 角色表
- **permissions**: 权限表
- **apis**: API接口表
- **user_roles**: 用户角色关联表
- **role_permissions**: 角色权限关联表
- **tokens**: 令牌表

详细表结构请参考 [docker/init.sql](docker/init.sql)

## 权限模型

### RBAC权限模型

- **用户(User)**: 系统中的具体用户
- **角色(Role)**: 权限的集合
- **权限(Permission)**: 具体的操作权限
- **资源(Resource)**: 被操作的对象
- **操作(Action)**: 对资源的操作类型

### 权限示例

```
用户: admin
├── 角色: 管理员
│   ├── 权限: user:manage
│   ├── 权限: user:read
│   ├── 权限: user:create
│   ├── 权限: user:update
│   └── 权限: user:delete
└── 角色: 系统管理员
    ├── 权限: system:manage
    └── 权限: app:manage
```

## 安全特性

### 密码安全

- 使用Argon2id算法哈希密码
- 支持密码强度验证
- 防止时序攻击

### 令牌安全

- JWT访问令牌
- 刷新令牌机制
- 令牌黑名单
- 令牌过期管理

### 权限安全

- 细粒度权限控制
- 权限缓存优化
- 权限审计日志

## 性能优化

### 缓存策略

- Redis缓存权限信息
- 缓存过期时间管理
- 缓存预热机制

### 数据库优化

- 索引优化
- 查询优化
- 连接池管理

## 监控告警

### 指标监控

- 认证成功率
- 权限检查性能
- 数据库连接状态
- Redis缓存命中率

### 日志管理

- 结构化日志
- 日志级别控制
- 日志轮转

## 扩展开发

### 添加新的认证方式

1. 实现认证接口
2. 添加认证中间件
3. 更新路由配置

### 添加新的权限类型

1. 扩展权限模型
2. 更新权限检查逻辑
3. 添加权限管理接口

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License

## 联系方式

- 项目地址: [GitHub Repository]
- 问题反馈: [GitHub Issues]
- 文档地址: [Documentation]

## 更新日志

### v1.0.0 (2024-01-01)

- 初始版本发布
- 支持用户认证和权限管理
- 提供Go SDK
- 支持Docker部署
