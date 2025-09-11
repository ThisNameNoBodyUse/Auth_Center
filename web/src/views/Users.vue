<template>
  <div class="users-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            {{ activeTab === 'system' ? '创建系统管理员' : '创建应用用户' }}
          </el-button>
        </div>
      </template>

      <!-- 标签页 -->
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="系统管理员" name="system" v-if="isSystemAdmin">
          <div class="tab-content">
            <!-- 系统管理员搜索栏 -->
            <div class="search-bar">
              <el-input
                v-model="systemSearchForm.username"
                placeholder="请输入用户名"
                style="width: 200px; margin-right: 10px"
                clearable
              />
              <el-input
                v-model="systemSearchForm.email"
                placeholder="请输入邮箱"
                style="width: 200px; margin-right: 10px"
                clearable
              />
              <el-select
                v-model="systemSearchForm.admin_type"
                placeholder="管理员类型"
                style="width: 150px; margin-right: 10px"
                clearable
              >
                <el-option label="全部类型" value="" />
                <el-option label="系统级" value="system" />
                <el-option label="应用级" value="app" />
              </el-select>
              <el-select
                v-model="systemSearchForm.app_id"
                placeholder="选择应用"
                style="width: 150px; margin-right: 10px"
                clearable
              >
                <el-option label="全部应用" value="" />
                <el-option
                  v-for="app in apps"
                  :key="app.app_id"
                  :label="app.name"
                  :value="app.app_id"
                />
              </el-select>
              <el-button type="primary" @click="handleSystemSearch">
                <el-icon><Search /></el-icon>
                搜索
              </el-button>
              <el-button @click="handleSystemReset">
                <el-icon><Refresh /></el-icon>
                重置
              </el-button>
            </div>

            <!-- 系统管理员列表 -->
            <el-table
              :data="systemAdmins"
              v-loading="systemLoading"
              style="width: 100%"
            >
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="username" label="用户名" />
              <el-table-column prop="email" label="邮箱" />
              <el-table-column prop="phone" label="手机号" />
              <el-table-column prop="admin_type" label="类型" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.admin_type === 'system' ? 'danger' : 'primary'">
                    {{ row.admin_type === 'system' ? '系统级' : '应用级' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="app_id" label="关联应用" />
              <el-table-column prop="is_active" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.is_active ? 'success' : 'danger'">
                    {{ row.is_active ? '启用' : '禁用' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="last_login_at" label="最后登录" width="180">
                <template #default="{ row }">
                  {{ row.last_login_at ? dayjs(row.last_login_at).format('YYYY-MM-DD HH:mm:ss') : '-' }}
                </template>
              </el-table-column>
              <el-table-column prop="created_at" label="创建时间" width="180">
                <template #default="{ row }">
                  {{ dayjs(row.created_at).format('YYYY-MM-DD HH:mm:ss') }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="200">
                <template #default="{ row }">
                  <el-button size="small" @click="showSystemEditDialog(row)">编辑</el-button>
                  <el-button size="small" type="warning" @click="showResetPasswordDialog(row)">重置密码</el-button>
                  <el-button size="small" type="danger" @click="handleSystemDelete(row)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>

            <!-- 系统管理员分页 -->
            <el-pagination
              v-model:current-page="systemPagination.page"
              v-model:page-size="systemPagination.size"
              :page-sizes="[10, 20, 50, 100]"
              :total="systemPagination.total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="loadSystemAdmins"
              @current-change="loadSystemAdmins"
              style="margin-top: 20px; text-align: right"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="应用用户" name="app">
          <div class="tab-content">
            <!-- 应用用户搜索栏 -->
            <div class="search-bar">
        <el-input
          v-model="searchForm.username"
          placeholder="请输入用户名"
          style="width: 200px; margin-right: 10px"
          clearable
        />
        <el-input
          v-model="searchForm.email"
          placeholder="请输入邮箱"
          style="width: 200px; margin-right: 10px"
          clearable
        />
        <el-select
          v-model="searchForm.app_id"
          placeholder="选择应用"
          style="width: 150px; margin-right: 10px"
          :disabled="isAppAdmin"
          clearable
        >
          <el-option v-if="!isAppAdmin" label="全部应用" value="" />
          <el-option
            v-for="app in apps"
            :key="app.app_id"
            :label="app.name"
            :value="app.app_id"
          />
        </el-select>
        <el-button type="primary" @click="handleSearch">
          <el-icon><Search /></el-icon>
          搜索
        </el-button>
        <el-button @click="handleReset">
          <el-icon><Refresh /></el-icon>
          重置
        </el-button>
      </div>

      <!-- 用户列表 -->
      <el-table
        :data="users"
        v-loading="loading"
        style="width: 100%"
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="email" label="邮箱" min-width="200" />
        <el-table-column prop="phone" label="手机号" min-width="120" />
        <el-table-column prop="app_id" label="所属应用" min-width="120">
          <template #default="{ row }">
            <el-tag type="info">{{ getAppName(row.app_id) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_super_admin" label="超级管理员" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_super_admin ? 'danger' : 'info'">
              {{ row.is_super_admin ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="handleManageRoles(row)"
            >
              管理角色
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

            <!-- 应用用户分页 -->
            <el-pagination
              v-model:current-page="pagination.page"
              v-model:page-size="pagination.size"
              :page-sizes="[10, 20, 50, 100]"
              :total="pagination.total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
              style="margin-top: 20px; text-align: right"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 创建/编辑用户对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑用户' : '创建用户'"
      width="500px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="所属应用" prop="app_id">
          <el-select v-model="form.app_id" placeholder="请选择应用" style="width: 100%" :disabled="isAppAdmin">
            <el-option
              v-for="app in apps"
              :key="app.app_id"
              :label="app.name"
              :value="app.app_id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="超级管理员" prop="is_super_admin">
          <el-radio-group v-model="form.is_super_admin">
            <el-radio :label="true">是</el-radio>
            <el-radio :label="false">否</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 系统管理员创建/编辑对话框 -->
    <el-dialog
      v-model="systemDialogVisible"
      :title="isEdit ? '编辑系统管理员' : '创建系统管理员'"
      width="600px"
    >
      <el-form
        ref="systemFormRef"
        :model="systemForm"
        :rules="systemFormRules"
        label-width="120px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="systemForm.username" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="systemForm.email" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="systemForm.phone" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="systemForm.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="管理员类型" prop="admin_type">
          <el-radio-group v-model="systemForm.admin_type">
            <el-radio label="system">系统级</el-radio>
            <el-radio label="app">应用级</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="关联应用" prop="app_id" v-if="systemForm.admin_type === 'app'">
          <el-select v-model="systemForm.app_id" placeholder="请选择应用">
            <el-option
              v-for="app in apps"
              :key="app.app_id"
              :label="app.name"
              :value="app.app_id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="is_active">
          <el-switch v-model="systemForm.is_active" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="systemDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSystemSubmit" :loading="submitLoading">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 重置密码对话框 -->
    <el-dialog
      v-model="resetPasswordDialogVisible"
      title="重置密码"
      width="400px"
    >
      <el-form
        ref="resetPasswordFormRef"
        :model="resetPasswordForm"
        :rules="resetPasswordFormRules"
        label-width="100px"
      >
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="resetPasswordForm.new_password" type="password" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirm_password">
          <el-input v-model="resetPasswordForm.confirm_password" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetPasswordDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleResetPasswordSubmit" :loading="submitLoading">
          重置
        </el-button>
      </template>
    </el-dialog>

    <!-- 角色管理对话框 -->
    <el-dialog
      v-model="roleDialogVisible"
      title="管理用户角色"
      width="600px"
    >
      <div class="role-management">
        <div class="user-info">
          <h4>用户信息</h4>
          <p><strong>用户名:</strong> {{ currentUser?.username }}</p>
          <p><strong>所属应用:</strong> {{ getAppName(currentUser?.app_id) }}</p>
        </div>
        
        <div class="role-assignment">
          <h4>角色分配</h4>
          <el-checkbox-group v-model="selectedRoles">
            <el-checkbox
              v-for="role in availableRoles"
              :key="role.id"
              :label="role.id"
            >
              {{ role.name }} ({{ role.code }})
            </el-checkbox>
          </el-checkbox-group>
        </div>
      </div>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveRoles" :loading="roleLoading">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getApps } from '@/api/apps'
import { getUsers, createUser, updateUser, deleteUser, assignUserRoles, getUserRoles, getSelfApp } from '@/api/app-resources'
import { getSystemAdmins, createSystemAdmin, updateSystemAdmin, deleteSystemAdmin, resetSystemAdminPassword } from '@/api/system-admins'
import dayjs from 'dayjs'

// 标签页状态
const activeTab = ref('system')

// 应用用户相关
const loading = ref(false)
const submitLoading = ref(false)
const roleLoading = ref(false)
const dialogVisible = ref(false)
const roleDialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

const searchForm = reactive({
  username: '',
  email: '',
  app_id: ''
})

const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

const users = ref([])
const apps = ref([])
const availableRoles = ref([])

// 系统管理员相关
const systemLoading = ref(false)
const systemDialogVisible = ref(false)
const resetPasswordDialogVisible = ref(false)
const systemFormRef = ref()
const resetPasswordFormRef = ref()

const systemSearchForm = reactive({
  username: '',
  email: '',
  admin_type: '',
  app_id: ''
})

const systemPagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

const systemAdmins = ref([])
const systemForm = reactive({
  id: null,
  username: '',
  email: '',
  phone: '',
  password: '',
  admin_type: 'system',
  app_id: '',
  is_active: true
})

const resetPasswordForm = reactive({
  id: null,
  new_password: '',
  confirm_password: ''
})
const selectedRoles = ref([])
const currentUser = ref(null)

const form = reactive({
  id: null,
  app_id: '',
  username: '',
  email: '',
  phone: '',
  password: '',
  is_super_admin: false,
  status: 1
})

const formRules = {
  app_id: [
    { required: true, message: '请选择所属应用', trigger: 'change' }
  ],
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

// 系统管理员表单验证规则
const systemFormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  admin_type: [
    { required: true, message: '请选择管理员类型', trigger: 'change' }
  ],
  app_id: [
    { 
      validator: (rule, value, callback) => {
        if (systemForm.admin_type === 'app' && !value) {
          callback(new Error('应用级管理员必须选择关联应用'))
        } else {
          callback()
        }
      }, 
      trigger: 'change' 
    }
  ]
}

// 重置密码表单验证规则
const resetPasswordFormRules = {
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { 
      validator: (rule, value, callback) => {
        if (value !== resetPasswordForm.new_password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ]
}

// 认证与身份
const authStore = useAuthStore()
const user = computed(() => authStore.user)
const isSystemAdmin = computed(() => authStore.loginType === 'system' || user.value?.admin_type === 'system')
const isAppAdmin = computed(() => authStore.loginType === 'app' || user.value?.admin_type === 'app')
const currentAppId = computed(() => user.value?.app_id || '')

// 获取应用名称
const getAppName = (appId) => {
  const app = apps.value.find(a => a.app_id === appId)
  return app ? app.name : appId
}

// 加载应用列表
const loadApps = async () => {
  try {
    if (isAppAdmin.value && currentAppId.value) {
      try {
        const app = await getSelfApp()
        apps.value = [{ app_id: app.app_id, name: app.name }]
      } catch (e) {
        apps.value = [{ app_id: currentAppId.value, name: currentAppId.value }]
      }
      if (!systemSearchForm.app_id) systemSearchForm.app_id = currentAppId.value
      if (!searchForm.app_id) searchForm.app_id = currentAppId.value
      if (!form.app_id) form.app_id = currentAppId.value
      activeTab.value = 'app'
      return
    }
    const response = await getApps()
    apps.value = response.apps || []
  } catch (error) {
    // 静默
  }
}

// 标签页切换
const handleTabChange = (tabName) => {
  if (tabName === 'system') {
    loadSystemAdmins()
  } else {
    loadUsers()
  }
}

// 加载系统管理员列表
const loadSystemAdmins = async () => {
  systemLoading.value = true
  try {
    const params = {
      page: systemPagination.page,
      page_size: systemPagination.size,
      ...systemSearchForm
    }
    const response = await getSystemAdmins(params)
    systemAdmins.value = response.data || []
    systemPagination.total = response.pagination?.total || 0
  } catch (error) {
    ElMessage.error('加载系统管理员列表失败')
  } finally {
    systemLoading.value = false
  }
}

// 系统管理员搜索
const handleSystemSearch = () => {
  systemPagination.page = 1
  loadSystemAdmins()
}

// 系统管理员重置
const handleSystemReset = () => {
  systemSearchForm.username = ''
  systemSearchForm.email = ''
  systemSearchForm.admin_type = ''
  systemSearchForm.app_id = ''
  systemPagination.page = 1
  loadSystemAdmins()
}

// 显示系统管理员创建/编辑对话框
const showSystemCreateDialog = () => {
  isEdit.value = false
  systemDialogVisible.value = true
  resetSystemForm()
}

const showSystemEditDialog = (admin) => {
  isEdit.value = true
  systemDialogVisible.value = true
  Object.assign(systemForm, {
    id: admin.id,
    username: admin.username,
    email: admin.email,
    phone: admin.phone,
    password: '',
    admin_type: admin.admin_type,
    app_id: admin.app_id,
    is_active: admin.is_active
  })
}

// 显示重置密码对话框
const showResetPasswordDialog = (admin) => {
  resetPasswordDialogVisible.value = true
  resetPasswordForm.id = admin.id
  resetPasswordForm.new_password = ''
  resetPasswordForm.confirm_password = ''
}

// 重置系统管理员表单
const resetSystemForm = () => {
  Object.assign(systemForm, {
    id: null,
    username: '',
    email: '',
    phone: '',
    password: '',
    admin_type: 'system',
    app_id: '',
    is_active: true
  })
  if (systemFormRef.value) {
    systemFormRef.value.resetFields()
  }
}

// 系统管理员提交
const handleSystemSubmit = async () => {
  if (!systemFormRef.value) return
  
  try {
    await systemFormRef.value.validate()
    submitLoading.value = true
    
    if (isEdit.value) {
      await updateSystemAdmin(systemForm.id, systemForm)
    } else {
      await createSystemAdmin(systemForm)
    }
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
    
    systemDialogVisible.value = false
    loadSystemAdmins()
  } catch (error) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  } finally {
    submitLoading.value = false
  }
}

// 系统管理员删除
const handleSystemDelete = async (admin) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除系统管理员 "${admin.username}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteSystemAdmin(admin.id)
    ElMessage.success('删除成功')
    loadSystemAdmins()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 重置密码提交
const handleResetPasswordSubmit = async () => {
  if (!resetPasswordFormRef.value) return
  
  try {
    await resetPasswordFormRef.value.validate()
    submitLoading.value = true
    
    await resetSystemAdminPassword(resetPasswordForm.id, {
      new_password: resetPasswordForm.new_password
    })
    ElMessage.success('密码重置成功')
    
    resetPasswordDialogVisible.value = false
  } catch (error) {
    ElMessage.error('密码重置失败')
  } finally {
    submitLoading.value = false
  }
}

// 加载用户列表
const loadUsers = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.size,
      ...searchForm
    }
    const response = await getUsers(params)
    users.value = response.data || []
    pagination.total = response.pagination?.total || 0
  } catch (error) {
    ElMessage.error('加载用户列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadUsers()
}

// 重置
const handleReset = () => {
  searchForm.username = ''
  searchForm.email = ''
  searchForm.app_id = isAppAdmin.value ? currentAppId.value : ''
  pagination.page = 1
  loadUsers()
}

// 显示创建对话框
const showCreateDialog = () => {
  if (activeTab.value === 'system') {
    showSystemCreateDialog()
  } else {
    isEdit.value = false
    dialogVisible.value = true
    resetForm()
    if (isAppAdmin.value && currentAppId.value) {
      form.app_id = currentAppId.value
    }
  }
}

// 编辑用户
const handleEdit = (row) => {
  isEdit.value = true
  dialogVisible.value = true
  Object.assign(form, { ...row })
}

// 管理角色
const handleManageRoles = async (row) => {
  currentUser.value = row
  roleDialogVisible.value = true
  
  // 加载可用角色
  try {
    // 这里应该调用角色API，暂时用模拟数据
    availableRoles.value = [
      { id: 1, name: '超级管理员', code: 'super_admin' },
      { id: 2, name: '管理员', code: 'admin' },
      { id: 3, name: '普通用户', code: 'user' }
    ]
    
    // 加载用户当前角色
    selectedRoles.value = [1] // 模拟数据
  } catch (error) {
    ElMessage.error('加载角色信息失败')
  }
}

// 保存角色
const handleSaveRoles = async () => {
  roleLoading.value = true
  try {
    // 这里应该调用保存角色API
    ElMessage.success('角色保存成功')
    roleDialogVisible.value = false
  } catch (error) {
    ElMessage.error('保存角色失败')
  } finally {
    roleLoading.value = false
  }
}

// 删除用户
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${row.username}" 吗？此操作不可恢复！`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteUser(row.id)
    ElMessage.success('删除成功')
    loadUsers()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitLoading.value = true
    
    if (isEdit.value) {
      await updateUser(form.id, form)
    } else {
      await createUser(form)
    }
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
    
    dialogVisible.value = false
    loadUsers()
  } catch (error) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  } finally {
    submitLoading.value = false
  }
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    id: null,
    app_id: isAppAdmin.value ? currentAppId.value : '',
    username: '',
    email: '',
    phone: '',
    password: '',
    is_super_admin: false,
    status: 1
  })
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

// 关闭对话框
const handleDialogClose = () => {
  resetForm()
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 分页处理
const handleSizeChange = (size) => {
  pagination.size = size
  pagination.page = 1
  loadUsers()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadUsers()
}

onMounted(() => {
  loadApps()
  if (isAppAdmin.value) {
    activeTab.value = 'app'
    systemSearchForm.app_id = currentAppId.value
    searchForm.app_id = currentAppId.value
    form.app_id = currentAppId.value
    loadUsers()
  } else {
    loadSystemAdmins() // 默认加载系统管理员
  }
})

watch(() => user.value, (val) => {
  if (val && isAppAdmin.value) {
    loadApps()
  }
})
</script>

<style scoped>
.users-page {
  padding: 0;
}

.page-card {
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 500;
}

.search-bar {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.role-management {
  max-height: 400px;
  overflow-y: auto;
}

.user-info {
  margin-bottom: 20px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 4px;
}

.user-info h4 {
  margin: 0 0 12px 0;
  color: #333;
}

.user-info p {
  margin: 4px 0;
  color: #666;
}

.role-assignment h4 {
  margin: 0 0 12px 0;
  color: #333;
}

.role-assignment .el-checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
