<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <img src="@/assets/logo.png" alt="Logo" class="login-logo" />
        <h2 class="login-title">AILAP 智能分析平台</h2>
        <p class="login-subtitle">欢迎回来，请登录您的账户</p>
      </div>
      
      <a-form :model="form" @submit-prevent="onSubmit" layout="vertical" class="login-form">
        <a-form-item field="username" hide-label>
          <a-input 
            v-model="form.username" 
            placeholder="用户名" 
            size="large"
            allow-clear
          >
            <template #prefix>
              <icon-user />
            </template>
          </a-input>
        </a-form-item>
        
        <a-form-item field="password" hide-label>
          <a-input-password 
            v-model="form.password" 
            placeholder="密码" 
            size="large"
            allow-clear
          >
            <template #prefix>
              <icon-lock />
            </template>
          </a-input-password>
        </a-form-item>
        
        <div class="form-actions">
          <a-button 
            type="primary" 
            html-type="submit" 
            long 
            size="large" 
            :loading="loading"
            class="login-button"
            @click="onSubmit"
          >
            登录
          </a-button>
        </div>
      </a-form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { login } from '@/api/auth'
import { IconUser, IconLock } from '@arco-design/web-vue/es/icon'

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

<style scoped>
.login-container {
  min-height: 100vh;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #ffffff;
  position: relative;
}

.login-card {
  width: 100%;
  max-width: 400px;
  padding: 40px;
  background: #ffffff;
  border-radius: 8px;
  /* Subtle shadow for separation */
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  z-index: 1;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  height: 48px;
  width: auto;
  margin-bottom: 16px;
  object-fit: contain;
}

.login-title {
  color: #1d2129;
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 8px;
}

.login-subtitle {
  color: #86909c;
  font-size: 14px;
  margin: 0;
}

/* Override Arco Input styles for better fit if needed, 
   but default Arco styles work well on white. 
   Just ensuring full width and spacing. */
.login-form :deep(.arco-form-item) {
  margin-bottom: 20px;
}

.login-button {
  background-color: #165dff;
  border: none;
  border-radius: 4px;
  height: 40px;
  font-size: 15px;
  font-weight: 500;
  transition: all 0.2s;
}

.login-button:hover {
  background-color: #4080ff;
}

.login-button:active {
  background-color: #0e42d2;
}

/* Responsive adjustments */
@media (max-width: 480px) {
  .login-card {
    box-shadow: none;
    padding: 20px;
  }
}
</style>

















