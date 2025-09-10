# 认证授权中心管理后台

基于 Vue 3 + Element Plus 构建的现代化管理后台界面。

## 功能特性

- 🎨 现代化 UI 设计，基于 Element Plus
- 📱 响应式布局，支持移动端
- 🔐 完整的认证授权流程
- 📊 数据可视化仪表盘
- 🛠️ 应用管理
- 👥 用户管理
- 🎭 角色管理
- 🔑 权限管理
- ⚙️ 系统设置

## 技术栈

- **框架**: Vue 3
- **构建工具**: Vite
- **UI 组件库**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP 客户端**: Axios
- **样式**: SCSS
- **图标**: Element Plus Icons

## 项目结构

```
web/
├── public/                 # 静态资源
├── src/
│   ├── api/               # API 接口
│   │   ├── index.js       # Axios 配置
│   │   ├── auth.js        # 认证相关 API
│   │   └── apps.js        # 应用管理 API
│   ├── components/        # 公共组件
│   ├── layouts/           # 布局组件
│   │   └── MainLayout.vue # 主布局
│   ├── router/            # 路由配置
│   │   └── index.js       # 路由定义
│   ├── stores/            # 状态管理
│   │   └── auth.js        # 认证状态
│   ├── styles/            # 样式文件
│   │   └── index.scss     # 全局样式
│   ├── views/             # 页面组件
│   │   ├── Login.vue      # 登录页
│   │   ├── Dashboard.vue  # 仪表盘
│   │   ├── Apps.vue       # 应用管理
│   │   ├── Users.vue      # 用户管理
│   │   ├── Roles.vue      # 角色管理
│   │   ├── Permissions.vue # 权限管理
│   │   └── Settings.vue   # 系统设置
│   ├── App.vue            # 根组件
│   └── main.js            # 入口文件
├── index.html             # HTML 模板
├── package.json           # 依赖配置
├── vite.config.js         # Vite 配置
└── README.md              # 项目说明
```

## 快速开始

### 环境要求

- Node.js >= 16.0.0
- npm >= 7.0.0

### 安装依赖

```bash
cd web
npm install
```

### 开发环境

```bash
npm run dev
```

访问 http://localhost:3000

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 配置说明

### 代理配置

开发环境下，API 请求会自动代理到后端服务：

```javascript
// vite.config.js
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true
    }
  }
}
```

### 环境变量

可以通过环境变量配置不同的后端地址：

```bash
# .env.development
VITE_API_BASE_URL=http://localhost:8080

# .env.production
VITE_API_BASE_URL=https://api.example.com
```

## 功能说明

### 1. 登录认证

- 支持超级管理员登录
- 自动保存登录状态
- 路由守卫保护

### 2. 仪表盘

- 系统统计信息
- 应用使用情况
- 最近创建的应用

### 3. 应用管理

- 创建/编辑/删除应用
- 应用密钥管理
- 应用状态控制

### 4. 用户管理

- 用户列表查看
- 用户信息编辑
- 角色分配管理

### 5. 角色管理

- 角色创建/编辑/删除
- 权限分配管理
- 角色状态控制

### 6. 权限管理

- 权限点管理
- API 权限关联
- 权限层级管理

### 7. 系统设置

- 基本设置
- 安全设置
- 数据库配置
- Redis 配置
- 日志设置

## API 接口

### 认证接口

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户登出
- `GET /api/v1/auth/user` - 获取用户信息

### 应用管理接口

- `GET /api/v1/apps` - 获取应用列表
- `POST /api/v1/apps` - 创建应用
- `GET /api/v1/apps/:id` - 获取应用详情
- `PUT /api/v1/apps/:id` - 更新应用
- `DELETE /api/v1/apps/:id` - 删除应用
- `POST /api/v1/apps/:id/regenerate-secret` - 重新生成密钥

## 开发指南

### 添加新页面

1. 在 `src/views/` 目录下创建新的 Vue 组件
2. 在 `src/router/index.js` 中添加路由配置
3. 在 `src/layouts/MainLayout.vue` 中添加导航菜单

### 添加新 API

1. 在 `src/api/` 目录下创建对应的 API 文件
2. 在组件中导入并使用 API 方法

### 状态管理

使用 Pinia 进行状态管理，在 `src/stores/` 目录下定义 store：

```javascript
import { defineStore } from 'pinia'

export const useExampleStore = defineStore('example', () => {
  const state = ref({})
  
  const action = () => {
    // 状态操作
  }
  
  return { state, action }
})
```

## 部署说明

### 构建部署

1. 执行构建命令：
   ```bash
   npm run build
   ```

2. 将 `dist` 目录下的文件部署到 Web 服务器

3. 配置 Nginx 代理 API 请求到后端服务

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # 前端静态文件
    location / {
        root /path/to/dist;
        try_files $uri $uri/ /index.html;
    }
    
    # API 代理
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88

## 许可证

MIT License
