<template>
  <page-container>
    <a-form ref="formRef" :model="form" :rules="rules" layout="vertical" style="max-width:900px; margin:0 auto;">
      <a-form-item :label="$t('common.name')" field="name" required>
        <a-input v-model="form.name" />
      </a-form-item>

      <a-divider>{{ $t('common.connection') }}</a-divider>
      <a-form-item :label="$t('common.url') + ' *'" field="endpoint" required validate-trigger="blur">
        <a-input v-model="form.endpoint" placeholder="http://localhost:3100" />
      </a-form-item>

      <a-divider>{{ $t('common.auth') }}</a-divider>
      <a-form-item :label="$t('common.authType')" field="authType">
        <a-select v-model="form.authType">
          <a-option value="none">无认证</a-option>
          <a-option value="basic">Basic 认证</a-option>
          <a-option value="bearer">Bearer Token</a-option>
        </a-select>
      </a-form-item>
      <a-grid :cols="24" :col-gap="12">
        <a-grid-item :span="12" v-if="form.authType==='basic'">
          <a-form-item :label="$t('common.username')" field="username"><a-input v-model="form.username" /></a-form-item>
        </a-grid-item>
        <a-grid-item :span="12" v-if="form.authType==='basic'">
          <a-form-item :label="$t('common.password')" field="password"><a-input-password v-model="form.password" /></a-form-item>
        </a-grid-item>
        <a-grid-item :span="24" v-if="form.authType==='bearer'">
          <a-form-item :label="$t('common.token')" field="token"><a-input v-model="form.token" /></a-form-item>
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
          <a-form-item :label="$t('common.caCert') + ' *'"><a-textarea v-model="form.tls.caCert" :auto-size="{minRows:6}" placeholder="Min --- BEGIN CERTIFICATE ---" /></a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="selfSignedVisible=false">{{ $t('common.testFail') }}</a-button> <!-- Hacky usage, better 'Close' -->
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
      <a-form-item :label="$t('common.maxLines')">
        <a-input-number v-model="form.query.maxLines" :min="1" :max="10000" />
      </a-form-item>

      <a-divider>{{ $t('common.derivedFields') }}</a-divider>
      <div>
        <a-button type="outline" size="small" @click="addDerived">+ {{ $t('common.new') }}</a-button>
        <div v-for="(f, i) in form.derivedFields" :key="i" style="margin-top:8px">
          <a-grid :cols="24" :col-gap="8">
            <a-grid-item :span="8"><a-input v-model="f.name" placeholder="字段名" /></a-grid-item>
            <a-grid-item :span="8"><a-input v-model="f.matcherRegex" placeholder="匹配正则" /></a-grid-item>
            <a-grid-item :span="6"><a-input v-model="f.url" placeholder="链接 URL" /></a-grid-item>
            <a-grid-item :span="2" style="display:flex;align-items:center"><a-button size="mini" status="danger" @click="removeDerived(i)">{{ $t('common.delete') }}</a-button></a-grid-item>
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

const formRef = ref(null)

const router = useRouter()
const route = useRoute()
const { t } = useI18n()

const urlPattern = /^(https?:)\/\//i
const selfSignedVisible = ref(false)
const clientAuthVisible = ref(false)

const testing = ref(false)
const saving = ref(false)

const form = reactive({
  name: 'loki',
  type: 'loki',
  endpoint: '',
  authType: 'none',
  username: '', password: '', token: '',
  tls: { addSelfSigned: false, clientAuth: false, skipVerify: false, caCert: '', serverName: '', clientCert: '', clientKey: '' },
  http: { allowedCookies: [], timeout: 30 },
  query: { maxLines: 1000 },
  derivedFields: [],
})

const rules = {
  name: [{ required: true, message: '名称不能为空' }],
  endpoint: [
    { required: true, message: 'URL 不能为空' },
    { match: urlPattern, message: '请输入合法的 URL（http/https）' },
  ],
}

function onToggleSelfSigned(v) { if (v) selfSignedVisible.value = true }
function onToggleClientAuth(v) { if (v) clientAuthVisible.value = true }

function addDerived() { form.derivedFields.push({ name: '', matcherRegex: '', url: '' }) }
function removeDerived(i) { form.derivedFields.splice(i, 1) }

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
