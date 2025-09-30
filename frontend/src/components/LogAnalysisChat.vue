<template>
  <div>
    <a-button type="primary" shape="circle" class="fab" @click="visible = true">
      <icon-robot />
    </a-button>

    <a-drawer v-model:visible="visible" :footer="false" :width="560" placement="right" title="日志智能分析">
      <div class="chat-container">
        <div class="options-row">
          <div class="opt">
            <span class="label">时间范围</span>
            <a-range-picker v-model="range" show-time style="width: 300px;" />
          </div>
          <div class="opt">
            <span class="label">采样条数</span>
            <a-input-number v-model="limit" :min="50" :max="1000" :step="50" />
          </div>
        </div>
        <div class="messages" ref="listRef">
          <div v-for="(m, i) in messages" :key="i" :class="['msg', m.role]">
            <div class="bubble" v-html="format(m.content)"></div>
          </div>
          <div v-if="loading" class="msg assistant">
            <div class="bubble">正在分析...</div>
          </div>
        </div>
        <a-divider style="margin: 8px 0;" />
        <div class="input-row">
          <a-input v-model="prompt" :disabled="loading" placeholder="描述你想分析的问题，例如：错误原因和排查步骤" allow-clear @press-enter="onSend" />
          <a-button type="primary" :loading="loading" @click="onSend">
            <template #icon><icon-send /></template>
            发送
          </a-button>
        </div>
        <div class="tips">将基于当前查询结果进行分析，仅发送部分日志片段参与推理。</div>
      </div>
    </a-drawer>
  </div>
  
</template>

<script setup>
import { ref, watch, nextTick, onMounted } from 'vue'
import { analyzeLogs } from '@/api/ai'
import { IconRobot, IconSend } from '@arco-design/web-vue/es/icon'

const props = defineProps({ 
  logs: { type: Array, default: () => [] },
  initialRange: { type: Object, default: () => ({ start: 0, end: 0 }) }, // ms
})

const visible = ref(false)
const messages = ref([])
const prompt = ref('')
const loading = ref(false)
const listRef = ref(null)
const range = ref([]) // [Date, Date]
const limit = ref(200)

watch(messages, async () => {
  await nextTick()
  if (listRef.value) {
    listRef.value.scrollTop = listRef.value.scrollHeight
  }
})

function format(text) {
  return (text || '').replace(/\n/g, '<br/>')
}

function toMs(ts) {
  if (!ts) return 0
  if (typeof ts === 'number') {
    if (ts > 1e15) return Math.floor(ts / 1e6) // ns -> ms
    if (ts > 1e12) return Math.floor(ts / 1)   // ms
    return ts
  }
  if (typeof ts === 'string' && /^\d+$/.test(ts)) {
    const n = parseInt(ts)
    return n > 1e15 ? Math.floor(n / 1e6) : n
  }
  const d = new Date(ts)
  return isNaN(d.getTime()) ? 0 : d.getTime()
}

function filterLogs() {
  let items = props.logs || []
  if (range.value && range.value.length === 2 && range.value[0] && range.value[1]) {
    const startMs = range.value[0].getTime()
    const endMs = range.value[1].getTime()
    items = items.filter(r => {
      const ms = toMs(r.timestamp || r.time || r.ts)
      return ms >= startMs && ms <= endMs
    })
  }
  return items.slice(0, Math.max(1, Number(limit.value) || 200))
}

async function onSend() {
  if (!prompt.value || loading.value) return
  const user = { role: 'user', content: prompt.value }
  messages.value.push(user)
  loading.value = true
  try {
    const compact = filterLogs()
    const timeHint = (range.value && range.value.length === 2 && range.value[0] && range.value[1])
      ? `（时间范围：${range.value[0].toLocaleString()} - ${range.value[1].toLocaleString()}）`
      : ''
    const { data } = await analyzeLogs({ prompt: `${prompt.value} ${timeHint}`, logs: compact })
    if (data?.code === 0) {
      messages.value.push({ role: 'assistant', content: data?.data?.reply || '（无内容）' })
    } else {
      messages.value.push({ role: 'assistant', content: data?.message || '分析失败' })
    }
  } catch (e) {
    messages.value.push({ role: 'assistant', content: e?.response?.data?.message || e?.message || '分析失败' })
  } finally {
    loading.value = false
    prompt.value = ''
  }
}

onMounted(() => {
  if (props.initialRange && props.initialRange.start && props.initialRange.end) {
    range.value = [new Date(props.initialRange.start), new Date(props.initialRange.end)]
  }
})
</script>

<style scoped>
.fab {
  position: fixed;
  right: 20px;
  bottom: 20px;
  z-index: 1000;
}
.options-row { display:flex; align-items:center; justify-content: space-between; gap: 12px; margin-bottom: 8px; }
.opt { display:flex; align-items:center; gap:8px; }
.label { color: var(--color-text-2); font-size: 12px; }
.chat-container { display: flex; flex-direction: column; min-height: 400px; }
.messages {
  flex: 1; overflow: auto; padding: 8px; border: 1px solid var(--color-border-2); border-radius: 6px; background: var(--color-bg-2);
}
.msg { display: flex; margin-bottom: 8px; }
.msg.user { justify-content: flex-end; }
.msg .bubble { max-width: 80%; padding: 8px 10px; border-radius: 8px; line-height: 1.6; }
.msg.user .bubble { background: var(--color-primary-1); color: var(--color-text-1); }
.msg.assistant .bubble { background: var(--color-fill-2); color: var(--color-text-1); }
.input-row { display:flex; gap:8px; margin-top:8px; }
.tips { color: var(--color-text-3); margin-top: 6px; font-size: 12px; }
</style>


