<template>
  <page-container>
    <a-card :bordered="false" :title="null">
      <a-table :data="alerts" :loading="loading" :pagination="{ pageSize: 20 }">
        <template #columns>
          <a-table-column :title="$t('monitor.time')" data-index="createdAt" :width="200">
             <template #cell="{ record }">
               {{ formatTime(record.createdAt) }}
             </template>
          </a-table-column>
          <a-table-column :title="$t('monitor.taskName')" data-index="monitor.name" :width="150">
             <template #cell="{ record }">
               {{ record.monitor?.name || '-' }}
             </template>
          </a-table-column>
          <a-table-column :title="$t('common.status')" data-index="status" :width="80">
            <template #cell="{ record }">
               <a-tag v-if="record.status === 'failed'" color="red">{{ $t('common.failed') }}</a-tag>
               <a-tag v-else color="green">{{ $t('common.success') }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column :title="$t('common.details')" data-index="content">
            <template #cell="{ record }">
               <div style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 600px; cursor: pointer;" @click="openDetail(record)">
                  {{ record.content }}
               </div>
            </template>
          </a-table-column>
          <a-table-column :title="$t('common.actions')" :width="100">
              <template #cell="{ record }">
                  <a-button size="mini" @click="openDetail(record)">{{ $t('common.view') }}</a-button>
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
import { listAlerts } from '@/api/monitor'

const router = useRouter()
const alerts = ref([])
const loading = ref(false)

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await listAlerts()
    alerts.value = data?.data?.items || []
  } finally {
    loading.value = false
  }
}

const openDetail = (record) => {
    router.push(`/monitors/alerts/${record.id}`)
}

const formatTime = (iso) => {
    if(!iso) return ''
    return new Date(iso).toLocaleString()
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
