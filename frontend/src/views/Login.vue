<template>
  <div class="login-page">
    <div class="login-card" :class="{ 'login-card-enter': mounted, 'login-card-shake': shaking }">
      <div class="login-logo">
        <svg width="32" height="32" viewBox="0 0 76 76" fill="none">
          <rect width="76" height="76" rx="12" fill="#171717"/>
          <path d="M49 26H27v24l22-24z" fill="white"/>
          <path d="M38 38L27 50h22L38 38z" fill="white" fill-opacity="0.5"/>
        </svg>
      </div>
      <h1 class="login-title">Tunnel Manager</h1>
      <p class="login-subtitle">登录以继续</p>

      <form class="login-form" @submit.prevent="handleLogin">
        <div class="field">
          <label class="field-label">用户名</label>
          <input
            v-model="form.username"
            type="text"
            placeholder="admin"
            class="vercel-input"
            :class="{ 'input-error': error }"
            autocomplete="username"
          />
        </div>
        <div class="field">
          <label class="field-label">密码</label>
          <input
            v-model="form.password"
            type="password"
            placeholder="••••••••"
            class="vercel-input"
            :class="{ 'input-error': error }"
            autocomplete="current-password"
          />
        </div>

        <div v-if="error" class="login-error">{{ error }}</div>

        <button type="submit" class="btn btn-primary login-btn" :disabled="loading">
          <span v-if="loading" class="spinner"></span>
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { login as loginApi } from '../api'
import { useConfigStore } from '../stores/config'

const router = useRouter()
const store = useConfigStore()

const form = reactive({ username: '', password: '' })
const loading = ref(false)
const error = ref('')
const mounted = ref(false)
const shaking = ref(false)

onMounted(() => {
  requestAnimationFrame(() => { mounted.value = true })
})

async function handleLogin() {
  if (!form.username || !form.password) {
    error.value = '请输入用户名和密码'
    triggerShake()
    return
  }
  loading.value = true
  error.value = ''
  try {
    const { data } = await loginApi(form.username, form.password)
    store.setAuth(data.token, data.username)
    router.push('/')
  } catch (e: any) {
    if (e.response?.status === 401) {
      error.value = '用户名或密码错误'
    } else {
      error.value = '登录失败: ' + (e.response?.data?.error || e.message)
    }
    triggerShake()
  } finally {
    loading.value = false
  }
}

function triggerShake() {
  shaking.value = true
  setTimeout(() => { shaking.value = false }, 500)
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-md);
  background: var(--color-canvas-soft);
  box-sizing: border-box;
}

.login-card {
  width: 100%;
  max-width: 380px;
  background: var(--color-canvas);
  border: 1px solid var(--color-hairline);
  border-radius: 12px;
  padding: var(--spacing-2xl);
  text-align: center;
  opacity: 0;
  transform: translateY(16px) scale(0.98);
  transition: opacity 0.5s ease-out, transform 0.5s cubic-bezier(0.16, 1, 0.3, 1);
}

.login-card-enter {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.login-card-shake {
  animation: shake 0.45s ease-out;
}

.login-logo {
  margin-bottom: var(--spacing-md);
}

.login-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-ink);
  margin: 0 0 4px 0;
}

.login-subtitle {
  font-size: 14px;
  color: var(--color-mute);
  margin: 0 0 var(--spacing-xl) 0;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  text-align: left;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.field-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-ink);
}

.login-error {
  font-size: 14px;
  color: var(--color-error);
  text-align: center;
  padding: 8px;
  background: var(--color-result-error-bg);
  border-radius: 6px;
}

.login-btn {
  width: 100%;
  justify-content: center;
  margin-top: var(--spacing-xs);
}

.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.2);
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

.login-error {
  animation: fadeIn 0.3s ease-out;
}

@keyframes spin { to { transform: rotate(360deg); } }
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  20% { transform: translateX(-4px); }
  40% { transform: translateX(4px); }
  60% { transform: translateX(-3px); }
  80% { transform: translateX(3px); }
}
@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
</style>