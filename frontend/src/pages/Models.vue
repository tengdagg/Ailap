<template>
  <page-container>
    <div v-if="!editingId && !creating && !preset" style="display:flex; justify-content:flex-start; align-items:center; margin-bottom: 16px;">
      <a-button type="primary" @click="openNew">
        <template #icon><icon-plus /></template>
        添加模型
      </a-button>
    </div>

    <a-grid v-if="!editingId && !creating && !preset" :cols="24" :col-gap="18" :row-gap="16">
      <a-grid-item v-for="m in models" :key="m.id" :span="6">
        <a-card hoverable style="position:relative">
          <template #title>
            <div style="display:flex;align-items:center;gap:8px">
              <img :src="getLogo(m.provider)" alt="logo" style="width:24px;height:24px;object-fit:contain" />
              <span>{{ m.name }}</span>
              <a-tag v-if="m.isDefault" color="arcoblue" size="small">默认</a-tag>
            </div>
          </template>
          <div style="color:var(--color-text-3); min-height:52px">
            <div>供应商：{{ m.provider }}</div>
            <div>模型：{{ m.model }}</div>
          </div>
          <div style="display:flex;justify-content:space-between;align-items:center;margin-top:8px">
            <a-space>
              <a-switch size="small" :model-value="!!m.enabled" @change="(v)=>onToggleEnabled(m, v)">启用</a-switch>
              <a-button size="mini" type="outline" @click="setDefault(m)" :disabled="m.isDefault">设为默认</a-button>
            </a-space>
            <a-space>
              <a-button size="mini" @click="startEdit(m)">编辑</a-button>
              <a-popconfirm content="确认删除？" @ok="remove(m)">
                <a-button size="mini" status="danger">删除</a-button>
              </a-popconfirm>
            </a-space>
          </div>
        </a-card>
      </a-grid-item>
    </a-grid>

    <div v-if="creating && !editingId && !preset">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px;">
        <a-button type="text" @click="cancelCreate">返回列表</a-button>
        <span style="font-weight:600">选择模型类型</span>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import PageContainer from '@/components/PageContainer.vue'
import { listModels, deleteModel, toggleModelEnabled, setModelDefault } from '@/api/models'
import ModelConfigInline from '@/pages/model/ModelConfig.vue'

const router = useRouter()
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
  if (data?.code === 0) { Message.success('已删除'); fetchList() }
  else { Message.error(data?.message || '删除失败') }
}

async function onToggleEnabled(m, val) {
  const { data } = await toggleModelEnabled(m.id, !!val)
  if (data?.code === 0) { m.enabled = !!val; Message.success(val ? '已启用' : '已停用') }
  else { Message.error(data?.message || '操作失败') }
}

async function setDefault(m) {
  const { data } = await setModelDefault(m.id)
  if (data?.code === 0) { Message.success('已设为默认'); fetchList() }
  else { Message.error(data?.message || '操作失败') }
}

onMounted(fetchList)
</script>

