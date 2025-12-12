<template>
  <div>
    <a-tabs v-model:active-key="tab" size="large">
      <a-tab-pane key="logs" title="Logs">
        <a-input v-model="lucene" :placeholder="$t('logs.enterLuceneQuery')" @keydown.shift.enter.prevent="run" />
        <a-collapse :default-active-key="['opt']" style="margin-top:8px">
          <a-collapse-item :header="$t('logs.options')" key="opt">
            <a-space>
              <span>Limit</span>
              <a-input-number v-model="limit" :min="1" />
            </a-space>
          </a-collapse-item>
        </a-collapse>
        <div style="margin-top:8px">
          <a-space>
            <a-button type="primary" @click="run">{{ $t('logs.runQuery') }}</a-button>
            <a-button>{{ $t('logs.addQuery') }}</a-button>
            <a-button @click="$emit('history')">{{ $t('logs.queryHistory') }}</a-button>
            <a-button @click="emitInspect">{{ $t('logs.queryInspector') }}</a-button>
          </a-space>
        </div>
      </a-tab-pane>
      <a-tab-pane key="raw-data" title="Raw Data">
        <a-input v-model="lucene" :placeholder="$t('logs.enterLuceneQuery')" @keydown.shift.enter.prevent="run" />
        <a-collapse :default-active-key="['opt']" style="margin-top:8px">
          <a-collapse-item :header="$t('logs.options')" key="opt">
            <a-space>
              <span>Limit</span>
              <a-input-number v-model="limit" :min="1" />
            </a-space>
          </a-collapse-item>
        </a-collapse>
        <div style="margin-top:8px">
          <a-space>
            <a-button type="primary" @click="run">{{ $t('logs.runQuery') }}</a-button>
            <a-button @click="$emit('history')">{{ $t('logs.queryHistory') }}</a-button>
            <a-button @click="emitInspect">{{ $t('logs.queryInspector') }}</a-button>
          </a-space>
        </div>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const emit = defineEmits(['run', 'history', 'inspect'])
const tab = ref('logs')
const lucene = ref('')
const limit = ref(500)

function run() { 
  emit('run', { 
    mode: tab.value === 'raw-data' ? 'raw' : 'code', 
    query: lucene.value, 
    lineLimit: limit.value 
  }) 
}

function emitInspect() {
  emit('inspect', lucene.value)
}
</script>


