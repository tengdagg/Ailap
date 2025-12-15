<template>
  <page-container>
    <a-card :bordered="false" :title="null">
      <a-table :data="items" :loading="loading" :pagination="{ pageSize: 20 }">
        <template #columns>
          <a-table-column :title="$t('logs.time')" data-index="createdAt" :width="200">
             <template #cell="{ record }">
               {{ formatTime(record.createdAt) }}
             </template>
          </a-table-column>
          <a-table-column :title="$t('logs.engine')" data-index="engine" :width="120">
             <template #cell="{ record }">
                <a-tag :color="getEngineColor(record.engine)">{{ record.engine }}</a-tag>
             </template>
          </a-table-column>
          <a-table-column :title="$t('logs.mode')" data-index="mode" :width="100">
             <template #cell="{ record }">
                <a-tag>{{ record.mode }}</a-tag>
             </template>
          </a-table-column>
          <a-table-column :title="$t('logs.query')" data-index="query">
            <template #cell="{ record }">
               <div style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 600px; cursor: pointer; font-family: monospace;" @click="runQuery(record)">
                  {{ record.query }}
               </div>
            </template>
          </a-table-column>
          <a-table-column :title="$t('common.actions')" :width="100">
              <template #cell="{ record }">
                  <a-button size="mini" @click="runQuery(record)">{{ $t('common.view') }}</a-button>
              </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>
  </page-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import PageContainer from '@/components/PageContainer.vue'
import { history } from '@/api/logs'

const router = useRouter()
const items = ref([])
const loading = ref(false)

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await history('recent')
    items.value = data?.data?.items || []
  } finally {
    loading.value = false
  }
}

const runQuery = (record) => {
    router.push({
        path: '/logs',
        query: {
            engine: record.engine,
            query: record.query,
            mode: record.mode
        }
    })
}

const formatTime = (iso) => {
    if(!iso) return ''
    return new Date(iso).toLocaleString()
}

const getEngineColor = (engine) => {
    switch(engine) {
        case 'loki': return 'blue'
        case 'elasticsearch': return 'green'
        case 'victorialogs': return 'orange'
        default: return 'gray'
    }
}

onMounted(fetchData)
</script>

<style scoped>
:deep(.arco-table-cell) {
    font-size: 13px;
    color: var(--color-text-2);
}
:deep(.arco-table-th) {
    background-color: var(--color-fill-2);
    font-weight: 600;
    font-size: 13px;
}
</style>
