<template>
  <page-container>
    <a-card :bordered="false" :title="null">
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
              <a-tag v-else color="orangered">{{ $t('common.triggered') }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item :label="$t('common.error')" v-if="alert.status === 'failed'">
              <span style="color: red">{{ alert.error }}</span>
            </a-descriptions-item>
          </a-descriptions>

          <div style="margin-top: 20px;">
            <div style="font-weight: bold; margin-bottom: 10px;">{{ $t('common.content') }}</div>
            <div class="content-box markdown-body" v-html="renderMarkdown(alert.content)"></div>
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
import { marked } from 'marked'
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
    loading.value = false
  }
}

const formatTime = (iso) => {
    if(!iso) return ''
    return new Date(iso).toLocaleString()
}

const renderMarkdown = (text) => {
    if (!text) return ''
    return marked.parse(text)
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
  color: var(--color-text-1);
}
:deep(.markdown-body h1), :deep(.markdown-body h2) {
    border-bottom: 1px solid var(--color-border-2);
    padding-bottom: 0.3em;
    margin-bottom: 16px;
    margin-top: 24px;
}
:deep(.markdown-body p) {
    margin-bottom: 16px;
    line-height: 1.6;
}
:deep(.markdown-body code) {
    background-color: var(--color-fill-3);
    padding: 2px 4px;
    border-radius: 4px;
    font-family: monospace;
}
:deep(.markdown-body pre) {
    background-color: var(--color-fill-3);
    padding: 16px;
    border-radius: 8px;
    overflow: auto;
    margin-bottom: 16px;
}
:deep(.markdown-body ul), :deep(.markdown-body ol) {
    padding-left: 20px;
    margin-bottom: 16px;
}
</style>

