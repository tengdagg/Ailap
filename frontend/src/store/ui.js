import { defineStore } from 'pinia'

export const useUiStore = defineStore('ui', {
  state: () => ({
    theme: localStorage.getItem('theme') || 'light',
    siderCollapsed: false,
  }),
  getters: {
    isDark: (s) => s.theme === 'dark',
  },
  actions: {
    toggleTheme() {
      this.theme = this.theme === 'dark' ? 'light' : 'dark'
      localStorage.setItem('theme', this.theme)
      const isDark = this.theme === 'dark'
      document.body.classList.toggle('arco-theme-dark', isDark)
      document.documentElement.classList.toggle('arco-theme-dark', isDark)
      document.body.setAttribute('arco-theme', isDark ? 'dark' : 'light')
      document.documentElement.setAttribute('arco-theme', isDark ? 'dark' : 'light')
    },
    initTheme() {
      const isDark = this.theme === 'dark'
      document.body.classList.toggle('arco-theme-dark', isDark)
      document.documentElement.classList.toggle('arco-theme-dark', isDark)
      document.body.setAttribute('arco-theme', isDark ? 'dark' : 'light')
      document.documentElement.setAttribute('arco-theme', isDark ? 'dark' : 'light')
    },
    setSiderCollapsed(v) { this.siderCollapsed = v },
  },
})













