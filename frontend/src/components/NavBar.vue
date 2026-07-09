<template>
  <div class="nav-bar">
    <div class="nav-inner">
      <div class="nav-left">
        <router-link to="/" class="logo">
          <svg width="20" height="20" viewBox="0 0 76 76" fill="none">
            <rect width="76" height="76" rx="12" fill="var(--color-ink)"/>
            <path d="M49 26H27v24l22-24z" fill="var(--color-canvas)"/>
            <path d="M38 38L27 50h22L38 38z" fill="var(--color-canvas)" fill-opacity="0.5"/>
          </svg>
          <span class="logo-text">Tunnel Manager</span>
        </router-link>
      </div>

      <div class="nav-center">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="nav-link"
          :class="{ active: route.path === item.path }"
        >
          {{ item.label }}
        </router-link>
      </div>

      <div class="nav-right">
        <button class="hamburger" @click="mobileOpen = !mobileOpen" :aria-label="mobileOpen ? '关闭菜单' : '打开菜单'">
          <span class="hamburger-line" :class="{ open: mobileOpen }"></span>
          <span class="hamburger-line" :class="{ open: mobileOpen }"></span>
          <span class="hamburger-line" :class="{ open: mobileOpen }"></span>
        </button>
        <button class="icon-button" @click="configStore.toggleDarkMode()" :title="configStore.darkMode ? '亮色模式' : '暗色模式'">
          <svg v-if="configStore.darkMode" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="5"/>
            <line x1="12" y1="1" x2="12" y2="3"/>
            <line x1="12" y1="21" x2="12" y2="23"/>
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
            <line x1="1" y1="12" x2="3" y2="12"/>
            <line x1="21" y1="12" x2="23" y2="12"/>
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
          </svg>
          <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
          </svg>
        </button>
        <button v-if="configStore.isAuthenticated" class="icon-button logout-btn" @click="handleLogout" title="退出登录">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
            <polyline points="16 17 21 12 16 7"/>
            <line x1="21" y1="12" x2="9" y2="12"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Mobile menu overlay -->
    <transition name="fade">
      <div v-if="mobileOpen" class="mobile-overlay" @click="mobileOpen = false"></div>
    </transition>

    <!-- Mobile menu panel -->
    <transition name="slide">
      <div v-if="mobileOpen" class="mobile-menu">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="mobile-nav-link"
          :class="{ active: route.path === item.path }"
          @click="mobileOpen = false"
        >
          {{ item.label }}
        </router-link>
        <div class="mobile-menu-divider"></div>
        <button class="mobile-logout-btn" @click="handleLogout">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
            <polyline points="16 17 21 12 16 7"/>
            <line x1="21" y1="12" x2="9" y2="12"/>
          </svg>
          退出登录
        </button>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useConfigStore } from '../stores/config'
import { logout as logoutApi } from '../api'

const route = useRoute()
const router = useRouter()
const configStore = useConfigStore()
const mobileOpen = ref(false)

const navItems = [
  { path: '/', label: '控制面板' },
  { path: '/tunnels', label: '隧道管理' },
  { path: '/domain', label: '域名绑定' },
  { path: '/settings', label: '全局设置' },
  { path: '/account', label: '账户' },
]

async function handleLogout() {
  try { await logoutApi() } catch (_) { /* ignore */ }
  configStore.clearAuth()
  mobileOpen.value = false
  router.push('/login')
}
</script>

<style scoped>
.nav-bar {
  position: sticky;
  top: 0;
  z-index: 100;
  height: var(--header-height);
  background: var(--color-canvas);
  border-bottom: 1px solid var(--color-hairline);
}

.nav-inner {
  max-width: var(--max-width);
  height: 100%;
  margin: 0 auto;
  padding: 0 var(--spacing-lg);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-left {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  color: var(--color-ink);
}

.logo-text {
  font-size: 16px;
  font-weight: 600;
  letter-spacing: -0.4px;
}


.nav-center {
  display: flex;
  align-items: center;
  gap: 4px;
}

.nav-link {
  padding: 6px 12px;
  border-radius: 9999px;
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
  color: var(--color-body);
  text-decoration: none;
  transition: all 0.15s ease;
}

.nav-link:hover {
  color: var(--color-ink);
  background: var(--color-canvas-soft-2);
}

.nav-link.active {
  color: var(--color-ink);
  background: var(--color-canvas-soft-2);
  font-weight: 500;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.icon-button {
  width: 32px;
  height: 32px;
  border-radius: 9999px;
  border: 1px solid var(--color-hairline);
  background: var(--color-canvas);
  color: var(--color-ink);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.icon-button:hover {
  border-color: var(--color-hairline-strong);
}

.logout-btn {
  color: var(--color-mute);
}
.logout-btn:hover {
  color: var(--color-error);
  border-color: var(--color-error);
}

/* Hamburger */
.hamburger {
  display: none;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  border: 1px solid var(--color-hairline);
  background: var(--color-canvas);
  cursor: pointer;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 0;
}

.hamburger-line {
  display: block;
  width: 16px;
  height: 2px;
  background: var(--color-ink);
  border-radius: 1px;
  transition: all 0.2s ease;
}

.hamburger-line.open:nth-child(1) {
  transform: translateY(6px) rotate(45deg);
}
.hamburger-line.open:nth-child(2) {
  opacity: 0;
}
.hamburger-line.open:nth-child(3) {
  transform: translateY(-6px) rotate(-45deg);
}

/* Mobile overlay */
.mobile-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.4);
  z-index: 99;
}

/* Mobile menu */
.mobile-menu {
  position: fixed;
  top: var(--header-height);
  left: 0;
  right: 0;
  background: var(--color-canvas);
  border-bottom: 1px solid var(--color-hairline);
  padding: var(--spacing-sm);
  z-index: 100;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.mobile-nav-link {
  display: block;
  padding: 12px var(--spacing-md);
  border-radius: 6px;
  font-size: 15px;
  font-weight: 500;
  color: var(--color-body);
  text-decoration: none;
  transition: all 0.15s ease;
}

.mobile-nav-link:hover,
.mobile-nav-link.active {
  color: var(--color-ink);
  background: var(--color-canvas-soft-2);
}

.mobile-menu-divider {
  height: 1px;
  background: var(--color-hairline);
  margin: var(--spacing-xs) 0;
}

.mobile-logout-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px var(--spacing-md);
  border-radius: 6px;
  font-size: 15px;
  font-weight: 500;
  color: var(--color-mute);
  background: transparent;
  border: none;
  cursor: pointer;
  width: 100%;
  text-align: left;
  transition: all 0.15s ease;
}

.mobile-logout-btn:hover {
  color: var(--color-error);
  background: var(--color-canvas-soft-2);
}

/* Transitions */
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

.slide-enter-active, .slide-leave-active { transition: all 0.25s ease; }
.slide-enter-from, .slide-leave-to { opacity: 0; transform: translateY(-12px); }

@media (max-width: 768px) {
  .nav-center { display: none; }
  .hamburger { display: flex; }
  .nav-inner { padding: 0 var(--spacing-md); }
}
</style>