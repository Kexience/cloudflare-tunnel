<script lang="ts">
import { onMount } from 'svelte'
import { navigate } from 'svelte-routing'
import { Layout } from '../lib/layout'
import { PrimaryButton, SecondaryButton } from '../lib/components'
import { getTunnel, getTunnelToken, listTunnelConnections, getTunnelConfig, updateTunnelConfig, startTunnel, stopTunnel, getTunnelStatus } from '../lib/tunnel/api'
import type { TunnelVO, TunnelTokenVO, TunnelConfigVO, TunnelConnectionVO, IngressRule } from '../lib/tunnel/types'

let { tunnelId, credentialId }: { tunnelId: string; credentialId: number } = $props()

let tunnel = $state<TunnelVO | null>(null)
let token = $state<TunnelTokenVO | null>(null)
let connections = $state<TunnelConnectionVO[]>([])
let config = $state<TunnelConfigVO | null>(null)
let tunnelStatus = $state<string | null>(null)
let loading = $state(false)
let error = $state<string | null>(null)
let activeTab = $state<'details' | 'config' | 'connections' | 'token'>('details')
let editingConfig = $state(false)
let savingConfig = $state(false)
let ingressRules = $state<IngressRule[]>([])
let newRule = $state<IngressRule>({ hostname: '', service: '' })
let configEditMode = $state<'form' | 'json'>('form')
let configText = $state('')
let draggedIndex = $state<number | null>(null)
let dragOverIndex = $state<number | null>(null)

function cleanErrorMessage(message: string): string {
  return message.replace(/http_status:\d+/gi, '').replace(/\s+/g, ' ').trim()
}

  async function loadTunnelDetails() {
    loading = true
    error = null
    try {
      const [tunnelRes, tokenRes, connectionsRes, configRes, statusRes] = await Promise.all([
        getTunnel(tunnelId, { credential_id: credentialId }),
        getTunnelToken(tunnelId, { credential_id: credentialId }),
        listTunnelConnections(tunnelId, { credential_id: credentialId }),
        getTunnelConfig(tunnelId, { credential_id: credentialId }),
        getTunnelStatus(tunnelId, { credential_id: credentialId })
      ])

      if (tunnelRes.code === 0 && tunnelRes.data) {
        tunnel = tunnelRes.data
      } else {
        error = cleanErrorMessage(tunnelRes.message) || '加载隧道详情失败'
        return
      }

      if (tokenRes.code === 0 && tokenRes.data) {
        token = tokenRes.data
      }

      if (connectionsRes.code === 0 && connectionsRes.data) {
        connections = connectionsRes.data.connections || []
      }

if (configRes.code === 0 && configRes.data) {
  config = configRes.data
  ingressRules = config.config.ingress || []
  configText = JSON.stringify(config.config, null, 2)
}

      if (statusRes.code === 0 && statusRes.data) {
        tunnelStatus = statusRes.data.status
      }
    } catch (err) {
      error = '加载隧道详情失败'
      console.error('Failed to load tunnel details:', err)
    } finally {
      loading = false
    }
  }

async function handleSaveConfig() {
  if (!config) return
  savingConfig = true
  error = null
  try {
    let newConfig: Record<string, unknown>
    
    if (configEditMode === 'json') {
      try {
        newConfig = JSON.parse(configText)
      } catch {
        error = 'JSON格式错误'
        return
      }
    } else {
      newConfig = {
        ingress: ingressRules,
        ...(config.config['warp-routing'] && { 'warp-routing': config.config['warp-routing'] }),
        ...(config.config.originRequest && { originRequest: config.config.originRequest })
      }
    }

    const response = await updateTunnelConfig(tunnelId, {
      credential_id: credentialId,
      config: newConfig
    })

    if (response.code === 0 && response.data) {
      config = response.data
      ingressRules = response.data.config.ingress || []
      configText = JSON.stringify(response.data.config, null, 2)
      editingConfig = false
    } else {
      error = cleanErrorMessage(response.message) || '保存配置失败'
    }
  } catch (err) {
    error = '保存配置失败'
    console.error('Failed to save config:', err)
  } finally {
    savingConfig = false
  }
}

