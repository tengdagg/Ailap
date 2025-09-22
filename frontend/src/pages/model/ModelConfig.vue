<template>
  <page-container>
    <a-form :model="form" layout="vertical" style="max-width:900px">
      <a-grid :cols="24" :col-gap="12">
        <a-grid-item :span="12">
          <a-form-item label="名称">
            <a-input v-model="form.name" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item label="供应商">
            <a-select v-model="form.provider" :disabled="true">
              <a-option value="openai">OpenAI</a-option>
              <a-option value="deepseek">Deepseek</a-option>
              <a-option value="qwen">Qwen</a-option>
            </a-select>
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item label="模型">
            <a-select v-model="form.model" :options="modelOptions" allow-search allow-create />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
            <a-form-item label="API Base">
            <a-input v-model="form.apiBase" :placeholder="apiBasePlaceholder" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="24">
          <a-form-item label="API Key">
            <a-input-password v-model="form.apiKey" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item label="Temperature">
            <a-input-number v-model="form.temperature" :min="0" :max="2" :step="0.1" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item label="Max Tokens">
            <a-input-number v-model="form.maxTokens" :min="1" :max="32000" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item label="是否启用">
            <a-switch v-model="form.enabled" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item label="设为默认">
            <a-switch v-model="form.isDefault" />
          </a-form-item>
        </a-grid-item>
      </a-grid>

      <a-divider>角色定义</a-divider>
      <a-space direction="vertical" fill>
        <div v-for="(r, idx) in roles" :key="idx" style="border:1px solid var(--color-border-2); padding:12px; border-radius:8px;">
          <a-grid :cols="24" :col-gap="8">
            <a-grid-item :span="8"><a-input v-model="r.name" placeholder="角色名，如：运维助手" /></a-grid-item>
            <a-grid-item :span="16"><a-input v-model="r.description" placeholder="角色描述" /></a-grid-item>
            <a-grid-item :span="24" style="margin-top:8px"><a-textarea v-model="r.systemPrompt" placeholder="系统提示词" :auto-size="{minRows:2, maxRows:6}" /></a-grid-item>
          </a-grid>
          <div style="display:flex; justify-content:flex-end; margin-top:8px">
            <a-button size="mini" status="danger" @click="removeRole(idx)">删除</a-button>
          </div>
        </div>
        <a-button type="outline" size="small" @click="addRole">+ 新增角色</a-button>
      </a-space>

      <a-space style="margin-top:16px">
        <a-button @click="$emit('back')">返回</a-button>
        <a-button @click="onTest" :loading="testing">测试</a-button>
        <a-button type="primary" @click="onSave" :loading="saving">保存</a-button>
      </a-space>
    </a-form>
  </page-container>
</template>
<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { Message } from '@arco-design/web-vue'
import PageContainer from '@/components/PageContainer.vue'
import { listModels, updateModel, createModel, testModel } from '@/api/models'
import request from '@/api/request'

const props = defineProps({
  modelId: { type: [String, Number], required: false },
  preset: { type: Object, required: false, default: null },
})
const emit = defineEmits(['back', 'saved'])

const form = ref({ name: '', provider: 'openai', model: '', apiBase: '', apiKey: '', temperature: 0.7, maxTokens: 2048, enabled: true, isDefault: false })
const roles = ref([])
const testing = ref(false)
const saving = ref(false)

const providerModels = {
  // OpenAI reference: https://platform.openai.com/docs/api-reference/introduction
  openai: [
    { label: 'gpt-4o', value: 'gpt-4o' },
    { label: 'gpt-4o-mini', value: 'gpt-4o-mini' },
    { label: 'gpt-4.1', value: 'gpt-4.1' },
    { label: 'gpt-4.1-mini', value: 'gpt-4.1-mini' },
    { label: 'gpt-3.5-turbo', value: 'gpt-3.5-turbo' },
  ],
  // DeepSeek reference: https://api-docs.deepseek.com/zh-cn/
  deepseek: [
    { label: 'deepseek-chat (V3.1 非思考)', value: 'deepseek-chat' },
    { label: 'deepseek-reasoner (V3.1 思考)', value: 'deepseek-reasoner' },
  ],
  // Qwen
  qwen: [
    { label: 'qwen2.5', value: 'qwen2.5' },
    { label: 'qwen2.5-instruct', value: 'qwen2.5-instruct' },
    { label: 'qwen2.5-coder', value: 'qwen2.5-coder' },
  ],
}

const providerApiBase = {
  openai: 'https://api.openai.com/v1',
  deepseek: 'https://api.deepseek.com',
  qwen: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
}

const modelOptions = computed(() => providerModels[form.value.provider] || [])
const apiBasePlaceholder = computed(() => providerApiBase[form.value.provider] || 'https://api.example.com/v1')

watch(() => form.value.provider, (p) => {
  // when provider changes, if current model not in list, reset
  const list = providerModels[p] || []
  if (!list.find(x => x.value === form.value.model)) {
    form.value.model = list[0]?.value || ''
  }
  // Do not auto-fill apiBase; show provider placeholder in grey instead
})

function addRole() { roles.value.push({ name: '', description: '', systemPrompt: '' }) }
function removeRole(i) { roles.value.splice(i, 1) }

async function load() {
  if (props.modelId) {
    const { data } = await listModels()
    const items = data?.data?.items || []
    const m = items.find(it => String(it.id) === String(props.modelId))
    if (m) {
      form.value = { name: m.name, provider: m.provider, model: m.model, apiBase: m.apiBase, apiKey: m.apiKey, temperature: m.temperature, maxTokens: m.maxTokens, enabled: !!m.enabled, isDefault: !!m.isDefault }
      try { roles.value = JSON.parse(m.roles || '[]') } catch { roles.value = [] }
    }
  } else if (props.preset) {
    form.value = {
      name: props.preset.name || '',
      provider: props.preset.provider || 'openai',
      model: props.preset.model || '',
      apiBase: '',
      apiKey: '',
      temperature: 0.7,
      maxTokens: 2048,
      enabled: true,
      isDefault: false,
    }
    roles.value = []
  }
}

async function onSave() {
  saving.value = true
  try {
    const payload = { ...form.value, roles: JSON.stringify(roles.value) }
    if (props.modelId) {
      const { data } = await updateModel(props.modelId, payload)
      if (data?.code === 0) { Message.success('已保存'); emit('saved', props.modelId) } else { Message.error(data?.message || '保存失败') }
    } else {
      const { data } = await createModel(payload)
      if (data?.code === 0) { Message.success('已创建并保存'); emit('saved', data?.data?.id) } else { Message.error(data?.message || '保存失败') }
    }
  } finally { saving.value = false }
}

async function onTest() {
  testing.value = true
  try {
    const payload = { ...form.value, roles: JSON.stringify(roles.value) }
    const { data } = await testModel(payload)
    if (data?.code === 0) Message.success('测试成功')
    else Message.error(data?.message || '测试失败')
  } catch (err) {
    const msg = err?.response?.data?.message || err?.message || '测试失败'
    Message.error(msg)
  } finally { testing.value = false }
}

onMounted(load)
</script>


