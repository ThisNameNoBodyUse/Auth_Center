<template>
  <div class="users-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            创建用户
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
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

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
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
          <el-select v-model="form.app_id" placeholder="请选择应用" style="width: 100%">
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
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getApps } from '@/api/apps'
import dayjs from 'dayjs'

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

// 获取应用名称
const getAppName = (appId) => {
  const app = apps.value.find(a => a.app_id === appId)
  return app ? app.name : appId
}

// 加载应用列表
const loadApps = async () => {
  try {
    const response = await getApps()
    apps.value = response.apps || []
  } catch (error) {
    console.error('Failed to load apps:', error)
  }
}

// 加载用户列表
const loadUsers = async () => {
  loading.value = true
  try {
    // 这里应该调用用户API，暂时用模拟数据
    const mockUsers = [
      {
        id: 1,
        username: 'superadmin',
        email: 'admin@auth-center.com',
        phone: '',
        app_id: 'system-admin',
        is_super_admin: true,
        status: 1,
        created_at: '2024-01-01T00:00:00Z'
      }
    ]
    users.value = mockUsers
    pagination.total = mockUsers.length
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
  searchForm.app_id = ''
  pagination.page = 1
  loadUsers()
}

// 显示创建对话框
const showCreateDialog = () => {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
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
    
    // 这里应该调用删除用户API
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
    
    // 这里应该调用创建/更新用户API
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
    app_id: '',
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
  loadUsers()
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
