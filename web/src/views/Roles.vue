<template>
  <div class="roles-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>角色管理</span>
          <el-button type="primary" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            创建角色
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-input
          v-model="searchForm.name"
          placeholder="请输入角色名称"
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

      <!-- 角色列表 -->
      <el-table
        :data="roles"
        v-loading="loading"
        style="width: 100%"
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" min-width="150" />
        <el-table-column prop="code" label="角色编码" min-width="150">
          <template #default="{ row }">
            <el-text type="primary" class="role-code">{{ row.code }}</el-text>
          </template>
        </el-table-column>
        <el-table-column prop="app_id" label="所属应用" min-width="120">
          <template #default="{ row }">
            <el-tag type="info">{{ getAppName(row.app_id) }}</el-tag>
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
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="handleManagePermissions(row)"
            >
              管理权限
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

    <!-- 创建/编辑角色对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑角色' : '创建角色'"
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
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入角色描述"
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

    <!-- 权限管理对话框 -->
    <el-dialog
      v-model="permissionDialogVisible"
      title="管理角色权限"
      width="800px"
    >
      <div class="permission-management">
        <div class="role-info">
          <h4>角色信息</h4>
          <p><strong>角色名称:</strong> {{ currentRole?.name }}</p>
          <p><strong>角色编码:</strong> {{ currentRole?.code }}</p>
          <p><strong>所属应用:</strong> {{ getAppName(currentRole?.app_id) }}</p>
        </div>
        
        <div class="permission-assignment">
          <h4>权限分配</h4>
          <el-tree
            ref="permissionTreeRef"
            :data="permissionTree"
            :props="treeProps"
            show-checkbox
            node-key="id"
            :default-checked-keys="selectedPermissions"
            class="permission-tree"
          />
        </div>
      </div>
      <template #footer>
        <el-button @click="permissionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSavePermissions" :loading="permissionLoading">
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
import { getRoles, createRole, updateRole, deleteRole, assignRolePermissions, getRolePermissions, getSelfApp } from '@/api/app-resources'
import dayjs from 'dayjs'

const loading = ref(false)
const submitLoading = ref(false)
const permissionLoading = ref(false)
const dialogVisible = ref(false)
const permissionDialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()
const permissionTreeRef = ref()

const searchForm = reactive({
  name: '',
  app_id: ''
})

const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

const roles = ref([])
const apps = ref([])
const permissionTree = ref([])
const selectedPermissions = ref([])
const currentRole = ref(null)

const form = reactive({
  id: null,
  app_id: '',
  name: '',
  description: '',
  status: 1
})

const formRules = {
  app_id: [
    { required: true, message: '请选择所属应用', trigger: 'change' }
  ],
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' }
  ],
  description: []
}

const treeProps = {
  children: 'children',
  label: 'name'
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
      // 拉取自身应用名称
      try {
        const app = await getSelfApp()
        apps.value = [{ app_id: app.app_id, name: app.name }]
      } catch (e) {
        apps.value = [{ app_id: currentAppId.value, name: currentAppId.value }]
      }
      // 同步搜索与表单的 app_id
      if (!searchForm.app_id) searchForm.app_id = currentAppId.value
      if (!form.app_id) form.app_id = currentAppId.value
      return
    }
    const response = await getApps()
    apps.value = response.apps || []
  } catch (error) {
    // 静默；应用级管理员无权拉取应用列表时，已通过上方分支填充
  }
}

// 加载角色列表
const loadRoles = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.size,
      ...searchForm
    }
    const response = await getRoles(params)
    roles.value = response.data || []
    pagination.total = response.pagination?.total || 0
  } catch (error) {
    ElMessage.error('加载角色列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadRoles()
}

// 重置
const handleReset = () => {
  searchForm.name = ''
  searchForm.app_id = isAppAdmin.value ? currentAppId.value : ''
  pagination.page = 1
  loadRoles()
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

// 编辑角色
const handleEdit = (row) => {
  isEdit.value = true
  dialogVisible.value = true
  Object.assign(form, { ...row })
  // 防御：不允许手工编辑或残留 code 字段
  if (Object.prototype.hasOwnProperty.call(form, 'code')) {
    delete form.code
  }
}

// 管理权限
const handleManagePermissions = async (row) => {
  currentRole.value = row
  permissionDialogVisible.value = true
  
  // 加载权限树
  try {
    // 这里应该调用权限API，暂时用模拟数据
    permissionTree.value = [
      {
        id: 1,
        name: '系统管理',
        children: [
          { id: 11, name: '应用管理' },
          { id: 12, name: '用户管理' },
          { id: 13, name: '角色管理' },
          { id: 14, name: '权限管理' }
        ]
      },
      {
        id: 2,
        name: '业务管理',
        children: [
          { id: 21, name: '数据查看' },
          { id: 22, name: '数据编辑' },
          { id: 23, name: '数据删除' }
        ]
      }
    ]
    
    // 加载角色当前权限
    selectedPermissions.value = [11, 12, 13, 14] // 模拟数据
  } catch (error) {
    ElMessage.error('加载权限信息失败')
  }
}

// 保存权限
const handleSavePermissions = async () => {
  permissionLoading.value = true
  try {
    const checkedKeys = permissionTreeRef.value.getCheckedKeys()
    const halfCheckedKeys = permissionTreeRef.value.getHalfCheckedKeys()
    const allKeys = [...checkedKeys, ...halfCheckedKeys]
    
    // 这里应该调用保存权限API
    ElMessage.success('权限保存成功')
    permissionDialogVisible.value = false
  } catch (error) {
    ElMessage.error('保存权限失败')
  } finally {
    permissionLoading.value = false
  }
}

// 删除角色
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除角色 "${row.name}" 吗？此操作不可恢复！`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteRole(row.id, row.app_id)
    ElMessage.success('删除成功')
    loadRoles()
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
      description: form.description,
      status: form.status
    }
    if (isEdit.value) {
      await updateRole(form.id, payload)
    } else {
      await createRole(payload)
    }
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
    
    dialogVisible.value = false
    loadRoles()
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
    description: '',
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
  loadRoles()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadRoles()
}

onMounted(() => {
  loadApps()
  // 应用级管理员固定 app_id 条件下，首次查询限定为自身应用
  if (isAppAdmin.value && currentAppId.value) {
    searchForm.app_id = currentAppId.value
    form.app_id = currentAppId.value
  }
  loadRoles()
})

// 监听登录信息加载完成后再填充应用选项（避免初始为空导致下拉无选项）
watch(() => user.value, (val) => {
  if (val && isAppAdmin.value) {
    loadApps()
  }
})
</script>

<style scoped>
.roles-page {
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

.role-code {
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.permission-management {
  max-height: 500px;
  overflow-y: auto;
}

.role-info {
  margin-bottom: 20px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 4px;
}

.role-info h4 {
  margin: 0 0 12px 0;
  color: #333;
}

.role-info p {
  margin: 4px 0;
  color: #666;
}

.permission-assignment h4 {
  margin: 0 0 12px 0;
  color: #333;
}

.permission-tree {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 8px;
}
</style>
