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
      document.body.classList.toggle('arco-theme-dark', this.theme === 'dark')
    },
    initTheme() {
      document.body.classList.toggle('arco-theme-dark', this.theme === 'dark')
    },
    setSiderCollapsed(v) { this.siderCollapsed = v },
  },
})






