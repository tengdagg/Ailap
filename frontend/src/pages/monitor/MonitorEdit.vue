<template>
  <div class="monitor-edit">
    <a-form :model="form" @submit="onSubmit" layout="vertical" class="form-content">
      <a-form-item field="name" label="任务名称" required>
        <a-input v-model="form.name" placeholder="请输入任务名称" />
      </a-form-item>
      
      <a-form-item field="datasourceId" label="数据源" required>
        <a-select v-model="form.datasourceId" placeholder="选择数据源" @change="onDatasourceChange">
          <a-option v-for="ds in datasources" :key="ds.id" :value="String(ds.id)" :label="ds.name + ' (' + ds.type + ')'" />
        </a-select>
      </a-form-item>

      <a-form-item field="cron" label="Cron表达式" required help="支持 @every 格式或标准 Cron。例如: @every 1m (每分钟), @every 1h30m (每1.5小时), 0 30 * * * * (每小时30分), 0 0 12 * * * (每天中午12点)">
        <a-input v-model="form.cron" placeholder="@every 1h" />
      </a-form-item>

      <a-form-item field="query" label="查询语句" help="基础查询语句，如 Loki 的 {app='nginx'} 或 ES 的 service:api">
        <a-textarea v-model="form.query" placeholder="输入基础查询..." />
      </a-form-item>

      <a-form-item field="keywords" label="监控关键词" help="逗号分隔，如: error, exception, 500. 命中任意关键词将触发告警">
        <a-input v-model="form.keywords" placeholder="error, 500" />
      </a-form-item>

      <a-form-item field="channelId" label="通知渠道" required>
        <a-select v-model="form.channelId" placeholder="选择通知渠道">
          <a-option v-for="ch in channels" :key="ch.id" :value="ch.id" :label="ch.name" />
        </a-select>
      </a-form-item>

      <a-form-item field="status" label="状态">
        <a-switch v-model="form.status" checked-value="active" unchecked-value="paused">
            <template #checked>运行中</template>
            <template #unchecked>暂停</template>
        </a-switch>
      </a-form-item>

      <a-form-item>
        <a-button type="primary" html-type="submit">提交</a-button>
        <a-button @click="$router.back()" style="margin-left: 10px">取消</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import request from '@/api/request'

const route = useRoute()
const router = useRouter()
const id = route.params.id

const isEdit = computed(() => id && id !== 'new')

const form = ref({
  name: '',
  datasourceId: '',
  engine: '',
  cron: '@every 1h',
  query: '',
  keywords: 'error',
  channelId: null,
  status: 'active'
})

const datasources = ref([])
const channels = ref([])

const loadMeta = async () => {
    // Load Datasources
    try {
        const { data: resDs } = await request.get('/datasources')
        if (resDs.code === 0) {
            datasources.value = resDs.data.items
        }
        // Load Channels
        const { data: resCh } = await request.get('/channels')
        if (resCh.code === 0) {
            channels.value = resCh.data.items
        }
    } catch (e) { console.error(e) }
}

const loadData = async () => {
    if (!isEdit.value) return
    try {
        const { data: res } = await request.get(`/monitors/${id}`)
        if (res.code === 0) {
            form.value = { 
                ...res.data.item, 
                datasourceId: String(res.data.item.datasourceId) // ensure string
            }
        }
    } catch (e) { console.error(e) }
}

const onDatasourceChange = (val) => {
    const ds = datasources.value.find(d => String(d.id) === val)
    if (ds) {
        form.value.engine = ds.type
    }
}

const onSubmit = async () => {
    try {
        let res
        if (isEdit.value) {
            res = await request.put(`/monitors/${id}`, form.value)
        } else {
            res = await request.post('/monitors', form.value)
        }
        
        const { data } = res
        if (data.code === 0) {
            Message.success('保存成功')
            router.push('/monitors')
        } else {
            Message.error(data.message)
        }
    } catch (e) {
        Message.error('保存失败')
    }
}

onMounted(async () => {
    await loadMeta()
    await loadData()
})
</script>

<style scoped>
.form-content {
  max-width: 600px;
  margin: 20px auto;
}
</style>
