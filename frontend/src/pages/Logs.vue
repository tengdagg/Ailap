<template>
  <page-container title="日志查询" subtitle="Grafana 风格查询编辑器">
    <div style="margin-bottom:12px; display:flex; gap:12px; align-items:center">
      <span>数据源</span>
      <a-select v-model="datasource" :options="dsOptions" style="width:200px" />
      <a-segmented v-model="mode" :options="['Builder','Code']" />
    </div>

    <loki-editor v-if="datasource==='loki'" @run="onRunLoki" />
    <elasticsearch-editor v-else @run="onRunES" />

    <a-table :data="rows" :loading="loading" row-key="id" style="margin-top:12px" :pagination="{ pageSize: 20 }">
      <a-table-column title="时间" data-index="timestamp" :width="220" />
      <a-table-column title="级别" data-index="level" :width="100" />
      <a-table-column title="内容" data-index="message" />
    </a-table>
  </page-container>
</template>
<script setup>
import { ref } from 'vue'
import PageContainer from '@/components/PageContainer.vue'
import LokiEditor from '@/components/logs/LokiEditor.vue'
import ElasticsearchEditor from '@/components/logs/ElasticsearchEditor.vue'
import { queryLogs } from '@/api/logs'

const datasource = ref('loki')
const dsOptions = [ { label: 'Loki', value: 'loki' }, { label: 'Elasticsearch', value: 'elasticsearch' } ]
const mode = ref('Builder')

const loading = ref(false)
const rows = ref([])

async function onRunLoki(payload) {
  await runQuery({ engine: 'loki', payload })
}
async function onRunES(payload) {
  await runQuery({ engine: 'elasticsearch', payload })
}

async function runQuery(params) {
  loading.value = true
  try {
    const { data } = await queryLogs({ engine: params.engine, ...params.payload })
    rows.value = data?.data?.items || []
  } finally {
    loading.value = false
  }
}
</script>

