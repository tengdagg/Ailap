<template>
  <page-container>
    <a-grid :cols="24" :col-gap="16" :row-gap="16">
      <a-grid-item :span="8">
        <a-card title="账户信息">
          <div class="profile-row"><span class="label">用户名</span><span class="value">{{ userName }}</span></div>
          <div class="profile-row"><span class="label">主题</span>
            <a-space>
              <a-tag :color="isDark ? 'arcoblue' : 'gray'">{{ isDark ? '深色' : '浅色' }}</a-tag>
              <a-switch :model-value="isDark" @change="toggleTheme" />
            </a-space>
          </div>
        </a-card>
      </a-grid-item>
      <a-grid-item :span="16">
        <a-card title="修改密码">
          <a-form :model="form" :rules="rules" ref="formRef" layout="vertical" @submit.prevent>
            <a-form-item field="oldPassword" label="当前密码">
              <a-input-password v-model="form.oldPassword" placeholder="请输入当前密码" allow-clear />
            </a-form-item>
            <a-form-item field="newPassword" label="新密码">
              <a-input-password v-model="form.newPassword" placeholder="至少 8 位，包含字母与数字" allow-clear />
            </a-form-item>
            <a-form-item field="confirmPassword" label="确认新密码">
              <a-input-password v-model="form.confirmPassword" placeholder="再次输入新密码" allow-clear />
            </a-form-item>
            <a-space>
              <a-button type="primary" :loading="submitting" @click="onSubmit">保存</a-button>
              <a-button @click="onReset">重置</a-button>
            </a-space>
          </a-form>
        </a-card>
      </a-grid-item>
    </a-grid>
  </page-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import PageContainer from '@/components/PageContainer.vue'
import { useUiStore } from '@/store/ui'
import { useAuthStore } from '@/store/auth'
import { profile, changePassword } from '@/api/auth'

const ui = useUiStore()
const auth = useAuthStore()

const isDark = computed(() => ui.isDark)
function toggleTheme() { ui.toggleTheme() }

const userName = ref('')
async function loadProfile() {
  try {
    const { data } = await profile()
    if (data?.code === 0) {
      userName.value = data?.data?.name || '用户'
      auth.setUser({ name: userName.value })
    }
  } catch (_) {}
}

const formRef = ref(null)
const submitting = ref(false)
const form = ref({ oldPassword: '', newPassword: '', confirmPassword: '' })
const rules = {
  oldPassword: [{ required: true, message: '请输入当前密码' }],
  newPassword: [
    { required: true, message: '请输入新密码' },
    { validator: (val, cb) => { if (val && val.length < 8) cb('至少 8 位'); else cb() } },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码' },
    { validator: (val, cb) => { if (val !== form.value.newPassword) cb('两次输入不一致'); else cb() } },
  ],
}

async function onSubmit() {
  const err = await formRef.value?.validate()
  if (err) return
  submitting.value = true
  try {
    const { data } = await changePassword({ oldPassword: form.value.oldPassword, newPassword: form.value.newPassword })
    if (data?.code === 0) {
      Message.success('密码已更新')
      form.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
    } else {
      Message.error(data?.message || '更新失败')
    }
  } catch (e) {
    Message.error(e?.response?.data?.message || e?.message || '更新失败')
  } finally {
    submitting.value = false
  }
}

function onReset() { form.value = { oldPassword: '', newPassword: '', confirmPassword: '' } }

onMounted(loadProfile)
</script>

<style scoped>
.profile-row { display:flex; align-items:center; justify-content:space-between; padding:8px 0; }
.label { color: var(--color-text-3); }
.value { font-weight: 600; }
</style>


