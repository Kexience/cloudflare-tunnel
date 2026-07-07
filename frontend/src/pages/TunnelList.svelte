<script lang="ts">
  import { navigate } from 'svelte-routing'
  import { Layout } from '../lib/layout'
  import { PrimaryButton, SecondaryButton, TunnelItem } from '../lib/components'
  import { listTunnels, createTunnel, deleteTunnel, startTunnel, stopTunnel, getTunnelStatus } from '../lib/tunnel/api'
  import { getCredentials } from '../lib/config/api'
  import type { TunnelVO, CreateTunnelRequest } from '../lib/tunnel/types'
  import type { Credential } from '../lib/config/types'

  let tunnels = $state<TunnelVO[]>([])
  let tunnelStatuses = $state<Record<string, string>>({})
  let credentials = $state<Credential[]>([])
  let loading = $state(false)
  let error = $state<string | null>(null)
  let selectedCredentialId = $state<number | null>(null)
  let showCreateForm = $state(false)
  let newTunnelName = $state('')

  async function loadCredentials() {
    try {
      const response = await getCredentials()
      if (response.code === 0 && response.data) {
        credentials = response.data
        if (credentials.length > 0) {
          selectedCredentialId = credentials[0].id
        }
      }
    } catch (err) {
      console.error('Failed to load credentials:', err)
    }
  }

  async function loadTunnels() {
    if (!selectedCredentialId) return
    loading = true
    error = null
    try {
      const response = await listTunnels({ credential_id: selectedCredentialId })
      if (response.code === 0 && response.data) {
        tunnels = response.data
        await loadTunnelStatuses()
      } else {
        error = response.message || '加载隧道列表失败'
      }
    } catch (err) {
      error = '加载隧道列表失败'
      console.error('Failed to load tunnels:', err)
    } finally {
      loading = false
    }
  }

  async function loadTunnelStatuses() {
    if (!selectedCredentialId) return
    const statuses: Record<string, string> = {}
    await Promise.all(
      tunnels.map(async (tunnel) => {
        try {
          const response = await getTunnelStatus(tunnel.id, { credential_id: selectedCredentialId! })
          if (response.code === 0 && response.data) {
            statuses[tunnel.id] = response.data.status
          }
        } catch (err) {
          console.error(`Failed to load status for tunnel ${tunnel.id}:`, err)
        }
      })
    )
    tunnelStatuses = statuses
  }

  async function handleCreateTunnel() {
    if (!selectedCredentialId || !newTunnelName.trim()) return
    loading = true
    error = null
    try {
      const request: CreateTunnelRequest = {
        credential_id: selectedCredentialId,
        name: newTunnelName.trim()
      }
      const response = await createTunnel(request)
      if (response.code === 0) {
        newTunnelName = ''
        showCreateForm = false
        await loadTunnels()
      } else {
        error = response.message || '创建隧道失败'
      }
    } catch (err) {
      error = '创建隧道失败'
      console.error('Failed to create tunnel:', err)
    } finally {
      loading = false
    }
  }

  async function handleDeleteTunnel(tunnelId: string) {
    if (!selectedCredentialId) return
    if (!confirm('确定要删除这个隧道吗？')) return
    loading = true
    error = null
    try {
      const response = await deleteTunnel(tunnelId, { credential_id: selectedCredentialId })
      if (response.code === 0) {
        await loadTunnels()
      } else {
        error = response.message || '删除隧道失败'
      }
    } catch (err) {
      error = '删除隧道失败'
      console.error('Failed to delete tunnel:', err)
    } finally {
      loading = false
    }
  }

  async function handleStartTunnel(tunnelId: string) {
    if (!selectedCredentialId) return
    loading = true
    error = null
    try {
      const response = await startTunnel(tunnelId, { credential_id: selectedCredentialId })
      if (response.code === 0) {
        await loadTunnelStatuses()
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

  async function handleStopTunnel(tunnelId: string) {
    if (!selectedCredentialId) return
    loading = true
    error = null
    try {
      const response = await stopTunnel(tunnelId, { credential_id: selectedCredentialId })
      if (response.code === 0) {
        await loadTunnelStatuses()
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

  $effect(() => {
    loadCredentials()
  })

  $effect(() => {
    if (selectedCredentialId) {
      loadTunnels()
    }
  })
</script>

<Layout title="隧道管理" subtitle="管理您的 Cloudflare Tunnel 配置和状态">
  <div class="flex items-center justify-between mb-8">
    <div></div>
    <PrimaryButton onclick={() => showCreateForm = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
      </svg>
      创建隧道
    </PrimaryButton>
  </div>

  <!-- 筛选条件 -->
  <div class="bg-white rounded-2xl shadow-lg border border-gray-100 p-6 mb-6">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label for="credential" class="block text-sm font-medium text-gray-700 mb-2">选择凭证</label>
        <select
          id="credential"
          bind:value={selectedCredentialId}
          class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
        >
          {#each credentials as credential}
            <option value={credential.id}>{credential.name}</option>
          {/each}
        </select>
      </div>
      <div class="flex items-end">
        <PrimaryButton
          onclick={loadTunnels}
          disabled={loading || !selectedCredentialId}
          loading={loading}
          loadingText="加载中..."
        >
          查询隧道
        </PrimaryButton>
      </div>
    </div>
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

  <!-- 隧道列表 -->
  <div class="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden">
    <div class="px-6 py-4 border-b border-gray-200 bg-gray-50/50">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-semibold text-gray-900">隧道列表</h3>
        <button 
          onclick={loadTunnels}
          class="text-sm text-indigo-600 hover:text-indigo-500 font-medium transition duration-200"
        >
          刷新
        </button>
      </div>
    </div>
    <div class="p-6">
      {#if loading}
        <div class="flex justify-center items-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
        </div>
      {:else if tunnels.length === 0}
        <div class="text-center py-8">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">暂无隧道</h3>
          <p class="mt-1 text-sm text-gray-500">点击"创建隧道"按钮开始使用</p>
        </div>
      {:else}
        <div class="space-y-4">
          {#each tunnels as tunnel}
            <TunnelItem
              {tunnel}
              status={tunnelStatuses[tunnel.id]}
              onView={(id) => navigate(`/tunnels/${id}/${selectedCredentialId}`)}
              onDelete={handleDeleteTunnel}
              onStart={handleStartTunnel}
              onStop={handleStopTunnel}
            />
          {/each}
        </div>
      {/if}
    </div>
  </div>

  <!-- 创建隧道模态框 -->
  {#if showCreateForm}
    <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <h3 class="text-lg font-medium text-gray-900 mb-4">创建新隧道</h3>
          <form onsubmit={handleCreateTunnel}>
            <div class="mb-4">
              <label for="credential" class="block text-sm font-medium text-gray-700 mb-2">选择凭证</label>
              <select
                id="credential"
                bind:value={selectedCredentialId}
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              >
                {#each credentials as credential}
                  <option value={credential.id}>{credential.name}</option>
                {/each}
              </select>
            </div>
            <div class="mb-4">
              <label for="tunnelName" class="block text-sm font-medium text-gray-700 mb-2">隧道名称</label>
              <input
                type="text"
                id="tunnelName"
                bind:value={newTunnelName}
                placeholder="请输入隧道名称"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                required
              />
            </div>
            <div class="flex justify-end space-x-3">
              <SecondaryButton onclick={() => showCreateForm = false} />
              <PrimaryButton type="submit" loading={loading} loadingText="创建中...">
                创建
              </PrimaryButton>
            </div>
          </form>
        </div>
      </div>
    </div>
  {/if}
</Layout>
