import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { 
  login, logout, getUserInfo, 
  systemLogin, systemLogout, getSystemAdminInfo 
} from '@/api/auth'
import Cookies from 'js-cookie'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(Cookies.get('access_token') || '')
  const user = ref(null)
  const appId = ref('system-admin')
  const appSecret = ref('system-admin-secret-key-change-in-production')
  const loginType = ref(Cookies.get('login_type') || 'system') // 'system' 或 'app'

  const isLoggedIn = computed(() => !!token.value)
  const isSystemAdmin = computed(() => loginType.value === 'system')

  // 系统管理员登录
  const systemLoginAction = async (credentials) => {
    try {
      const response = await systemLogin({
        username: credentials.username,
        password: credentials.password
      })
      
      token.value = response.access_token
      user.value = response.admin
      loginType.value = 'system'
      
      // 保存到cookie
      Cookies.set('access_token', token.value, { expires: 7 })
      Cookies.set('login_type', 'system', { expires: 7 })
      
      return response
    } catch (error) {
      throw error
    }
  }

  // 应用用户登录
  const appLoginAction = async (credentials) => {
    try {
      const response = await login({
        app_id: appId.value,
        app_secret: appSecret.value,
        username: credentials.username,
        password: credentials.password
      })
      
      token.value = response.access_token
      user.value = response.user
      loginType.value = 'app'
      
      // 保存到cookie
      Cookies.set('access_token', token.value, { expires: 7 })
      Cookies.set('login_type', 'app', { expires: 7 })
      
      return response
    } catch (error) {
      throw error
    }
  }

  // 通用登录方法（根据登录类型选择）
  const loginAction = async (credentials, type = 'system') => {
    if (type === 'system') {
      return await systemLoginAction(credentials)
    } else {
      return await appLoginAction(credentials)
    }
  }

  // 登出
  const logoutAction = async () => {
    try {
      if (token.value) {
        if (loginType.value === 'system') {
          await systemLogout(token.value)
        } else {
          await logout(token.value)
        }
      }
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      token.value = ''
      user.value = null
      loginType.value = 'system'
      Cookies.remove('access_token')
      Cookies.remove('login_type')
    }
  }

  // 检查认证状态
  const checkAuth = async () => {
    if (token.value) {
      try {
        let response
        if (loginType.value === 'system') {
          response = await getSystemAdminInfo()
        } else {
          response = await getUserInfo()
        }
        user.value = response
      } catch (error) {
        // token可能已过期，清除本地状态
        token.value = ''
        user.value = null
        loginType.value = 'system'
        Cookies.remove('access_token')
        Cookies.remove('login_type')
      }
    }
  }

  return {
    token,
    user,
    appId,
    appSecret,
    loginType,
    isLoggedIn,
    isSystemAdmin,
    loginAction,
    systemLoginAction,
    appLoginAction,
    logoutAction,
    checkAuth
  }
})
