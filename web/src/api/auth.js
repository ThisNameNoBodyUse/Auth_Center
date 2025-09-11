import api from './index'

// 应用用户登录
export const login = (data) => {
  return api.post('/auth/login', data)
}

// 应用用户登出
export const logout = (token) => {
  return api.post('/auth/logout', { token })
}

// 获取应用用户信息
export const getUserInfo = () => {
  return api.get('/auth/user')
}

// 刷新应用用户令牌
export const refreshToken = (refreshToken) => {
  return api.post('/auth/refresh', { refresh_token: refreshToken })
}

// 系统管理员登录
export const systemLogin = (data) => {
  return api.post('/system/login', data)
}

// 系统管理员登出
export const systemLogout = (token) => {
  return api.post('/system/logout', { token })
}

// 获取系统管理员信息
export const getSystemAdminInfo = () => {
  return api.get('/system/admin/info')
}

// 刷新系统管理员令牌
export const systemRefreshToken = (refreshToken) => {
  return api.post('/system/refresh', { refresh_token: refreshToken })
}
