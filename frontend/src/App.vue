<template>
  <n-message-provider>
    <n-config-provider :theme="naiveTheme" :theme-overrides="themeOverrides">
      <n-layout class="app-layout">
        <nav-bar />
        <main class="main-content">
          <router-view />
        </main>
      </n-layout>
    </n-config-provider>
  </n-message-provider>
</template>

<script setup lang="ts">
import { darkTheme } from 'naive-ui'
import { NMessageProvider, NConfigProvider, NLayout } from 'naive-ui'
import { computed } from 'vue'
import NavBar from './components/NavBar.vue'
import { useConfigStore } from './stores/config'
import { vercelThemeOverrides } from './theme'

const configStore = useConfigStore()
const naiveTheme = computed(() => configStore.darkMode ? darkTheme : null)
const themeOverrides = computed(() => configStore.darkMode ? {} : vercelThemeOverrides)

// Sync dark mode to data-theme attribute on mount
if (configStore.darkMode) {
  document.documentElement.setAttribute('data-theme', 'dark')
}
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
  background-color: var(--color-canvas-soft);
}

.main-content {
  width: 100%;
}
</style>