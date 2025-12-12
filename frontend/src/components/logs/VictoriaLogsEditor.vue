<template>
  <div style="background:var(--color-bg-2); padding:16px; border-radius:4px; margin-bottom:16px">
    <a-alert v-if="checks.length > 0" type="warning" closable style="margin-bottom:12px">
      <div v-for="(c, i) in checks" :key="i">{{ c }}</div>
    </a-alert>

    <div style="margin-bottom:12px; display:flex; justify-content:space-between; align-items:center">
      <a-typography-text type="secondary">
        LogsQL 查询 (例如: <code>_time:5m</code> or <code>error</code>)
        <a href="https://docs.victoriametrics.com/victorialogs/logsql/" target="_blank" style="margin-left:8px">文档</a>
      </a-typography-text>
      <a-space>
         <a-button size="small" @click="$emit('history')">历史记录 <icon-history /></a-button>
         <a-button size="small" @click="$emit('inspect', query)">检查 Query <icon-code /></a-button>
      </a-space>
    </div>

    <!-- Code Editor -->
    <div style="border:1px solid var(--color-border-3); border-radius:4px;">
      <a-textarea 
        v-model="query" 
        :auto-size="{minRows:3, maxRows:10}" 
        style="border:none; background:var(--color-bg-1); font-family:monospace" 
        placeholder="Enter LogsQL query..." 
        @keydown.enter.prevent="onRun"
      />
    </div>

    <div style="margin-top:12px; display:flex; justify-content:flex-end">
      <a-button type="primary" @click="onRun">
        <template #icon><icon-play-arrow /></template>
        运行查询
      </a-button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { IconPlayArrow, IconHistory, IconCode } from '@arco-design/web-vue/es/icon'

const props = defineProps({
  datasourceId: String,
})

const emit = defineEmits(['run', 'history', 'inspect'])

const query = ref('*')
const checks = ref([])

function onRun() {
  emit('run', {
    mode: 'code',
    query: query.value
  })
}

// Ensure query is not empty initially
watch(() => query.value, (val) => {
  if (!val) query.value = ''
})
</script>
