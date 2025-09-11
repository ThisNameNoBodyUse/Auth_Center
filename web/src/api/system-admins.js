import api from './index'

// 系统管理员管理
export const getSystemAdmins = (params = {}) => {
  return api.get('/system-admins', { params })
}

export const createSystemAdmin = (data) => {
  return api.post('/system-admins', data)
}

export const updateSystemAdmin = (id, data) => {
  return api.put(`/system-admins/${id}`, data)
}

export const deleteSystemAdmin = (id) => {
  return api.delete(`/system-admins/${id}`)
}

export const resetSystemAdminPassword = (id, data) => {
  return api.post(`/system-admins/${id}/reset-password`, data)
}
