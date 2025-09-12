<template>
  <page-container title="仪表盘" subtitle="欢迎使用 AILAP">
    <template #extra>
      <a-space>
        <a-button type="primary" @click="$router.push('/logs')">日志</a-button>
        <a-button @click="$router.push('/models')">模型</a-button>
        <a-button @click="$router.push('/datasources')">数据源</a-button>
      </a-space>
    </template>
    <a-grid :cols="24" :col-gap="12" :row-gap="12">
      <a-grid-item :span="8"><a-card title="系统状态">占位</a-card></a-grid-item>
      <a-grid-item :span="8"><a-card title="近期活动">占位</a-card></a-grid-item>
      <a-grid-item :span="8"><a-card title="快速提示">占位</a-card></a-grid-item>
    </a-grid>
  </page-container>
</template>

<script setup>
import PageContainer from '@/components/PageContainer.vue'
</script>

