<template>
  <a-modal v-model:visible="visible" title="模型">
    <a-form :model="model" layout="vertical">
      <a-form-item field="name" label="名称">
        <a-input v-model="model.name" />
      </a-form-item>
      <a-form-item field="version" label="版本">
        <a-input v-model="model.version" />
      </a-form-item>
    </a-form>
    <template #footer>
      <a-button @click="visible=false">取消</a-button>
      <a-button type="primary" @click="$emit('save', model)">保存</a-button>
    </template>
  </a-modal>
</template>
<script setup>
import { ref, reactive, watch } from 'vue'
const props = defineProps({ value: Object })
const visible = ref(false)
const model = reactive({ name: '', version: '' })
watch(() => props.value, (v) => { Object.assign(model, v || {}) })
</script>

