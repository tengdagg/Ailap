<template>
  <page-container title="选择数据源类型" subtitle="请选择需要添加的数据源">
    <a-grid :cols="24" :col-gap="16" :row-gap="16">
      <a-grid-item :span="12">
        <data-source-type-card :icon="icons.loki" title="Loki" description="Grafana Loki 日志系统" @click="go('loki')" />
      </a-grid-item>
      <a-grid-item :span="12">
        <data-source-type-card :icon="icons.es" title="Elasticsearch" description="Elastic 搜索与日志" @click="go('elasticsearch')" />
      </a-grid-item>
    </a-grid>
  </page-container>
</template>
<script setup>
import PageContainer from '@/components/PageContainer.vue'
import DataSourceTypeCard from '@/components/DataSourceTypeCard.vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const icons = {
  loki: '<svg width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg"><rect width="32" height="32" fill="none"/><path d="M8 6l2 10 2-1-2-10-2 1z" fill="#F59E0B"/><path d="M12 5l2 10 2-1-2-10-2 1z" fill="#F97316"/><path d="M16 4l2 10 2-1-2-10-2 1z" fill="#EA580C"/><rect x="6" y="20" width="12" height="2" fill="#FBBF24"/></svg>',
  es: '<svg width="36" height="36" viewBox="0 0 36 36" xmlns="http://www.w3.org/2000/svg"><circle cx="18" cy="18" r="16" fill="#1F2937"/><path d="M9 18a9 9 0 0 0 9 9 9 9 0 0 0 6.9-3.2l-7.1-3.8H9.6A8.96 8.96 0 0 1 9 18z" fill="#34D399"/><path d="M27 18a9 9 0 0 0-9-9 9 9 0 0 0-6.9 3.2l7.1 3.8h8.2c.4.98.6 2.05.6 3z" fill="#F59E0B"/><rect x="10" y="17" width="16" height="2" rx="1" fill="#38BDF8"/></svg>'
}

function go(type) { router.push(`/datasources/new/${type}`) }
</script>


