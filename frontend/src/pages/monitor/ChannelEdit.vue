<template>
  <div class="channel-edit">
    <a-form :model="form" @submit="onSubmit" layout="vertical" class="form-content">
      <a-form-item field="name" label="渠道名称" required>
        <a-input v-model="form.name" placeholder="请输入渠道名称" />
      </a-form-item>
      
      <a-form-item field="type" label="类型" required>
        <a-radio-group v-model="form.type" type="button">
            <a-radio value="webhook">Webhook</a-radio>
            <a-radio value="email">邮件 (SMTP)</a-radio>
        </a-radio-group>
      </a-form-item>

      <!-- Webhook Config -->
      <template v-if="form.type === 'webhook'">
          <a-form-item field="config.url" label="Webhook URL" required>
              <a-input v-model="form.config.url" placeholder="https://example.com/webhook" />
          </a-form-item>
      </template>

      <!-- Email Config -->
      <template v-if="form.type === 'email'">
           <a-form-item field="config.smtp_host" label="SMTP Host" required>
              <a-input v-model="form.config.smtp_host" placeholder="smtp.example.com" />
          </a-form-item>
           <a-form-item field="config.smtp_port" label="SMTP Port" required>
              <a-input v-model="form.config.smtp_port" placeholder="587" />
          </a-form-item>
           <a-form-item field="config.username" label="用户名" required>
              <a-input v-model="form.config.username" />
          </a-form-item>
           <a-form-item field="config.password" label="密码" required>
              <a-input-password v-model="form.config.password" />
          </a-form-item>
           <a-form-item field="config.to" label="收件人" required help="多个收件人用逗号分隔">
              <a-input v-model="form.config.to" placeholder="admin@example.com" />
          </a-form-item>
      </template>

      <a-form-item>
        <a-button type="primary" html-type="submit">提交</a-button>
        <a-button @click="onTest" status="success" style="margin-left: 10px">测试</a-button>
        <a-button @click="$router.back()" style="margin-left: 10px">取消</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import request from '@/api/request'

const route = useRoute()
const router = useRouter()
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
            Message.success('测试消息发送成功')
        } else {
            Message.error('测试失败: ' + data.message)
        }
    } catch (e) {
        // e is usually the axios error object; the interceptor might have handled it or rejected it
        // If 404/500, interceptor prints console.error but we want to show message
        console.error(e)
        Message.error('请求失败，请检查后端服务是否启动')
    }
}

const onSubmit = async () => {
    if (!form.value.name) {
        return Message.warning('请输入渠道名称')
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
            Message.success('保存成功')
            router.push('/channels')
        } else {
            Message.error(data.message)
        }
    } catch (e) {
        console.error(e)
        // Message is already shown by interceptor if 401/403. For others:
        Message.error('保存失败')
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
