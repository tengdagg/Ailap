<template>
  <div>
    <a-tabs v-model:active-key="mode" type="line" size="large">
      <a-tab-pane key="builder" title="Builder">
        <a-space direction="vertical" fill>
          <div>
            <div style="margin-bottom:6px">{{ $t('logs.labelFilters') }}</div>
            <div v-for="(f,i) in form.labelFilters" :key="i" style="display:flex; gap:8px; margin-bottom:8px; align-items:center">
              <a-select v-model="f.label" style="width:180px" :placeholder="$t('logs.selectLabel')" allow-search :options="labelOptions" :loading="loadingLabels" @focus="ensureLabels" @change="onLabelChange(i)" />
              <a-select v-model="f.op" :options="ops" style="width:60px" />
              <a-select v-model="valuesDraft[i]" multiple allow-search style="width:160px" :placeholder="$t('logs.selectValue')" :options="valueOptions[i]" :loading="loadingValues[i]" @focus="ensureLabelValues(i)" @change="commitValues(i)" />
              
              <!-- + 号用于添加新行，只在最后一行显示 -->
              <a-button v-if="i === form.labelFilters.length - 1" @click="addLabelFilter" size="mini" type="outline" style="color: #1890ff; border-color: #1890ff; width: 28px; height: 28px; padding: 0;">+</a-button>
              
              <!-- X 号用于删除行，当有多行时显示 -->
              <a-button v-if="form.labelFilters.length > 1" @click="removeLabelFilter(i)" size="mini" type="outline" status="danger" style="color: #f53f3f; border-color: #f53f3f; width: 28px; height: 28px; padding: 0;">×</a-button>
            </div>
          </div>

          <div>
            <div style="margin-bottom:6px">{{ $t('logs.lineContains') }}</div>
            <div style="display:flex; gap:8px; align-items:center">
              <a-input v-model="form.contains" :placeholder="$t('logs.textToFind')" style="max-width:360px" />
            </div>
          </div>

          <a-collapse :default-active-key="['opt']">
            <a-collapse-item :header="$t('logs.options')" key="opt">
              <a-space wrap>
                <a-descriptions :column="3" size="small" :bordered="false">
                  <a-descriptions-item :label="$t('logs.type')">
                    <a-segmented v-model="form.type" :options="['Range','Instant']" />
                  </a-descriptions-item>
                  <a-descriptions-item :label="$t('logs.lineLimit')">
                    <a-input-number v-model="form.lineLimit" :min="1" />
                  </a-descriptions-item>
                </a-descriptions>
              </a-space>
            </a-collapse-item>
          </a-collapse>

          <a-space>
            <a-button type="primary" @click="run">{{ $t('logs.runQuery') }}</a-button>
            <a-button @click="$emit('history')">{{ $t('logs.queryHistory') }}</a-button>
            <a-button @click="emitInspectBuilder">{{ $t('logs.queryInspector') }}</a-button>
          </a-space>
        </a-space>
      </a-tab-pane>

      <a-tab-pane key="code" title="Code">
        <a-input
          v-model="code"
          :placeholder="$t('logs.enterLokiQuery')"
          @keydown.shift.enter.prevent="runCode"
        />
        <div style="margin-top:8px">
          <a-space>
            <a-segmented v-model="form.type" :options="['Range','Instant']" />
            <a-input-number v-model="form.lineLimit" :min="1" />
            <a-button type="primary" @click="runCode">{{ $t('logs.runQuery') }}</a-button>
            <a-button @click="emitInspectCode">{{ $t('logs.queryInspector') }}</a-button>
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
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const emit = defineEmits(['run','history','inspect'])
const props = defineProps({ datasourceId: { type: [String, Number], default: '' } })

const mode = ref('builder')
const code = ref('')

const ops = [
  { label: '=', value: '=' },
  { label: '!=', value: '!=' },
  { label: '=~', value: '=~' },
  { label: '!~', value: '!~' },
]

const form = reactive({
  labelFilters: [{ label: '', op: '=', values: [] }], // 默认有一行
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

function onLabelChange(idx) {
  valuesDraft.value[idx] = []
  form.labelFilters[idx].values = []
  valueOptions.value[idx] = []
}

function addLabelFilter() { form.labelFilters.push({ label: '', op: '=', values: [] }) }
function removeLabelFilter(i) { form.labelFilters.splice(i, 1) }

function run() {
  emit('run', { mode: 'builder', builder: { ...form }, lineLimit: form.lineLimit, explain: false })
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


