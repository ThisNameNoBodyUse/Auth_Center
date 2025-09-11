<template>
  <div class="apps-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>应用管理</span>
          <el-button type="primary" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            创建应用
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-input
          v-model="searchForm.name"
          placeholder="请输入应用名称"
          style="width: 200px; margin-right: 10px"
          clearable
        />
        <el-select
          v-model="searchForm.status"
          placeholder="选择状态"
          style="width: 120px; margin-right: 10px"
          clearable
        >
          <el-option label="全部" value="" />
          <el-option label="启用" :value="1" />
          <el-option label="禁用" :value="0" />
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

      <!-- 应用列表 -->
      <el-table
        :data="apps"
        v-loading="loading"
        style="width: 100%"
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="应用名称" min-width="150" />
        <el-table-column prop="app_id" label="应用ID" min-width="150">
          <template #default="{ row }">
            <el-text type="primary" class="app-id">{{ row.app_id }}</el-text>
          </template>
        </el-table-column>
        <el-table-column prop="app_secret" label="应用密钥" min-width="200">
          <template #default="{ row }">
            <div class="secret-container">
              <el-text v-if="!row.showSecret" class="secret-text">
                {{ maskSecret(row.app_secret) }}
              </el-text>
              <el-text v-else class="secret-text">{{ row.app_secret }}</el-text>
              <el-button
                type="text"
                size="small"
                @click="toggleSecret(row)"
              >
                {{ row.showSecret ? '隐藏' : '显示' }}
              </el-button>
            </div>
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
              @click="handleRegenerateSecret(row)"
            >
              重新生成密钥
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

    <!-- 创建/编辑应用对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑应用' : '创建应用'"
      width="500px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="应用名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入应用名称" />
        </el-form-item>
        <el-form-item label="应用ID" prop="app_id" v-if="isEdit">
          <el-input v-model="form.app_id" disabled />
        </el-form-item>
        <el-form-item label="应用ID" v-else>
          <el-input value="系统自动生成" disabled />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入应用描述"
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getApps, createApp, updateApp, deleteApp, regenerateAppSecret } from '@/api/apps'
import dayjs from 'dayjs'

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

const searchForm = reactive({
  name: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

const apps = ref([])

const form = reactive({
  id: null,
  name: '',
  app_id: '',
  description: '',
  status: 1
})

const formRules = {
  name: [
    { required: true, message: '请输入应用名称', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入应用描述', trigger: 'blur' }
  ]
}

// 加载应用列表
const loadApps = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      size: pagination.size,
      ...searchForm
    }
    const response = await getApps(params)
    apps.value = response.apps || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('加载应用列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadApps()
}

// 重置
const handleReset = () => {
  searchForm.name = ''
  searchForm.status = ''
  pagination.page = 1
  loadApps()
}

// 显示创建对话框
const showCreateDialog = () => {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
}

// 编辑应用
const handleEdit = (row) => {
  isEdit.value = true
  dialogVisible.value = true
  Object.assign(form, { ...row })
}

// 删除应用
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除应用 "${row.name}" 吗？此操作不可恢复！`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteApp(row.app_id)
    ElMessage.success('删除成功')
    loadApps()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 重新生成密钥
const handleRegenerateSecret = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要重新生成应用 "${row.name}" 的密钥吗？原密钥将失效！`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await regenerateAppSecret(row.app_id)
    ElMessage.success('密钥重新生成成功')
    loadApps()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('重新生成密钥失败')
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
      await updateApp(form.app_id, form)
      ElMessage.success('更新成功')
    } else {
      await createApp(form)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    loadApps()
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
    name: '',
    app_id: '',
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

// 切换密钥显示
const toggleSecret = (row) => {
  row.showSecret = !row.showSecret
}

// 掩码密钥
const maskSecret = (secret) => {
  if (!secret) return ''
  return secret.substring(0, 8) + '****' + secret.substring(secret.length - 4)
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 分页处理
const handleSizeChange = (size) => {
  pagination.size = size
  pagination.page = 1
  loadApps()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadApps()
}

onMounted(() => {
  loadApps()
})
</script>

<style scoped>
.apps-page {
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

.app-id {
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.secret-container {
  display: flex;
  align-items: center;
  gap: 8px;
}

.secret-text {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  flex: 1;
  word-break: break-all;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
