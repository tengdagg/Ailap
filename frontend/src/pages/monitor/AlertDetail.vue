<template>
  <page-container>
    <a-card :bordered="false" :title="$t('monitor.alertDetails')">
      <template #extra>
        <a-button type="secondary" @click="$router.back()">
          {{ $t('common.back') }}
        </a-button>
      </template>

      <a-spin :loading="loading" style="width: 100%">
        <div v-if="alert">
          <a-descriptions :column="1" bordered size="large" class="alert-info">
            <a-descriptions-item :label="$t('monitor.time')">
              {{ formatTime(alert.createdAt) }}
            </a-descriptions-item>
            <a-descriptions-item :label="$t('monitor.taskName')">
              {{ alert.monitor?.name }}
            </a-descriptions-item>
            <a-descriptions-item :label="$t('common.status')">
              <a-tag v-if="alert.status === 'failed'" color="red">{{ $t('common.failed') }}</a-tag>
              <a-tag v-else color="green">{{ $t('common.success') }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item :label="$t('common.error')" v-if="alert.status === 'failed'">
              <span style="color: red">{{ alert.error }}</span>
            </a-descriptions-item>
          </a-descriptions>

          <div style="margin-top: 20px;">
            <div style="font-weight: bold; margin-bottom: 10px;">{{ $t('common.content') }}</div>
            <div class="content-box">
              {{ alert.content }}
            </div>
          </div>
        </div>
        <a-empty v-else-if="!loading" />
      </a-spin>
    </a-card>
  </page-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import PageContainer from '@/components/PageContainer.vue'
import { getAlert } from '@/api/monitor'

const route = useRoute()
const loading = ref(false)
const alert = ref(null)

const fetchData = async () => {
  const id = route.params.id
  if (!id) return
  
  loading.value = true
  try {
    const { data } = await getAlert(id)
    alert.value = data?.data?.item
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

const formatTime = (iso) => {
    if(!iso) return ''
    return new Date(iso).toLocaleString()
}

onMounted(fetchData)
</script>

<style scoped>
.alert-info {
  margin-bottom: 20px;
}
.content-box {
  background: var(--color-fill-2);
  padding: 16px;
  border-radius: 4px;
  font-family: monospace;
  white-space: pre-wrap;
  color: var(--color-text-1);
}
</style>
