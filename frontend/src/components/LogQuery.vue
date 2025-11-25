<template>
  <a-space>
    <a-input v-model="query" placeholder="输入查询语句" style="width:340px" />
    <a-button type="primary" @click="$emit('search', query)">查询</a-button>
  </a-space>
</template>
<script setup>
import { ref } from 'vue'
const query = ref('')
</script>

















