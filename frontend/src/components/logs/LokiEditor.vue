<template>
  <div>
    <a-tabs v-model:active-key="mode" type="line" size="large">
      <a-tab-pane key="builder" title="Builder">
        <a-space direction="vertical" fill>
          <a-space>
            <a-button @click="addLabelFilter">新增标签过滤</a-button>
            <a-button @click="addOperation">+ Operations</a-button>
            <a-switch v-model="explain" size="small">解释查询</a-switch>
          </a-space>

          <div>
            <div style="margin-bottom:6px">Label filters</div>
            <div v-for="(f,i) in form.labelFilters" :key="i" style="display:flex; gap:8px; margin-bottom:8px">
              <a-select v-model="f.label" style="width:140px" placeholder="Select label" allow-search :options="labelOptions" :loading="loadingLabels" @focus="ensureLabels" />
              <a-select v-model="f.op" :options="ops" style="width:80px" />
              <a-select v-model="valuesDraft[i]" multiple allow-search style="min-width:220px" placeholder="Select value" :options="valueOptions[i]" :loading="loadingValues[i]" @focus="ensureLabelValues(i)" @change="commitValues(i)" />
              <a-button @click="removeLabelFilter(i)" size="mini">-</a-button>
            </div>
          </div>

          <div>
            <div style="margin-bottom:6px">Line contains</div>
            <div style="display:flex; gap:8px; align-items:center">
              <a-input v-model="form.contains" placeholder="Text to find" style="max-width:360px" />
              <a-button size="mini" @click="addContains">+</a-button>
            </div>
          </div>

          <a-collapse :default-active-key="['opt']">
            <a-collapse-item header="Options" key="opt">
              <a-space wrap>
                <a-descriptions :column="3" size="small" :bordered="false">
                  <a-descriptions-item label="Type">
                    <a-segmented v-model="form.type" :options="['Range','Instant']" />
                  </a-descriptions-item>
                  <a-descriptions-item label="Line limit">
                    <a-input-number v-model="form.lineLimit" :min="1" />
                  </a-descriptions-item>
                </a-descriptions>
              </a-space>
            </a-collapse-item>
          </a-collapse>

          <a-space>
            <a-button type="primary" @click="run">运行查询</a-button>
            <a-button @click="$emit('history')">查询历史记录</a-button>
            <a-button @click="emitInspectBuilder">查询检查器</a-button>
          </a-space>
        </a-space>
      </a-tab-pane>

      <a-tab-pane key="code" title="Code">
        <a-input
          v-model="code"
          placeholder="Enter a Loki query (Shift+Enter 执行)"
          @keydown.shift.enter.prevent="runCode"
        />
        <div style="margin-top:8px">
          <a-space>
            <a-segmented v-model="form.type" :options="['Range','Instant']" />
            <a-input-number v-model="form.lineLimit" :min="1" />
            <a-button type="primary" @click="runCode">运行查询</a-button>
            <a-button @click="emitInspectCode">查询检查器</a-button>
          </a-space>
        </div>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>
<script setup>
import { reactive, ref } from 'vue'
import { suggestions, labelValues } from '@/api/logs'
import { Message } from '@arco-design/web-vue'

const emit = defineEmits(['run','history','inspect'])
const props = defineProps({ datasourceId: { type: [String, Number], default: '' } })

const mode = ref('builder')
const explain = ref(false)
const code = ref('')

const ops = [
  { label: '=', value: '=' },
  { label: '!=', value: '!=' },
  { label: '=~', value: '=~' },
  { label: '!~', value: '!~' },
]

const form = reactive({
  labelFilters: [],
  contains: '',
  type: 'Range',
  lineLimit: 1000,
})

const labelOptions = ref([])
const loadingLabels = ref(false)
const valueOptions = ref([])
const loadingValues = ref([])
const valuesDraft = ref([])

async function ensureLabels() {
  if (labelOptions.value.length || loadingLabels.value) return
  loadingLabels.value = true
  try {
    const dsId = String(props.datasourceId || localStorage.getItem('last_loki_ds_id') || '')
    const { data } = await suggestions({ engine: 'loki', datasourceId: dsId })
    const items = data?.data?.items || []
    labelOptions.value = items.map((x) => ({ label: String(x), value: String(x) }))
  } catch (e) {
    Message.warning('无法加载标签，可稍后重试')
  } finally {
    loadingLabels.value = false
  }
}

function ensureArraySize(arr, n, init) {
  while (arr.length <= n) arr.push(typeof init === 'function' ? init() : init)
}

async function ensureLabelValues(idx) {
  const label = form.labelFilters[idx]?.label
  if (!label) return
  ensureArraySize(valueOptions.value, idx, [])
  ensureArraySize(loadingValues.value, idx, false)
  if (valueOptions.value[idx]?.length || loadingValues.value[idx]) return
  loadingValues.value[idx] = true
  try {
    const dsId = String(props.datasourceId || localStorage.getItem('last_loki_ds_id') || '')
    const { data } = await labelValues({ engine: 'loki', label, datasourceId: dsId })
    const items = data?.data?.items || []
    valueOptions.value[idx] = items.map((x) => ({ label: String(x), value: String(x) }))
  } catch (_) {
  } finally {
    loadingValues.value[idx] = false
  }
}

function commitValues(idx) {
  ensureArraySize(valuesDraft.value, idx, [])
  const vals = valuesDraft.value[idx] || []
  form.labelFilters[idx].values = Array.isArray(vals) ? vals.slice(0, 1) : []
}

function addLabelFilter() { form.labelFilters.push({ label: '', op: '=', values: [] }) }
function removeLabelFilter(i) { form.labelFilters.splice(i, 1) }
function addOperation() {}
function addContains() {}

function run() {
  emit('run', { mode: 'builder', builder: { ...form }, explain: explain.value })
}
function runCode() {
  emit('run', { mode: 'code', query: code.value, type: form.type, lineLimit: form.lineLimit })
}

function emitInspectBuilder() {
  const q = buildQueryFromBuilder()
  emit('inspect', q)
}
function emitInspectCode() {
  const q = code.value || ''
  emit('inspect', q)
}

function buildQueryFromBuilder() {
  // mirror backend buildLokiQuery for UI inspector
  const parts = []
  const labels = []
  for (const f of form.labelFilters) {
    if (!f.label || !Array.isArray(f.values) || f.values.length === 0) continue
    const op = f.op || '='
    const v = String(f.values[0]).replaceAll('"', '\\"')
    labels.push(`${f.label}${op}"${v}"`)
  }
  parts.push(`{${labels.join(',')}}`)
  if (form.contains) {
    const text = String(form.contains).replaceAll('"', '\\"')
    parts.push(`|~ "${text}"`)
  }
  return parts.join(' ')
}
</script>


