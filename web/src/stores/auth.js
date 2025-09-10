import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, logout, getUserInfo } from '@/api/auth'
import Cookies from 'js-cookie'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(Cookies.get('access_token') || '')
  const user = ref(null)
  const appId = ref('system-admin')
  const appSecret = ref('system-admin-secret-key-change-in-production')

  const isLoggedIn = computed(() => !!token.value)

  // 登录
  const loginAction = async (credentials) => {
    try {
      const response = await login({
        app_id: appId.value,
        username: credentials.username,
        password: credentials.password
      })
      
      token.value = response.access_token
      user.value = response.user
      
      // 保存token到cookie
      Cookies.set('access_token', token.value, { expires: 7 })
      
      return response
    } catch (error) {
      throw error
    }
  }

  // 登出
  const logoutAction = async () => {
    try {
      if (token.value) {
        await logout(token.value)
      }
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      token.value = ''
      user.value = null
      Cookies.remove('access_token')
    }
  }

  // 检查认证状态
  const checkAuth = async () => {
    if (token.value) {
      try {
        const response = await getUserInfo()
        user.value = response
      } catch (error) {
        // token可能已过期，清除本地状态
        token.value = ''
        user.value = null
        Cookies.remove('access_token')
      }
    }
  }

  return {
    token,
    user,
    appId,
    appSecret,
    isLoggedIn,
    loginAction,
    logoutAction,
    checkAuth
  }
})
