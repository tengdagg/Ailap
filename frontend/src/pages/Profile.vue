<template>
  <page-container>
    <a-grid :cols="24" :col-gap="16" :row-gap="16">
      <a-grid-item :span="8">
        <a-card :title="$t('profile.accountInfo')">
          <div class="profile-row"><span class="label">{{ $t('common.username') }}</span><span class="value">{{ userName }}</span></div>
          <div class="profile-row"><span class="label">{{ $t('profile.theme') }}</span>
            <a-space>
              <a-tag :color="isDark ? 'arcoblue' : 'gray'">{{ isDark ? $t('profile.dark') : $t('profile.light') }}</a-tag>
              <a-switch :model-value="isDark" @change="toggleTheme" />
            </a-space>
          </div>
        </a-card>
      </a-grid-item>
      <a-grid-item :span="16">
        <a-card :title="$t('profile.changePassword')">
          <a-form :model="form" :rules="rules" ref="formRef" layout="vertical" @submit.prevent>
            <a-form-item field="oldPassword" :label="$t('profile.currentPassword')">
              <a-input-password v-model="form.oldPassword" :placeholder="$t('profile.placeCurrentPassword')" allow-clear />
            </a-form-item>
            <a-form-item field="newPassword" :label="$t('profile.newPassword')">
              <a-input-password v-model="form.newPassword" :placeholder="$t('profile.placeNewPassword')" allow-clear />
            </a-form-item>
            <a-form-item field="confirmPassword" :label="$t('profile.confirmPassword')">
              <a-input-password v-model="form.confirmPassword" :placeholder="$t('profile.placeConfirmPassword')" allow-clear />
            </a-form-item>
            <a-space>
              <a-button type="primary" :loading="submitting" @click="onSubmit">{{ $t('profile.save') }}</a-button>
              <a-button @click="onReset">{{ $t('profile.reset') }}</a-button>
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
import { useI18n } from 'vue-i18n'
import PageContainer from '@/components/PageContainer.vue'
import { useUiStore } from '@/store/ui'
import { useAuthStore } from '@/store/auth'
import { profile, changePassword } from '@/api/auth'

const ui = useUiStore()
const auth = useAuthStore()
const { t } = useI18n()

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
const rules = computed(() => ({
  oldPassword: [{ required: true, message: t('profile.placeCurrentPassword') }],
  newPassword: [
    { required: true, message: t('profile.placeNewPassword') },
    { validator: (val, cb) => { if (val && val.length < 8) cb(t('profile.placeNewPassword')); else cb() } },
  ],
  confirmPassword: [
    { required: true, message: t('profile.placeConfirmPassword') },
    { validator: (val, cb) => { if (val !== form.value.newPassword) cb(t('profile.passwordMismatch')); else cb() } },
  ],
}))

async function onSubmit() {
  const err = await formRef.value?.validate()
  if (err) return
  submitting.value = true
  try {
    const { data } = await changePassword({ oldPassword: form.value.oldPassword, newPassword: form.value.newPassword })
    if (data?.code === 0) {
      Message.success(t('profile.updateSuccess'))
      form.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
    } else {
      Message.error(data?.message || t('profile.updateFail'))
    }
  } catch (e) {
    Message.error(e?.response?.data?.message || e?.message || t('profile.updateFail'))
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


