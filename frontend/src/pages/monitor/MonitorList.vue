<template>
  <div class="monitor-list">
    <a-space direction="vertical" fill>
      <div class="header">
        <a-button type="primary" @click="$router.push('/monitors/new')" size="small">
          <template #icon><icon-plus /></template>
          新建任务
        </a-button>
      </div>

      <a-table :data="items" :loading="loading" row-key="id">
        <template #columns>
          <a-table-column title="名称" data-index="name" />
          <a-table-column title="引擎" data-index="engine">
             <template #cell="{ record }">
               <a-tag color="blue" v-if="record.engine==='loki'">Loki</a-tag>
               <a-tag color="green" v-else-if="record.engine==='elasticsearch'">ES</a-tag>
               <a-tag color="orange" v-else-if="record.engine==='victorialogs'">VictoriaLogs</a-tag>
               <a-tag v-else>{{ record.engine }}</a-tag>
             </template>
          </a-table-column>
          <a-table-column title="Cron表达式" data-index="cron" />
          <a-table-column title="关键词" data-index="keywords" />
          <a-table-column title="状态" data-index="status">
             <template #cell="{ record }">
               <a-badge :status="record.status === 'active' ? 'success' : 'warning'" :text="record.status === 'active' ? '运行中' : '已暂停'" />
             </template>
          </a-table-column>
          <a-table-column title="上次运行" data-index="lastRunAt">
             <template #cell="{ record }">
               {{ record.lastRunAt ? new Date(record.lastRunAt).toLocaleString() : '-' }}
             </template>
          </a-table-column>
          <a-table-column title="操作">
            <template #cell="{ record }">
              <a-space>
                <a-button size="small" @click="$router.push(`/monitors/${record.id}`)">编辑</a-button>
                <a-popconfirm content="确定删除吗?" @ok="doDelete(record.id)">
                  <a-button size="small" status="danger">删除</a-button>
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
import request from '@/api/request'

const items = ref([])
const loading = ref(false)

const loadData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/monitors')
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
    const { data } = await request.delete(`/monitors/${id}`)
    if (data.code === 0) {
        Message.success('删除成功')
        loadData()
    } else {
        Message.error(data.message)
    }
  } catch (e) {
      Message.error('删除失败')
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
