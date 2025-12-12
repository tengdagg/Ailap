<template>
  <page-container>
    <div style="margin-bottom:12px; display:flex; gap:12px; align-items:center">
      <span>{{ $t('logs.datasource') }}</span>
      <a-select v-model="datasource" :options="dsOptions" style="width:200px" />
      <a-select v-if="datasource==='loki' && lokiDsOptions.length > 1" v-model="selectedLokiId" :options="lokiDsOptions" style="width:200px" :placeholder="$t('logs.selectLoki')" />
      <a-select v-if="datasource==='elasticsearch' && esDsOptions.length > 1" v-model="selectedEsId" :options="esDsOptions" style="width:200px" :placeholder="$t('logs.selectES')" />
      <a-select v-if="datasource==='victorialogs' && vlDsOptions.length > 1" v-model="selectedVlId" :options="vlDsOptions" style="width:200px" :placeholder="$t('logs.selectVL')" />
      <a-segmented v-model="mode" :options="['Builder','Code']" v-if="datasource!=='victorialogs'" />
      <a-segmented v-model="mode" :options="['Code']" v-else disabled />
      <span>{{ $t('logs.range') }}</span>
      <a-select v-model="range" :options="rangeOptions" style="width:140px" />
      <span>{{ $t('logs.step') }}</span>
      <a-input v-model="step" placeholder="60s" style="width:100px" />
      <span>{{ $t('logs.direction') }}</span>
      <a-select v-model="direction" :options="['BACKWARD','FORWARD']" style="width:120px" />
    </div>

    <loki-editor v-if="datasource==='loki'" :datasource-id="selectedLokiId" @run="onRunLoki" @history="openHistory" @inspect="openInspector" />
    <elasticsearch-editor v-else-if="datasource==='elasticsearch'" @run="onRunES" @history="openHistory" @inspect="openInspector" />
    <victoria-logs-editor v-else-if="datasource==='victorialogs'" :datasource-id="selectedVlId" @run="onRunVL" @history="openHistory" @inspect="openInspector" />

    <div v-if="rows.length > 0 && viewMode==='logs'" style="margin-top:12px">
      <div style="margin-bottom:8px; color: var(--color-text-3);">{{ $t('logs.queryResults', { count: rows.length }) }}</div>

      <div style="border: 1px solid var(--color-border-2); border-radius: 4px; overflow: auto; max-height: calc(100vh - 360px);">
        <table style="width: 100%; border-collapse: collapse; font-size: 13px;">
          <thead style="background: var(--color-fill-2); border-bottom: 1px solid var(--color-border-2); position: sticky; top: 0; z-index: 1;">
            <tr>
              <th style="padding: 8px 12px; text-align: left; font-weight: 500; width: 180px; border-right: 1px solid var(--color-border-2);">
                <div style="display:flex; align-items:center; justify-content:space-between">
                  {{ $t('logs.time') }}
                  <a-popover trigger="click" position="bottom">
                    <icon-filter :style="{ color: filters.time ? 'rgb(var(--primary-6))' : 'var(--color-text-3)', cursor: 'pointer' }" />
                    <template #content>
                      <div style="width:200px">
                        <a-input v-model="filters.time" :placeholder="$t('logs.searchTime')" size="small" allow-clear />
                      </div>
                    </template>
                  </a-popover>
                </div>
              </th>
              <th style="padding: 8px 12px; text-align: left; font-weight: 500;">
                <div style="display:flex; align-items:center; justify-content:space-between">
                  {{ $t('logs.content') }}
                  <a-popover trigger="click" position="bottom">
                    <icon-filter :style="{ color: filters.content ? 'rgb(var(--primary-6))' : 'var(--color-text-3)', cursor: 'pointer' }" />
                    <template #content>
                      <div style="width:250px">
                        <a-input v-model="filters.content" :placeholder="$t('logs.searchContent')" size="small" allow-clear />
                      </div>
                    </template>
                  </a-popover>
                </div>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(record, index) in paginatedRows" :key="index" 
                :style="{ backgroundColor: index % 2 === 0 ? 'var(--color-bg-1)' : 'var(--color-fill-1)' }"
                style="border-bottom: 1px solid var(--color-border-2);">
              <td style="padding: 8px 12px; border-right: 1px solid var(--color-border-2); font-family: monospace; vertical-align: top; color: var(--color-text-3);">
                {{ record.parsed.time }}
              </td>
              <td style="padding: 8px 12px; word-break: break-all; font-family: monospace; vertical-align: top;">
                <!-- Tags for structured logs -->
                <div v-if="record.parsed.method !== '-' || record.parsed.status !== '-'" style="margin-bottom: 4px; display: flex; gap: 6px; align-items: center;">
                   <a-tag v-if="record.parsed.method !== '-'" size="small" :color="getMethodColor(record.parsed.method)">{{ record.parsed.method }}</a-tag>
                   <a-tag v-if="record.parsed.status !== '-'" size="small" :color="getStatusColor(record.parsed.status)">{{ record.parsed.status }}</a-tag>
                   <span v-if="record.parsed.host !== '-'" style="color: var(--color-text-3); font-size: 12px;">{{ record.parsed.host }}</span>
                   <span v-if="record.parsed.source !== '-'" style="color: var(--color-text-3); font-size: 12px;">From: {{ record.parsed.source }}</span>
                </div>
                <!-- Raw Message -->
                <div style="line-height: 1.5;">{{ record.message }}</div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-if="rows.length > pageSize" style="margin-top: 16px; text-align: center;">
        <a-space>
          <a-button @click="prevPage" :disabled="currentPage === 1" size="small">{{ $t('logs.prevPage') }}</a-button>
          <span style="margin: 0 16px; font-size: 14px;">
            {{ $t('logs.pageInfo', { current: currentPage, total: totalPages, count: rows.length }) }}
          </span>
          <a-button @click="nextPage" :disabled="currentPage === totalPages" size="small">{{ $t('logs.nextPage') }}</a-button>
        </a-space>
      </div>
    </div>

    <div v-else-if="rows.length > 0 && viewMode==='raw'" style="margin-top:12px">
      <div style="margin-bottom:8px; color: var(--color-text-3);">{{ $t('logs.rawData', { count: rows.length }) }}</div>
      <div style="border: 1px solid var(--color-border-2); border-radius: 4px; overflow: auto; max-height: calc(100vh - 320px);">
        <table style="width: 100%; border-collapse: collapse; font-size: 13px;">
          <thead style="background: var(--color-fill-2); border-bottom: 1px solid var(--color-border-2); position: sticky; top: 0; z-index: 1;">
            <tr>
              <th v-for="(col, cidx) in rawColumns" :key="cidx" style="padding: 8px; text-align: left; font-weight: 500; border-right: 1px solid var(--color-border-2); white-space: nowrap;">{{ col }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(doc, ridx) in paginatedRows" :key="ridx" :style="{ backgroundColor: ridx % 2 === 0 ? 'var(--color-bg-1)' : 'var(--color-fill-1)' }" style="border-bottom: 1px solid var(--color-border-2);">
              <td v-for="(col, cidx) in rawColumns" :key="cidx" style="padding: 8px; border-right: 1px solid var(--color-border-2); font-family: monospace;">
                {{ formatRawCell(doc.__raw, col) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div v-else-if="!loading" style="margin-top:12px; padding:20px; text-align:center; color:var(--color-text-3); border:1px dashed var(--color-border-2); border-radius:4px">
      {{ $t('logs.noResults') }}
    </div>

     <!-- 历史记录抽屉 -->
    <a-drawer v-model:visible="historyVisible" title="" width="600px" placement="bottom" :height="450" :footer="false">
      <!-- 搜索框 -->
      <div style="margin-bottom: 12px;">
        <a-input
          v-model="searchKeyword"
          :placeholder="$t('logs.searchHistory')"
          allow-clear
          @input="onSearchInput"
        >
          <template #prefix>
            <icon-search />
          </template>
        </a-input>
      </div>
      
      <div style="margin-bottom: 16px;">
        <a-tabs v-model:active-key="historyTab" type="line">
          <a-tab-pane key="recent" :title="$t('logs.historyTitle')" />
          <a-tab-pane key="favorite" :title="$t('logs.favoriteTitle')" />
        </a-tabs>
      </div>

      <div v-if="historyItems.length === 0" style="text-align: center; padding: 40px; color: var(--color-text-3);">
        {{ historyTab === 'favorite' ? $t('logs.noFavorites') : $t('logs.noHistory') }}
      </div>
      
      <div v-else style="max-height: 280px; overflow-y: auto;">
        <div v-for="item in historyItems" :key="item.id" 
             style="border: 1px solid var(--color-border-2); border-radius: 6px; padding: 12px; margin-bottom: 8px; background: var(--color-fill-1);">
          
          <!-- 时间和操作按钮在同一行 -->
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; font-size: 12px; color: var(--color-text-3);">
            <div style="display: flex; gap: 12px; align-items: center;">
              <span>{{ new Date(item.createdAt).toLocaleString() }}</span>
              <a-tag :color="item.engine === 'loki' ? 'blue' : (item.engine === 'elasticsearch' ? 'green' : 'orange')" size="small">{{ item.engine }}</a-tag>
              <a-tag color="gray" size="small">{{ item.mode }}</a-tag>
            </div>
            
            <!-- 操作按钮组 -->
            <div style="display: flex; gap: 4px;">
              <a-tooltip :content="$t('logs.editNote')">
                <a-button 
                  size="mini" 
                  type="text" 
                  @click="showNoteModal(item)"
                  style="width: 20px; height: 20px; padding: 0; display: flex; align-items: center; justify-content: center;"
                >
                  <icon-tag style="color: var(--color-primary-6) !important;" />
                </a-button>
              </a-tooltip>
              <a-tooltip :content="$t('common.delete')">
                <a-button 
                  size="mini" 
                  type="text" 
                  @click="confirmDelete(item)"
                  style="width: 20px; height: 20px; padding: 0; display: flex; align-items: center; justify-content: center;"
                >
                  <icon-delete style="color: #f53f3f !important;" />
                </a-button>
              </a-tooltip>
              <a-tooltip :content="item.isFavorite ? $t('logs.removeFavorite') : $t('logs.addFavorite')">
                <a-button 
                  size="mini" 
                  type="text" 
                  style="width: 20px; height: 20px; padding: 0; display: flex; align-items: center; justify-content: center;"
                  @click="toggleQueryFavorite(item)"
                >
                  <icon-star-fill v-if="item.isFavorite" style="color: #ff7d00 !important;" />
                  <icon-star v-else style="color: var(--color-text-3);" />
                </a-button>
              </a-tooltip>
              <a-tooltip :content="$t('logs.runQuery')">
                <a-button 
                  size="mini" 
                  type="text" 
                  @click="executeQuery(item)"
                  style="width: 20px; height: 20px; padding: 0; display: flex; align-items: center; justify-content: center;"
                >
                  <icon-send style="color: var(--color-primary-6) !important;" />
                </a-button>
              </a-tooltip>
            </div>
          </div>
          
          <!-- 备注显示 -->
          <div v-if="item.note" style="margin-bottom: 8px; font-size: 12px; color: var(--color-text-3); font-style: italic;">
            💬 {{ item.note }}
          </div>
          
          <div style="font-family: monospace; font-size: 13px; word-break: break-all; background: var(--color-fill-2); padding: 8px; border-radius: 4px; border: 1px solid var(--color-border-2);">
            {{ item.query || '-' }}
          </div>
        </div>
      </div>
    </a-drawer>

    <!-- 备注编辑模态框 -->
    <a-modal v-model:visible="noteModalVisible" :title="$t('logs.editNote')" @ok="saveNote" @cancel="cancelNote">
      <a-textarea 
        v-model="noteContent" 
        :placeholder="$t('logs.notePlaceholder')"
        :rows="3"
        :max-length="200"
        show-word-limit
      />
    </a-modal>

    <a-modal v-model:visible="inspectVisible" :title="$t('logs.inspector')" :footer="false">
      <a-typography-paragraph copyable>
        {{ inspectUrl }}
      </a-typography-paragraph>
      <div style="margin-top: 8px; color: var(--color-text-3);">Body</div>
      <pre style="white-space: pre; background: var(--color-fill-2); padding: 12px; border-radius: 4px; border: 1px solid var(--color-border-2); max-height: 300px; overflow: auto;">
{{ inspectBody }}
      </pre>
    </a-modal>

    <!-- 智能分析悬浮按钮与对话框 -->
    <log-analysis-chat v-if="rows.length > 0" :logs="rows" :initial-range="{ start: lastRangeStartMs, end: lastRangeEndMs }" />
  </page-container>
</template>
<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import PageContainer from '@/components/PageContainer.vue'
import LokiEditor from '@/components/logs/LokiEditor.vue'
import ElasticsearchEditor from '@/components/logs/ElasticsearchEditor.vue'
import VictoriaLogsEditor from '@/components/logs/VictoriaLogsEditor.vue'
import { queryLogs, history as apiHistory, inspect, toggleFavorite, updateNote, deleteHistory } from '@/api/logs'
import { listDataSources } from '@/api/datasources'
import { Message, Modal } from '@arco-design/web-vue'
import { useI18n } from 'vue-i18n'
import { IconTag, IconDelete, IconStar, IconStarFill, IconSend, IconSearch, IconFilter } from '@arco-design/web-vue/es/icon'
import LogAnalysisChat from '@/components/LogAnalysisChat.vue'

  const { t } = useI18n()
  const datasource = ref('')
watch(datasource, () => {
  rows.value = []
  viewMode.value = 'logs'
  // When switching engines, ensure we are on page 1
  currentPage.value = 1
})
const lokiDsOptions = ref([])
const esDsOptions = ref([])
const vlDsOptions = ref([])

// Computed property to dynamically show only configured data sources
const dsOptions = computed(() => {
  const options = []
  if (lokiDsOptions.value.length > 0) {
    options.push({ label: 'Loki', value: 'loki' })
  }
  if (esDsOptions.value.length > 0) {
    options.push({ label: 'Elasticsearch', value: 'elasticsearch' })
  }
  if (vlDsOptions.value.length > 0) {
    options.push({ label: 'VictoriaLogs', value: 'victorialogs' })
  }
  return options
})
const selectedLokiId = ref('')
const selectedEsId = ref('')
const selectedVlId = ref('')
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
const lastRangeStartMs = ref(0)
const lastRangeEndMs = ref(0)

const historyVisible = ref(false)
const historyTab = ref('recent')
const historyItems = ref([])
const allHistoryItems = ref([]) // 存储所有历史记录
const searchKeyword = ref('')
const inspectVisible = ref(false)
const inspectUrl = ref('')
const inspectBody = ref('')

// 备注编辑相关
const noteModalVisible = ref(false)
const noteContent = ref('')
const currentEditItem = ref(null)

const loading = ref(false)
const rows = ref([])
const viewMode = ref('logs') // 'logs' | 'raw'
const rawColumns = ref([])

// 分页相关
const currentPage = ref(1)
const pageSize = ref(20)

const filters = ref({
  source: '',
  time: '',
  host: '',
  method: '',
  status: '',
  duration: '',
  upstream: '',
  agent: '',
  content: ''
})

const parsedRows = computed(() => {
  return rows.value.map(record => {
    const parsed = parseLogMessage(record.message)
    return {
      ...record,
      parsed: parsed || {
        source: '-',
        time: formatTimestamp(record.timestamp),
        host: '-',
        path: '-',
        method: '-',
        status: '-',
        duration: '-',
        upstream: '-',
        agent: record.message
      }
    }
  })
})

const filteredRows = computed(() => {
  return parsedRows.value.filter(row => {
    const p = row.parsed
    if (filters.value.source && !p.source.toLowerCase().includes(filters.value.source.toLowerCase())) return false
    if (filters.value.time && !p.time.toLowerCase().includes(filters.value.time.toLowerCase())) return false
    if (filters.value.host && !p.host.toLowerCase().includes(filters.value.host.toLowerCase()) && !p.path.toLowerCase().includes(filters.value.host.toLowerCase())) return false
    if (filters.value.method && !p.method.toLowerCase().includes(filters.value.method.toLowerCase())) return false
    if (filters.value.status && !p.status.toLowerCase().includes(filters.value.status.toLowerCase())) return false
    if (filters.value.duration && !p.duration.toLowerCase().includes(filters.value.duration.toLowerCase())) return false
    if (filters.value.upstream && !p.upstream.toLowerCase().includes(filters.value.upstream.toLowerCase())) return false
    if (filters.value.agent && !p.agent.toLowerCase().includes(filters.value.agent.toLowerCase())) return false
    if (filters.value.content && !row.message.toLowerCase().includes(filters.value.content.toLowerCase())) return false
    return true
  })
})

// 获取列的唯一值
function getUniqueValues(key) {
  const values = new Set()
  parsedRows.value.forEach(row => {
    if (row.parsed[key]) {
      values.add(row.parsed[key])
    }
  })
  return Array.from(values).sort()
}

function setFilter(key, val) {
  filters.value[key] = val
}

// 计算属性
const totalPages = computed(() => Math.ceil(filteredRows.value.length / pageSize.value))
const paginatedRows = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredRows.value.slice(start, end)
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
    // Helper function to format date to YYYY-MM-DD HH:mm:ss
    const formatDate = (date) => {
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      const hours = String(date.getHours()).padStart(2, '0')
      const minutes = String(date.getMinutes()).padStart(2, '0')
      const seconds = String(date.getSeconds()).padStart(2, '0')
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
    }

    // Numeric timestamps (number or numeric string)
    if (typeof timestamp === 'number' || (typeof timestamp === 'string' && /^\d+$/.test(timestamp))) {
      let ts = typeof timestamp === 'number' ? timestamp : parseInt(timestamp)
      // Convert nanoseconds to milliseconds if needed
      if (ts > 1e15) {
        ts = Math.floor(ts / 1e6)
      }
      const dateNum = new Date(ts)
      if (!isNaN(dateNum.getTime())) {
        return formatDate(dateNum)
      }
    }

    // ISO8601 string timestamps (e.g., 2025-09-16T14:29:15.609+08:00)
    if (typeof timestamp === 'string') {
      const dateIso = new Date(timestamp)
      if (!isNaN(dateIso.getTime())) {
        return formatDate(dateIso)
      }
    }

    // Fallback
    return String(timestamp)
  } catch (_) {
    return String(timestamp)
  }
}

async function onRunLoki(payload) {
  // Validation
  if (payload.mode === 'code') {
    if (!payload.query || !payload.query.trim()) {
      Message.warning(t('logs.inputQuery'))
      return
    }
  } else if (payload.mode === 'builder') {
    const hasFilter = payload.builder.labelFilters.some(f => f.label && f.values && f.values.length > 0)
    const hasContains = payload.builder.contains && payload.builder.contains.trim()
    if (!hasFilter && !hasContains) {
      Message.warning(t('logs.selectCondition'))
      return
    }
  }
  await runQuery({ engine: 'loki', payload })
}

async function onRunES(payload) {
  // Validation
  if (!payload.query || !payload.query.trim()) {
    Message.warning(t('logs.inputQuery'))
    return
  }
  viewMode.value = payload?.mode === 'raw' ? 'raw' : 'logs'
  await runQuery({ engine: 'elasticsearch', payload })
}

async function onRunVL(payload) {
    if (!payload.query || !payload.query.trim()) {
        Message.warning(t('logs.inputQuery'))
        return
    }
    viewMode.value = 'logs'
    await runQuery({ engine: 'victorialogs', payload })
}

function computeTimeRange() {
  const now = Date.now()
  const map = { m: 60*1000, h: 60*60*1000 }
  const m = range.value.endsWith('m') ? map.m : map.h
  const num = parseInt(range.value)
  const startMs = now - num * m
  return { start: String(startMs * 1e6), end: String(now * 1e6), startMs, nowMs: now }
}

async function runQuery(params) {
  loading.value = true
  try {
    const { start, end, startMs, nowMs } = computeTimeRange()
    lastRangeStartMs.value = startMs
    lastRangeEndMs.value = nowMs
    const dsId = params.engine === 'loki' 
      ? (selectedLokiId.value || localStorage.getItem('last_loki_ds_id') || '') 
      : (params.engine === 'elasticsearch' 
         ? (selectedEsId.value || localStorage.getItem('last_es_ds_id') || '')
         : (selectedVlId.value || localStorage.getItem('last_vl_ds_id') || ''))
    
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
    const items = data?.data?.items || []
    if (viewMode.value === 'raw') {
      // Expect backend sends each row with __raw being full _source
      // Fallback: build rawColumns from keys present
      const cols = new Set()
      items.forEach(it => {
        if (it.__raw && typeof it.__raw === 'object') {
          Object.keys(it.__raw).forEach(k => cols.add(k))
        }
      })
      rawColumns.value = Array.from(cols)
    }
    rows.value = items
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
  await loadHistoryData()
}

async function loadHistoryData() {
  try {
    const { data } = await apiHistory(historyTab.value)
    allHistoryItems.value = data?.data?.items || []
    filterHistoryItems() // 应用搜索过滤
  } catch (error) {
    console.error('Failed to load history:', error)
    allHistoryItems.value = []
    historyItems.value = []
  }
}

// 根据搜索关键词过滤历史记录
function filterHistoryItems() {
  if (!searchKeyword.value.trim()) {
    historyItems.value = allHistoryItems.value
    return
  }
  
  const keyword = searchKeyword.value.toLowerCase()
  historyItems.value = allHistoryItems.value.filter(item => {
    return (
      item.query?.toLowerCase().includes(keyword) ||
      item.note?.toLowerCase().includes(keyword) ||
      item.engine?.toLowerCase().includes(keyword) ||
      item.mode?.toLowerCase().includes(keyword)
    )
  })
}

// 搜索输入处理
function onSearchInput() {
  filterHistoryItems()
}

async function toggleQueryFavorite(item) {
  try {
    const { data } = await toggleFavorite(item.id)
    if (data?.code === 0) {
      // 更新本地状态
      item.isFavorite = data.data.item.isFavorite
      // 如果当前在收藏页面且取消收藏，则重新加载数据
      if (historyTab.value === 'favorite' && !item.isFavorite) {
        await loadHistoryData()
      }
    }
  } catch (error) {
    console.error('Failed to toggle favorite:', error)
  }
}

function useQuery(item) {
  // 根据引擎类型设置对应的查询内容
  if (item.engine === 'loki') {
    datasource.value = 'loki'
    // 这里可以进一步设置 LokiEditor 的查询内容
  } else if (item.engine === 'elasticsearch') {
    datasource.value = 'elasticsearch'
    // 这里可以进一步设置 ElasticsearchEditor 的查询内容
  } else if (item.engine === 'victorialogs') {
    datasource.value = 'victorialogs'
  }
  historyVisible.value = false
}

// 显示备注编辑模态框
function showNoteModal(item) {
  currentEditItem.value = item
  noteContent.value = item.note || ''
  noteModalVisible.value = true
}

// 保存备注
async function saveNote() {
  if (!currentEditItem.value) return
  
  try {
    const { data } = await updateNote(currentEditItem.value.id, noteContent.value)
    if (data?.code === 0) {
      currentEditItem.value.note = noteContent.value
      Message.success(t('common.saveSuccess'))
      noteModalVisible.value = false
    }
  } catch (error) {
    console.error('Failed to save note:', error)
    Message.error(t('common.saveFail'))
  }
}

// 取消编辑备注
function cancelNote() {
  noteModalVisible.value = false
  noteContent.value = ''
  currentEditItem.value = null
}

// 确认删除
function confirmDelete(item) {
  Modal.confirm({
    title: t('common.deleteConfirm'),
    content: t('logs.confirmDeleteHist'),
    onOk: () => deleteHistoryItem(item)
  })
}

// 删除历史记录
async function deleteHistoryItem(item) {
  try {
    const { data } = await deleteHistory(item.id)
    if (data?.code === 0) {
      Message.success(t('common.deleteSuccess'))
      await loadHistoryData() // 重新加载数据
    }
  } catch (error) {
    console.error('Failed to delete history:', error)
    Message.error(t('common.deleteFail'))
  }
}

// 执行查询
async function executeQuery(item) {
  try {
    // 设置数据源
    if (item.engine === 'loki') {
      datasource.value = 'loki'
    } else if (item.engine === 'elasticsearch') {
      datasource.value = 'elasticsearch'
    } else if (item.engine === 'victorialogs') {
      datasource.value = 'victorialogs'
    }
    
    // 关闭抽屉
    historyVisible.value = false
    
    // 构造查询参数并执行
    const { start, end } = computeTimeRange()
    let dsId = ''
    if (item.engine === 'loki') {
      dsId = selectedLokiId.value || localStorage.getItem('last_loki_ds_id') || ''
    } else if (item.engine === 'elasticsearch') {
      dsId = selectedEsId.value || localStorage.getItem('last_es_ds_id') || ''
    } else {
      dsId = selectedVlId.value || localStorage.getItem('last_vl_ds_id') || ''
    }
    
    const params = {
      engine: item.engine,
      datasourceId: dsId,
      start,
      end,
      step: step.value,
      direction: direction.value,
      mode: item.mode,
      query: item.query,
      lineLimit: item.lineLimit || 1000
    }
    
    loading.value = true
    const { data } = await queryLogs(params)
    rows.value = data?.data?.items || []
    currentPage.value = 1
    
    Message.success(t('logs.querySuccess'))
  } catch (error) {
    console.error('Execute query error:', error)
    Message.error(t('logs.queryFail'))
  } finally {
    loading.value = false
  }
}

// 监听 tab 切换
watch(historyTab, async () => {
  if (historyVisible.value) {
    await loadHistoryData()
  }
})

// 清空搜索关键词当抽屉关闭时
watch(historyVisible, (visible) => {
  if (!visible) {
    searchKeyword.value = ''
  }
})

async function openInspector(queryStr = '') {
  inspectVisible.value = true
  const { start, end } = computeTimeRange()
  
  let params
  if (datasource.value === 'loki') {
    const dsId = selectedLokiId.value || localStorage.getItem('last_loki_ds_id') || ''
    params = { engine: 'loki', datasourceId: dsId, start, end, step: step.value, direction: direction.value }
  } else if (datasource.value === 'elasticsearch') {
    const dsId = selectedEsId.value || localStorage.getItem('last_es_ds_id') || ''
    params = { engine: 'elasticsearch', datasourceId: dsId, start, end }
  } else {
    const dsId = selectedVlId.value || localStorage.getItem('last_vl_ds_id') || ''
    params = { engine: 'victorialogs', datasourceId: dsId, start, end }
  }
  
  if (queryStr) params.query = queryStr
  const { data } = await inspect(params)
  inspectUrl.value = data?.data?.url || ''
  inspectBody.value = data?.data?.body || ''
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
    lokiDsOptions.value = items.filter(x => x.type === 'loki').map(x => ({ label: x.name, value: String(x.id) }))
    esDsOptions.value = items.filter(x => x.type === 'elasticsearch').map(x => ({ label: x.name, value: String(x.id) }))
    vlDsOptions.value = items.filter(x => x.type === 'victorialogs').map(x => ({ label: x.name, value: String(x.id) }))
    
    console.log('Loki datasources:', lokiDsOptions.value)
    console.log('ES datasources:', esDsOptions.value)
    console.log('VL datasources:', vlDsOptions.value)
    
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
    if (!selectedVlId.value && vlDsOptions.value.length) {
      selectedVlId.value = vlDsOptions.value[0].value
      localStorage.setItem('last_vl_ds_id', selectedVlId.value)
      console.log('Selected VL datasource:', selectedVlId.value)
    }
    
    // Auto-select first available data source type
    if (!datasource.value) {
      if (lokiDsOptions.value.length > 0) {
        datasource.value = 'loki'
      } else if (esDsOptions.value.length > 0) {
        datasource.value = 'elasticsearch'
      } else if (vlDsOptions.value.length > 0) {
        datasource.value = 'victorialogs'
      }
    }
  } catch (e) {
    console.error('Failed to load datasources:', e)
  }
})

