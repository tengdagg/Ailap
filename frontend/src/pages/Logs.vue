<template>
  <page-container title="日志查询" subtitle="Grafana 风格查询编辑器">
    <div style="margin-bottom:12px; display:flex; gap:12px; align-items:center">
      <span>数据源</span>
      <a-select v-model="datasource" :options="dsOptions" style="width:200px" />
      <a-select v-if="datasource==='loki' && lokiDsOptions.length > 1" v-model="selectedLokiId" :options="lokiDsOptions" style="width:200px" placeholder="选择 Loki 数据源" />
      <a-select v-if="datasource==='elasticsearch' && esDsOptions.length > 1" v-model="selectedEsId" :options="esDsOptions" style="width:200px" placeholder="选择 ES 数据源" />
      <a-segmented v-model="mode" :options="['Builder','Code']" />
      <span>Range</span>
      <a-select v-model="range" :options="rangeOptions" style="width:140px" />
      <span>Step</span>
      <a-input v-model="step" placeholder="60s" style="width:100px" />
      <span>Direction</span>
      <a-select v-model="direction" :options="['BACKWARD','FORWARD']" style="width:120px" />
    </div>

    <loki-editor v-if="datasource==='loki'" :datasource-id="selectedLokiId" @run="onRunLoki" @history="openHistory" @inspect="openInspector" />
    <elasticsearch-editor v-else @run="onRunES" />

    <div v-if="rows.length > 0" style="margin-top:12px">
      <div style="margin-bottom:8px; color: #666;">查询结果: {{ rows.length }} 条记录</div>
      
      <!-- Debug: Show first few records -->
      <details style="margin-bottom:8px; font-size:12px; color:#666;">
        <summary>调试信息 (点击展开)</summary>
        <pre>{{ JSON.stringify(rows.slice(0, 2), null, 2) }}</pre>
      </details>
      
      <!-- 使用原生表格替代 Arco 表格 -->
      <div style="border: 1px solid #e5e6eb; border-radius: 4px; overflow: hidden;">
        <table style="width: 100%; border-collapse: collapse; font-size: 14px;">
          <thead style="background: #f7f8fa; border-bottom: 1px solid #e5e6eb;">
            <tr>
              <th style="padding: 12px; text-align: left; font-weight: 500; width: 220px; border-right: 1px solid #e5e6eb;">时间</th>
              <th style="padding: 12px; text-align: left; font-weight: 500; width: 100px; border-right: 1px solid #e5e6eb;">级别</th>
              <th style="padding: 12px; text-align: left; font-weight: 500;">内容</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(record, index) in paginatedRows" :key="index" 
                :style="{ backgroundColor: index % 2 === 0 ? '#fff' : '#fafafa' }"
                style="border-bottom: 1px solid #f0f0f0;">
              <td style="padding: 12px; border-right: 1px solid #f0f0f0; font-family: monospace; font-size: 12px;">
                {{ formatTimestamp(record.timestamp) }}
              </td>
              <td style="padding: 12px; border-right: 1px solid #f0f0f0;">
                {{ record.level || '-' }}
              </td>
              <td style="padding: 12px; word-break: break-all; max-width: 600px;">
                {{ record.message || '-' }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- 分页控制 -->
      <div v-if="rows.length > pageSize" style="margin-top: 16px; text-align: center;">
        <a-space>
          <a-button @click="prevPage" :disabled="currentPage === 1" size="small">上一页</a-button>
          <span style="margin: 0 16px; font-size: 14px;">
            第 {{ currentPage }} / {{ totalPages }} 页，共 {{ rows.length }} 条
          </span>
          <a-button @click="nextPage" :disabled="currentPage === totalPages" size="small">下一页</a-button>
        </a-space>
      </div>
    </div>
    <div v-else-if="!loading" style="margin-top:12px; padding:20px; text-align:center; color:#999; border:1px dashed #d9d9d9; border-radius:4px">
      暂无查询结果，请点击"运行查询"执行查询
    </div>

    <a-modal v-model:visible="historyVisible" title="查询历史记录" :footer="false" width="720px">
      <a-table :data="historyItems" row-key="id" :pagination="false">
        <a-table-column title="ID" data-index="id" :width="80" />
        <a-table-column title="Mode" data-index="mode" :width="100" />
        <a-table-column title="Engine" data-index="engine" :width="120" />
        <a-table-column title="Query" data-index="query" />
      </a-table>
    </a-modal>

    <a-modal v-model:visible="inspectVisible" title="查询检查器" :footer="false">
      <a-typography-paragraph copyable>
        {{ inspectUrl }}
      </a-typography-paragraph>
    </a-modal>
  </page-container>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import PageContainer from '@/components/PageContainer.vue'
import LokiEditor from '@/components/logs/LokiEditor.vue'
import ElasticsearchEditor from '@/components/logs/ElasticsearchEditor.vue'
import { queryLogs, history as apiHistory, inspect } from '@/api/logs'
import { listDataSources } from '@/api/datasources'

const datasource = ref('loki')
const dsOptions = [ { label: 'Loki', value: 'loki' }, { label: 'Elasticsearch', value: 'elasticsearch' } ]
const lokiDsOptions = ref([])
const esDsOptions = ref([])
const selectedLokiId = ref('')
const selectedEsId = ref('')
const mode = ref('Builder')
const rangeOptions = [
  { label: 'Last 5m', value: '5m' },
  { label: 'Last 15m', value: '15m' },
  { label: 'Last 1h', value: '1h' },
  { label: 'Last 6h', value: '6h' },
  { label: 'Last 24h', value: '24h' },
]
const range = ref('1h')
const step = ref('60s')
const direction = ref('BACKWARD')

const historyVisible = ref(false)
const historyItems = ref([])
const inspectVisible = ref(false)
const inspectUrl = ref('')

const loading = ref(false)
const rows = ref([])

// 分页相关
const currentPage = ref(1)
const pageSize = ref(20)

// 计算属性
const totalPages = computed(() => Math.ceil(rows.value.length / pageSize.value))
const paginatedRows = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return rows.value.slice(start, end)
})

