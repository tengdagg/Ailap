<template>
  <a-layout style="height:100%">
    <a-layout-header class="app-header">
      <div class="logo">AILAP</div>
      <a-breadcrumb class="crumb">
        <a-breadcrumb-item v-for="(m, i) in crumbs" :key="i">{{ m.meta?.title || m.path }}</a-breadcrumb-item>
      </a-breadcrumb>
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
function toggleTheme() { ui.toggleTheme() }
function onCollapse(v) { ui.setSiderCollapsed(v) }

function goProfile() { router.push('/dashboard') }
function onLogout() { auth.clear(); router.replace('/login') }
</script>

<style scoped>
.app-header { display:flex; align-items:center; gap:12px; height:56px; }
.logo { font-weight:700; font-size:16px; letter-spacing:.5px; }
.crumb { margin-left:8px; flex: 1; }
.header-actions { margin-left:auto; }
.app-content { padding:16px; overflow:auto; background: var(--color-bg-1); }
</style>
