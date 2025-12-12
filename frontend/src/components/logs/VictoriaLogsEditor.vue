<template>
  <div style="background:var(--color-bg-2); padding:16px; border-radius:4px; margin-bottom:16px">
    <a-alert v-if="checks.length > 0" type="warning" closable style="margin-bottom:12px">
      <div v-for="(c, i) in checks" :key="i">{{ c }}</div>
    </a-alert>

    <div style="margin-bottom:12px; display:flex; justify-content:space-between; align-items:center">
      <a-typography-text type="secondary">
        <span v-html="$t('logs.logsQLHelp')" />
        <a href="https://docs.victoriametrics.com/victorialogs/logsql/" target="_blank" style="margin-left:8px">{{ $t('logs.docs') }}</a>
      </a-typography-text>
      <a-space>
         <a-button size="small" @click="$emit('history')">{{ $t('logs.history') }} <icon-history /></a-button>
         <a-button size="small" @click="$emit('inspect', query)">{{ $t('logs.inspectQuery') }} <icon-code /></a-button>
      </a-space>
    </div>

    <!-- Code Editor -->
    <div style="border:1px solid var(--color-border-3); border-radius:4px;">
      <a-textarea 
        v-model="query" 
        :auto-size="{minRows:3, maxRows:10}" 
        style="border:none; background:var(--color-bg-1); font-family:monospace" 
        :placeholder="$t('logs.enterLogsQLQuery')" 
        @keydown.enter.prevent="onRun"
      />
    </div>

    <div style="margin-top:12px; display:flex; justify-content:flex-end">
      <a-button type="primary" @click="onRun">
        <template #icon><icon-play-arrow /></template>
        {{ $t('logs.runQuery') }}
      </a-button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { IconPlayArrow, IconHistory, IconCode } from '@arco-design/web-vue/es/icon'

const { t } = useI18n()

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