// 分页方法
function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

function formatTimestamp(timestamp) {
  if (!timestamp) return '-'
  try {
    // Handle both string and number timestamps
    let ts = timestamp
    if (typeof ts === 'string') {
      ts = parseInt(ts)
    }
    // Convert nanoseconds to milliseconds if needed
    if (ts > 1e15) {
      ts = ts / 1e6
    }
    const date = new Date(ts)
    if (isNaN(date.getTime())) {
      return String(timestamp) // fallback to raw value
    }
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (e) {
    console.error('Error formatting timestamp:', timestamp, e)
    return String(timestamp)
  }
}

async function onRunLoki(payload) {
  await runQuery({ engine: 'loki', payload })
}
async function onRunES(payload) {
  await runQuery({ engine: 'elasticsearch', payload })
}

function computeTimeRange() {
  const now = Date.now()
  const map = { m: 60*1000, h: 60*60*1000 }
  const m = range.value.endsWith('m') ? map.m : map.h
  const num = parseInt(range.value)
  const startMs = now - num * m
  return { start: String(startMs * 1e6), end: String(now * 1e6) }
}

async function runQuery(params) {
  loading.value = true
  try {
    const { start, end } = computeTimeRange()
    const dsId = params.engine === 'loki' 
      ? (selectedLokiId.value || localStorage.getItem('last_loki_ds_id') || '') 
      : (selectedEsId.value || localStorage.getItem('last_es_ds_id') || '')
    
    console.log('Running query with:', { 
      engine: params.engine, 
      datasourceId: dsId, 
      start, 
      end, 
      step: step.value, 
      direction: direction.value,
      payload: params.payload 
    })
    
    const { data } = await queryLogs({ engine: params.engine, datasourceId: dsId, start, end, step: step.value, direction: direction.value, ...params.payload })
    console.log('API Response:', data)
    rows.value = data?.data?.items || []
    currentPage.value = 1 // 重置到第一页
    console.log('Rows after setting:', rows.value.length, 'items, first few:', rows.value.slice(0, 2))
  } catch (error) {
    console.error('Query error:', error)
  } finally {
    loading.value = false
  }
}

async function openHistory() {
  historyVisible.value = true
  const { data } = await apiHistory()
  historyItems.value = data?.data?.items || []
}

async function openInspector(queryStr = '') {
  inspectVisible.value = true
  const { start, end } = computeTimeRange()
  const dsId = selectedLokiId.value || localStorage.getItem('last_loki_ds_id') || ''
  const params = { engine: 'loki', datasourceId: dsId, start, end, step: step.value, direction: direction.value }
  if (queryStr) params.query = queryStr
  const { data } = await inspect(params)
  inspectUrl.value = data?.data?.url || ''
}

onMounted(async () => {
  try {
    const { data } = await listDataSources()
    const items = data?.data?.items || []
    console.log('Loaded datasources:', items)
    
    lokiDsOptions.value = items.filter(x => x.type === 'loki').map(x => ({ label: x.name, value: String(x.id) }))
    esDsOptions.value = items.filter(x => x.type === 'elasticsearch').map(x => ({ label: x.name, value: String(x.id) }))
    
    console.log('Loki datasources:', lokiDsOptions.value)
    console.log('ES datasources:', esDsOptions.value)
    
    if (!selectedLokiId.value && lokiDsOptions.value.length) {
      selectedLokiId.value = lokiDsOptions.value[0].value
      localStorage.setItem('last_loki_ds_id', selectedLokiId.value)
      console.log('Selected Loki datasource:', selectedLokiId.value)
    }
    if (!selectedEsId.value && esDsOptions.value.length) {
      selectedEsId.value = esDsOptions.value[0].value
      localStorage.setItem('last_es_ds_id', selectedEsId.value)
      console.log('Selected ES datasource:', selectedEsId.value)
    }
  } catch (e) {
    console.error('Failed to load datasources:', e)
  }
})
</script>

