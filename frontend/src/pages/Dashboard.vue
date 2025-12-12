<template>
  <page-container>
    <!-- Welcome Section -->


    <!-- Overview Cards -->
    <a-grid :cols="{ xs: 1, sm: 2, md: 4 }" :col-gap="16" :row-gap="16" class="overview-cards">
      <a-grid-item>
        <a-card class="overview-card blue-card" :bordered="false">
          <div class="card-content">
            <div class="card-icon blue">
              <icon-storage />
            </div>
            <div class="card-info">
              <div class="card-label">数据源</div>
              <div class="card-value">{{ stats.datasources }}</div>
            </div>
          </div>
        </a-card>
      </a-grid-item>
      <a-grid-item>
        <a-card class="overview-card purple-card" :bordered="false">
          <div class="card-content">
            <div class="card-icon purple">
              <icon-robot />
            </div>
            <div class="card-info">
              <div class="card-label">AI 模型</div>
              <div class="card-value">{{ stats.models }}</div>
            </div>
          </div>
        </a-card>
      </a-grid-item>
      <a-grid-item>
        <a-card class="overview-card green-card" :bordered="false">
          <div class="card-content">
            <div class="card-icon green">
              <icon-history />
            </div>
            <div class="card-info">
              <div class="card-label">近期查询</div>
              <div class="card-value">{{ stats.queries }}</div>
            </div>
          </div>
        </a-card>
      </a-grid-item>
      <a-grid-item>
        <a-card class="overview-card orange-card" :bordered="false">
          <div class="card-content">
            <div class="card-icon orange">
              <icon-star />
            </div>
            <div class="card-info">
              <div class="card-label">已收藏</div>
              <div class="card-value">{{ stats.favorites }}</div>
            </div>
          </div>
        </a-card>
      </a-grid-item>
    </a-grid>

    <a-grid :cols="{ xs: 1, md: 24 }" :col-gap="16" :row-gap="16" style="margin-top: 16px;">
      <!-- Recent Activity -->
      <a-grid-item :span="{ xs: 24, md: 16 }">
        <a-card title="近期活动" :bordered="false" class="activity-card" :head-style="{ fontSize: '14px' }">
          <template #extra>
            <a-link @click="$router.push('/logs')">查看全部</a-link>
          </template>
          <a-list :bordered="false" :split="false">
            <a-list-item v-for="item in recentActivity" :key="item.id" class="activity-item">
              <a-list-item-meta
                :title="item.query || 'No query'"
                :description="formatTime(item.createdAt)"
              >
                <template #avatar>
                  <a-avatar shape="square" :style="{ backgroundColor: item.engine === 'loki' ? '#165dff' : (item.engine === 'elasticsearch' ? '#00b42a' : '#ff7d00') }">
                    {{ item.engine === 'loki' ? 'L' : (item.engine === 'elasticsearch' ? 'E' : 'V') }}
                  </a-avatar>
                </template>
              </a-list-item-meta>
              <template #actions>
                <a-tag size="small">{{ item.mode }}</a-tag>
              </template>
            </a-list-item>
            <div v-if="recentActivity.length === 0" class="empty-activity">
              暂无近期活动
            </div>
          </a-list>
        </a-card>
      </a-grid-item>

      <!-- Quick Actions -->
      <a-grid-item :span="{ xs: 24, md: 8 }">
        <a-card title="快速开始" :bordered="false" class="quick-actions-card">
          <a-space direction="vertical" fill size="large">
            <a-button long size="large" @click="$router.push('/datasources')">
              <template #icon><icon-plus /></template>
              添加数据源
            </a-button>
            <a-button long size="large" @click="$router.push('/models')">
              <template #icon><icon-settings /></template>
              配置 AI 模型
            </a-button>
            <a-button type="primary" long size="large" @click="$router.push('/logs')">
              <template #icon><icon-search /></template>
              开始日志分析
            </a-button>
          </a-space>
        </a-card>
      </a-grid-item>
    </a-grid>
  </page-container>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import PageContainer from '@/components/PageContainer.vue'
import { IconStorage, IconRobot, IconHistory, IconSearch, IconPlus, IconSettings, IconStar } from '@arco-design/web-vue/es/icon'
import { listDataSources } from '@/api/datasources'
import { listModels } from '@/api/models'
import { history } from '@/api/logs'

const stats = reactive({
  datasources: 0,
  models: 0,
  queries: 0,
  favorites: 0
})

const recentActivity = ref([])

onMounted(async () => {
  try {
    // Fetch Data Sources
    const { data: dsData } = await listDataSources()
    stats.datasources = dsData?.data?.items?.length || 0

    // Fetch Models
    const { data: modelData } = await listModels()
    stats.models = modelData?.data?.items?.length || 0

    // Fetch History
    const { data: historyData } = await history('recent')
    const items = historyData?.data?.items || []
    stats.queries = items.length
    recentActivity.value = items.slice(0, 5) // Top 5

    // Fetch Favorites
    const { data: favData } = await history('favorite')
    stats.favorites = favData?.data?.items?.length || 0
  } catch (e) {
    console.error('Failed to load dashboard data:', e)
  }
})

