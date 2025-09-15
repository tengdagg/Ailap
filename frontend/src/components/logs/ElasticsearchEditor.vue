<template>
  <div>
    <a-tabs v-model:active-key="tab" size="large">
      <a-tab-pane key="logs" title="Logs">
        <a-input v-model="lucene" placeholder="输入 Lucene 查询语句 (Shift+Enter 执行)" @keydown.shift.enter.prevent="run" />
        <a-collapse :default-active-key="['opt']" style="margin-top:8px">
          <a-collapse-item header="Options" key="opt">
            <a-space>
              <span>Limit</span>
              <a-input-number v-model="limit" :min="1" />
            </a-space>
          </a-collapse-item>
        </a-collapse>
        <div style="margin-top:8px">
          <a-space>
            <a-button type="primary" @click="run">运行查询</a-button>
            <a-button>添加查询</a-button>
            <a-button @click="$emit('history')">查询历史记录</a-button>
            <a-button @click="emitInspect">查询检查器</a-button>
          </a-space>
        </div>
      </a-tab-pane>
      <a-tab-pane key="metrics" title="Metrics" disabled />
      <a-tab-pane key="raw-data" title="Raw Data" disabled />
      <a-tab-pane key="raw-doc" title="Raw Document" disabled />
    </a-tabs>
  </div>
</template>
<script setup>
import { ref } from 'vue'

const emit = defineEmits(['run', 'history', 'inspect'])
const tab = ref('logs')
const lucene = ref('')
const limit = ref(500)

function run() { 
  emit('run', { 
    mode: 'code', 
    query: lucene.value, 
    lineLimit: limit.value 
  }) 
}

function emitInspect() {
  emit('inspect', lucene.value)
}
</script>


