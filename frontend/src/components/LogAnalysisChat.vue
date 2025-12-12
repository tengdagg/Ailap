<template>
  <div class="log-analysis-chat">
    <a-tooltip :content="$t('chat.tooltip')" position="left">
      <a-button
        class="chat-trigger"
        type="primary"
        shape="circle"
        size="large"
        @click="visible = true"
      >
        <img v-if="defaultModelLogo" :src="defaultModelLogo" style="width: 24px; height: 24px; object-fit: contain;" />
        <icon-robot v-else />
      </a-button>
    </a-tooltip>

    <a-drawer
      v-model:visible="visible"
      :title="$t('chat.title')"
      placement="right"
      width="450px"
      :drawer-style="{ top: '10%', height: '90%', borderRadius: '16px 0 0 16px' }"
      :footer="false"
      unmount-on-close
    >
      <div class="chat-container">
        <div class="messages" ref="messagesRef">
          <div v-if="messages.length === 0" class="empty-state">
            <img v-if="defaultModelLogo" :src="defaultModelLogo" style="width: 48px; height: 48px; object-fit: contain; margin-bottom: 16px;" />
            <icon-robot v-else :style="{ fontSize: '48px', color: 'var(--color-text-3)' }" />
            <p>{{ $t('chat.intro') }}</p>
            <p>{{ $t('chat.help', { count: logs.length }) }}</p>
            <div class="quick-actions">
              <a-tag clickable @click="sendPrompt($t('chat.promptAnomaly'))">{{ $t('chat.actionAnomaly') }}</a-tag>
              <a-tag clickable @click="sendPrompt($t('chat.promptSummary'))">{{ $t('chat.actionSummary') }}</a-tag>
              <a-tag clickable @click="sendPrompt($t('chat.promptError'))">{{ $t('chat.actionError') }}</a-tag>
            </div>
          </div>

          <div
            v-for="(msg, index) in messages"
            :key="index"
            class="message-item"
            :class="{ 'is-me': msg.role === 'user' }"
          >
            <div class="avatar">
              <icon-user v-if="msg.role === 'user'" />
              <img v-else-if="defaultModelLogo" :src="defaultModelLogo" style="width: 20px; height: 20px; object-fit: contain;" />
              <icon-robot v-else />
            </div>
            <div class="content">
              <div class="bubble">
                <div v-if="msg.role === 'assistant'" v-html="formatContent(msg.content)"></div>
                <div v-else>{{ msg.content }}</div>
              </div>
              <div class="time" v-if="msg.time">{{ msg.time }}</div>
            </div>
          </div>

          <div v-if="loading" class="message-item">
            <div class="avatar">
              <img v-if="defaultModelLogo" :src="defaultModelLogo" style="width: 20px; height: 20px; object-fit: contain;" />
              <icon-robot v-else />
            </div>
            <div class="content">
              <div class="bubble loading">
                <icon-loading /> {{ $t('chat.loading') }}
              </div>
            </div>
          </div>
        </div>

        <div class="input-area">
          <a-textarea
            v-model="inputContent"
            :placeholder="$t('chat.placeholder')"
            :auto-size="{ minRows: 3, maxRows: 6 }"
            @keydown.enter.prevent="handleEnter"
          />
          <a-button type="primary" @click="handleSend" :loading="loading">
            <template #icon><icon-send /></template>
          </a-button>
        </div>
      </div>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, watch, nextTick, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { IconRobot, IconUser, IconSend, IconLoading } from '@arco-design/web-vue/es/icon'
import { Message } from '@arco-design/web-vue'
import { analyzeLogs } from '@/api/ai'
import { listModels } from '@/api/models'


const props = defineProps({
  logs: {
    type: Array,
    default: () => []
  },
  initialRange: {
    type: Object,
    default: () => ({ start: 0, end: 0 })
  }
})

const { t } = useI18n()
const visible = ref(false)
const inputContent = ref('')
const loading = ref(false)
const messages = ref([])
const messagesRef = ref(null)
const defaultModelLogo = ref('')
let currentAbortController = null

function getLogo(provider) {
  try {
    return new URL(`../assets/${provider}.png`, import.meta.url).href
  } catch (_) {
    return new URL(`../assets/logo.png`, import.meta.url).href
  }
}

