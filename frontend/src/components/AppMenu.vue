<template>
  <div class="menu-container">
    <!-- Menu Items -->
    <div class="menu-items">
      <div
        v-for="item in menuItems"
        :key="item.key"
        :class="['menu-item', { 
          'active': selected === item.key,
          'collapsed': isCollapsed 
        }]"
        @click="onClick(item.key)"
      >
        <div class="menu-item-content">
          <div class="menu-icon">
            <component :is="item.icon" />
          </div>
          <transition name="fade">
            <span v-if="!isCollapsed" class="menu-label">{{ item.label }}</span>
          </transition>
        </div>
        <div v-if="selected === item.key" class="active-indicator"></div>
      </div>
    </div>

    <!-- Toggle Button -->
    <div class="menu-toggle" @click="toggleMenu">
      <icon-menu-fold v-if="!isCollapsed" />
      <icon-menu-unfold v-else />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  IconApps, 
  IconList, 
  IconExperiment, 
  IconStorage,
  IconRobot,
  IconNotification,
  IconMenuFold,
  IconMenuUnfold
} from '@arco-design/web-vue/es/icon'
import { useUiStore } from '@/store/ui'

const route = useRoute()
const router = useRouter()
const ui = useUiStore()

const selected = computed(() => '/' + (route.path.split('/')[1] || 'dashboard'))
const isCollapsed = computed(() => ui.siderCollapsed)

const menuItems = [
  { key: '/dashboard', label: '仪表盘', icon: IconApps },
  { key: '/logs', label: '日志分析', icon: IconList },
  { key: '/models', label: '模型', icon: IconExperiment },
  { key: '/datasources', label: '数据源', icon: IconStorage },
  { key: '/monitors', label: '智能监控任务', icon: IconRobot },
  { key: '/channels', label: '通知渠道', icon: IconNotification }
]

function onClick(key) {
  router.push(key)
}

function toggleMenu() {
  ui.setSiderCollapsed(!isCollapsed.value)
}
</script>

<style scoped>
.menu-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-2);
  position: relative;
}

/* Toggle Button */
.menu-toggle {
  padding: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  color: var(--color-text-3);
  font-size: 18px;
  border-top: 1px solid var(--color-border-2);
}

/* Menu Items */
.menu-items {
  flex: 1;
  padding: 12px 8px;
  overflow-y: auto;
  overflow-x: hidden;
}

.menu-item {
  position: relative;
  margin-bottom: 8px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.menu-item-content {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  gap: 12px;
  position: relative;
  z-index: 1;
}

.menu-item.collapsed .menu-item-content {
  justify-content: center;
  padding: 12px 8px;
}

.menu-icon {
  font-size: 20px;
  color: var(--color-text-2);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
}

.menu-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-1);
  white-space: nowrap;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Active State */
.menu-item.active {
  background: linear-gradient(
    135deg,
    rgba(var(--primary-6), 0.1) 0%,
    rgba(var(--primary-6), 0.05) 100%
  );
  box-shadow: 0 2px 8px rgba(var(--primary-6), 0.15);
}

.menu-item.active .menu-icon {
  color: rgb(var(--primary-6));
  transform: scale(1.15);
}

.menu-item.active .menu-label {
  color: rgb(var(--primary-6));
  font-weight: 600;
}

/* Active Indicator */
.active-indicator {
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 60%;
  background: linear-gradient(
    180deg,
    rgb(var(--primary-6)) 0%,
    rgb(var(--primary-5)) 100%
  );
  border-radius: 0 4px 4px 0;
  animation: slideIn 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Collapsed State */
.menu-item.collapsed {
  margin-bottom: 12px;
}

/* Animations */
@keyframes slideIn {
  from {
    height: 0;
    opacity: 0;
  }
  to {
    height: 60%;
    opacity: 1;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Scrollbar Styling */
.menu-items::-webkit-scrollbar {
  width: 4px;
}

.menu-items::-webkit-scrollbar-track {
  background: transparent;
}

.menu-items::-webkit-scrollbar-thumb {
  background: var(--color-fill-3);
  border-radius: 2px;
}

.menu-items::-webkit-scrollbar-thumb:hover {
  background: var(--color-fill-4);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .menu-item-content {
    padding: 10px 12px;
  }
  
  .menu-icon {
    font-size: 18px;
  }
  
  .menu-label {
    font-size: 13px;
  }
}
</style>












