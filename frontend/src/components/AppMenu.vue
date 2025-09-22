<template>
  <a-menu :selected-keys="[selected]" @menu-item-click="onClick" breakpoint="xl">
    <a-menu-item key="/dashboard"><icon-apps /> 仪表盘</a-menu-item>
    <a-menu-item key="/logs"><icon-list /> 日志</a-menu-item>
    <a-menu-item key="/models"><icon-experiment /> 模型</a-menu-item>
    <a-menu-item key="/datasources"><icon-storage /> 数据源</a-menu-item>
  </a-menu>
</template>
<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { IconApps, IconList, IconExperiment, IconStorage } from '@arco-design/web-vue/es/icon'

const route = useRoute()
const router = useRouter()
const selected = computed(() => '/' + (route.path.split('/')[1] || 'dashboard'))

function onClick(key) { router.push(key) }
</script>










