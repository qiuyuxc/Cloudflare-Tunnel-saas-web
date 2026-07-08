<template>
  <div class="page-container" style="padding-top: 0;">
    <div class="page-header">
      <h2>全局设置</h2>
      <p>管理全局配置参数</p>
    </div>

    <div class="settings-list section">
      <div class="settings-card">
        <div class="settings-card-header">
          <div class="settings-card-title">全局优选 CNAME</div>
          <div class="settings-card-desc">该 CNAME 将用于主域名的 DNS 解析指向，以实现优选 IP 加速。</div>
        </div>
        <div class="settings-input-row">
          <div class="input-wrapper">
            <input v-model="preferredCNAME" placeholder="cf.090227.xyz" class="vercel-input" />
          </div>
          <button class="btn btn-primary" :disabled="savingCNAME" @click="savePreferredCNAME">
            {{ savingCNAME ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>

      <div class="settings-card">
        <div class="settings-card-header">
          <div class="settings-card-title">回退源设置</div>
          <div class="settings-card-desc">设置 Custom Hostnames 的回退源（Fallback Origin），用于 SaaS 模块。</div>
        </div>
        <div class="settings-input-row">
          <div class="input-wrapper">
            <input v-model="fallbackDomain" placeholder="例如: fallback.169977.xyz" class="vercel-input" />
          </div>
          <button class="btn btn-secondary" :disabled="savingFallback" @click="saveFallback">
            {{ savingFallback ? '设置中...' : '设置' }}
          </button>
        </div>
      </div>

      <div class="settings-card">
        <div class="settings-card-header">
          <div class="settings-card-title">转发地址</div>
          <div class="settings-card-desc">本地服务的 URL 地址，隧道将流量转发至此地址。</div>
        </div>
        <div class="settings-input-row">
          <div class="input-wrapper">
            <input v-model="serviceURL" placeholder="http://localhost:3000" class="vercel-input" />
          </div>
          <button class="btn btn-primary" :disabled="savingService" @click="saveServiceURL">
            {{ savingService ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>

      <div class="settings-card">
        <div class="settings-card-header">
          <div class="settings-card-title">隧道 ID</div>
          <div class="settings-card-desc">当前锁定的 Cloudflare Tunnel ID。</div>
        </div>
        <div class="settings-input-row">
          <div class="input-wrapper">
            <input v-model="tunnelID" placeholder="隧道 ID" class="vercel-input" />
          </div>
          <button class="btn btn-primary" :disabled="savingTunnel" @click="saveTunnelID">
            {{ savingTunnel ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { setPreferredCNAME, setFallbackOrigin, setServiceURL, setTunnelID } from '../api'
import { useConfigStore } from '../stores/config'

const message = useMessage()
const store = useConfigStore()
const config = store.config

const preferredCNAME = ref(config.preferred_cname)
const fallbackDomain = ref('')
const serviceURL = ref(config.service_url)
const tunnelID = ref(config.tunnel_id)

const savingCNAME = ref(false)
const savingFallback = ref(false)
const savingService = ref(false)
const savingTunnel = ref(false)

async function savePreferredCNAME() {
  savingCNAME.value = true
  try {
    await setPreferredCNAME(preferredCNAME.value)
    config.preferred_cname = preferredCNAME.value
    message.success('优选 CNAME 已更新')
  } catch (e: any) {
    message.error('保存失败: ' + (e.response?.data?.error || e.message))
  } finally {
    savingCNAME.value = false
  }
}

async function saveFallback() {
  savingFallback.value = true
  try {
    const { data } = await setFallbackOrigin(fallbackDomain.value)
    message.success(data.message || '回退源已设置')
  } catch (e: any) {
    message.error('设置失败: ' + (e.response?.data?.error || e.message))
  } finally {
    savingFallback.value = false
  }
}

async function saveServiceURL() {
  savingService.value = true
  try {
    await setServiceURL(serviceURL.value)
    config.service_url = serviceURL.value
    message.success('转发地址已更新')
  } catch (e: any) {
    message.error('保存失败: ' + (e.response?.data?.error || e.message))
  } finally {
    savingService.value = false
  }
}

async function saveTunnelID() {
  savingTunnel.value = true
  try {
    await setTunnelID(tunnelID.value)
    config.tunnel_id = tunnelID.value
    message.success('隧道 ID 已更新')
  } catch (e: any) {
    message.error('保存失败: ' + (e.response?.data?.error || e.message))
  } finally {
    savingTunnel.value = false
  }
}

onMounted(() => { store.fetchConfig() })
</script>

<style scoped>
.page-header { margin-bottom: var(--spacing-lg); }
.section { margin-bottom: 0; }

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