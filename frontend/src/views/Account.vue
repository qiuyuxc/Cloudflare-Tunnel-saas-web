<template>
  <div class="page-container" style="padding-top: 0;">
    <div class="page-header">
      <router-link to="/" class="back-link">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"/></svg>
        返回控制面板
      </router-link>
      <h2>账户设置</h2>
      <p>管理登录凭据</p>
    </div>

    <div class="settings-list">
      <div class="settings-card card-transition" :class="{ 'stagger-item': visible }" :style="{ animationDelay: '0.08s' }">
        <div class="settings-card-header">
          <div class="settings-card-title">修改用户名</div>
          <div class="settings-card-desc">当前用户名: <strong>{{ store.username }}</strong></div>
        </div>
        <div class="settings-input-row">
          <div class="input-wrapper">
            <input v-model="newUsername" placeholder="新用户名" class="vercel-input" />
          </div>
        </div>
        <div class="settings-input-row" style="margin-top: var(--spacing-sm);">
          <div class="input-wrapper">
            <input v-model="usernamePassword" type="password" placeholder="输入当前密码确认" class="vercel-input" />
          </div>
          <button class="btn btn-secondary" :disabled="savingUsername" @click="saveUsername">
            {{ savingUsername ? '保存中...' : '更新用户名' }}
          </button>
        </div>
      </div>

      <div class="settings-card card-transition" :class="{ 'stagger-item': visible }" :style="{ animationDelay: '0.16s' }">
        <div class="settings-card-header">
          <div class="settings-card-title">修改密码</div>
          <div class="settings-card-desc">密码长度不少于 6 位</div>
        </div>
        <div class="settings-input-row">
          <div class="input-wrapper">
            <input v-model="currentPassword" type="password" placeholder="当前密码" class="vercel-input" />
          </div>
        </div>
        <div class="settings-input-row" style="margin-top: var(--spacing-sm);">
          <div class="input-wrapper">
            <input v-model="newPassword" type="password" placeholder="新密码" class="vercel-input" />
          </div>
          <button class="btn btn-primary" :disabled="savingPassword" @click="savePassword">
            {{ savingPassword ? '保存中...' : '更新密码' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { changePassword, changeUsername } from '../api'
import { useConfigStore } from '../stores/config'

const message = useMessage()
const store = useConfigStore()
const visible = ref(false)

const newUsername = ref('')
const usernamePassword = ref('')
const currentPassword = ref('')
const newPassword = ref('')
const savingUsername = ref(false)
const savingPassword = ref(false)

async function saveUsername() {
  if (!newUsername.value || !usernamePassword.value) {
    message.error('请填写新用户名和当前密码')
    return
  }
  savingUsername.value = true
  try {
    await changeUsername(usernamePassword.value, newUsername.value)
    store.clearAuth()
    message.success('用户名已更新，请重新登录')
  } catch (e: any) {
    message.error('更新失败: ' + (e.response?.data?.error || e.message))
  } finally {
    savingUsername.value = false
  }
}

onMounted(() => {
  requestAnimationFrame(() => { visible.value = true })
})

async function savePassword() {
  if (!currentPassword.value || !newPassword.value) {
    message.error('请填写当前密码和新密码')
    return
  }
  if (newPassword.value.length < 6) {
    message.error('新密码长度不能少于 6 位')
    return
  }
  savingPassword.value = true
  try {
    await changePassword(currentPassword.value, newPassword.value)
    store.clearAuth()
    message.success('密码已更新，请重新登录')
  } catch (e: any) {
    message.error('更新失败: ' + (e.response?.data?.error || e.message))
  } finally {
    savingPassword.value = false
  }
}
</script>

<style scoped>
.page-header { margin-bottom: var(--spacing-lg); }

.settings-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.settings-card {
  background: var(--color-canvas);
  border: 1px solid var(--color-hairline);
  border-radius: 8px;
  padding: var(--spacing-lg);
  transition: border-color 0.25s ease, box-shadow 0.25s ease, transform 0.25s ease-out;
}
.settings-card:hover {
  border-color: var(--color-hairline-strong);
  box-shadow: 0px 0px 0px 1px var(--color-hairline-strong), 0px 1px 1px rgba(0,0,0,0.02), 0px 2px 2px rgba(0,0,0,0.04);
  transform: translateY(-1px);
}

.settings-card-header { margin-bottom: var(--spacing-md); }
.settings-card-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--color-ink);
  margin-bottom: 4px;
}
.settings-card-desc {
  font-size: 14px;
  color: var(--color-mute);
  line-height: 20px;
}

.settings-input-row {
  display: flex;
  gap: 10px;
  align-items: flex-start;
}

.input-wrapper { flex: 1; }

@media (max-width: 768px) {
  .settings-input-row {
    flex-direction: column;
  }
  .settings-input-row .btn {
    width: 100%;
    max-width: 100%;
    justify-content: center;
    box-sizing: border-box;
  }
}
</style>
