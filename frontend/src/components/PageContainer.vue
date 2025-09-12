<template>
  <div class="page-container">
    <a-page-header :title="title" :subtitle="subtitle">
      <template #extra>
        <slot name="extra" />
      </template>
    </a-page-header>
    <a-card :bordered="false" style="margin-top:12px">
      <slot />
    </a-card>
  </div>
</template>
<script setup>
defineProps({ title: String, subtitle: String })
</script>
<style scoped>
.page-container { display:block; }
</style>
