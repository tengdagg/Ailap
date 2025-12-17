<template>
  <page-container>
    <a-form :model="form" layout="vertical" style="max-width:900px; margin:0 auto;">
      <a-grid :cols="24" :col-gap="12">
        <a-grid-item :span="12">
          <a-form-item :label="$t('common.name')">
            <a-input v-model="form.name" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item :label="$t('models.provider')">
            <a-select v-model="form.provider" :disabled="true">
              <a-option value="openai">OpenAI</a-option>
              <a-option value="deepseek">Deepseek</a-option>
              <a-option value="qwen">Qwen</a-option>
            </a-select>
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item :label="$t('models.model')">
            <a-select v-model="form.model" :options="modelOptions" allow-search allow-create />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
            <a-form-item :label="$t('common.apiBase')">
            <a-input v-model="form.apiBase" :placeholder="apiBasePlaceholder" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="24">
          <a-form-item :label="$t('common.apiKey')">
            <a-input-password v-model="form.apiKey" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item :label="$t('common.temp')">
            <a-input-number v-model="form.temperature" :min="0" :max="2" :step="0.1" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item :label="$t('common.maxTokens')">
            <a-input-number v-model="form.maxTokens" :min="1" :max="32000" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item :label="$t('common.enabled')">
            <a-switch v-model="form.enabled" />
          </a-form-item>
        </a-grid-item>
        <a-grid-item :span="12">
          <a-form-item :label="$t('common.setDefault')">
            <a-switch v-model="form.isDefault" />
          </a-form-item>
        </a-grid-item>
      </a-grid>

      <a-divider>{{ $t('common.roleDef') }}</a-divider>
      <a-space direction="vertical" fill>
        <div v-for="(r, idx) in roles" :key="idx" style="border:1px solid var(--color-border-2); padding:12px; border-radius:8px;">
          <a-grid :cols="24" :col-gap="8">
            <a-grid-item :span="8"><a-input v-model="r.name" placeholder="角色名，如：运维助手" /></a-grid-item>
            <a-grid-item :span="16"><a-input v-model="r.description" placeholder="角色描述" /></a-grid-item>
            <a-grid-item :span="24" style="margin-top:8px"><a-textarea v-model="r.systemPrompt" :placeholder="$t('common.sysPrompt')" :auto-size="{minRows:2, maxRows:6}" /></a-grid-item>
          </a-grid>
          <div style="display:flex; justify-content:flex-end; margin-top:8px">
            <a-button size="mini" status="danger" @click="removeRole(idx)">{{ $t('common.delete') }}</a-button>
          </div>
        </div>
        <a-button type="outline" size="small" @click="addRole">+ {{ $t('common.addRole') }}</a-button>
      </a-space>

      <a-space style="margin-top:16px">
        <a-button @click="$emit('back')">{{ $t('common.return') }}</a-button>
        <a-button @click="onTest" :loading="testing">{{ $t('common.test') }}</a-button>
        <a-button type="primary" @click="onSave" :loading="saving">{{ $t('common.save') }}</a-button>
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
import { useI18n } from 'vue-i18n'

const props = defineProps({
  modelId: { type: [String, Number], required: false },
  preset: { type: Object, required: false, default: null },
})
const emit = defineEmits(['back', 'saved'])
const { t } = useI18n()

const form = ref({ name: '', provider: 'openai', model: '', apiBase: '', apiKey: '', temperature: 0.7, maxTokens: 2048, enabled: true, isDefault: false })
const roles = ref([])
const testing = ref(false)
const saving = ref(false)

const providerModels = {
  // OpenAI
  openai: [
    { label: 'GPT-4o', value: 'gpt-4o' },
    { label: 'GPT-4o Mini', value: 'gpt-4o-mini' },
    { label: 'O1', value: 'o1' },
    { label: 'GPT-4.5 Preview', value: 'gpt-4.5-preview' },
    { label: 'GPT-4 Turbo', value: 'gpt-4-turbo' },
    { label: 'GPT-3.5 Turbo', value: 'gpt-3.5-turbo' },
  ],
  // DeepSeek
  deepseek: [
    { label: 'DeepSeek-V3 (Chat)', value: 'deepseek-chat' },
    { label: 'DeepSeek-R1 (Reasoner)', value: 'deepseek-reasoner' },
  ],
  // Qwen
  qwen: [
    { label: 'Qwen-Max (旗舰)', value: 'qwen-max' },
    { label: 'Qwen-Plus (增强)', value: 'qwen-plus' },
    { label: 'Qwen-Turbo (快速)', value: 'qwen-turbo' },
    { label: 'Qwen-Long (长文本)', value: 'qwen-long' },
  ],
}

const providerApiBase = {
  openai: 'https://api.openai.com/v1',
  deepseek: 'https://api.deepseek.com',
  qwen: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
}

const modelOptions = computed(() => providerModels[form.value.provider] || [])
const apiBasePlaceholder = computed(() => providerApiBase[form.value.provider] || 'https://api.example.com/v1')

watch(() => form.value.provider, (p, oldP) => {
  if (!p) return
  // If switching providers manually
  if (p !== oldP && oldP) {
    // Check if the current model value belongs to the OLD provider's list.
    // If it does, it's a stale value from the previous selection => reset it.
    // If it does NOT, it implies it was either just loaded (custom/valid for new provider) or is a custom value we shouldn't clobber.
    const oldList = providerModels[oldP] || []
    const isStale = oldList.some(item => item.value === form.value.model)
    
    // Also reset if empty, to provide a nice default
    if (!form.value.model || isStale) {
      const list = providerModels[p] || []
      form.value.model = list[0]?.value || ''
    }
  } else {
     // Initial load or same provider: trust the form value unless empty
     if (!form.value.model) {
        const list = providerModels[p] || []
        form.value.model = list[0]?.value || ''
     }
  }
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


