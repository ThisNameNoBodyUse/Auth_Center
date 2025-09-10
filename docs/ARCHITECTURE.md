# 认证授权中心架构设计

## 项目结构

```
auth_center/
├── main.go                 # 应用入口
├── go.mod                  # Go模块文件
├── go.sum                  # 依赖版本锁定
├── Makefile               # 构建脚本
├── README.md              # 项目说明
├── config/                # 配置文件
│   ├── app.ini           # 应用配置
│   ├── app.ini.example   # 配置示例
│   └── config.go         # 配置管理
├── models/                # 数据模型
│   └── models.go         # 数据库模型定义
├── controllers/           # 控制器层
│   ├── auth_controller.go    # 认证控制器
│   ├── app_controller.go     # 应用管理控制器
│   └── permission_controller.go # 权限控制器
├── service/               # 业务逻辑层
│   ├── auth_service.go       # 认证服务
│   ├── app_service.go        # 应用管理服务
│   └── permission_service.go # 权限服务
├── middleware/            # 中间件
│   └── auth.go           # 认证中间件
├── utils/                 # 工具函数
│   ├── jwt.go            # JWT工具
│   ├── redis.go          # Redis工具
│   ├── password.go       # 密码工具
│   └── swagger.go        # Swagger工具
├── routers/               # 路由配置
│   └── routes.go         # 路由定义
├── sdk/                   # 客户端SDK
│   └── go/               # Go SDK
│       ├── go.mod        # SDK模块文件
│       ├── auth_client.go # 认证客户端
│       └── example/      # 使用示例
│           └── main.go   # 示例代码
├── docker/                # Docker配置
│   ├── Dockerfile        # Docker镜像构建
│   ├── docker-compose.yml # Docker Compose配置
│   └── init.sql          # 数据库初始化脚本
├── docs/                  # 文档
│   ├── API.md            # API文档
│   ├── ARCHITECTURE.md   # 架构文档
│   └── swagger/          # Swagger文档
├── scripts/               # 脚本
│   ├── start.sh          # Linux启动脚本
│   └── start.bat         # Windows启动脚本
├── test/                  # 测试文件
│   └── auth_test.go      # 认证测试
└── bin/                   # 构建输出目录
```

## 架构层次

### 1. 表现层 (Presentation Layer)
- **控制器 (Controllers)**: 处理HTTP请求和响应
- **中间件 (Middleware)**: 请求预处理和后处理
- **路由 (Routers)**: URL路由配置

### 2. 业务逻辑层 (Business Logic Layer)
- **服务 (Services)**: 核心业务逻辑实现
- **模型 (Models)**: 数据模型定义

### 3. 数据访问层 (Data Access Layer)
- **数据库**: MySQL存储持久化数据
- **缓存**: Redis存储临时数据和缓存

### 4. 工具层 (Utility Layer)
- **JWT工具**: 令牌生成和验证
- **密码工具**: 密码哈希和验证
- **Redis工具**: 缓存操作封装

## 核心组件

### 1. 认证服务 (Auth Service)
- 用户注册/登录
- 令牌生成和验证
- 用户信息管理

### 2. 权限服务 (Permission Service)
- 权限检查
- 角色管理
- API权限控制

### 3. 应用管理服务 (App Service)
- 应用注册
- 应用配置管理
- 应用凭据管理

### 4. 中间件 (Middleware)
- JWT认证中间件
- 权限检查中间件
- 应用认证中间件

## 数据流

### 1. 用户认证流程
```
客户端 → 认证控制器 → 认证服务 → 数据库/Redis → 返回令牌
```

### 2. 权限检查流程
```
请求 → 认证中间件 → 权限中间件 → 权限服务 → Redis缓存 → 返回结果
```

### 3. 令牌验证流程
```
请求头 → JWT解析 → 黑名单检查 → 权限验证 → 放行/拒绝
```

## 安全设计

### 1. 密码安全
- 使用Argon2id算法哈希密码
- 防止时序攻击
- 密码强度验证

### 2. 令牌安全
- JWT访问令牌
- 刷新令牌机制
- 令牌黑名单
- 令牌过期管理

### 3. 权限安全
- 基于RBAC的权限模型
- 细粒度权限控制
- 权限缓存优化

## 性能优化

### 1. 缓存策略
- Redis缓存权限信息
- 缓存过期时间管理
- 缓存预热机制

### 2. 数据库优化
- 索引优化
- 查询优化
- 连接池管理

### 3. 并发处理
- Goroutine并发处理
- 连接池复用
- 异步处理

## 扩展性设计

### 1. 微服务架构
- 独立部署
- 服务发现
- 负载均衡

### 2. 多租户支持
- 应用隔离
- 数据隔离
- 配置隔离

### 3. 插件化设计
- 认证方式插件
- 权限模型插件
- 存储后端插件

## 监控和运维

### 1. 日志管理
- 结构化日志
- 日志级别控制
- 日志轮转

### 2. 指标监控
- 认证成功率
- 权限检查性能
- 数据库连接状态
- Redis缓存命中率

### 3. 健康检查
- 服务健康检查
- 数据库连接检查
- Redis连接检查

## 部署架构

### 1. 单机部署
```
应用 → MySQL → Redis
```

### 2. 集群部署
```
负载均衡器 → 应用集群 → MySQL主从 → Redis集群
```

### 3. 容器化部署
```
Docker Compose → 应用容器 → MySQL容器 → Redis容器
```

## 开发规范

### 1. 代码规范
- Go标准格式
- 错误处理规范
- 注释规范

### 2. 测试规范
- 单元测试
- 集成测试
- 性能测试

### 3. 文档规范
- API文档
- 架构文档
- 部署文档
