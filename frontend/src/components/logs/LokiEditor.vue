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
              <a-select v-model="f.label" style="width:140px" placeholder="Select label" allow-search />
              <a-select v-model="f.op" :options="ops" style="width:80px" />
              <a-input-tag v-model="f.values" style="min-width:220px" placeholder="Select value" />
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
            <a-button>添加查询</a-button>
            <a-button>查询历史记录</a-button>
            <a-button>查询检查器</a-button>
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
          </a-space>
        </div>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>
<script setup>
import { reactive, ref } from 'vue'

const emit = defineEmits(['run'])

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
</script>
