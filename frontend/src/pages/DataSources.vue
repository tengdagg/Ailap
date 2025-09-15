<template>
  <page-container>
    <div style="display: flex; justify-content: flex-end; margin-bottom: 16px;">
      <a-button type="primary" @click="$router.push('/datasources/new')">新建数据源</a-button>
    </div>
    <a-alert v-if="errorMsg" type="error" :content="errorMsg" style="margin-bottom:8px" />

    <a-table v-if="rows.length" :key="listKey" :data="rows" :columns="columns" :loading="loading" row-key="id" :bordered="true" :pagination="false">
      <template #actions="{ record }">
        <a-space>
          <a-button size="mini" @click="edit(record)">编辑</a-button>
          <a-button size="mini" @click="test(record)">测试连接</a-button>
          <a-popconfirm content="确认删除？" @ok="remove(record)">
            <a-button size="mini" status="danger">删除</a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </a-table>
    <a-empty v-else description="暂无数据" />
  </page-container>
</template>
<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import PageContainer from '@/components/PageContainer.vue'
import { listDataSources, updateDataSource, deleteDataSource, testConnection } from '@/api/datasources'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const rows = ref([])
const errorMsg = ref('')
const listKey = computed(() => (route.query.ts || '') + ':' + rows.value.length)

const columns = [
  { title: '名称', dataIndex: 'name' },
  { title: '类型', dataIndex: 'type', width: 140 },
  { title: '地址', dataIndex: 'endpoint' },
  { title: '操作', slotName: 'actions', width: 260 },
]

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

async function remove(record) { await deleteDataSource(record.id); await fetchList() }
async function test(record) {
  loading.value = true
  try {
    const { data } = await testConnection(record.id)
    if (data?.code === 0) {
      // 使用全局消息提示
      // eslint-disable-next-line no-alert
      console.log('连接成功', data?.data)
    } else {
      console.error('连接失败', data?.message)
    }
  } finally {
    loading.value = false
  }
}

onMounted(fetchList)
watch(() => route.query.ts, () => fetchList())
</script>

