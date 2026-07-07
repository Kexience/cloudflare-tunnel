<script lang="ts">
  import { currentUser } from '../lib/auth/store'
  import { navigate } from 'svelte-routing'
  import { Layout } from '../lib/layout'
  import { listTunnels, createTunnel, deleteTunnel } from '../lib/tunnel/api'
  import { getCredentials } from '../lib/config/api'
  import type { TunnelVO, CreateTunnelRequest } from '../lib/tunnel/types'
  import type { Credential } from '../lib/config/types'

  let tunnels = $state<TunnelVO[]>([])
  let credentials = $state<Credential[]>([])
  let loading = $state(false)
  let error = $state<string | null>(null)
  let selectedCredentialId = $state<number | null>(null)
  let showCreateForm = $state(false)
  let newTunnelName = $state('')

  let healthyTunnels = $derived(tunnels.filter(t => t.status === 'healthy').length)

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

  $effect(() => {
    loadCredentials()
  })

  $effect(() => {
    if (selectedCredentialId) {
      loadTunnels()
    }
  })
</script>

<Layout title="欢迎回来，{$currentUser?.nickname || '用户'}！" subtitle="管理您的 Cloudflare Tunnel 配置和状态">
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
    <div class="bg-white rounded-2xl shadow-lg p-6 border border-gray-100 hover:shadow-xl transition duration-300">
      <div class="flex items-center justify-between mb-4">
        <div class="w-12 h-12 bg-linear-to-r from-green-500 to-emerald-500 rounded-xl flex items-center justify-center">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <span class="px-3 py-1 bg-green-100 text-green-800 text-xs font-medium rounded-full">{healthyTunnels} 运行中</span>
      </div>
      <h3 class="text-lg font-semibold text-gray-900 mb-1">隧道状态</h3>
      <p class="text-gray-600 text-sm">所有隧道正常运行</p>
    </div>

    <div class="bg-white rounded-2xl shadow-lg p-6 border border-gray-100 hover:shadow-xl transition duration-300">
      <div class="flex items-center justify-between mb-4">
        <div class="w-12 h-12 bg-linear-to-r from-blue-500 to-cyan-500 rounded-xl flex items-center justify-center">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
        </div>
        <span class="text-2xl font-bold text-gray-900">{tunnels.length}</span>
      </div>
      <h3 class="text-lg font-semibold text-gray-900 mb-1">隧道数量</h3>
      <p class="text-gray-600 text-sm">已配置的隧道总数</p>
    </div>

    <div class="bg-white rounded-2xl shadow-lg p-6 border border-gray-100 hover:shadow-xl transition duration-300">
      <div class="flex items-center justify-between mb-4">
        <div class="w-12 h-12 bg-linear-to-r from-purple-500 to-pink-500 rounded-xl flex items-center justify-center">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
        </div>
        <span class="text-2xl font-bold text-gray-900">0 GB</span>
      </div>
      <h3 class="text-lg font-semibold text-gray-900 mb-1">流量使用</h3>
      <p class="text-gray-600 text-sm">本月已使用流量</p>
    </div>
  </div>

  <div class="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden">
    <div class="px-6 py-4 border-b border-gray-200 bg-gray-50/50">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-semibold text-gray-900">快速操作</h3>
        <button class="text-sm text-indigo-600 hover:text-indigo-500 font-medium transition duration-200">
          查看全部
        </button>
      </div>
    </div>
    <div class="p-6">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <button 
          onclick={() => navigate('/tunnels')}
          class="flex flex-col items-center p-4 bg-linear-to-br from-indigo-50 to-purple-50 rounded-xl border border-indigo-100 hover:border-indigo-300 hover:shadow-md transition duration-200"
        >
          <div class="w-10 h-10 bg-linear-to-r from-indigo-500 to-purple-500 rounded-lg flex items-center justify-center mb-3">
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
            </svg>
          </div>
          <span class="text-sm font-medium text-gray-900">隧道列表</span>
        </button>

        <button
          onclick={() => navigate('/config')}
          class="flex flex-col items-center p-4 bg-linear-to-br from-green-50 to-emerald-50 rounded-xl border border-green-100 hover:border-green-300 hover:shadow-md transition duration-200"
        >
          <div class="w-10 h-10 bg-linear-to-r from-green-500 to-emerald-500 rounded-lg flex items-center justify-center mb-3">
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
          <span class="text-sm font-medium text-gray-900">配置管理</span>
        </button>

        <button 
          onclick={() => navigate('/dns')}
          class="flex flex-col items-center p-4 bg-linear-to-br from-blue-50 to-cyan-50 rounded-xl border border-blue-100 hover:border-blue-300 hover:shadow-md transition duration-200"
        >
          <div class="w-10 h-10 bg-linear-to-r from-blue-500 to-cyan-500 rounded-lg flex items-center justify-center mb-3">
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
            </svg>
          </div>
          <span class="text-sm font-medium text-gray-900">DNS管理</span>
        </button>

        <button class="flex flex-col items-center p-4 bg-linear-to-br from-orange-50 to-red-50 rounded-xl border border-orange-100 hover:border-orange-300 hover:shadow-md transition duration-200">
          <div class="w-10 h-10 bg-linear-to-r from-orange-500 to-red-500 rounded-lg flex items-center justify-center mb-3">
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
          </div>
          <span class="text-sm font-medium text-gray-900">日志查看</span>
        </button>
      </div>
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
              <button
                type="button"
                onclick={() => showCreateForm = false}
                class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-gray-300 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                取消
              </button>
              <button
                type="submit"
                disabled={loading}
                class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
              >
                {loading ? '创建中...' : '创建'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  {/if}

  <!-- 隧道列表 -->
  <div id="tunnel-list" class="mt-8 bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden">
    <div class="px-6 py-4 border-b border-gray-200 bg-gray-50/50">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-semibold text-gray-900">隧道列表</h3>
        <div class="flex space-x-2">
          <button 
            onclick={() => showCreateForm = true}
            class="text-sm text-white bg-indigo-600 hover:bg-indigo-700 font-medium transition duration-200 px-3 py-1 rounded-md"
          >
            创建隧道
          </button>
          <button 
            onclick={() => loadTunnels()}
            class="text-sm text-indigo-600 hover:text-indigo-500 font-medium transition duration-200"
          >
            刷新
          </button>
        </div>
      </div>
    </div>
    <div class="p-6">
      {#if error}
        <div class="mb-4 p-4 bg-red-50 border border-red-200 rounded-md">
          <p class="text-sm text-red-600">{error}</p>
        </div>
      {/if}

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
            <div class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition duration-200">
              <div class="flex items-center justify-between">
                <div class="flex items-center space-x-3">
                  <div class="w-10 h-10 bg-linear-to-r from-green-500 to-emerald-500 rounded-lg flex items-center justify-center">
                    <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                    </svg>
                  </div>
                  <div>
                    <h4 class="text-sm font-medium text-gray-900">{tunnel.name}</h4>
                    <p class="text-xs text-gray-500">ID: {tunnel.id}</p>
                  </div>
                </div>
                <div class="flex items-center space-x-2">
                  <span class="px-2 py-1 text-xs font-medium rounded-full {tunnel.status === 'healthy' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'}">
                    {tunnel.status}
                  </span>
                  <button
                    onclick={() => navigate(`/tunnels/${tunnel.id}/${selectedCredentialId}`)}
                    class="p-2 text-blue-600 hover:text-blue-900 hover:bg-blue-50 rounded-md transition duration-200"
                    title="查看详情"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                  </button>
                  <button
                    onclick={() => handleDeleteTunnel(tunnel.id)}
                    class="p-2 text-red-600 hover:text-red-900 hover:bg-red-50 rounded-md transition duration-200"
                    title="删除隧道"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </div>
              {#if tunnel.created_at}
                <div class="mt-2 text-xs text-gray-500">
                  创建时间: {new Date(tunnel.created_at).toLocaleString()}
                </div>
              {/if}
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>
</Layout>
