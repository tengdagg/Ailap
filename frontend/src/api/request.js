import axios from 'axios'
import { useAuthStore } from '@/store/auth'
import { Message } from '@arco-design/web-vue'
import router from '@/router'

const request = axios.create({
  baseURL: '/api',
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  try {
    const auth = useAuthStore()
    if (auth?.token) {
      config.headers = config.headers || {}
      config.headers.Authorization = `Bearer ${auth.token}`
    }
  } catch (_) {}
  return config
})

request.interceptors.response.use(
  (resp) => resp,
  (error) => {
    const status = error?.response?.status
    const message = error?.response?.data?.message || error.message
    if (status === 401 || status === 403) {
      try {
        const auth = useAuthStore()
        auth.clear()
      } catch (_) {}
      Message.error('登录已过期，请重新登录')
      router.replace('/login')
    } else {
      console.error('API Error:', message)
    }
    return Promise.reject(error)
  },
)

export default request

