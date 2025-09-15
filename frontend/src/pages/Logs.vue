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

    <div v-if="rows.length > 0" style="margin-top:12px">
      <div style="margin-bottom:8px; color: #666;">查询结果: {{ rows.length }} 条记录</div>
      
      
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

      <div v-if="historyItems.length === 0" style="text-align: center; padding: 40px; color: #999;">
        {{ historyTab === 'favorite' ? '暂无收藏的查询' : '暂无查询历史' }}
      </div>
      
      <div v-else style="max-height: 280px; overflow-y: auto;">
        <div v-for="item in historyItems" :key="item.id" 
             style="border: 1px solid #e5e6eb; border-radius: 6px; padding: 12px; margin-bottom: 8px; background: #fafafa;">
          
          <!-- 时间和操作按钮在同一行 -->
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; font-size: 12px; color: #666;">
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
                  <icon-tag />
                </a-button>
              </a-tooltip>
              <a-tooltip content="删除记录">
                <a-button 
                  size="mini" 
                  type="text" 
                  @click="confirmDelete(item)"
                  style="width: 20px; height: 20px; padding: 0; display: flex; align-items: center; justify-content: center; color: #f53f3f;"
                >
                  <icon-delete />
                </a-button>
              </a-tooltip>
              <a-tooltip :content="item.isFavorite ? '取消收藏' : '添加收藏'">
                <a-button 
                  size="mini" 
                  type="text" 
                  :style="{ 
                    width: '20px', 
                    height: '20px', 
                    padding: '0', 
                    display: 'flex', 
                    alignItems: 'center', 
                    justifyContent: 'center',
                    color: item.isFavorite ? '#faad14' : '#8c8c8c'
                  }"
                  @click="toggleQueryFavorite(item)"
                >
                  <icon-star-fill v-if="item.isFavorite" />
                  <icon-star v-else />
                </a-button>
              </a-tooltip>
              <a-tooltip content="执行查询">
                <a-button 
                  size="mini" 
                  type="text" 
                  @click="executeQuery(item)"
                  style="width: 20px; height: 20px; padding: 0; display: flex; align-items: center; justify-content: center; color: #1890ff;"
                >
                  <icon-send />
                </a-button>
              </a-tooltip>
            </div>
          </div>
          
          <!-- 备注显示 -->
          <div v-if="item.note" style="margin-bottom: 8px; font-size: 12px; color: #666; font-style: italic;">
            💬 {{ item.note }}
          </div>
          
          <div style="font-family: monospace; font-size: 13px; word-break: break-all; background: #f8f9fa; padding: 8px; border-radius: 4px; border: 1px solid #e9ecef;">
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
    </a-modal>
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
import { IconTag, IconDelete, IconStar, IconStarFill, IconSend, IconSearch } from '@arco-design/web-vue/es/icon'

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
const historyTab = ref('recent')
const historyItems = ref([])
const allHistoryItems = ref([]) // 存储所有历史记录
const searchKeyword = ref('')
const inspectVisible = ref(false)
const inspectUrl = ref('')

// 备注编辑相关
const noteModalVisible = ref(false)
const noteContent = ref('')
const currentEditItem = ref(null)

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
      query: item.query
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

