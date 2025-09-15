<template>
  <page-container>
    <div style="display: flex; justify-content: flex-end; margin-bottom: 16px;">
      <a-button type="primary" @click="openCreate">新建模型</a-button>
    </div>
    <a-table :data="rows" :loading="loading" row-key="id">
      <a-table-column title="名称" data-index="name" />
      <a-table-column title="版本" data-index="version" />
      <a-table-column title="操作" :width="200">
        <template #cell="{ record }">
          <a-space>
            <a-button size="mini" @click="edit(record)">编辑</a-button>
            <a-popconfirm content="确认删除？" @ok="remove(record)">
              <a-button size="mini" status="danger">删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </a-table-column>
    </a-table>
    <model-dialog v-model:visible="dialogVisible" :value="current" @save="save" />
  </page-container>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import PageContainer from '@/components/PageContainer.vue'
import { listModels, createModel, updateModel, deleteModel } from '@/api/models'
import ModelDialog from '@/components/ModelDialog.vue'

const loading = ref(false)
const rows = ref([])
const dialogVisible = ref(false)
const current = ref(null)

async function fetchList() {
  loading.value = true
  try {
    const { data } = await listModels()
    rows.value = data?.data?.items || []
  } finally {
    loading.value = false
  }
}

function openCreate() { current.value = { name: '', version: '' }; dialogVisible.value = true }
function edit(record) { current.value = { ...record }; dialogVisible.value = true }

async function save(model) {
  if (model.id) await updateModel(model.id, model)
  else await createModel(model)
  dialogVisible.value = false
  await fetchList()
}

async function remove(record) {
  await deleteModel(record.id)
  await fetchList()
}

onMounted(fetchList)
</script>