async function fetchDefaultModel() {
  try {
    const { data } = await listModels()
    const models = data?.data?.items || []
    const defaultModel = models.find(m => m.isDefault)
    if (defaultModel) {
      defaultModelLogo.value = getLogo(defaultModel.provider)
    }
  } catch (e) {
    console.error('Failed to fetch default model:', e)
  }
}

onMounted(() => {
  fetchDefaultModel()
})

watch(visible, (val) => {
  if (!val && loading.value && currentAbortController) {
    currentAbortController.abort()
    loading.value = false
    messages.value.push({
      role: 'assistant',
      content: t('chat.cancel'),
      time: new Date().toLocaleTimeString()
    })
  }
})

// Simple markdown formatter if marked is not installed, but usually it is or we can just display text
// For now, let's assume we just display text with line breaks
function formatContent(text) {
  // If marked is available globally or we can import it. 
  // Since I don't know if marked is installed, I'll do simple formatting:
  // Convert newlines to <br> and bold to <b>
  if (!text) return ''
  return text
    .replace(/\n/g, '<br/>')
    .replace(/\*\*(.*?)\*\*/g, '<b>$1</b>')
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight
    }
  })
}

function handleEnter(e) {
  if (!e.shiftKey) {
    handleSend()
  }
}

async function sendPrompt(text) {
  if (loading.value) return
  
  const content = text || inputContent.value.trim()
  if (!content) return

  // Add user message
  messages.value.push({
    role: 'user',
    content: content,
    time: new Date().toLocaleTimeString()
  })
  
  inputContent.value = ''
  scrollToBottom()
  loading.value = true
  currentAbortController = new AbortController()

  try {
    // Prepare logs - limit to last 50 or so to avoid huge payload if not handled by backend
    // Backend handler seems to handle truncation, but let's be safe and send a reasonable amount
    // The backend handler says: "limit := 8000 // characters limit"
    // So we should probably limit the number of logs we send.
    const logsToSend = props.logs.slice(0, 100) // Send up to 100 logs

    const { data } = await analyzeLogs({
      prompt: content,
      logs: logsToSend
    }, { signal: currentAbortController.signal })

    if (data.code === 0) {
      messages.value.push({
        role: 'assistant',
        content: data.data.reply,
        time: new Date().toLocaleTimeString()
      })
    } else {
      messages.value.push({
        role: 'assistant',
        content: `Sorry, something went wrong: ${data.message}`,
        time: new Date().toLocaleTimeString()
      })
    }
  } catch (error) {
    if (error.name === 'CanceledError' || error.code === 'ERR_CANCELED') {
      // Request canceled, do nothing or log
      console.log('Request canceled')
      return
    }
    console.error(error)
    const errorMsg = error.response?.data?.message || error.message || 'Network error or server unavailable.'
    messages.value.push({
      role: 'assistant',
      content: `Error: ${errorMsg}`,
      time: new Date().toLocaleTimeString()
    })
  } finally {
    loading.value = false
    currentAbortController = null
    scrollToBottom()
  }
}

function handleSend() {
  sendPrompt()
}
</script>

<style scoped>
.log-analysis-chat {
  position: fixed;
  bottom: 40px;
  right: 40px;
  z-index: 1000;
}

.chat-trigger {
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
}

.chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.empty-state {
  text-align: center;
  margin-top: 40px;
  color: var(--color-text-3);
}

.quick-actions {
  display: flex;
  gap: 8px;
  justify-content: center;
  flex-wrap: wrap;
  margin-top: 16px;
}

.message-item {
  display: flex;
  gap: 12px;
  max-width: 90%;
}

.message-item.is-me {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--color-fill-3);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.bubble {
  background: var(--color-fill-2);
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.5;
  white-space: pre-wrap;
}

.is-me .bubble {
  background: rgb(var(--primary-6));
  color: #fff;
}

.time {
  font-size: 12px;
  color: var(--color-text-4);
  margin-top: 4px;
  text-align: right;
}

.is-me .time {
  text-align: left;
}

.input-area {
  padding: 16px;
  border-top: 1px solid var(--color-border-2);
  display: flex;
  gap: 12px;
  align-items: flex-start;
}
</style>
