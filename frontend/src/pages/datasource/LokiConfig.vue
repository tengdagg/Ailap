<template>
  <page-container title="Loki 配置" subtitle="配置 Loki 数据源连接信息">
    <a-form :model="form" :rules="rules" layout="vertical" style="max-width:900px">
      <a-grid :cols="24" :col-gap="16">
        <a-grid-item :span="18">
          <a-form-item label="名称" field="name" required>
            <a-input v-model="form.name" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="6" style="display:flex;align-items:flex-end">
          <a-switch v-model="form.isDefault">默认</a-switch>
        </a-grid-item>
      </a-grid>

      <a-divider>连接</a-divider>
      <a-form-item label="URL *" field="endpoint" required validate-trigger="blur">
        <a-input v-model="form.endpoint" placeholder="http://localhost:3100" />
      </a-form-item>

      <a-divider>认证</a-divider>
      <a-form-item label="认证方式" field="authType">
        <a-select v-model="form.authType">
          <a-option value="none">无认证</a-option>
          <a-option value="basic">Basic 认证</a-option>
          <a-option value="bearer">Bearer Token</a-option>
        </a-select>
      </a-form-item>
      <a-grid :cols="24" :col-gap="12">
        <a-grid-item :span="12" v-if="form.authType==='basic'">
          <a-form-item label="用户名" field="username"><a-input v-model="form.username" /></a-form-item>
        </a-grid-item>
        <a-grid-item :span="12" v-if="form.authType==='basic'">
          <a-form-item label="密码" field="password"><a-input-password v-model="form.password" /></a-form-item>
        </a-grid-item>
        <a-grid-item :span="24" v-if="form.authType==='bearer'">
          <a-form-item label="Token" field="token"><a-input v-model="form.token" /></a-form-item>
        </a-grid-item>
      </a-grid>

      <a-divider>TLS 设置</a-divider>
      <a-space direction="vertical">
        <a-checkbox v-model="form.tls.addSelfSigned" @change="onToggleSelfSigned">接受自签名证书</a-checkbox>
        <a-checkbox v-model="form.tls.clientAuth" @change="onToggleClientAuth">TLS 客户端认证</a-checkbox>
        <a-checkbox v-model="form.tls.skipVerify">跳过证书校验</a-checkbox>
      </a-space>

      <a-modal v-model:visible="selfSignedVisible" title="自签名证书">
        <a-form layout="vertical">
          <a-form-item label="CA 证书 *"><a-textarea v-model="form.tls.caCert" :auto-size="{minRows:6}" placeholder="以 --- BEGIN CERTIFICATE --- 开头" /></a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="selfSignedVisible=false">关闭</a-button>
          <a-button type="primary" @click="selfSignedVisible=false">确定</a-button>
        </template>
      </a-modal>

      <a-modal v-model:visible="clientAuthVisible" title="TLS 客户端认证">
        <a-form layout="vertical">
          <a-form-item label="ServerName *"><a-input v-model="form.tls.serverName" placeholder="domain.example.com" /></a-form-item>
          <a-form-item label="客户端证书 *"><a-textarea v-model="form.tls.clientCert" :auto-size="{minRows:6}" placeholder="以 --- BEGIN CERTIFICATE --- 开头" /></a-form-item>
          <a-form-item label="客户端私钥 *"><a-textarea v-model="form.tls.clientKey" :auto-size="{minRows:6}" placeholder="以 --- RSA PRIVATE KEY CERTIFICATE --- 开头" /></a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="clientAuthVisible=false">关闭</a-button>
          <a-button type="primary" @click="clientAuthVisible=false">确定</a-button>
        </template>
      </a-modal>

      <a-divider>高级设置</a-divider>
      <a-form-item label="允许的 Cookie"><a-input-tag v-model="form.http.allowedCookies" placeholder="回车新增" /></a-form-item>
      <a-form-item label="超时时间 (秒)"><a-input-number v-model="form.http.timeout" placeholder="单位：秒" :min="0" /></a-form-item>

      <a-form-item label="告警"><a-switch v-model="form.alerting.manageInUI">在 UI 中管理告警规则</a-switch></a-form-item>
      <a-form-item label="查询最大行数"><a-input-number v-model="form.query.maxLines" :min="1" /></a-form-item>

      <a-divider>派生字段（Derived fields）</a-divider>
      <div>
        <a-button type="outline" size="small" @click="addDerived">+ 新增</a-button>
        <div v-for="(f, i) in form.derivedFields" :key="i" style="margin-top:8px">
          <a-grid :cols="24" :col-gap="8">
            <a-grid-item :span="8"><a-input v-model="f.name" placeholder="字段名" /></a-grid-item>
            <a-grid-item :span="8"><a-input v-model="f.matcherRegex" placeholder="匹配正则" /></a-grid-item>
            <a-grid-item :span="6"><a-input v-model="f.url" placeholder="链接 URL" /></a-grid-item>
            <a-grid-item :span="2" style="display:flex;align-items:center"><a-button size="mini" status="danger" @click="removeDerived(i)">删除</a-button></a-grid-item>
          </a-grid>
        </div>
      </div>

      <a-space style="margin-top:16px">
        <a-button @click="$router.back()">返回</a-button>
        <a-button @click="onTest" :loading="testing">测试连接</a-button>
        <a-button type="primary" @click="onSave" :loading="saving">保存</a-button>
      </a-space>
    </a-form>
  </page-container>
</template>
<script setup>
import { reactive, ref } from 'vue'
import PageContainer from '@/components/PageContainer.vue'
import { createDataSource, testConnectionPayload } from '@/api/datasources'
import { Message } from '@arco-design/web-vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const urlPattern = /^(https?:)\/\//i
const selfSignedVisible = ref(false)
const clientAuthVisible = ref(false)

const testing = ref(false)
const saving = ref(false)

const form = reactive({
  name: 'loki',
  type: 'loki',
  isDefault: false,
  endpoint: '',
  authType: 'none',
  username: '', password: '', token: '',
  tls: { addSelfSigned: false, clientAuth: false, skipVerify: false, caCert: '', serverName: '', clientCert: '', clientKey: '' },
  http: { allowedCookies: [], timeout: undefined },
  alerting: { manageInUI: true },
  query: { maxLines: 1000 },
  derivedFields: [],
})

const rules = {
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
    const { data } = await createDataSource(form)
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
</script>
