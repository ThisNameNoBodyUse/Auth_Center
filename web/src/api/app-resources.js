import api from './index'

// 角色管理
export const getRoles = (params = {}) => {
  return api.get('/app/roles', { params })
}

export const createRole = (data) => {
  const params = data?.app_id ? { app_id: data.app_id } : undefined
  return api.post('/app/roles', data, { params })
}

export const updateRole = (id, data) => {
  const params = data?.app_id ? { app_id: data.app_id } : undefined
  return api.put(`/app/roles/${id}`, data, { params })
}

export const deleteRole = (id, appId) => {
  const params = appId ? { app_id: appId } : undefined
  return api.delete(`/app/roles/${id}`, { params })
}

export const assignRolePermissions = (id, data) => {
  const params = data?.app_id ? { app_id: data.app_id } : undefined
  return api.post(`/app/roles/${id}/permissions`, data, { params })
}

export const getRolePermissions = (id, appId) => {
  const params = appId ? { app_id: appId } : undefined
  return api.get(`/app/roles/${id}/permissions`, { params })
}

// 权限管理
export const getPermissions = (params = {}) => {
  return api.get('/app/permissions', { params })
}

export const createPermission = (data) => {
  const params = data?.app_id ? { app_id: data.app_id } : undefined
  return api.post('/app/permissions', data, { params })
}

export const updatePermission = (id, data) => {
  const params = data?.app_id ? { app_id: data.app_id } : undefined
  return api.put(`/app/permissions/${id}`, data, { params })
}

export const deletePermission = (id, appId) => {
  const params = appId ? { app_id: appId } : undefined
  return api.delete(`/app/permissions/${id}`, { params })
}

// 用户管理
export const getUsers = (params = {}) => {
  return api.get('/app/users', { params })
}

export const createUser = (data) => {
  const params = data?.app_id ? { app_id: data.app_id } : undefined
  return api.post('/app/users', data, { params })
}

export const updateUser = (id, data) => {
  const params = data?.app_id ? { app_id: data.app_id } : undefined
  return api.put(`/app/users/${id}`, data, { params })
}

export const deleteUser = (id, appId) => {
  const params = appId ? { app_id: appId } : undefined
  return api.delete(`/app/users/${id}`, { params })
}

export const assignUserRoles = (id, data) => {
  const params = data?.app_id ? { app_id: data.app_id } : undefined
  return api.post(`/app/users/${id}/roles`, data, { params })
}

export const getUserRoles = (id) => {
  return api.get(`/app/users/${id}/roles`)
}

// 当前应用信息（仅应用级管理员）
export const getSelfApp = () => {
  return api.get('/app/self').then((resp) => resp.data)
}
