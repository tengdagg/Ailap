<template>
  <page-container>
    <a-form ref="formRef" :model="form" :rules="rules" layout="vertical" style="max-width:900px; margin:0 auto;">
      <a-form-item :label="$t('common.name')" field="name" required>
        <a-input v-model="form.name" />
      </a-form-item>

      <a-divider>{{ $t('common.connection') }}</a-divider>
      <a-form-item :label="$t('common.url') + ' *'" field="endpoint" required validate-trigger="blur">
        <a-input v-model="form.endpoint" placeholder="http://localhost:9200" />
      </a-form-item>

      <a-divider>{{ $t('common.auth') }}</a-divider>
      <a-form-item :label="$t('common.authType')" field="authType">
        <a-select v-model="form.authType">
          <a-option value="none">无认证</a-option>
          <a-option value="basic">Basic 认证</a-option>
          <a-option value="apiKey">API Key</a-option>
        </a-select>
      </a-form-item>
      <a-grid :cols="24" :col-gap="12">
        <a-grid-item :span="12" v-if="form.authType==='basic'">
          <a-form-item :label="$t('common.username')" field="username">
            <a-input v-model="form.username" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12" v-if="form.authType==='basic'">
          <a-form-item :label="$t('common.password')" field="password">
            <a-input-password v-model="form.password" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="24" v-if="form.authType==='apiKey'">
          <a-form-item :label="$t('common.apiKey')" field="apiKey">
            <a-input v-model="form.apiKey" />
          </a-form-item>
        </a-grid-item>
      </a-grid>

      <a-divider>{{ $t('common.tls') }}</a-divider>
      <a-space direction="vertical">
        <a-checkbox v-model="form.tls.addSelfSigned" @change="onToggleSelfSigned">{{ $t('common.addSelfSigned') }}</a-checkbox>
        <a-checkbox v-model="form.tls.clientAuth" @change="onToggleClientAuth">{{ $t('common.clientAuth') }}</a-checkbox>
        <a-checkbox v-model="form.tls.skipVerify">{{ $t('common.skipVerify') }}</a-checkbox>
      </a-space>

      <a-modal v-model:visible="selfSignedVisible" :title="$t('common.addSelfSigned')">
        <a-form layout="vertical">
          <a-form-item :label="$t('common.caCert') + ' *'">
            <a-textarea v-model="form.tls.caCert" :auto-size="{minRows:6}" placeholder="Min --- BEGIN CERTIFICATE ---" />
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="selfSignedVisible=false">{{ $t('common.cancel') }}</a-button>
          <a-button type="primary" @click="selfSignedVisible=false">{{ $t('common.confirm') }}</a-button>
        </template>
      </a-modal>

      <a-modal v-model:visible="clientAuthVisible" :title="$t('common.clientAuth')">
        <a-form layout="vertical">
          <a-form-item :label="$t('common.serverName') + ' *'"><a-input v-model="form.tls.serverName" placeholder="domain.example.com" /></a-form-item>
          <a-form-item :label="$t('common.clientCert') + ' *'"><a-textarea v-model="form.tls.clientCert" :auto-size="{minRows:6}" placeholder="Min --- BEGIN CERTIFICATE ---" /></a-form-item>
          <a-form-item :label="$t('common.clientKey') + ' *'"><a-textarea v-model="form.tls.clientKey" :auto-size="{minRows:6}" placeholder="Min --- RSA PRIVATE KEY CERTIFICATE ---" /></a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="clientAuthVisible=false">{{ $t('common.cancel') }}</a-button>
          <a-button type="primary" @click="clientAuthVisible=false">{{ $t('common.confirm') }}</a-button>
        </template>
      </a-modal>

      <a-divider>{{ $t('common.advanced') }}</a-divider>
      <a-form-item :label="$t('common.allowedCookies')">
        <a-input-tag v-model="form.http.allowedCookies" placeholder="Enter..." />
      </a-form-item>
      <a-form-item :label="$t('common.timeout')">
        <a-input-number v-model="form.http.timeout" placeholder="30" :min="1" :max="300" />
      </a-form-item>

      <a-divider>Elasticsearch</a-divider>
      <a-grid :cols="24" :col-gap="12">
        <a-grid-item :span="12">
          <a-form-item :label="$t('common.index')" field="index">
            <a-input v-model="form.es.index" placeholder="es-index-name" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item :label="$t('common.pattern')" field="pattern">
            <a-select v-model="form.es.pattern">
              <a-option value="none">无</a-option>
              <a-option value="daily">按天</a-option>
              <a-option value="weekly">按周</a-option>
              <a-option value="monthly">按月</a-option>
            </a-select>
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item :label="$t('common.timeField')" field="timeField">
            <a-input v-model="form.es.timeField" placeholder="@timestamp" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item :label="$t('common.maxShardReq')" field="maxShardReq">
            <a-input-number v-model="form.es.maxShardRequests" :min="1" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item :label="$t('common.minInterval')" field="minInterval">
            <a-input v-model="form.es.minTimeInterval" placeholder="10s" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12" style="display:flex;align-items:flex-end">
          <a-form-item :label="$t('common.xpack')">
            <a-switch v-model="form.es.xpack" />
          </a-form-item>
        </a-grid-item>
      </a-grid>

      <a-divider>日志字段</a-divider>
      <a-grid :cols="24" :col-gap="12">
        <a-grid-item :span="12">
          <a-form-item :label="$t('common.messageField')">
            <a-input v-model="form.logs.messageField" placeholder="_source" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item :label="$t('common.levelField')">
            <a-input v-model="form.logs.levelField" />
          </a-form-item>
        </a-grid-item>
      </a-grid>

      <a-divider>{{ $t('common.dataLinks') }}</a-divider>
      <div>
        <a-button type="outline" size="small" @click="addDataLink">+ {{ $t('common.new') }}</a-button>
        <div v-for="(link, i) in form.dataLinks" :key="i" style="margin-top:8px">
          <a-grid :cols="24" :col-gap="8">
            <a-grid-item :span="10">
              <a-input v-model="link.field" placeholder="字段名" />
            </a-grid-item>
            <a-grid-item :span="12">
              <a-input v-model="link.url" placeholder="链接 URL" />
            </a-grid-item>
            <a-grid-item :span="2" style="display:flex;align-items:center">
              <a-button size="mini" status="danger" @click="removeDataLink(i)">{{ $t('common.delete') }}</a-button>
            </a-grid-item>
          </a-grid>
        </div>
      </div>

      <a-space style="margin-top:16px">
        <a-button @click="$router.back()">{{ $t('common.return') }}</a-button>
        <a-button @click="onTest" :loading="testing">{{ $t('common.testConnection') }}</a-button>
        <a-button type="primary" @click="onSave" :loading="saving">{{ $t('common.save') }}</a-button>
      </a-space>
    </a-form>
  </page-container>
