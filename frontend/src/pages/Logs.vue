<template>
  <page-container>
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
    <elasticsearch-editor v-else @run="onRunES" @history="openHistory" @inspect="openInspector" />

    <div v-if="rows.length > 0 && viewMode==='logs'" style="margin-top:12px">
      <div style="margin-bottom:8px; color: var(--color-text-3);">查询结果: {{ rows.length }} 条记录</div>
      <div style="border: 1px solid var(--color-border-2); border-radius: 4px; overflow: auto; max-height: calc(100vh - 320px);">
        <table style="width: 100%; border-collapse: collapse; font-size: 13px;">
          <thead style="background: var(--color-fill-2); border-bottom: 1px solid var(--color-border-2); position: sticky; top: 0; z-index: 1;">
            <tr>
              <th style="padding: 8px 12px; text-align: left; font-weight: 500; width: 140px; border-right: 1px solid var(--color-border-2);">
                <div style="display:flex; align-items:center; justify-content:space-between">
                  源地址
                  <a-popover trigger="click" position="bottom">
                    <icon-filter :style="{ color: filters.source ? 'rgb(var(--primary-6))' : 'var(--color-text-3)', cursor: 'pointer' }" />
                    <template #content>
                      <div style="width:200px">
                        <a-input v-model="filters.source" placeholder="搜索源地址..." size="small" allow-clear />
                        <div style="margin-top:8px; max-height:150px; overflow-y:auto; border-top:1px solid var(--color-border-1); padding-top:4px">
                          <div v-for="val in getUniqueValues('source')" :key="val" 
                               @click="setFilter('source', val)"
                               style="padding:4px 8px; cursor:pointer; font-size:12px; border-radius:4px"
                               class="filter-item">
                            {{ val }}
                          </div>
                        </div>
                      </div>
                    </template>
                  </a-popover>
                </div>
              </th>
              <th style="padding: 8px 12px; text-align: left; font-weight: 500; width: 170px; border-right: 1px solid var(--color-border-2);">
                <div style="display:flex; align-items:center; justify-content:space-between">
                  时间
                  <a-popover trigger="click" position="bottom">
                    <icon-filter :style="{ color: filters.time ? 'rgb(var(--primary-6))' : 'var(--color-text-3)', cursor: 'pointer' }" />
                    <template #content>
                      <div style="width:200px">
                        <a-input v-model="filters.time" placeholder="搜索时间..." size="small" allow-clear />
                        <div style="margin-top:8px; max-height:150px; overflow-y:auto; border-top:1px solid var(--color-border-1); padding-top:4px">
                          <div v-for="val in getUniqueValues('time')" :key="val" 
                               @click="setFilter('time', val)"
                               style="padding:4px 8px; cursor:pointer; font-size:12px; border-radius:4px"
                               class="filter-item">
                            {{ val }}
                          </div>
                        </div>
                      </div>
                    </template>
                  </a-popover>
                </div>
              </th>
              <th style="padding: 8px 12px; text-align: left; font-weight: 500; width: 200px; border-right: 1px solid var(--color-border-2);">
                <div style="display:flex; align-items:center; justify-content:space-between">
                  请求地址
                  <a-popover trigger="click" position="bottom">
                    <icon-filter :style="{ color: filters.host ? 'rgb(var(--primary-6))' : 'var(--color-text-3)', cursor: 'pointer' }" />
                    <template #content>
                      <div style="width:200px">
                        <a-input v-model="filters.host" placeholder="搜索地址..." size="small" allow-clear />
                        <div style="margin-top:8px; max-height:150px; overflow-y:auto; border-top:1px solid var(--color-border-1); padding-top:4px">
                          <div v-for="val in getUniqueValues('host')" :key="val" 
                               @click="setFilter('host', val)"
                               style="padding:4px 8px; cursor:pointer; font-size:12px; border-radius:4px"
                               class="filter-item">
                            {{ val }}
                          </div>
                        </div>
                      </div>
                    </template>
                  </a-popover>
                </div>
              </th>
              <th style="padding: 8px 12px; text-align: left; font-weight: 500; width: 80px; border-right: 1px solid var(--color-border-2);">
                <div style="display:flex; align-items:center; justify-content:space-between">
                  方式
                  <a-popover trigger="click" position="bottom">
                    <icon-filter :style="{ color: filters.method ? 'rgb(var(--primary-6))' : 'var(--color-text-3)', cursor: 'pointer' }" />
                    <template #content>
                      <div style="width:120px">
                        <a-input v-model="filters.method" placeholder="搜索方式..." size="small" allow-clear />
                        <div style="margin-top:8px; max-height:150px; overflow-y:auto; border-top:1px solid var(--color-border-1); padding-top:4px">
                          <div v-for="val in getUniqueValues('method')" :key="val" 
                               @click="setFilter('method', val)"
                               style="padding:4px 8px; cursor:pointer; font-size:12px; border-radius:4px"
                               class="filter-item">
                            {{ val }}
                          </div>
                        </div>
                      </div>
                    </template>
                  </a-popover>
                </div>
              </th>
              <th style="padding: 8px 12px; text-align: left; font-weight: 500; width: 80px; border-right: 1px solid var(--color-border-2);">
                <div style="display:flex; align-items:center; justify-content:space-between">
                  状态
                  <a-popover trigger="click" position="bottom">
                    <icon-filter :style="{ color: filters.status ? 'rgb(var(--primary-6))' : 'var(--color-text-3)', cursor: 'pointer' }" />
                    <template #content>
                      <div style="width:100px">
                        <a-input v-model="filters.status" placeholder="搜索状态..." size="small" allow-clear />
                        <div style="margin-top:8px; max-height:150px; overflow-y:auto; border-top:1px solid var(--color-border-1); padding-top:4px">
                          <div v-for="val in getUniqueValues('status')" :key="val" 
                               @click="setFilter('status', val)"
                               style="padding:4px 8px; cursor:pointer; font-size:12px; border-radius:4px"
                               class="filter-item">
                            {{ val }}
                          </div>
                        </div>
                      </div>
                    </template>
                  </a-popover>
                </div>
              </th>
              <th style="padding: 8px 12px; text-align: left; font-weight: 500; width: 100px; border-right: 1px solid var(--color-border-2);">
                耗时
              </th>
              <th style="padding: 8px 12px; text-align: left; font-weight: 500; width: 150px; border-right: 1px solid var(--color-border-2);">
                后端
              </th>
              <th style="padding: 8px 12px; text-align: left; font-weight: 500; border-right: 1px solid var(--color-border-2);">
                客户端
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(record, index) in paginatedRows" :key="index" 
                :style="{ backgroundColor: index % 2 === 0 ? 'var(--color-bg-1)' : 'var(--color-fill-1)' }"
                style="border-bottom: 1px solid var(--color-border-2);">
              <td style="padding: 8px 12px; border-right: 1px solid var(--color-border-2); font-family: monospace;">
                {{ record.parsed.source }}
              </td>
              <td style="padding: 8px 12px; border-right: 1px solid var(--color-border-2); font-family: monospace;">
                {{ record.parsed.time }}
              </td>
              <td style="padding: 8px 12px; border-right: 1px solid var(--color-border-2); word-break: break-all;">
                <div style="font-weight:500">{{ record.parsed.host }}</div>
                <div style="color:var(--color-text-3); font-size:12px">{{ record.parsed.path }}</div>
              </td>
              <td style="padding: 8px 12px; border-right: 1px solid var(--color-border-2);">
                <a-tag size="small" :color="getMethodColor(record.parsed.method)">{{ record.parsed.method }}</a-tag>
              </td>
              <td style="padding: 8px 12px; border-right: 1px solid var(--color-border-2);">
                <a-tag size="small" :color="getStatusColor(record.parsed.status)">{{ record.parsed.status }}</a-tag>
              </td>
              <td style="padding: 8px 12px; border-right: 1px solid var(--color-border-2); font-family: monospace;">
                {{ record.parsed.duration }}s
              </td>
              <td style="padding: 8px 12px; border-right: 1px solid var(--color-border-2); word-break: break-all; font-size: 12px;">
                {{ record.parsed.upstream }}
              </td>
              <td style="padding: 8px 12px; word-break: break-all;">
                {{ record.parsed.agent }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
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

    <div v-else-if="rows.length > 0 && viewMode==='raw'" style="margin-top:12px">
      <div style="margin-bottom:8px; color: var(--color-text-3);">Raw Data: {{ rows.length }} 条记录</div>
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
      暂无查询结果，请点击"运行查询"执行查询
    </div>

     <!-- 历史记录抽屉 -->
    <a-drawer v-model:visible="historyVisible" title="" width="600px" placement="bottom" :height="450" :footer="false">
      <!-- 搜索框 -->
      <div style="margin-bottom: 12px;">
        <a-input
          v-model="searchKeyword"
          placeholder="搜索查询历史..."
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
          <a-tab-pane key="recent" title="查询历史记录" />
          <a-tab-pane key="favorite" title="已收藏查询" />
        </a-tabs>
      </div>

      <div v-if="historyItems.length === 0" style="text-align: center; padding: 40px; color: var(--color-text-3);">
        {{ historyTab === 'favorite' ? '暂无收藏的查询' : '暂无查询历史' }}
      </div>
      
      <div v-else style="max-height: 280px; overflow-y: auto;">
        <div v-for="item in historyItems" :key="item.id" 
             style="border: 1px solid var(--color-border-2); border-radius: 6px; padding: 12px; margin-bottom: 8px; background: var(--color-fill-1);">
          
          <!-- 时间和操作按钮在同一行 -->
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; font-size: 12px; color: var(--color-text-3);">
            <div style="display: flex; gap: 12px; align-items: center;">
              <span>{{ new Date(item.createdAt).toLocaleString() }}</span>
              <a-tag :color="item.engine === 'loki' ? 'blue' : 'green'" size="small">{{ item.engine }}</a-tag>
              <a-tag color="gray" size="small">{{ item.mode }}</a-tag>
            </div>
            
            <!-- 操作按钮组 -->
            <div style="display: flex; gap: 4px;">
              <a-tooltip content="编辑备注">
                <a-button 
                  size="mini" 
                  type="text" 
                  @click="showNoteModal(item)"
                  style="width: 20px; height: 20px; padding: 0; display: flex; align-items: center; justify-content: center;"
                >
                  <icon-tag style="color: var(--color-primary-6) !important;" />
                </a-button>
              </a-tooltip>
              <a-tooltip content="删除记录">
                <a-button 
                  size="mini" 
                  type="text" 
                  @click="confirmDelete(item)"
                  style="width: 20px; height: 20px; padding: 0; display: flex; align-items: center; justify-content: center;"
                >
                  <icon-delete style="color: #f53f3f !important;" />
                </a-button>
              </a-tooltip>
              <a-tooltip :content="item.isFavorite ? '取消收藏' : '添加收藏'">
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
              <a-tooltip content="执行查询">
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
    <a-modal v-model:visible="noteModalVisible" title="编辑备注" @ok="saveNote" @cancel="cancelNote">
      <a-textarea 
        v-model="noteContent" 
        placeholder="为这个查询添加备注..."
        :rows="3"
        :max-length="200"
        show-word-limit
      />
    </a-modal>

    <a-modal v-model:visible="inspectVisible" title="查询检查器" :footer="false">
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
import { queryLogs, history as apiHistory, inspect, toggleFavorite, updateNote, deleteHistory } from '@/api/logs'
import { listDataSources } from '@/api/datasources'
import { Message, Modal } from '@arco-design/web-vue'
import { IconTag, IconDelete, IconStar, IconStarFill, IconSend, IconSearch, IconFilter } from '@arco-design/web-vue/es/icon'
import LogAnalysisChat from '@/components/LogAnalysisChat.vue'

const datasource = ref('')
const lokiDsOptions = ref([])
const esDsOptions = ref([])

// Computed property to dynamically show only configured data sources
const dsOptions = computed(() => {
  const options = []
  if (lokiDsOptions.value.length > 0) {
    options.push({ label: 'Loki', value: 'loki' })
  }
  if (esDsOptions.value.length > 0) {
    options.push({ label: 'Elasticsearch', value: 'elasticsearch' })
  }
  return options
})
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
  agent: ''
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
    // Numeric timestamps (number or numeric string)
    if (typeof timestamp === 'number' || (typeof timestamp === 'string' && /^\d+$/.test(timestamp))) {
      let ts = typeof timestamp === 'number' ? timestamp : parseInt(timestamp)
      // Convert nanoseconds to milliseconds if needed
      if (ts > 1e15) {
        ts = Math.floor(ts / 1e6)
      }
      const dateNum = new Date(ts)
      if (!isNaN(dateNum.getTime())) {
        return dateNum.toLocaleString('zh-CN', {
          year: 'numeric', month: '2-digit', day: '2-digit',
          hour: '2-digit', minute: '2-digit', second: '2-digit'
        })
      }
    }

    // ISO8601 string timestamps (e.g., 2025-09-16T14:29:15.609+08:00)
    if (typeof timestamp === 'string') {
      const dateIso = new Date(timestamp)
      if (!isNaN(dateIso.getTime())) {
        return dateIso.toLocaleString('zh-CN', {
          year: 'numeric', month: '2-digit', day: '2-digit',
          hour: '2-digit', minute: '2-digit', second: '2-digit'
        })
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
      Message.warning('请输入查询语句')
      return
    }
  } else if (payload.mode === 'builder') {
    const hasFilter = payload.builder.labelFilters.some(f => f.label && f.values && f.values.length > 0)
    const hasContains = payload.builder.contains && payload.builder.contains.trim()
    if (!hasFilter && !hasContains) {
      Message.warning('请选择查询条件')
      return
    }
  }
  await runQuery({ engine: 'loki', payload })
}

