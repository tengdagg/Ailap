<template>
  <a-card class="ds-type-card" hoverable @click="$emit('click')">
    <div class="icon" v-html="icon" />
    <div class="meta">
      <div class="title">{{ title }}</div>
      <div class="desc">{{ description }}</div>
    </div>
  </a-card>
</template>
<script setup>
defineProps({ title: String, description: String, icon: String })
</script>
<style scoped>
.ds-type-card { display:flex; align-items:center; gap:12px; cursor:pointer; background: var(--color-fill-2); }
.icon { width:40px; height:40px; display:flex; align-items:center; justify-content:center; }
.title { font-weight:600; }
.desc { opacity:.7; font-size:12px; }
</style>








