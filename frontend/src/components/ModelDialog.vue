<template>
  <a-modal v-model:visible="innerVisible" :title="$t('models.model')">
    <a-form :model="model" layout="vertical">
      <a-form-item field="name" :label="$t('common.name')">
        <a-input v-model="model.name" />
      </a-form-item>
      <a-form-item field="version" label="版本">
        <a-input v-model="model.version" />
      </a-form-item>
    </a-form>
    <template #footer>
      <a-button @click="innerVisible=false">{{ $t('common.cancel') }}</a-button>
      <a-button type="primary" @click="$emit('save', model)">{{ $t('common.save') }}</a-button>
    </template>
  </a-modal>
</template>
<script setup>
import { reactive, watch, computed } from 'vue'
const props = defineProps({ visible: Boolean, value: Object })
const emit = defineEmits(['update:visible', 'save'])
const innerVisible = computed({ get: () => props.visible, set: (v) => emit('update:visible', v) })
const model = reactive({ id: undefined, name: '', version: '' })
watch(() => props.value, (v) => { Object.assign(model, { id: undefined, name: '', version: '' }, v || {}) }, { immediate: true })
</script>



