<template>
  <a-layout style="height:100%">
    <a-layout-header class="app-header">
      <div class="logo">
        <img src="@/assets/logo.png" alt="AILAP" :class="['logo-img', { 'logo-dark': isDark }]" />
      </div>
      <div class="header-content">
        <a-breadcrumb class="crumb">
          <a-breadcrumb-item v-for="(m, i) in crumbs" :key="i">{{ m.meta?.locale ? $t(m.meta.locale) : m.path }}</a-breadcrumb-item>
        </a-breadcrumb>
        <div v-if="pageSubtitle" class="page-info">
          <span class="page-subtitle">{{ $t(pageSubtitle) }}</span>
        </div>
      </div>
      <div class="header-actions">
        <a-space>
          <a-button type="text" class="theme-btn" @click="toggleLocale">
            {{ currentLocaleText }}
          </a-button>
          <a-button type="text" class="theme-btn" @click="toggleTheme">
            <icon-moon v-if="!isDark" />
            <icon-sun v-else />
          </a-button>
          <a-dropdown>
            <a-avatar style="cursor:pointer">A</a-avatar>
            <template #content>
              <a-doption @click="goProfile">个人中心</a-doption>
              <a-doption @click="onLogout">退出登录</a-doption>
            </template>
          </a-dropdown>
        </a-space>
      </div>
    </a-layout-header>
    <a-layout class="main-layout">
      <a-layout-sider
        class="app-sider"
        :collapsed="collapsed"
        :width="208"
        :collapsed-width="56"
        breakpoint="xl"
      >
        <app-menu />
      </a-layout-sider>
      <a-layout-content class="app-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { useUiStore } from '@/store/ui'
import AppMenu from '@/components/AppMenu.vue'
import { IconSun, IconMoon } from '@arco-design/web-vue/es/icon'

import { useI18n } from 'vue-i18n'
import enUS from '@arco-design/web-vue/es/locale/lang/en-us'
import zhCN from '@arco-design/web-vue/es/locale/lang/zh-cn'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const ui = useUiStore()
const { locale } = useI18n()

const collapsed = computed({ get:() => ui.siderCollapsed, set:(v)=> ui.setSiderCollapsed(v) })
const crumbs = computed(() => route.matched.filter(m => m.path !== '/'))
const isDark = computed(() => ui.isDark)

function toggleLocale() {
  const next = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
  locale.value = next
  localStorage.setItem('locale', next)
  // Also update Arco locale in App.vue via store/event if needed, but for now standard vue-i18n
}
const currentLocaleText = computed(() => locale.value === 'zh-CN' ? '中' : 'En')

const pageTitle = computed(() => route.meta?.locale || '')
const pageSubtitle = computed(() => route.meta?.localeSubtitle || '')
function toggleTheme() { ui.toggleTheme() }

function goProfile() { router.push('/profile') }
function onLogout() { auth.clear(); router.replace('/login') }
</script>

<style scoped>
.app-header { display:flex; align-items:center; gap:12px; height:56px; }
.main-layout { flex: 1; min-height: 0; }
.app-sider { overflow: hidden; }
.app-sider :deep(.arco-layout-sider-children) { height: 100%; }
.app-sider:not(:deep(.arco-layout-sider-collapsed)) { overflow-y: auto; }
.logo { 
  display: flex; 
  align-items: center; 
  gap: 8px;
  height: 100%;
  width: 200px;
  justify-content: center;
}
.logo-img { 
  height: 32px; 
  width: auto; 
  object-fit: contain; 
}
.logo-dark { filter: brightness(0) invert(1); }
.header-content { display: flex; align-items: center; gap: 16px; flex: 1; margin-left:8px; }
.page-info { display: flex; align-items: center; gap: 8px; }
.page-title { font-weight: 600; font-size: 16px; color: var(--color-text-1); }
.page-subtitle { font-size: 14px; color: var(--color-text-3); }
.header-actions { margin-left:auto; padding-right:12px; }
.app-content { padding:16px; overflow:auto; background: var(--color-bg-1); }
.theme-btn { padding: 0; line-height: 1; color: var(--color-text-1); }
.theme-btn :deep(svg) { font-size: 18px; }
</style>
