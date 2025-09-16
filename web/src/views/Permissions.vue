<template>
  <div class="permissions-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>权限管理</span>
          <el-button type="primary" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            创建权限
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-input
          v-model="searchForm.name"
          placeholder="请输入权限名称"
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
        <el-select
          v-model="searchForm.resource"
          placeholder="选择资源类型"
          style="width: 150px; margin-right: 10px"
          clearable
        >
          <el-option label="全部资源" value="" />
          <el-option label="菜单" value="menu" />
          <el-option label="按钮" value="button" />
          <el-option label="API" value="api" />
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

      <!-- 权限列表 -->
      <el-table
        :data="permissions"
        v-loading="loading"
        style="width: 100%"
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="权限名称" min-width="150" />
        <el-table-column prop="code" label="权限编码" min-width="150">
          <template #default="{ row }">
            <el-text type="primary" class="permission-code">{{ row.code }}</el-text>
          </template>
        </el-table-column>
        <el-table-column prop="app_id" label="所属应用" min-width="120">
          <template #default="{ row }">
            <el-tag type="info">{{ getAppName(row.app_id) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource" label="资源类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getResourceType(row.resource)">{{ row.resource }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action" label="操作类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getActionType(row.action)">{{ row.action }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" />
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
              @click="handleManageAPIs(row)"
            >
              管理API
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

    <!-- 创建/编辑权限对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑权限' : '创建权限'"
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
        <el-form-item label="权限名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入权限名称" />
        </el-form-item>
        
        <el-form-item label="资源类型" prop="resource">
          <el-select v-model="form.resource" placeholder="请选择资源类型" style="width: 100%">
            <el-option label="菜单" value="menu" />
            <el-option label="按钮" value="button" />
            <el-option label="API" value="api" />
            <el-option label="数据" value="data" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作类型" prop="action">
          <el-select v-model="form.action" placeholder="请选择操作类型" style="width: 100%">
            <el-option label="查看" value="read" />
            <el-option label="创建" value="create" />
            <el-option label="更新" value="update" />
            <el-option label="删除" value="delete" />
            <el-option label="管理" value="manage" />
            <el-option label="全部" value="all" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入权限描述"
          />
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

    <!-- API管理对话框 -->
    <el-dialog
      v-model="apiDialogVisible"
      title="管理权限API"
      width="800px"
    >
      <div class="api-management">
        <div class="permission-info">
          <h4>权限信息</h4>
          <p><strong>权限名称:</strong> {{ currentPermission?.name }}</p>
          <p><strong>权限编码:</strong> {{ currentPermission?.code }}</p>
          <p><strong>所属应用:</strong> {{ getAppName(currentPermission?.app_id) }}</p>
        </div>
        
        <div class="api-list">
          <div class="api-header">
            <h4>关联的API</h4>
            <el-button type="primary" size="small" @click="showAddAPIDialog">
              <el-icon><Plus /></el-icon>
              添加API
            </el-button>
          </div>
          
          <el-table :data="permissionAPIs" style="width: 100%">
            <el-table-column prop="path" label="API路径" min-width="200" />
            <el-table-column prop="method" label="请求方法" width="100">
              <template #default="{ row }">
                <el-tag :type="getMethodType(row.method)">{{ row.method }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="描述" min-width="200" />
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button type="danger" size="small" @click="handleRemoveAPI(row)">
                  移除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
      <template #footer>
        <el-button @click="apiDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 添加API对话框 -->
    <el-dialog
      v-model="addAPIDialogVisible"
      title="添加API"
      width="500px"
    >
      <el-form
        ref="apiFormRef"
        :model="apiForm"
        :rules="apiFormRules"
        label-width="100px"
      >
        <el-form-item label="API路径" prop="path">
          <el-input v-model="apiForm.path" placeholder="请输入API路径" />
        </el-form-item>
        <el-form-item label="请求方法" prop="method">
          <el-select v-model="apiForm.method" placeholder="请选择请求方法" style="width: 100%">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
            <el-option label="PATCH" value="PATCH" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="apiForm.description"
            type="textarea"
            :rows="2"
            placeholder="请输入API描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addAPIDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAddAPI" :loading="apiLoading">
          添加
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
import { getPermissions, createPermission, updatePermission, deletePermission, getSelfApp } from '@/api/app-resources'
import dayjs from 'dayjs'

const loading = ref(false)
const submitLoading = ref(false)
const apiLoading = ref(false)
const dialogVisible = ref(false)
const apiDialogVisible = ref(false)
const addAPIDialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()
const apiFormRef = ref()

const searchForm = reactive({
  name: '',
  app_id: '',
  resource: ''
})

const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

const permissions = ref([])
const apps = ref([])
const permissionAPIs = ref([])
const currentPermission = ref(null)

const form = reactive({
  id: null,
  app_id: '',
  name: '',
  resource: '',
  action: '',
  description: '',
  status: 1
})

const apiForm = reactive({
  path: '',
  method: '',
  description: ''
})

const formRules = {
  app_id: [
    { required: true, message: '请选择所属应用', trigger: 'change' }
  ],
  name: [
    { required: true, message: '请输入权限名称', trigger: 'blur' }
  ],
  resource: [
    { required: true, message: '请选择资源类型', trigger: 'change' }
  ],
  action: [
    { required: true, message: '请选择操作类型', trigger: 'change' }
  ],
  description: []
}

const apiFormRules = {
  path: [
    { required: true, message: '请输入API路径', trigger: 'blur' }
  ],
  method: [
    { required: true, message: '请选择请求方法', trigger: 'change' }
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

// 获取资源类型标签类型
const getResourceType = (resource) => {
  const types = {
    menu: 'primary',
    button: 'success',
    api: 'warning',
    data: 'info'
  }
  return types[resource] || 'info'
}

// 获取操作类型标签类型
const getActionType = (action) => {
  const types = {
    read: 'info',
    create: 'success',
    update: 'warning',
    delete: 'danger',
    manage: 'primary',
    all: 'danger'
  }
  return types[action] || 'info'
}

// 获取请求方法标签类型
const getMethodType = (method) => {
  const types = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger',
    PATCH: 'info'
  }
  return types[method] || 'info'
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
      if (!searchForm.app_id) searchForm.app_id = currentAppId.value
      if (!form.app_id) form.app_id = currentAppId.value
      return
    }
    const response = await getApps()
    apps.value = response.apps || []
  } catch (error) {
    // 静默
  }
}

// 加载权限列表
const loadPermissions = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.size,
      ...searchForm
    }
    const response = await getPermissions(params)
    permissions.value = response.data || []
    pagination.total = response.pagination?.total || 0
  } catch (error) {
    ElMessage.error('加载权限列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadPermissions()
}

// 重置
const handleReset = () => {
  searchForm.name = ''
  searchForm.app_id = isAppAdmin.value ? currentAppId.value : ''
  searchForm.resource = ''
  pagination.page = 1
  loadPermissions()
}

// 显示创建对话框
const showCreateDialog = () => {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
  if (isAppAdmin.value && currentAppId.value) {
    form.app_id = currentAppId.value
  }
}

// 编辑权限
const handleEdit = (row) => {
  isEdit.value = true
  dialogVisible.value = true
  Object.assign(form, { ...row })
}

// 管理API
const handleManageAPIs = async (row) => {
  currentPermission.value = row
  apiDialogVisible.value = true
  
  // 加载权限关联的API
  try {
    // 这里应该调用API，暂时用模拟数据
    permissionAPIs.value = [
      {
        id: 1,
        path: '/api/v1/apps',
        method: 'POST',
        description: '创建应用'
      },
      {
        id: 2,
        path: '/api/v1/apps',
        method: 'GET',
        description: '获取应用列表'
      }
    ]
  } catch (error) {
    ElMessage.error('加载API信息失败')
  }
}

// 显示添加API对话框
const showAddAPIDialog = () => {
  addAPIDialogVisible.value = true
  resetAPIForm()
}

// 添加API
const handleAddAPI = async () => {
  if (!apiFormRef.value) return
  
  try {
    await apiFormRef.value.validate()
    apiLoading.value = true
    
    // 这里应该调用添加API的API
    ElMessage.success('API添加成功')
    addAPIDialogVisible.value = false
    handleManageAPIs(currentPermission.value) // 重新加载API列表
  } catch (error) {
    ElMessage.error('添加API失败')
  } finally {
    apiLoading.value = false
  }
}

// 移除API
const handleRemoveAPI = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要移除API "${row.path}" 吗？`,
      '确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 这里应该调用移除API的API
    ElMessage.success('API移除成功')
    handleManageAPIs(currentPermission.value) // 重新加载API列表
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('移除API失败')
    }
  }
}

// 删除权限
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除权限 "${row.name}" 吗？此操作不可恢复！`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deletePermission(row.id, row.app_id)
    ElMessage.success('删除成功')
    loadPermissions()
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
    
    const payload = {
      app_id: form.app_id,
      name: form.name,
      resource: form.resource,
      action: form.action,
      description: form.description,
      status: form.status
    }
    if (isEdit.value) {
      await updatePermission(form.id, payload)
    } else {
      await createPermission(payload)
    }
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
    
    dialogVisible.value = false
    loadPermissions()
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
    name: '',
    resource: '',
    action: '',
    description: '',
    status: 1
  })
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

// 重置API表单
const resetAPIForm = () => {
  Object.assign(apiForm, {
    path: '',
    method: '',
    description: ''
  })
  if (apiFormRef.value) {
    apiFormRef.value.resetFields()
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
  loadPermissions()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadPermissions()
}

onMounted(() => {
  loadApps()
  if (isAppAdmin.value && currentAppId.value) {
    searchForm.app_id = currentAppId.value
    form.app_id = currentAppId.value
  }
  loadPermissions()
})

watch(() => user.value, (val) => {
  if (val && isAppAdmin.value) {
    loadApps()
  }
})
</script>

<style scoped>
.permissions-page {
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

.permission-code {
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.api-management {
  max-height: 500px;
  overflow-y: auto;
}

.permission-info {
  margin-bottom: 20px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 4px;
}

.permission-info h4 {
  margin: 0 0 12px 0;
  color: #333;
}

.permission-info p {
  margin: 4px 0;
  color: #666;
}

.api-list h4 {
  margin: 0 0 12px 0;
  color: #333;
}

.api-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
</style>
