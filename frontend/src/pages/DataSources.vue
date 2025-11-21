<template>
  <page-container>
    <div style="display: flex; justify-content: flex-start; margin-bottom: 16px;">
      <a-button type="primary" @click="$router.push('/datasources/new')">添加数据源</a-button>
    </div>
    <a-alert v-if="errorMsg" type="error" :content="errorMsg" style="margin-bottom:8px" />

    <div v-if="rows.length">
      <div v-for="item in rows" :key="item.id" class="ds-item">
        <div class="ds-left">
          <img :src="getTypeLogo(item.type)" alt="logo" class="ds-logo" />
          <div class="ds-meta">
            <div class="ds-name">{{ item.name }}</div>
            <div class="ds-sub">{{ item.type }} ｜ {{ item.endpoint }}</div>
          </div>
        </div>
        <div class="ds-actions">
          <a-space>
            <a-button size="mini" @click="edit(item)">编辑</a-button>
            <a-button size="mini" @click="test(item)">测试连接</a-button>
            <a-popconfirm content="确认删除？" @ok="remove(item)">
              <a-button size="mini" status="danger">删除</a-button>
            </a-popconfirm>
          </a-space>
        </div>
      </div>
    </div>
    <a-empty v-else description="暂无数据" />
  </page-container>
</template>
<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import PageContainer from '@/components/PageContainer.vue'
import { listDataSources, updateDataSource, deleteDataSource, testConnection } from '@/api/datasources'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const rows = ref([])
const errorMsg = ref('')

function getTypeLogo(type) {
  try {
    return new URL(`../assets/${type}.png`, import.meta.url).href
  } catch (_) {
    return new URL(`../assets/logo.png`, import.meta.url).href
  }
}

async function fetchList() {
  loading.value = true
  errorMsg.value = ''
  try {
    const { data } = await listDataSources()
    if (data?.code !== 0) {
      errorMsg.value = data?.message || '加载失败'
      rows.value = []
    } else {
      rows.value = data?.data?.items || []
      console.log('datasources items', rows.value)
    }
  } catch (e) {
    errorMsg.value = (e && e.message) || '网络错误'
    rows.value = []
  } finally {
    loading.value = false
  }
}

function edit(record) {
  if (record.type === 'loki') {
    router.push(`/datasources/new/loki?id=${record.id}`)
    return
  }
  if (record.type === 'elasticsearch') {
    router.push(`/datasources/new/elasticsearch?id=${record.id}`)
  }
}

async function remove(record) { await deleteDataSource(record.id); Message.success('删除成功'); await fetchList() }
async function test(record) {
  loading.value = true
  try {
    const { data } = await testConnection(record.id)
    if (data?.code === 0) {
      Message.success('连接成功')
    } else {
      Message.error(data?.message || '连接失败')
    }
  } catch (error) {
    Message.error(error?.response?.data?.message || error?.message || '连接失败')
  } finally {
    loading.value = false
  }
}

onMounted(fetchList)
watch(() => route.query.ts, () => fetchList())
</script>

<style scoped>
.ds-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 15px 16px;
  border: 1px solid var(--color-border-2);
  border-radius: 8px;
  background: var(--color-bg-2);
  margin-bottom: 12px;
}
.ds-left { display: flex; align-items: center; gap: 12px; }
.ds-logo { width: 45px; height: 45px; object-fit: contain; }
.ds-meta { display: flex; flex-direction: column; }
.ds-name { font-weight: 600; }
.ds-sub { color: var(--color-text-3); font-size: 12px; }
</style>

