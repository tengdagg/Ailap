<template>
  <div class="channel-edit">
    <a-form :model="form" @submit="onSubmit" layout="vertical" class="form-content">
      <a-form-item field="name" :label="$t('channel.channelName')" required>
        <a-input v-model="form.name" :placeholder="$t('channel.placeName')" />
      </a-form-item>
      
      <a-form-item field="type" :label="$t('common.type')" required>
        <a-radio-group v-model="form.type" type="button">
            <a-radio value="webhook">Webhook</a-radio>
            <a-radio value="email">{{ $t('channel.emailType') }}</a-radio>
        </a-radio-group>
      </a-form-item>

      <!-- Webhook Config -->
      <template v-if="form.type === 'webhook'">
          <a-form-item field="config.url" :label="$t('channel.webhookUrl')" required>
              <a-input v-model="form.config.url" placeholder="https://example.com/webhook" />
          </a-form-item>
      </template>

      <!-- Email Config -->
      <template v-if="form.type === 'email'">
           <a-form-item field="config.smtp_host" :label="$t('channel.smtpHost')" required>
              <a-input v-model="form.config.smtp_host" placeholder="smtp.example.com" />
          </a-form-item>
           <a-form-item field="config.smtp_port" :label="$t('channel.smtpPort')" required>
              <a-input v-model="form.config.smtp_port" placeholder="587" />
          </a-form-item>
           <a-form-item field="config.username" :label="$t('channel.username')" required>
              <a-input v-model="form.config.username" />
          </a-form-item>
           <a-form-item field="config.password" :label="$t('channel.password')" required>
              <a-input-password v-model="form.config.password" />
          </a-form-item>
           <a-form-item field="config.to" :label="$t('channel.recipients')" required :help="$t('channel.helpRecipients')">
              <a-input v-model="form.config.to" placeholder="admin@example.com" />
          </a-form-item>
      </template>

      <a-form-item>
        <a-button type="primary" html-type="submit">{{ $t('common.submit') }}</a-button>
        <a-button @click="onTest" status="success" style="margin-left: 10px">{{ $t('common.test') }}</a-button>
        <a-button @click="$router.back()" style="margin-left: 10px">{{ $t('common.cancel') }}</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { useI18n } from 'vue-i18n'
import request from '@/api/request'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const id = route.params.id

const isEdit = computed(() => id && id !== 'new')

const form = ref({
  name: '',
  type: 'webhook',
  config: {
    url: '',
    smtp_host: '',
    smtp_port: '587',
    username: '',
    password: '',
    to: ''
  }
})

const loadData = async () => {
    if (!isEdit.value) return
    try {
        const { data } = await request.get(`/channels/${id}`)
        if (data.code === 0) {
            const item = data.data.item
            try {
                const parsed = JSON.parse(item.config || '{}')
                // Merge parsed config into form.config ensure all fields exist
                form.value = {
                    ...item,
                    config: { ...form.value.config, ...parsed }
                }
            } catch (e) {
                console.error(e)
                form.value = { ...item, config: { ...form.value.config } }
            }
        }
    } catch (e) { console.error(e) }
}

const onTest = async () => {
    // Pack config for test
    const payload = {
        name: form.value.name,
        type: form.value.type,
        config: JSON.stringify(form.value.config)
    }
    try {
        const { data } = await request.post('/channels/test', payload)
        if (data.code === 0) {
            Message.success(t('common.testSuccess'))
        } else {
            Message.error(t('common.testFail') + ': ' + data.message)
        }
    } catch (e) {
        // e is usually the axios error object; the interceptor might have handled it or rejected it
        // If 404/500, interceptor prints console.error but we want to show message
        console.error(e)
        Message.error(t('common.testFail'))
    }
}

const onSubmit = async () => {
    if (!form.value.name) {
        return Message.warning(t('channel.placeName'))
    }
    
    // Pack config for submission
    const payload = {
        ...form.value,
        config: JSON.stringify(form.value.config)
    }
    
    try {
        let res
        if (isEdit.value) {
            res = await request.put(`/channels/${id}`, payload)
        } else {
            res = await request.post('/channels', payload)
        }
        const { data } = res
        
        if (data.code === 0) {
            Message.success(t('common.saveSuccess'))
            router.push('/channels')
        } else {
            Message.error(data.message)
        }
    } catch (e) {
        console.error(e)
        // Message is already shown by interceptor if 401/403. For others:
        Message.error(t('common.saveFail'))
    }
}

onMounted(loadData)
</script>

<style scoped>
.form-content {
  max-width: 600px;
  margin: 20px auto;
}
</style>
