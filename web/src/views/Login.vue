<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <img src="/logo.svg" alt="Logo" class="logo" />
        <h1>认证授权中心</h1>
        <p>{{ loginType === 'system' ? '系统级超级管理员登录' : '应用级超级管理员登录' }}</p>
      </div>
      
      <div class="login-type-selector">
        <el-radio-group v-model="loginType" @change="handleLoginTypeChange">
          <el-radio-button label="system">系统级超级管理员</el-radio-button>
          <el-radio-button label="app">应用级超级管理员</el-radio-button>
        </el-radio-group>
      </div>
      
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @submit.prevent="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            size="large"
            prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleLogin"
            class="login-btn"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="login-footer">
        <template v-if="loginType === 'system'">
          <p>系统级超级管理员默认用户名: superadmin</p>
          <p>默认密码: admin123</p>
        </template>
        <template v-else>
          <p>应用级超级管理员默认用户名: appadmin</p>
          <p>默认密码: admin123</p>
        </template>
        <p class="warning">⚠️ 生产环境请修改默认密码</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()

const loginFormRef = ref()
const loading = ref(false)

const loginType = ref('system')

const loginForm = reactive({
  username: 'superadmin',
  password: 'admin123'
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

const handleLoginTypeChange = () => {
  // 切换登录类型时重置表单
  loginForm.username = loginType.value === 'system' ? 'superadmin' : 'appadmin'
  loginForm.password = 'admin123'
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    loading.value = true
    
    await authStore.loginAction(loginForm, loginType.value)
    ElMessage.success('登录成功')
    router.push('/')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-box {
  width: 400px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.logo {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
}

.login-header h1 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 24px;
  font-weight: 600;
}

.login-header p {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.login-type-selector {
  margin-bottom: 20px;
  text-align: center;
}

.login-form {
  margin-bottom: 20px;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
}

.login-footer {
  text-align: center;
  color: #666;
  font-size: 12px;
  line-height: 1.6;
}

.login-footer p {
  margin: 4px 0;
}

.warning {
  color: #ff6b6b;
  font-weight: 500;
}
</style>