async function onRunES(payload) {
  // Validation
  if (!payload.query || !payload.query.trim()) {
    Message.warning('请输入查询语句')
    return
  }
  viewMode.value = payload?.mode === 'raw' ? 'raw' : 'logs'
  await runQuery({ engine: 'elasticsearch', payload })
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
      Message.success('备注保存成功')
      noteModalVisible.value = false
    }
  } catch (error) {
    console.error('Failed to save note:', error)
    Message.error('保存备注失败')
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
    title: '确认删除',
    content: '确定要删除这条查询记录吗？此操作不可恢复。',
    onOk: () => deleteHistoryItem(item)
  })
}

// 删除历史记录
async function deleteHistoryItem(item) {
  try {
    const { data } = await deleteHistory(item.id)
    if (data?.code === 0) {
      Message.success('删除成功')
      await loadHistoryData() // 重新加载数据
    }
  } catch (error) {
    console.error('Failed to delete history:', error)
    Message.error('删除失败')
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
    }
    
    // 关闭抽屉
    historyVisible.value = false
    
    // 构造查询参数并执行
    const { start, end } = computeTimeRange()
    let dsId = ''
    if (item.engine === 'loki') {
      dsId = selectedLokiId.value || localStorage.getItem('last_loki_ds_id') || ''
    } else {
      dsId = selectedEsId.value || localStorage.getItem('last_es_ds_id') || ''
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
    
    Message.success('查询执行成功')
  } catch (error) {
    console.error('Execute query error:', error)
    Message.error('查询执行失败')
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
  } else {
    const dsId = selectedEsId.value || localStorage.getItem('last_es_ds_id') || ''
    params = { engine: 'elasticsearch', datasourceId: dsId, start, end }
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
    if (!selectedEsId.value && esDsOptions.value.length) {
      selectedEsId.value = esDsOptions.value[0].value
      localStorage.setItem('last_es_ds_id', selectedEsId.value)
      console.log('Selected ES datasource:', selectedEsId.value)
    }
    
    // Auto-select first available data source type
    if (!datasource.value) {
      if (lokiDsOptions.value.length > 0) {
        datasource.value = 'loki'
      } else if (esDsOptions.value.length > 0) {
        datasource.value = 'elasticsearch'
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
  
  return {
    source,
    time,
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

