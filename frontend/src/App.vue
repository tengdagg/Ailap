<template>
  <a-config-provider>
    <router-view />
  </a-config-provider>
</template>

<script setup>
import { onMounted } from 'vue'
import { useUiStore } from '@/store/ui'

const ui = useUiStore()

onMounted(() => {
  ui.initTheme()
})
</script>

<style>
html, body, #app { height: 100%; }
</style>

