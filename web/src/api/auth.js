import api from './index'

// 登录
export const login = (data) => {
  return api.post('/auth/login', data)
}

// 登出
export const logout = (token) => {
  return api.post('/auth/logout', { token })
}

// 获取用户信息
export const getUserInfo = () => {
  return api.get('/auth/user')
}

// 刷新令牌
export const refreshToken = (refreshToken) => {
  return api.post('/auth/refresh', { refresh_token: refreshToken })
}
