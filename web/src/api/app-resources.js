import api from './index'

// 角色管理
export const getRoles = (params = {}) => {
  return api.get('/app/roles', { params })
}

export const createRole = (data) => {
  return api.post('/app/roles', data)
}

export const updateRole = (id, data) => {
  return api.put(`/app/roles/${id}`, data)
}

export const deleteRole = (id) => {
  return api.delete(`/app/roles/${id}`)
}

export const assignRolePermissions = (id, data) => {
  return api.post(`/app/roles/${id}/permissions`, data)
}

export const getRolePermissions = (id) => {
  return api.get(`/app/roles/${id}/permissions`)
}

// 权限管理
export const getPermissions = (params = {}) => {
  return api.get('/app/permissions', { params })
}

export const createPermission = (data) => {
  return api.post('/app/permissions', data)
}

export const updatePermission = (id, data) => {
  return api.put(`/app/permissions/${id}`, data)
}

export const deletePermission = (id) => {
  return api.delete(`/app/permissions/${id}`)
}

// 用户管理
export const getUsers = (params = {}) => {
  return api.get('/app/users', { params })
}

export const createUser = (data) => {
  return api.post('/app/users', data)
}

export const updateUser = (id, data) => {
  return api.put(`/app/users/${id}`, data)
}

export const deleteUser = (id) => {
  return api.delete(`/app/users/${id}`)
}

export const assignUserRoles = (id, data) => {
  return api.post(`/app/users/${id}/roles`, data)
}

export const getUserRoles = (id) => {
  return api.get(`/app/users/${id}/roles`)
}
