<template>
  <div class="settings-page">
    <el-card class="page-card">
      <template #header>
        <span>系统设置</span>
      </template>

      <el-tabs v-model="activeTab" class="settings-tabs">
        <!-- 基本设置 -->
        <el-tab-pane label="基本设置" name="basic">
          <el-form
            ref="basicFormRef"
            :model="basicForm"
            :rules="basicRules"
            label-width="120px"
            class="settings-form"
          >
            <el-form-item label="系统名称" prop="systemName">
              <el-input v-model="basicForm.systemName" placeholder="请输入系统名称" />
            </el-form-item>
            <el-form-item label="系统描述" prop="systemDescription">
              <el-input
                v-model="basicForm.systemDescription"
                type="textarea"
                :rows="3"
                placeholder="请输入系统描述"
              />
            </el-form-item>
            <el-form-item label="系统版本" prop="systemVersion">
              <el-input v-model="basicForm.systemVersion" placeholder="请输入系统版本" />
            </el-form-item>
            <el-form-item label="管理员邮箱" prop="adminEmail">
              <el-input v-model="basicForm.adminEmail" placeholder="请输入管理员邮箱" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveBasic" :loading="basicLoading">
                保存设置
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 安全设置 -->
        <el-tab-pane label="安全设置" name="security">
          <el-form
            ref="securityFormRef"
            :model="securityForm"
            :rules="securityRules"
            label-width="120px"
            class="settings-form"
          >
            <el-form-item label="JWT密钥" prop="jwtSecret">
              <el-input
                v-model="securityForm.jwtSecret"
                type="password"
                placeholder="请输入JWT密钥"
                show-password
              />
              <div class="form-tip">用于生成和验证JWT令牌的密钥</div>
            </el-form-item>
            <el-form-item label="令牌有效期" prop="tokenTTL">
              <el-input-number
                v-model="securityForm.tokenTTL"
                :min="300"
                :max="86400"
                placeholder="令牌有效期（秒）"
              />
              <div class="form-tip">访问令牌的有效期，单位：秒</div>
            </el-form-item>
            <el-form-item label="刷新令牌有效期" prop="refreshTokenTTL">
              <el-input-number
                v-model="securityForm.refreshTokenTTL"
                :min="3600"
                :max="604800"
                placeholder="刷新令牌有效期（秒）"
              />
              <div class="form-tip">刷新令牌的有效期，单位：秒</div>
            </el-form-item>
            <el-form-item label="密码最小长度" prop="minPasswordLength">
              <el-input-number
                v-model="securityForm.minPasswordLength"
                :min="6"
                :max="32"
                placeholder="密码最小长度"
              />
              <div class="form-tip">用户密码的最小长度要求</div>
            </el-form-item>
            <el-form-item label="登录失败锁定" prop="loginLockEnabled">
              <el-switch v-model="securityForm.loginLockEnabled" />
              <div class="form-tip">启用后，连续登录失败将锁定账户</div>
            </el-form-item>
            <el-form-item label="最大失败次数" prop="maxLoginAttempts" v-if="securityForm.loginLockEnabled">
              <el-input-number
                v-model="securityForm.maxLoginAttempts"
                :min="3"
                :max="10"
                placeholder="最大失败次数"
              />
            </el-form-item>
            <el-form-item label="锁定时间" prop="lockoutDuration" v-if="securityForm.loginLockEnabled">
              <el-input-number
                v-model="securityForm.lockoutDuration"
                :min="300"
                :max="3600"
                placeholder="锁定时间（秒）"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveSecurity" :loading="securityLoading">
                保存设置
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 数据库设置 -->
        <el-tab-pane label="数据库设置" name="database">
          <el-form
            ref="databaseFormRef"
            :model="databaseForm"
            :rules="databaseRules"
            label-width="120px"
            class="settings-form"
          >
            <el-form-item label="数据库类型" prop="dbType">
              <el-select v-model="databaseForm.dbType" placeholder="请选择数据库类型" style="width: 100%">
                <el-option label="MySQL" value="mysql" />
                <el-option label="PostgreSQL" value="postgresql" />
                <el-option label="SQLite" value="sqlite" />
              </el-select>
            </el-form-item>
            <el-form-item label="数据库地址" prop="dbHost">
              <el-input v-model="databaseForm.dbHost" placeholder="请输入数据库地址" />
            </el-form-item>
            <el-form-item label="数据库端口" prop="dbPort">
              <el-input-number
                v-model="databaseForm.dbPort"
                :min="1"
                :max="65535"
                placeholder="数据库端口"
              />
            </el-form-item>
            <el-form-item label="数据库名称" prop="dbName">
              <el-input v-model="databaseForm.dbName" placeholder="请输入数据库名称" />
            </el-form-item>
            <el-form-item label="用户名" prop="dbUsername">
              <el-input v-model="databaseForm.dbUsername" placeholder="请输入数据库用户名" />
            </el-form-item>
            <el-form-item label="密码" prop="dbPassword">
              <el-input
                v-model="databaseForm.dbPassword"
                type="password"
                placeholder="请输入数据库密码"
                show-password
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleTestConnection" :loading="testLoading">
                测试连接
              </el-button>
              <el-button type="success" @click="handleSaveDatabase" :loading="databaseLoading">
                保存设置
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- Redis设置 -->
        <el-tab-pane label="Redis设置" name="redis">
          <el-form
            ref="redisFormRef"
            :model="redisForm"
            :rules="redisRules"
            label-width="120px"
            class="settings-form"
          >
            <el-form-item label="Redis地址" prop="redisHost">
              <el-input v-model="redisForm.redisHost" placeholder="请输入Redis地址" />
            </el-form-item>
            <el-form-item label="Redis端口" prop="redisPort">
              <el-input-number
                v-model="redisForm.redisPort"
                :min="1"
                :max="65535"
                placeholder="Redis端口"
              />
            </el-form-item>
            <el-form-item label="Redis密码" prop="redisPassword">
              <el-input
                v-model="redisForm.redisPassword"
                type="password"
                placeholder="请输入Redis密码"
                show-password
              />
            </el-form-item>
            <el-form-item label="数据库索引" prop="redisDB">
              <el-input-number
                v-model="redisForm.redisDB"
                :min="0"
                :max="15"
                placeholder="Redis数据库索引"
              />
            </el-form-item>
            <el-form-item label="连接超时" prop="redisTimeout">
              <el-input-number
                v-model="redisForm.redisTimeout"
                :min="1000"
                :max="30000"
                placeholder="连接超时（毫秒）"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleTestRedis" :loading="redisTestLoading">
                测试连接
              </el-button>
              <el-button type="success" @click="handleSaveRedis" :loading="redisLoading">
                保存设置
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 日志设置 -->
        <el-tab-pane label="日志设置" name="logging">
          <el-form
            ref="loggingFormRef"
            :model="loggingForm"
            :rules="loggingRules"
            label-width="120px"
            class="settings-form"
          >
            <el-form-item label="日志级别" prop="logLevel">
              <el-select v-model="loggingForm.logLevel" placeholder="请选择日志级别" style="width: 100%">
                <el-option label="DEBUG" value="debug" />
                <el-option label="INFO" value="info" />
                <el-option label="WARN" value="warn" />
                <el-option label="ERROR" value="error" />
              </el-select>
            </el-form-item>
            <el-form-item label="日志文件路径" prop="logFilePath">
              <el-input v-model="loggingForm.logFilePath" placeholder="请输入日志文件路径" />
            </el-form-item>
            <el-form-item label="最大文件大小" prop="maxFileSize">
              <el-input-number
                v-model="loggingForm.maxFileSize"
                :min="1"
                :max="1000"
                placeholder="最大文件大小（MB）"
              />
            </el-form-item>
            <el-form-item label="保留文件数" prop="maxFiles">
              <el-input-number
                v-model="loggingForm.maxFiles"
                :min="1"
                :max="100"
                placeholder="保留文件数"
              />
            </el-form-item>
            <el-form-item label="启用控制台输出" prop="consoleOutput">
              <el-switch v-model="loggingForm.consoleOutput" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveLogging" :loading="loggingLoading">
                保存设置
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'