function formatTime(timeStr) {
  if (!timeStr) return ''
  return new Date(timeStr).toLocaleString()
}
</script>

<style scoped>
.overview-cards {
  margin-bottom: 16px;
}

.overview-card {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 12px;
  overflow: hidden;
  position: relative;
}

.overview-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.08);
}

.overview-card.blue-card {
  background: linear-gradient(135deg, rgba(22, 93, 255, 0.03) 0%, rgba(22, 93, 255, 0.01) 100%);
  border: 1px solid rgba(22, 93, 255, 0.08);
}

.overview-card.purple-card {
  background: linear-gradient(135deg, rgba(114, 46, 209, 0.03) 0%, rgba(114, 46, 209, 0.01) 100%);
  border: 1px solid rgba(114, 46, 209, 0.08);
}

.overview-card.green-card {
  background: linear-gradient(135deg, rgba(0, 180, 42, 0.03) 0%, rgba(0, 180, 42, 0.01) 100%);
  border: 1px solid rgba(0, 180, 42, 0.08);
}

.overview-card.orange-card {
  background: linear-gradient(135deg, rgba(255, 125, 0, 0.03) 0%, rgba(255, 125, 0, 0.01) 100%);
  border: 1px solid rgba(255, 125, 0, 0.08);
}

.overview-card :deep(.arco-card-body) {
  padding: 20px;
}

.card-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.card-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26px;
  position: relative;
  transition: all 0.3s;
}

.overview-card:hover .card-icon {
  transform: scale(1.05);
}

.card-icon.blue {
  background: linear-gradient(135deg, rgba(22, 93, 255, 0.15) 0%, rgba(22, 93, 255, 0.05) 100%);
  color: #165dff;
  box-shadow: 0 4px 12px rgba(22, 93, 255, 0.15);
}

.card-icon.purple {
  background: linear-gradient(135deg, rgba(114, 46, 209, 0.15) 0%, rgba(114, 46, 209, 0.05) 100%);
  color: #722ed1;
  box-shadow: 0 4px 12px rgba(114, 46, 209, 0.15);
}

.card-icon.green {
  background: linear-gradient(135deg, rgba(0, 180, 42, 0.15) 0%, rgba(0, 180, 42, 0.05) 100%);
  color: #00b42a;
  box-shadow: 0 4px 12px rgba(0, 180, 42, 0.15);
}

.card-icon.orange {
  background: linear-gradient(135deg, rgba(255, 125, 0, 0.15) 0%, rgba(255, 125, 0, 0.05) 100%);
  color: #ff7d00;
  box-shadow: 0 4px 12px rgba(255, 125, 0, 0.15);
}

.card-info {
  flex: 1;
}

.card-label {
  font-size: 13px;
  color: var(--color-text-3);
  margin-bottom: 6px;
  font-weight: 500;
  letter-spacing: 0.3px;
}

.card-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-1);
  line-height: 1;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
}

.activity-card, .quick-actions-card {
  height: 100%;
  border-radius: 12px;
}

.activity-card :deep(.arco-card-header),
.quick-actions-card :deep(.arco-card-header) {
  border-bottom: 1px solid var(--color-border-2);
  padding: 16px 20px;
}

.activity-card :deep(.arco-card-body),
.quick-actions-card :deep(.arco-card-body) {
  padding: 20px;
}

.activity-item {
  padding: 14px 0;
  border-bottom: 1px solid var(--color-border-1);
  transition: background-color 0.2s;
}

.activity-item:hover {
  background-color: var(--color-fill-1);
  margin: 0 -12px;
  padding: 14px 12px;
  border-radius: 8px;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-item :deep(.arco-list-item-meta-title) {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-1);
}

.activity-item :deep(.arco-list-item-meta-description) {
  font-size: 12px;
  color: var(--color-text-3);
  margin-top: 4px;
}

.empty-activity {
  text-align: center;
  padding: 40px 20px;
  color: var(--color-text-3);
  font-size: 14px;
}

.quick-actions-card :deep(.arco-btn) {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s;
}

.quick-actions-card :deep(.arco-btn:not(.arco-btn-primary)) {
  border: 1px solid var(--color-border-2);
}

.quick-actions-card :deep(.arco-btn:not(.arco-btn-primary):hover) {
  border-color: rgb(var(--primary-6));
  color: rgb(var(--primary-6));
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(22, 93, 255, 0.1);
}

.quick-actions-card :deep(.arco-btn-primary) {
  background: linear-gradient(135deg, #165dff 0%, #4080ff 100%);
  border: none;
}

.quick-actions-card :deep(.arco-btn-primary:hover) {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(22, 93, 255, 0.25);
}
</style>