async function handleStartTunnel() {
  loading = true
  error = null
  try {
    const response = await startTunnel(tunnelId, { credential_id: credentialId })
    if (response.code === 0) {
      tunnelStatus = 'running'
    } else {
      error = response.message || '启动隧道失败'
    }
  } catch (err) {
    error = '启动隧道失败'
    console.error('Failed to start tunnel:', err)
  } finally {
    loading = false
  }
}

async function handleStopTunnel() {
  loading = true
  error = null
  try {
    const response = await stopTunnel(tunnelId, { credential_id: credentialId })
    if (response.code === 0) {
      tunnelStatus = 'stopped'
    } else {
      error = response.message || '停止隧道失败'
    }
  } catch (err) {
    error = '停止隧道失败'
    console.error('Failed to stop tunnel:', err)
  } finally {
    loading = false
  }
}

function addIngressRule() {
  if (!newRule.hostname || !newRule.service) return
  ingressRules = [...ingressRules, { ...newRule }]
  newRule = { hostname: '', service: '' }
}

function removeIngressRule(index: number) {
  ingressRules = ingressRules.filter((_, i) => i !== index)
}

function updateIngressRule(index: number, field: keyof IngressRule, value: string) {
  ingressRules = ingressRules.map((rule, i) => 
    i === index ? { ...rule, [field]: value } : rule
  )
}

function handleDragStart(index: number) {
  draggedIndex = index
}

function handleDragOver(event: DragEvent, index: number) {
  event.preventDefault()
  dragOverIndex = index
}

function handleDragEnd() {
  if (draggedIndex !== null && dragOverIndex !== null && draggedIndex !== dragOverIndex) {
    const newRules = [...ingressRules]
    const draggedRule = newRules[draggedIndex]
    newRules.splice(draggedIndex, 1)
    newRules.splice(dragOverIndex, 0, draggedRule)
    ingressRules = newRules
  }
  draggedIndex = null
  dragOverIndex = null
}

function handleDrop(event: DragEvent, index: number) {
  event.preventDefault()
  handleDragEnd()
}

  function copyToClipboard(text: string) {
    navigator.clipboard.writeText(text).then(() => {
      alert('已复制到剪贴板')
    }).catch(() => {
      alert('复制失败')
    })
  }

  onMount(loadTunnelDetails)
</script>

