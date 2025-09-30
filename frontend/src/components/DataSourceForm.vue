<template>
  <a-form :model="form" layout="vertical">
    <a-form-item field="name" label="名称">
      <a-input v-model="form.name" />
    </a-form-item>
    <a-form-item field="type" label="类型">
      <a-select v-model="form.type" :options="types" />
    </a-form-item>
    <a-form-item field="endpoint" label="地址">
      <a-input v-model="form.endpoint" />
    </a-form-item>
  </a-form>
</template>
<script setup>
import { reactive } from 'vue'
const props = defineProps({ value: Object })
const form = reactive({ name: '', type: 'loki', endpoint: '' })
const types = [
  { label: 'Loki', value: 'loki' },
  { label: 'ELK', value: 'elk' },
  { label: 'Database', value: 'db' },
]
</script>














