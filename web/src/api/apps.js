import api from './index'

// 获取应用列表
export const getApps = (params = {}) => {
  return api.get('/apps', { params }).then((resp) => {
    return {
      apps: resp?.data || [],
      total: resp?.pagination?.total || 0,
    }
  })
}

// 获取应用详情
export const getApp = (appId) => {
  return api.get(`/apps/${appId}`)
}

// 创建应用
export const createApp = (data) => {
  return api.post('/apps', data)
}

// 更新应用
export const updateApp = (appId, data) => {
  return api.put(`/apps/${appId}`, data)
}

// 删除应用
export const deleteApp = (appId) => {
  return api.delete(`/apps/${appId}`)
}

// 重新生成应用密钥
export const regenerateAppSecret = (appId) => {
  return api.post(`/apps/${appId}/regenerate-secret`)
}