const activeTab = ref('basic')
const basicFormRef = ref()
const securityFormRef = ref()
const databaseFormRef = ref()
const redisFormRef = ref()
const loggingFormRef = ref()

const basicLoading = ref(false)
const securityLoading = ref(false)
const databaseLoading = ref(false)
const testLoading = ref(false)
const redisLoading = ref(false)
const redisTestLoading = ref(false)
const loggingLoading = ref(false)

// 基本设置
const basicForm = reactive({
  systemName: '认证授权中心',
  systemDescription: '企业级认证授权中心管理系统',
  systemVersion: '1.0.0',
  adminEmail: 'admin@auth-center.com'
})

const basicRules = {
  systemName: [
    { required: true, message: '请输入系统名称', trigger: 'blur' }
  ],
  systemDescription: [
    { required: true, message: '请输入系统描述', trigger: 'blur' }
  ],
  systemVersion: [
    { required: true, message: '请输入系统版本', trigger: 'blur' }
  ],
  adminEmail: [
    { required: true, message: '请输入管理员邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

// 安全设置
const securityForm = reactive({
  jwtSecret: '',
  tokenTTL: 3600,
  refreshTokenTTL: 604800,
  minPasswordLength: 6,
  loginLockEnabled: true,
  maxLoginAttempts: 5,
  lockoutDuration: 900
})

const securityRules = {
  jwtSecret: [
    { required: true, message: '请输入JWT密钥', trigger: 'blur' },
    { min: 32, message: 'JWT密钥长度不能少于32位', trigger: 'blur' }
  ],
  tokenTTL: [
    { required: true, message: '请输入令牌有效期', trigger: 'blur' }
  ],
  refreshTokenTTL: [
    { required: true, message: '请输入刷新令牌有效期', trigger: 'blur' }
  ],
  minPasswordLength: [
    { required: true, message: '请输入密码最小长度', trigger: 'blur' }
  ]
}

// 数据库设置
const databaseForm = reactive({
  dbType: 'mysql',
  dbHost: 'localhost',
  dbPort: 3306,
  dbName: 'auth_center',
  dbUsername: 'root',
  dbPassword: ''
})

const databaseRules = {
  dbType: [
    { required: true, message: '请选择数据库类型', trigger: 'change' }
  ],
  dbHost: [
    { required: true, message: '请输入数据库地址', trigger: 'blur' }
  ],
  dbPort: [
    { required: true, message: '请输入数据库端口', trigger: 'blur' }
  ],
  dbName: [
    { required: true, message: '请输入数据库名称', trigger: 'blur' }
  ],
  dbUsername: [
    { required: true, message: '请输入数据库用户名', trigger: 'blur' }
  ]
}

// Redis设置
const redisForm = reactive({
  redisHost: 'localhost',
  redisPort: 6379,
  redisPassword: '',
  redisDB: 0,
  redisTimeout: 5000
})

const redisRules = {
  redisHost: [
    { required: true, message: '请输入Redis地址', trigger: 'blur' }
  ],
  redisPort: [
    { required: true, message: '请输入Redis端口', trigger: 'blur' }
  ],
  redisDB: [
    { required: true, message: '请输入Redis数据库索引', trigger: 'blur' }
  ]
}

// 日志设置
const loggingForm = reactive({
  logLevel: 'info',
  logFilePath: './logs/auth-center.log',
  maxFileSize: 100,
  maxFiles: 10,
  consoleOutput: true
})

const loggingRules = {
  logLevel: [
    { required: true, message: '请选择日志级别', trigger: 'change' }
  ],
  logFilePath: [
    { required: true, message: '请输入日志文件路径', trigger: 'blur' }
  ],
  maxFileSize: [
    { required: true, message: '请输入最大文件大小', trigger: 'blur' }
  ],
  maxFiles: [
    { required: true, message: '请输入保留文件数', trigger: 'blur' }
  ]
}

// 保存基本设置
const handleSaveBasic = async () => {
  if (!basicFormRef.value) return
  
  try {
    await basicFormRef.value.validate()
    basicLoading.value = true
    
    // 这里应该调用保存基本设置的API
    await new Promise(resolve => setTimeout(resolve, 1000)) // 模拟API调用
    
    ElMessage.success('基本设置保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    basicLoading.value = false
  }
}

// 保存安全设置
const handleSaveSecurity = async () => {
  if (!securityFormRef.value) return
  
  try {
    await securityFormRef.value.validate()
    securityLoading.value = true
    
    // 这里应该调用保存安全设置的API
    await new Promise(resolve => setTimeout(resolve, 1000)) // 模拟API调用
    
    ElMessage.success('安全设置保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    securityLoading.value = false
  }
}

// 测试数据库连接
const handleTestConnection = async () => {
  if (!databaseFormRef.value) return
  
  try {
    await databaseFormRef.value.validate()
    testLoading.value = true
    
    // 这里应该调用测试数据库连接的API
    await new Promise(resolve => setTimeout(resolve, 2000)) // 模拟API调用
    
    ElMessage.success('数据库连接测试成功')
  } catch (error) {
    ElMessage.error('数据库连接测试失败')
  } finally {
    testLoading.value = false
  }
}

// 保存数据库设置
const handleSaveDatabase = async () => {
  if (!databaseFormRef.value) return
  
  try {
    await databaseFormRef.value.validate()
    databaseLoading.value = true
    
    // 这里应该调用保存数据库设置的API
    await new Promise(resolve => setTimeout(resolve, 1000)) // 模拟API调用
    
    ElMessage.success('数据库设置保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    databaseLoading.value = false
  }
}

// 测试Redis连接
const handleTestRedis = async () => {
  if (!redisFormRef.value) return
  
  try {
    await redisFormRef.value.validate()
    redisTestLoading.value = true
    
    // 这里应该调用测试Redis连接的API
    await new Promise(resolve => setTimeout(resolve, 2000)) // 模拟API调用
    
    ElMessage.success('Redis连接测试成功')
  } catch (error) {
    ElMessage.error('Redis连接测试失败')
  } finally {
    redisTestLoading.value = false
  }
}

// 保存Redis设置
const handleSaveRedis = async () => {
  if (!redisFormRef.value) return
  
  try {
    await redisFormRef.value.validate()
    redisLoading.value = true
    
    // 这里应该调用保存Redis设置的API
    await new Promise(resolve => setTimeout(resolve, 1000)) // 模拟API调用
    
    ElMessage.success('Redis设置保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    redisLoading.value = false
  }
}

// 保存日志设置
const handleSaveLogging = async () => {
  if (!loggingFormRef.value) return
  
  try {
    await loggingFormRef.value.validate()
    loggingLoading.value = true
    
    // 这里应该调用保存日志设置的API
    await new Promise(resolve => setTimeout(resolve, 1000)) // 模拟API调用
    
    ElMessage.success('日志设置保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    loggingLoading.value = false
  }
}
</script>

<style scoped>
.settings-page {
  padding: 0;
}

.page-card {
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.settings-tabs {
  margin-top: 20px;
}

.settings-form {
  max-width: 600px;
  margin: 20px 0;
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
  line-height: 1.4;
}
</style>
