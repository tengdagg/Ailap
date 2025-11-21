<template>
  <div class="app-container" style="max-width:360px;margin:10% auto;">
    <a-card title="登录 AILAP">
      <a-form :model="form" @submit-prevent="onSubmit" layout="vertical">
        <a-form-item field="username" label="用户名">
          <a-input v-model="form.username" placeholder="admin" />
        </a-form-item>
        <a-form-item field="password" label="密码">
          <a-input-password v-model="form.password" placeholder="••••••" />
        </a-form-item>
        <a-space>
          <a-button type="primary" @click="onSubmit" :loading="loading">登录</a-button>
        </a-space>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { login } from '@/api/auth'

const router = useRouter()
const auth = useAuthStore()

const form = reactive({ username: '', password: '' })
const loading = ref(false)

async function onSubmit() {
  if (loading.value) return
  loading.value = true
  try {
    const { data } = await login(form)
    const token = data?.data?.token || data?.token
    if (token) auth.setToken(token)
    router.replace('/dashboard')
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}
</script>
















