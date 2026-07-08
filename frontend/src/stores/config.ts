import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { getConfig, type Config } from '../api'

export const useConfigStore = defineStore('config', () => {
  const config = ref<Config>({
    tunnel_id: '',
    service_url: '',
    preferred_cname: 'cf.090227.xyz',
  })
  const darkMode = ref(localStorage.getItem('dark_mode') === 'true')
  const loading = ref(false)

  // Auth state
  const token = ref(localStorage.getItem('auth_token') || '')
  const username = ref(localStorage.getItem('auth_username') || '')
  const isAuthenticated = ref(!!token.value)

  // Persist darkMode
  watch(darkMode, (val) => {
    localStorage.setItem('dark_mode', String(val))
    document.documentElement.setAttribute('data-theme', val ? 'dark' : '')
  }, { immediate: true })

  async function fetchConfig() {
    loading.value = true
    try {
      const { data } = await getConfig()
      config.value = data
    } catch (e) {
      // Config might not be loaded yet
    } finally {
      loading.value = false
    }
  }

  function toggleDarkMode() {
    darkMode.value = !darkMode.value
  }

  function setAuth(tokenVal: string, usernameVal: string) {
    token.value = tokenVal
    username.value = usernameVal
    isAuthenticated.value = true
    localStorage.setItem('auth_token', tokenVal)
    localStorage.setItem('auth_username', usernameVal)
  }

  function clearAuth() {
    token.value = ''
    username.value = ''
    isAuthenticated.value = false
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_username')
  }

  return {
    config, darkMode, loading, token, username, isAuthenticated,
    fetchConfig, toggleDarkMode, setAuth, clearAuth,
  }
})