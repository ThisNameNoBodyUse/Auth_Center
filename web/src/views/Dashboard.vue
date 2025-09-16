<template>
  <div class="dashboard">
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon apps">
              <el-icon><Grid /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.apps }}</div>
              <div class="stat-label">应用总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon users">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.users }}</div>
              <div class="stat-label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon roles">
              <el-icon><UserFilled /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.roles }}</div>
              <div class="stat-label">角色总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon permissions">
              <el-icon><Key /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.permissions }}</div>
              <div class="stat-label">权限总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="content-row">
      <el-col :span="16">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>应用使用统计</span>
              <el-button type="text" @click="refreshData">刷新</el-button>
            </div>
          </template>
          <div class="chart-container">
            <el-empty v-if="!chartData.length" description="暂无数据" />
            <div v-else class="chart-placeholder">
              <p>图表数据加载中...</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="8">
        <el-card class="recent-card">
          <template #header>
            <span>最近创建的应用</span>
          </template>
          <div class="recent-list">
            <div v-for="app in recentApps" :key="app.id" class="recent-item">
              <div class="app-info">
                <div class="app-name">{{ app.name }}</div>
                <div class="app-id">{{ app.app_id }}</div>
              </div>
              <div class="app-status">
                <el-tag :type="app.status === 1 ? 'success' : 'danger'">
                  {{ app.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </div>
            </div>
            <el-empty v-if="!recentApps.length" description="暂无数据" />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { getApps } from '@/api/apps'
import { getSelfApp, getUsers, getRoles, getPermissions } from '@/api/app-resources'

const stats = ref({
  apps: 0,
  users: 0,
  roles: 0,
  permissions: 0
})

const chartData = ref([])
const recentApps = ref([])

const authStore = useAuthStore()
const user = computed(() => authStore.user)
const isAppAdmin = computed(() => authStore.loginType === 'app' || user.value?.admin_type === 'app')
const currentAppId = computed(() => user.value?.app_id || '')

const loadStats = async () => {
  try {
    if (isAppAdmin.value && currentAppId.value) {
      // 应用级管理员：应用=1，最近应用显示自己
      try {
        const app = await getSelfApp()
        stats.value.apps = 1
        recentApps.value = [{ id: app.id || 1, name: app.name, app_id: app.app_id, status: app.status ?? 1 }]
      } catch (e) {
        stats.value.apps = 1
        recentApps.value = [{ id: 1, name: currentAppId.value, app_id: currentAppId.value, status: 1 }]
      }
      // 统计本应用用户/角色/权限
      const scoped = { page: 1, size: 1, app_id: currentAppId.value }
      const [u, r, p] = await Promise.all([
        getUsers(scoped),
        getRoles(scoped),
        getPermissions(scoped)
      ])
      stats.value.users = u?.pagination?.total || 0
      stats.value.roles = r?.pagination?.total || 0
      stats.value.permissions = p?.pagination?.total || 0
      return
    }
    // 系统级管理员：统计全局
    const response = await getApps()
    const appList = response.apps || []
    stats.value.apps = appList.length
    recentApps.value = appList.slice(0, 5)

    // 用户总数：/app/users 需要 app_id，因此对每个应用取 total 后累加
    const base = { page: 1, size: 1 }
    const usersTotals = await Promise.all(
      appList.map(app => getUsers({ ...base, app_id: app.app_id }).then(r => r?.pagination?.total || 0).catch(() => 0))
    )
    stats.value.users = usersTotals.reduce((a, b) => a + b, 0)

    // 角色/权限总数：后端已支持全局统计
    const [r, p] = await Promise.all([
      getRoles(base),
      getPermissions(base)
    ])
    stats.value.roles = r?.pagination?.total || 0
    stats.value.permissions = p?.pagination?.total || 0
  } catch (error) {
    // 静默
  }
}

const refreshData = () => {
  loadStats()
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
  margin-right: 16px;
}

.stat-icon.apps {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.users {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.roles {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.permissions {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.content-row {
  margin-bottom: 20px;
}

.chart-card,
.recent-card {
  height: 400px;
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-placeholder {
  text-align: center;
  color: #999;
}

.recent-list {
  max-height: 300px;
  overflow-y: auto;
}

.recent-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.recent-item:last-child {
  border-bottom: none;
}

.app-info {
  flex: 1;
}

.app-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.app-id {
  font-size: 12px;
  color: #999;
}

.app-status {
  margin-left: 12px;
}
</style>