<Layout title="隧道详情" subtitle="查看和管理隧道配置">
  <div class="mb-6">
    <button
      onclick={() => navigate('/tunnels')}
      class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition duration-200"
      aria-label="返回隧道列表"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
    </button>
  </div>

  {#if error}
    <div class="mb-6 bg-red-50 border border-red-200 rounded-xl p-4 flex items-center">
      <svg class="w-5 h-5 text-red-500 mr-3 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span class="text-red-700 text-sm">{error}</span>
      <button onclick={() => error = null} aria-label="关闭" class="ml-auto text-red-400 hover:text-red-600">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  {/if}

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto mb-4"></div>
        <p class="text-gray-500">加载中...</p>
      </div>
    </div>
  {:else if tunnel}
    <!-- 隧道基本信息卡片 -->
    <div class="bg-white rounded-2xl shadow-lg border border-gray-100 p-6 mb-6">
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center space-x-4">
          <div class="w-16 h-16 bg-linear-to-r from-green-500 to-emerald-500 rounded-2xl flex items-center justify-center">
            <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
            </svg>
          </div>
          <div>
            <h3 class="text-xl font-semibold text-gray-900">{tunnel.name}</h3>
            <p class="text-sm text-gray-500">ID: {tunnel.id}</p>
          </div>
        </div>
        <div class="flex items-center space-x-3">
          <span class="px-3 py-1 text-sm font-medium rounded-full {tunnel.status === 'healthy' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'}">
            {tunnel.status}
          </span>
          {#if tunnelStatus}
            <span class="px-3 py-1 text-sm font-medium rounded-full {tunnelStatus === 'running' ? 'bg-blue-100 text-blue-800' : 'bg-gray-100 text-gray-800'}">
              {tunnelStatus === 'running' ? '运行中' : '已停止'}
            </span>
          {/if}
          {#if tunnelStatus === 'running'}
            <PrimaryButton onclick={handleStopTunnel} disabled={loading} class="bg-orange-600 hover:bg-orange-700">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6" />
              </svg>
              停止
            </PrimaryButton>
          {:else if tunnelStatus === 'stopped'}
            <PrimaryButton onclick={handleStartTunnel} disabled={loading} class="bg-green-600 hover:bg-green-700">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              启动
            </PrimaryButton>
          {/if}
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="bg-gray-50 rounded-lg p-4">
          <p class="text-sm text-gray-500 mb-1">创建时间</p>
          <p class="text-gray-900">{tunnel.created_at ? new Date(tunnel.created_at).toLocaleString() : '未知'}</p>
        </div>
        <div class="bg-gray-50 rounded-lg p-4">
          <p class="text-sm text-gray-500 mb-1">连接数</p>
          <p class="text-gray-900">{connections.length}</p>
        </div>
        <div class="bg-gray-50 rounded-lg p-4">
          <p class="text-sm text-gray-500 mb-1">配置版本</p>
          <p class="text-gray-900">{config?.version || '未知'}</p>
        </div>
      </div>
    </div>

    <!-- 标签页导航 -->
    <div class="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden">
      <div class="border-b border-gray-200">
        <nav class="flex -mb-px">
          <button
            onclick={() => activeTab = 'details'}
            class="py-4 px-6 text-center border-b-2 font-medium text-sm {activeTab === 'details' ? 'border-indigo-500 text-indigo-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
          >
            详情
          </button>
          <button
            onclick={() => activeTab = 'config'}
            class="py-4 px-6 text-center border-b-2 font-medium text-sm {activeTab === 'config' ? 'border-indigo-500 text-indigo-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
          >
            配置
          </button>
          <button
            onclick={() => activeTab = 'connections'}
            class="py-4 px-6 text-center border-b-2 font-medium text-sm {activeTab === 'connections' ? 'border-indigo-500 text-indigo-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
          >
            连接
          </button>
          <button
            onclick={() => activeTab = 'token'}
            class="py-4 px-6 text-center border-b-2 font-medium text-sm {activeTab === 'token' ? 'border-indigo-500 text-indigo-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
          >
            Token
          </button>
        </nav>
      </div>

      <div class="p-6">
        {#if activeTab === 'details'}
          <div class="space-y-4">
            <h4 class="text-lg font-semibold text-gray-900">基本信息</h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-gray-500">隧道ID</p>
                <p class="text-gray-900 font-mono">{tunnel.id}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500">名称</p>
                <p class="text-gray-900">{tunnel.name}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500">状态</p>
                <p class="text-gray-900">{tunnel.status}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500">创建时间</p>
                <p class="text-gray-900">{tunnel.created_at ? new Date(tunnel.created_at).toLocaleString() : '未知'}</p>
              </div>
            </div>
          </div>

        {:else if activeTab === 'config'}
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <h4 class="text-lg font-semibold text-gray-900">隧道配置</h4>
              <div class="flex space-x-2">
                {#if editingConfig}
                  <SecondaryButton onclick={() => editingConfig = false} />
                  <PrimaryButton onclick={handleSaveConfig} loading={savingConfig} loadingText="保存中...">
                    保存
                  </PrimaryButton>
                {:else}
                  <PrimaryButton onclick={() => editingConfig = true}>
                    编辑配置
                  </PrimaryButton>
                {/if}
              </div>
            </div>

            {#if editingConfig}
              <!-- 子标签切换 -->
              <div class="border-b border-gray-200">
                <nav class="flex -mb-px">
                  <button
                    onclick={() => configEditMode = 'form'}
                    class="py-2 px-4 text-center border-b-2 font-medium text-sm {configEditMode === 'form' ? 'border-indigo-500 text-indigo-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
                  >
                    表单编辑
                  </button>
                  <button
                    onclick={() => configEditMode = 'json'}
                    class="py-2 px-4 text-center border-b-2 font-medium text-sm {configEditMode === 'json' ? 'border-indigo-500 text-indigo-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
                  >
                    JSON编辑
                  </button>
                </nav>
              </div>

              {#if configEditMode === 'form'}
                <div class="space-y-4">
                  <div class="bg-gray-50 rounded-lg p-4">
                    <h5 class="text-sm font-medium text-gray-700 mb-3">Ingress 规则</h5>
                    <div class="space-y-3">
                      {#each ingressRules as rule, index}
                        <div
                          class="flex items-start space-x-3 p-3 bg-white rounded-md border border-gray-200 {dragOverIndex === index ? 'border-indigo-500 bg-indigo-50' : ''}"
                          draggable="true"
                          role="listitem"
                          ondragstart={() => handleDragStart(index)}
                          ondragover={(e) => handleDragOver(e, index)}
                          ondragend={handleDragEnd}
                          ondrop={(e) => handleDrop(e, index)}
                        >
                          <div class="flex items-center justify-center w-8 h-8 text-gray-400 hover:text-gray-600 cursor-grab active:cursor-grabbing">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h16M4 16h16" />
                            </svg>
                          </div>
                          <div class="flex-1 grid grid-cols-1 md:grid-cols-2 gap-3">
                            <div>
                              <label for={`hostname-${index}`} class="block text-xs font-medium text-gray-500 mb-1">主机名</label>
                              <input
                                id={`hostname-${index}`}
                                type="text"
                                value={rule.hostname}
                                oninput={(e) => updateIngressRule(index, 'hostname', e.currentTarget.value)}
                                placeholder="example.com"
                                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                              />
                            </div>
                            <div>
                              <label for={`service-${index}`} class="block text-xs font-medium text-gray-500 mb-1">服务地址</label>
                              <input
                                id={`service-${index}`}
                                type="text"
                                value={rule.service}
                                oninput={(e) => updateIngressRule(index, 'service', e.currentTarget.value)}
                                placeholder="http://localhost:8080"
                                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                              />
                            </div>
                          </div>
                          <button
                            onclick={() => removeIngressRule(index)}
                            class="mt-5 p-2 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-md transition duration-200"
                            title="删除规则"
                          >
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                            </svg>
                          </button>
                        </div>
                      {/each}
                    </div>
                  </div>

                  <div class="bg-gray-50 rounded-lg p-4">
                    <h5 class="text-sm font-medium text-gray-700 mb-3">添加新规则</h5>
                    <div class="flex items-start space-x-3">
                      <div class="flex-1 grid grid-cols-1 md:grid-cols-2 gap-3">
                        <div>
                          <label for="new-hostname" class="block text-xs font-medium text-gray-500 mb-1">主机名</label>
                          <input
                            id="new-hostname"
                            type="text"
                            bind:value={newRule.hostname}
                            placeholder="example.com"
                            class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                          />
                        </div>
                        <div>
                          <label for="new-service" class="block text-xs font-medium text-gray-500 mb-1">服务地址</label>
                          <input
                            id="new-service"
                            type="text"
                            bind:value={newRule.service}
                            placeholder="http://localhost:8080"
                            class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                          />
                        </div>
                      </div>
                      <PrimaryButton
                        onclick={addIngressRule}
                        disabled={!newRule.hostname || !newRule.service}
                      >
                        添加
                      </PrimaryButton>
                    </div>
                  </div>
                </div>
              {:else}
                <div class="space-y-4">
                  <div class="bg-gray-50 rounded-lg p-4">
                    <h5 class="text-sm font-medium text-gray-700 mb-3">JSON 配置</h5>
                    <textarea
                      bind:value={configText}
                      class="w-full h-96 px-4 py-3 font-mono text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                      placeholder="输入JSON配置"
                    ></textarea>
                  </div>
                </div>
              {/if}
            {:else if config}
              <div class="space-y-3">
                {#each ingressRules as rule}
                  <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                    <div class="flex items-center space-x-4">
                      <div class="w-10 h-10 bg-linear-to-r from-blue-500 to-cyan-500 rounded-lg flex items-center justify-center">
                        <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                        </svg>
                      </div>
                      <div>
                        <p class="text-sm font-medium text-gray-900">{rule.hostname}</p>
                        <p class="text-xs text-gray-500">{rule.service}</p>
                      </div>
                    </div>
                    <div class="flex items-center space-x-2">
                      <span class="px-2 py-1 text-xs font-medium rounded-full bg-green-100 text-green-800">
                        活跃
                      </span>
                    </div>
                  </div>
                {/each}
              </div>
            {:else}
              <p class="text-gray-500">暂无配置信息</p>
            {/if}
          </div>

        {:else if activeTab === 'connections'}
          <div class="space-y-4">
            <h4 class="text-lg font-semibold text-gray-900">活跃连接</h4>
            {#if connections.length === 0}
              <p class="text-gray-500">暂无活跃连接</p>
            {:else}
              <div class="space-y-3">
                {#each connections as connection}
                  <div class="bg-gray-50 rounded-lg p-4">
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <p class="text-sm text-gray-500">连接ID</p>
                        <p class="text-gray-900 font-mono">{connection.id}</p>
                      </div>
                      <div>
                        <p class="text-sm text-gray-500">数据中心</p>
                        <p class="text-gray-900">{connection.colo_name}</p>
                      </div>
                      <div>
                        <p class="text-sm text-gray-500">客户端ID</p>
                        <p class="text-gray-900 font-mono">{connection.client_id}</p>
                      </div>
                      <div>
                        <p class="text-sm text-gray-500">客户端版本</p>
                        <p class="text-gray-900">{connection.client_version}</p>
                      </div>
                      <div>
                        <p class="text-sm text-gray-500">打开时间</p>
                        <p class="text-gray-900">{new Date(connection.opened_at).toLocaleString()}</p>
                      </div>
                      <div>
                        <p class="text-sm text-gray-500">源IP</p>
                        <p class="text-gray-900">{connection.origin_ip}</p>
                      </div>
                    </div>
                  </div>
                {/each}
              </div>
            {/if}
          </div>

        {:else if activeTab === 'token'}
          <div class="space-y-4">
            <h4 class="text-lg font-semibold text-gray-900">隧道Token</h4>
            <p class="text-gray-600">用于连接隧道的认证Token，请妥善保管。</p>
            
            {#if token}
              <div class="bg-gray-50 rounded-lg p-4">
                <div class="flex items-center justify-between">
                  <div class="flex-1 min-w-0">
                    <p class="text-sm text-gray-500 mb-2">Token</p>
                    <code class="block font-mono text-sm text-gray-900 break-all">{token.token}</code>
                  </div>
                  <button
                    onclick={() => token && copyToClipboard(token.token)}
                    class="ml-4 p-2 text-gray-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition duration-200"
                    title="复制Token"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
                    </svg>
                  </button>
                </div>
              </div>
            {:else}
              <p class="text-gray-500">无法获取Token</p>
            {/if}
          </div>
        {/if}
      </div>
    </div>
  {:else}
    <div class="bg-white rounded-2xl shadow-lg border border-gray-100 p-12 text-center">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">隧道不存在</h3>
      <p class="mt-1 text-sm text-gray-500">无法找到指定的隧道</p>
      <div class="mt-6">
        <PrimaryButton onclick={() => navigate('/tunnels')}>
          返回隧道列表
        </PrimaryButton>
      </div>
    </div>
  {/if}
</Layout>
