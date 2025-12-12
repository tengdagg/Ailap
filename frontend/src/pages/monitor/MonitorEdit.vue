<template>
  <div class="monitor-edit">
    <a-form :model="form" @submit="onSubmit" layout="vertical" class="form-content">
      <a-form-item field="name" :label="$t('monitor.taskName')" required>
        <a-input v-model="form.name" :placeholder="$t('monitor.placeName')" />
      </a-form-item>
      
      <a-form-item field="datasourceId" :label="$t('monitor.datasource')" required>
        <a-select v-model="form.datasourceId" :placeholder="$t('monitor.placeDs')" @change="onDatasourceChange">
          <a-option v-for="ds in datasources" :key="ds.id" :value="String(ds.id)" :label="ds.name + ' (' + ds.type + ')'" />
        </a-select>
      </a-form-item>

      <a-form-item field="cron" :label="$t('monitor.cron')" required :help="$t('monitor.helpCron')">
        <a-input v-model="form.cron" placeholder="@every 1h" />
      </a-form-item>

      <a-form-item field="query" :label="$t('monitor.query')" :help="$t('monitor.helpQuery')">
        <a-textarea v-model="form.query" :placeholder="$t('monitor.placeQuery')" />
      </a-form-item>

      <a-form-item field="keywords" :label="$t('monitor.keywords')" :help="$t('monitor.helpKw')">
        <a-input v-model="form.keywords" :placeholder="$t('monitor.placeKw')" />
      </a-form-item>

      <a-form-item field="channelId" :label="$t('monitor.channel')" required>
        <a-select v-model="form.channelId" :placeholder="$t('monitor.placeCh')">
          <a-option v-for="ch in channels" :key="ch.id" :value="ch.id" :label="ch.name" />
        </a-select>
      </a-form-item>

      <a-form-item field="status" :label="$t('common.status')">
        <a-switch v-model="form.status" checked-value="active" unchecked-value="paused">
            <template #checked>{{ $t('monitor.active') }}</template>
            <template #unchecked>{{ $t('monitor.paused') }}</template>
        </a-switch>
      </a-form-item>

      <a-form-item>
        <a-button type="primary" html-type="submit">{{ $t('common.submit') }}</a-button>
        <a-button @click="$router.back()" style="margin-left: 10px">{{ $t('common.cancel') }}</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
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
            Message.success(t('common.saveSuccess'))
            router.push('/monitors')
        } else {
            Message.error(data.message)
        }
    } catch (e) {
        Message.error(t('common.saveFail'))
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
