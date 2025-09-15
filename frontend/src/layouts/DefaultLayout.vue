<template>
  <a-layout style="height:100%">
    <a-layout-header class="app-header">
      <div class="logo">
        <img src="@/assets/logo.png" alt="AILAP" class="logo-img" />
      </div>
      <div class="header-content">
        <a-breadcrumb class="crumb">
          <a-breadcrumb-item v-for="(m, i) in crumbs" :key="i">{{ m.meta?.title || m.path }}</a-breadcrumb-item>
        </a-breadcrumb>
        <div v-if="pageSubtitle" class="page-info">
          <span class="page-subtitle">{{ pageSubtitle }}</span>
        </div>
      </div>
      <div class="header-actions">
        <a-space>
          <span>深色</span>
          <a-switch :model-value="isDark" @change="toggleTheme" />
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
    <a-layout>
      <a-layout-sider :collapsed="collapsed" collapsible @collapse="onCollapse">
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

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const ui = useUiStore()

const collapsed = computed({ get:() => ui.siderCollapsed, set:(v)=> ui.setSiderCollapsed(v) })
const crumbs = computed(() => route.matched.filter(m => m.path !== '/'))
const isDark = computed(() => ui.isDark)

// 获取当前页面的标题和副标题
const pageTitle = computed(() => route.meta?.pageTitle || '')
const pageSubtitle = computed(() => route.meta?.pageSubtitle || '')
function toggleTheme() { ui.toggleTheme() }
function onCollapse(v) { ui.setSiderCollapsed(v) }

function goProfile() { router.push('/dashboard') }
function onLogout() { auth.clear(); router.replace('/login') }
</script>

<style scoped>
.app-header { display:flex; align-items:center; gap:12px; height:56px; }
.logo { 
  display: flex; 
  align-items: center; 
  gap: 8px;
  height: 100%;
  width: 200px;
  justify-content: space-around;
}
.logo-img { 
  height: 32px; 
  width: auto; 
  object-fit: contain; 
}
.logo-text {
  font-weight: 700; 
  font-size: 18px; 
  letter-spacing: 0.5px;
  color: var(--color-text-1);
  line-height: 1;
}
.header-content { display: flex; align-items: center; gap: 16px; flex: 1; margin-left:8px; }
.page-info { display: flex; align-items: center; gap: 8px; }
.page-title { font-weight: 600; font-size: 16px; color: var(--color-text-1); }
.page-subtitle { font-size: 14px; color: var(--color-text-3); }
.header-actions { margin-left:auto; }
.app-content { padding:16px; overflow:auto; background: var(--color-bg-1); }
</style>
