<template>
  <div class="channel-list">
    <a-space direction="vertical" fill>
      <div class="header">
        <a-button type="primary" @click="$router.push('/channels/new')" size="small">
          <template #icon><icon-plus /></template>
          {{ $t('channel.newChannel') }}
        </a-button>
      </div>

      <a-table :data="items" :loading="loading" row-key="id">
        <template #columns>
          <a-table-column :title="$t('common.name')" data-index="name" />
          <a-table-column :title="$t('common.type')" data-index="type">
              <template #cell="{ record }">
                  <a-tag v-if="record.type === 'webhook'" color="blue">Webhook</a-tag>
                  <a-tag v-else-if="record.type === 'email'" color="arcoblue">{{ $t('channel.title') }}</a-tag>
                  <a-tag v-else>{{ record.type }}</a-tag>
              </template>
          </a-table-column>
          <a-table-column :title="$t('common.actions')">
            <template #cell="{ record }">
              <a-space>
                <a-button size="small" @click="$router.push(`/channels/${record.id}`)">{{ $t('common.edit') }}</a-button>
                <a-popconfirm :content="$t('common.confirm') + '?'" @ok="doDelete(record.id)">
                  <a-button size="small" status="danger">{{ $t('common.delete') }}</a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-space>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'

const items = ref([])
const loading = ref(false)
const { t } = useI18n()

const loadData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/channels')
    if (data.code === 0) {
      items.value = data.data.items
    }
  } catch (e) {
      console.error(e)
  } finally {
    loading.value = false
  }
}

const doDelete = async (id) => {
  try {
    const { data } = await request.delete(`/channels/${id}`)
    if (data.code === 0) {
        Message.success(t('common.deleteSuccess'))
        loadData()
    } else {
        Message.error(data.message)
    }
  } catch (e) {
      Message.error(t('common.deleteFail'))
  }
}

onMounted(loadData)
</script>

<style scoped>
.header {
  display: flex;
  justify-content: flex-end; /* Align button to right */
  align-items: center;
  margin-bottom: 16px;
}
:deep(.arco-table-th) {
  background-color: var(--color-fill-2);
  font-weight: 600;
  font-size: 13px; /* Reduced font size */
}
</style>
