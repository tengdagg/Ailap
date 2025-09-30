import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ArcoVue from '@arco-design/web-vue'
import '@arco-design/web-vue/dist/arco.css'
import './styles/global.css'
import App from './App.vue'
import router from './router'

// initialize theme ASAP to avoid flash
try {
  const savedTheme = localStorage.getItem('theme') || 'light'
  const isDark = savedTheme === 'dark'
  document.body.classList.toggle('arco-theme-dark', isDark)
  document.documentElement.classList.toggle('arco-theme-dark', isDark)
  document.body.setAttribute('arco-theme', isDark ? 'dark' : 'light')
  document.documentElement.setAttribute('arco-theme', isDark ? 'dark' : 'light')
} catch (_) {}

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(ArcoVue)
app.mount('#app')


