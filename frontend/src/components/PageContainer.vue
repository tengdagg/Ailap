<template>
  <div class="page-container">
    <a-card :bordered="false">
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