function formatRawCell(raw, key) {
  const val = raw && raw[key]
  if (val == null) return ''
  if (typeof val === 'object') return JSON.stringify(val)
  return String(val)
}

function parseLogMessage(msg) {
  if (!msg) return null
  
  // Clean up message if it starts with message="...
  let cleanMsg = msg
  if (cleanMsg.startsWith('message="')) {
    cleanMsg = cleanMsg.substring(9)
    if (cleanMsg.endsWith('"')) {
      cleanMsg = cleanMsg.substring(0, cleanMsg.length - 1)
    }
  } else if (cleanMsg.startsWith('message=')) {
    cleanMsg = cleanMsg.substring(8)
  }

  // Example: 172.21.8.88 - [26/Nov/2025:09:51:51 +0800] "zabbix6.aaa.com" "GET /path HTTP/1.1" 200 28748 "referer" "ua" "-" unix:/run/php-fpm/zabbix.sock 0.331 0.330
  // Regex to capture: IP, Time, Host, Request, Status, Size, Referer, UA, XFF, Upstream, Duration, UpstreamTime
  const regex = /^(\S+) \S+ \[(.*?)\] "(\S+)" "(.*?)" (\d+) (\d+) "(.*?)" "(.*?)" "(.*?)" (\S+) (\S+) (\S+)/
  const match = cleanMsg.match(regex)
  if (!match) return null
  
  const [_, source, time, host, request, status, size, referer, agent, xff, upstream, duration, upstreamTime] = match
  const [method, path] = request.split(' ')
  
  // Convert Nginx time format (e.g., "03/Dec/2025:14:01:40 +0800") to YYYY-MM-DD HH:mm:ss
  const formatNginxTime = (nginxTime) => {
    try {
      // Parse format: DD/MMM/YYYY:HH:mm:ss +ZZZZ
      const match = nginxTime.match(/(\d{2})\/(\w{3})\/(\d{4}):(\d{2}):(\d{2}):(\d{2})/)
      if (!match) return nginxTime
      
      const [, day, monthStr, year, hours, minutes, seconds] = match
      const monthMap = {
        'Jan': '01', 'Feb': '02', 'Mar': '03', 'Apr': '04',
        'May': '05', 'Jun': '06', 'Jul': '07', 'Aug': '08',
        'Sep': '09', 'Oct': '10', 'Nov': '11', 'Dec': '12'
      }
      const month = monthMap[monthStr] || '01'
      
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
    } catch (e) {
      return nginxTime
    }
  }
  
  return {
    source,
    time: formatNginxTime(time),
    host,
    method: method || '-',
    path: path || '-',
    agent: agent === '-' ? '' : agent,
    status,
    upstream: upstream === '-' ? '' : upstream,
    duration
  }
}

function getMethodColor(method) {
  if (!method) return 'gray'
  const m = method.toUpperCase()
  if (m === 'GET') return 'blue'
  if (m === 'POST') return 'green'
  if (m === 'PUT') return 'orange'
  if (m === 'DELETE') return 'red'
  return 'gray'
}

function getStatusColor(status) {
  if (!status) return 'gray'
  const s = parseInt(status)
  if (s >= 200 && s < 300) return 'green'
  if (s >= 300 && s < 400) return 'blue'
  if (s >= 400 && s < 500) return 'orange'
  if (s >= 500) return 'red'
  return 'gray'
}
</script>

<style scoped>
.filter-item:hover {
  background-color: var(--color-fill-2);
}
</style>

