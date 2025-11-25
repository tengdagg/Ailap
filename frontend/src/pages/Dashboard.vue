<template>
  <page-container>
    <!-- Welcome Section -->


    <!-- Overview Cards -->
    <a-grid :cols="{ xs: 1, sm: 2, md: 4 }" :col-gap="16" :row-gap="16" class="overview-cards">
      <a-grid-item>
        <a-card class="overview-card" :bordered="false">
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
        <a-card class="overview-card" :bordered="false">
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
        <a-card class="overview-card" :bordered="false">
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
        <a-card class="overview-card" :bordered="false">
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
                  <a-avatar shape="square" :style="{ backgroundColor: item.engine === 'loki' ? '#165dff' : '#00b42a' }">
                    {{ item.engine === 'loki' ? 'L' : 'E' }}
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


.overview-card {
  transition: all 0.3s;
}

.overview-card:hover {
  transform: translateY(-2px);
}

.card-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.card-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.card-icon.blue {
  background-color: rgba(22, 93, 255, 0.1);
  color: #165dff;
}

.card-icon.purple {
  background-color: rgba(114, 46, 209, 0.1);
  color: #722ed1;
}

.card-icon.green {
  background-color: rgba(0, 180, 42, 0.1);
  color: #00b42a;
}

.card-icon.orange {
  background-color: rgba(255, 125, 0, 0.1);
  color: #ff7d00;
}

.card-info {
  flex: 1;
}

.card-label {
  font-size: 14px;
  color: var(--color-text-3);
  margin-bottom: 4px;
}

.card-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-1);
}

.activity-card, .quick-actions-card {
  height: 100%;
}

.activity-item {
  padding: 12px 0;
  border-bottom: 1px solid var(--color-border-1);
}

.activity-item:last-child {
  border-bottom: none;
}

.empty-activity {
  text-align: center;
  padding: 20px;
  color: var(--color-text-3);
}
</style>