</template>
<script setup>
import { reactive, ref } from 'vue'
import PageContainer from '@/components/PageContainer.vue'
import { createDataSource, testConnectionPayload, getDataSourceById, updateDataSource } from '@/api/datasources'
import { Message } from '@arco-design/web-vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const formRef = ref(null)
const route = useRoute()
const { t } = useI18n()

const urlPattern = /^(https?:)\/\//i

const selfSignedVisible = ref(false)
const clientAuthVisible = ref(false)

const form = reactive({
  name: 'elasticsearch',
  type: 'elasticsearch',
  endpoint: '',
  authType: 'none',
  username: '',
  password: '',
  apiKey: '',
  tls: { addSelfSigned: false, clientAuth: false, skipVerify: false, caCert: '', serverName: '', clientCert: '', clientKey: '' },
  http: { allowedCookies: [], timeout: 30 },
  es: {
    index: '',
    pattern: 'none',
    timeField: '@timestamp',
    maxShardRequests: 5,
    minTimeInterval: '10s',
    xpack: false,
  },
  logs: { messageField: '_source', levelField: '' },
  dataLinks: [],
})

const rules = {
  name: [{ required: true, message: '名称不能为空' }],
  endpoint: [
    { required: true, message: 'URL 不能为空' },
    { match: urlPattern, message: '请输入合法的 URL（http/https）' },
  ],
}

const testing = ref(false)
const saving = ref(false)

function onToggleSelfSigned(v) { if (v) selfSignedVisible.value = true }
function onToggleClientAuth(v) { if (v) clientAuthVisible.value = true }

function addDataLink() { form.dataLinks.push({ field: '', url: '' }) }
function removeDataLink(i) { form.dataLinks.splice(i, 1) }

async function onTest() {
  testing.value = true
  try {
    const { data } = await testConnectionPayload(form)
    if (data?.code === 0) Message.success('连接成功')
    else Message.error(data?.message || '连接失败')
  } finally {
    testing.value = false
  }
}

async function onSave() {
  saving.value = true
  try {
    const res = await formRef.value?.validate()
    if (res) return

    const id = route.query.id
    let resp
    if (id) {
      resp = await updateDataSource(id, form)
    } else {
      resp = await createDataSource(form)
    }
    const { data } = resp
    if (data?.code === 0) {
      Message.success('已保存')
      router.replace({ path: '/datasources', query: { ts: Date.now().toString() } })
    } else {
      Message.error(data?.message || '保存失败')
    }
  } finally {
    saving.value = false
  }
}

// edit mode: load existing
if (route.query.id) {
  ;(async () => {
    const existing = await getDataSourceById(route.query.id)
    if (existing && existing.config) {
      try { Object.assign(form, JSON.parse(existing.config)) } catch (_) {}
      form.name = existing.name || form.name
      form.type = existing.type || form.type
      form.endpoint = existing.endpoint || form.endpoint
    }
  })()
}
</script>
