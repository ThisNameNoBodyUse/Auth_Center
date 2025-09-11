# 前端更新说明

## 概述

为了支持新的系统管理员认证架构，前端需要进行以下更新：

## 主要更新内容

### 1. 认证API更新 (`/web/src/api/auth.js`)

- **新增系统级超级管理员登录接口**：
  - `systemLogin()` - 系统级超级管理员登录
  - `systemLogout()` - 系统级超级管理员登出
  - `getSystemAdminInfo()` - 获取系统级超级管理员信息
  - `systemRefreshToken()` - 刷新系统级超级管理员令牌

- **保留应用级超级管理员登录接口**：
  - `login()` - 应用级超级管理员登录（已更新支持 `app_secret`）
  - `logout()` - 应用级超级管理员登出
  - `getUserInfo()` - 获取应用级超级管理员信息
  - `refreshToken()` - 刷新应用级超级管理员令牌

### 2. 认证Store更新 (`/web/src/stores/auth.js`)

- **新增状态管理**：
  - `loginType` - 登录类型（'system' 或 'app'）
  - `isSystemAdmin` - 是否为系统管理员

- **新增登录方法**：
  - `systemLoginAction()` - 系统级超级管理员登录
  - `appLoginAction()` - 应用级超级管理员登录
  - `loginAction(credentials, type)` - 通用登录方法

- **更新认证逻辑**：
  - 根据登录类型选择不同的API接口
  - 保存登录类型到Cookie
  - 根据登录类型获取用户信息

### 3. API拦截器更新 (`/web/src/api/index.js`)

- **智能认证头添加**：
  - 系统级超级管理员：只添加 `Authorization` 头
  - 应用级超级管理员：添加 `Authorization`、`X-App-Id`、`X-App-Secret` 头

### 4. 登录页面更新 (`/web/src/views/Login.vue`)

- **新增登录类型选择器**：
  - 系统级超级管理员登录
  - 应用级超级管理员登录

- **动态表单和提示**：
  - 根据登录类型显示不同的默认用户名
  - 根据登录类型显示不同的提示信息

### 5. 新增应用资源管理API (`/web/src/api/app-resources.js`)

- **角色管理**：
  - `getRoles()` - 获取角色列表
  - `createRole()` - 创建角色
  - `updateRole()` - 更新角色
  - `deleteRole()` - 删除角色
  - `assignRolePermissions()` - 分配角色权限
  - `getRolePermissions()` - 获取角色权限

- **权限管理**：
  - `getPermissions()` - 获取权限列表
  - `createPermission()` - 创建权限
  - `updatePermission()` - 更新权限
  - `deletePermission()` - 删除权限

- **用户管理**：
  - `getUsers()` - 获取用户列表
  - `createUser()` - 创建用户
  - `updateUser()` - 更新用户
  - `deleteUser()` - 删除用户
  - `assignUserRoles()` - 分配用户角色
  - `getUserRoles()` - 获取用户角色

## 使用说明

### 系统级超级管理员登录

```javascript
// 使用系统级超级管理员登录
await authStore.systemLoginAction({
  username: 'superadmin',
  password: 'admin123'
})
```

### 应用级超级管理员登录

```javascript
// 使用应用级超级管理员登录
await authStore.appLoginAction({
  username: 'appadmin',
  password: 'admin123'
})
```

### 通用登录

```javascript
// 根据类型选择登录方式
await authStore.loginAction(credentials, 'system') // 系统级超级管理员
await authStore.loginAction(credentials, 'app')    // 应用级超级管理员
```

## 默认账户

### 系统级超级管理员
- 用户名：`superadmin`
- 密码：`admin123`

### 应用级超级管理员
- 用户名：`appadmin`
- 密码：`admin123`

## 注意事项

1. **生产环境**：请务必修改默认密码
2. **Cookie管理**：登录类型会保存到Cookie中，确保页面刷新后状态正确
3. **API调用**：系统级超级管理员和应用级超级管理员使用不同的API端点
4. **权限控制**：前端会根据登录类型显示不同的功能模块

## 兼容性

- 现有的应用级超级管理员登录功能完全兼容
- 新增的系统级超级管理员功能不影响现有功能
- 支持两种登录模式的无缝切换
