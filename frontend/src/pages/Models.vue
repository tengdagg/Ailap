<template>
  <page-container>
    <div v-if="!editingId && !creating && !preset" style="display:flex; justify-content:flex-start; align-items:center; margin-bottom: 16px;">
      <a-button type="primary" @click="openNew">
        <template #icon><icon-plus /></template>
        {{ $t('models.add') }}
      </a-button>
    </div>

    <a-grid v-if="!editingId && !creating && !preset" :cols="24" :col-gap="18" :row-gap="16">
      <a-grid-item v-for="m in models" :key="m.id" :span="6">
        <a-card hoverable class="model-card" :bordered="true">
          <div class="model-card-header">
            <div class="model-info">
              <img :src="getLogo(m.provider)" alt="logo" class="model-logo" />
              <span class="model-name">{{ m.name }}</span>
            </div>
            <a-tag v-if="m.isDefault" color="arcoblue" size="small">{{ $t('common.default') }}</a-tag>
          </div>
          
          <div class="model-card-body">
            <div class="info-row">
              <span class="label">{{ $t('models.provider') }}：</span>
              <span class="value">{{ m.provider }}</span>
            </div>
            <div class="info-row">
              <span class="label">{{ $t('models.model') }}：</span>
              <span class="value">{{ m.model }}</span>
            </div>
          </div>

          <div class="model-card-footer">
            <div class="status-switch">
              <a-switch size="small" :model-value="!!m.enabled" @change="(v)=>onToggleEnabled(m, v)" />
              <span :class="['status-text', { enabled: m.enabled }]">{{ m.enabled ? $t('common.enabled') : $t('common.disabled') }}</span>
            </div>
            
            <div class="actions">
               <a-button size="mini" type="text" @click="setDefault(m)" :disabled="m.isDefault" v-if="!m.isDefault">{{ $t('common.setDefault') }}</a-button>
               <a-button size="mini" @click="startEdit(m)">{{ $t('common.edit') }}</a-button>
               <a-popconfirm :content="$t('common.deleteConfirm')" @ok="remove(m)">
                 <a-button size="mini" status="danger">{{ $t('common.delete') }}</a-button>
               </a-popconfirm>
            </div>
          </div>
        </a-card>
      </a-grid-item>
    </a-grid>

    <div v-if="creating && !editingId && !preset">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px;">
        <a-button type="text" @click="cancelCreate">{{ $t('models.backToList') }}</a-button>
        <span style="font-weight:600">{{ $t('models.selectType') }}</span>
        <div />
      </div>
      <a-grid :cols="24" :col-gap="12" :row-gap="12">
        <a-grid-item v-for="item in presets" :key="item.provider" :span="8">
          <a-card hoverable @click="choosePreset(item)" style="cursor:pointer; padding:10px;">
            <div style="display:flex;align-items:center;gap:8px">
              <img :src="getLogo(item.provider)" alt="logo" style="width:45px;height:45px;object-fit:contain" />
              <div>
                <div style="font-weight:600">{{ item.name }}</div>
                <div style="color:var(--color-text-3);font-size:12px">{{ item.desc }}</div>
              </div>
            </div>
          </a-card>
        </a-grid-item>
      </a-grid>
    </div>

    <div v-else-if="editingId || preset">
      <model-config-inline :model-id="editingId" :preset="preset" @back="stopEditOrCancel" @saved="onSaved" />
    </div>
  </page-container>
  
</template>
<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { useI18n } from 'vue-i18n'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import PageContainer from '@/components/PageContainer.vue'
import { listModels, deleteModel, toggleModelEnabled, setModelDefault } from '@/api/models'
import ModelConfigInline from '@/pages/model/ModelConfig.vue'

const router = useRouter()
const { t } = useI18n()
const models = ref([])
const loading = ref(false)
const creating = ref(false)
const editingId = ref('')
const preset = ref(null)

const presets = [
  { provider: 'deepseek', name: 'Deepseek', model: 'deepseek-chat', desc: 'Deepseek 对话模型' },
  { provider: 'openai', name: 'OpenAI', model: 'gpt-4o-mini', desc: 'OpenAI GPT 系列模型' },
  { provider: 'qwen', name: 'Qwen', model: 'qwen2.5' , desc: '通义千问系列模型'},
]

function getLogo(provider) {
  try {
    return new URL(`../assets/${provider}.png`, import.meta.url).href
  } catch (_) {
    return new URL(`../assets/logo.png`, import.meta.url).href
  }
}

async function fetchList() {
  loading.value = true
  try {
    const { data } = await listModels()
    models.value = data?.data?.items || []
  } finally { loading.value = false }
}

function openNew() { creating.value = true; preset.value = null }
function cancelCreate() { creating.value = false; preset.value = null }

function choosePreset(p) { preset.value = p; creating.value = false }

function startEdit(m) { editingId.value = String(m.id); preset.value = null }
function stopEditOrCancel() { editingId.value = ''; preset.value = null; creating.value = false }
async function onSaved(id) { await fetchList(); editingId.value = String(id); preset.value = null }

async function remove(m) {
  const { data } = await deleteModel(m.id)
  if (data?.code === 0) { Message.success(t('common.deleteSuccess')); fetchList() }
  else { Message.error(data?.message || t('common.deleteFail')) }
}

async function onToggleEnabled(m, val) {
  const { data } = await toggleModelEnabled(m.id, !!val)
  if (data?.code === 0) { m.enabled = !!val; Message.success(val ? t('common.enabled') : t('common.disabled')) }
  else { Message.error(data?.message || t('common.error')) }
}

async function setDefault(m) {
  const { data } = await setModelDefault(m.id)
  if (data?.code === 0) { Message.success(t('common.success')); fetchList() }
  else { Message.error(data?.message || t('common.error')) }
}

onMounted(fetchList)
</script>

<style scoped>
.model-card {
  border-radius: 8px;
  transition: all 0.3s;
}
.model-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(0,0,0,0.1);
}
.model-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.model-info {
  display: flex;
  align-items: center;
  gap: 10px;
}
.model-logo {
  width: 32px;
  height: 32px;
  object-fit: contain;
}
.model-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-1);
}
.model-card-body {
  margin-bottom: 20px;
  color: var(--color-text-2);
}
.info-row {
  display: flex;
  align-items: center;
  margin-bottom: 4px;
  font-size: 13px;
}
.label {
  color: var(--color-text-3);
  width: 60px;
}
.model-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid var(--color-border-1);
  padding-top: 12px;
}
.status-switch {
  display: flex;
  align-items: center;
  gap: 8px;
}
.status-text {
  font-size: 12px;
  color: var(--color-text-3);
}
.status-text.enabled {
  color: rgb(var(--green-6));
}
.actions {
  display: flex;
  gap: 8px;
}
</style>

