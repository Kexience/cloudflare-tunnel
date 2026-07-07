<script lang="ts">
  import { onMount } from 'svelte'
  import { navigate } from 'svelte-routing'
  import { currentUser } from '../lib/auth/store'
  import { authStore } from '../lib/auth/store'
  import type { Credential, CreateCredentialRequest } from '../lib/config/types'
  import * as configApi from '../lib/config/api'
  import { validateCredential } from '../lib/credential'

  let credentials: Credential[] = $state([])
  let isLoading = $state(true)
  let error = $state<string | null>(null)
  let showModal = $state(false)
  let editingId = $state<number | null>(null)
  let saving = $state(false)
  let deletingId = $state<number | null>(null)
  let testingId = $state<number | null>(null)
  let testResult = $state<{ id: number; success: boolean; message: string } | null>(null)

  let formData: CreateCredentialRequest = $state({
    name: '',
    api_token: '',
    account_id: '',
    is_default: false
  })

  onMount(loadCredentials)

  async function loadCredentials() {
    isLoading = true
    error = null
    try {
      const res = await configApi.getCredentials()
      if (res.code === 0 && res.data) {
        credentials = res.data
      } else {
        error = res.message
      }
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : '加载失败'
    } finally {
      isLoading = false
    }
  }

  function openCreate() {
    editingId = null
    formData = { name: '', api_token: '', account_id: '', is_default: false }
    showModal = true
  }

  function openEdit(cred: Credential) {
    editingId = cred.id
    formData = {
      name: cred.name,
      api_token: cred.api_token,
      account_id: cred.account_id,
      is_default: cred.is_default
    }
    showModal = true
  }

  function closeModal() {
    showModal = false
    editingId = null
  }

  async function handleSubmit() {
    saving = true
    try {
      if (editingId !== null) {
        const res = await configApi.updateCredential(editingId, formData)
        if (res.code !== 0) {
          error = res.message
          return
        }
      } else {
        const res = await configApi.createCredential(formData)
        if (res.code !== 0) {
          error = res.message
          return
        }
      }
      closeModal()
      await loadCredentials()
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : '保存失败'
    } finally {
      saving = false
    }
  }

  async function handleDelete(id: number) {
    if (!confirm('确定要删除此凭证吗？')) return
    deletingId = id
    try {
      const res = await configApi.deleteCredential(id)
      if (res.code !== 0) {
        error = res.message
        return
      }
      await loadCredentials()
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : '删除失败'
    } finally {
      deletingId = null
    }
  }

  async function handleSetDefault(id: number) {
    try {
      const res = await configApi.setDefaultCredential(id)
      if (res.code !== 0) {
        error = res.message
        return
      }
      await loadCredentials()
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : '设置失败'
    }
  }

  async function handleTestCredential(cred: Credential) {
    testingId = cred.id
    testResult = null
    try {
      const res = await validateCredential({
        api_token: cred.api_token,
        account_id: cred.account_id
      })
      if (res.code === 0) {
        testResult = { id: cred.id, success: true, message: '凭证验证成功' }
      } else {
        testResult = { id: cred.id, success: false, message: res.message }
      }
    } catch (e: unknown) {
      testResult = { id: cred.id, success: false, message: e instanceof Error ? e.message : '验证失败' }
    } finally {
      testingId = null
    }
  }

  function handleLogout() {
    authStore.logout()
    navigate('/login')
  }

  function maskToken(token: string): string {
    if (token.length <= 8) return '****'
    return token.slice(0, 4) + '****' + token.slice(-4)
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50">
  <nav class="bg-white/80 backdrop-blur-md shadow-sm border-b border-gray-200/50 sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16 items-center">
        <div class="flex items-center space-x-3">
          <button onclick={() => navigate('/dashboard')} class="flex items-center space-x-3">
            <div class="w-10 h-10 bg-gradient-to-r from-indigo-600 to-purple-600 rounded-xl flex items-center justify-center shadow-lg">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
              </svg>
            </div>
            <div>
              <h1 class="text-xl font-bold text-gray-900">Cloudflare Tunnel</h1>
              <p class="text-xs text-gray-500">配置管理</p>
            </div>
          </button>
        </div>

        <div class="flex items-center space-x-4">
          <div class="hidden sm:flex items-center space-x-2 bg-gray-100 rounded-xl px-4 py-2">
            <div class="w-8 h-8 bg-gradient-to-r from-indigo-500 to-purple-500 rounded-lg flex items-center justify-center">
              <span class="text-white text-sm font-medium">{$currentUser?.nickname?.[0] || 'U'}</span>
            </div>
            <div class="text-sm">
              <p class="font-medium text-gray-900">{$currentUser?.nickname || '用户'}</p>
            </div>
          </div>
          <button
            onclick={handleLogout}
            class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-xl text-white bg-gradient-to-r from-red-500 to-pink-500 hover:from-red-600 hover:to-pink-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition duration-200 shadow-lg hover:shadow-xl transform hover:scale-[1.02]"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            退出登录
          </button>
        </div>
      </div>
    </div>
  </nav>

  <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div class="flex items-center justify-between mb-8">
      <div>
        <h2 class="text-2xl font-bold text-gray-900 mb-2">凭证管理</h2>
        <p class="text-gray-600">管理您的 Cloudflare API 凭证，用于隧道的创建和管理</p>
      </div>
      <button
        onclick={openCreate}
        class="inline-flex items-center px-5 py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-medium rounded-xl shadow-lg hover:shadow-xl transition duration-200 transform hover:scale-[1.02]"
      >
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
        </svg>
        添加凭证
      </button>
    </div>

    {#if error}
      <div class="mb-6 bg-red-50 border border-red-200 rounded-xl p-4 flex items-center">
        <svg class="w-5 h-5 text-red-500 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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

    {#if isLoading}
      <div class="flex items-center justify-center py-20">
        <div class="text-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto mb-4"></div>
          <p class="text-gray-500">加载中...</p>
        </div>
      </div>
    {:else if credentials.length === 0}
      <div class="bg-white rounded-2xl shadow-lg border border-gray-100 p-12 text-center">
        <div class="w-16 h-16 bg-gradient-to-r from-gray-100 to-gray-200 rounded-2xl flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 mb-2">暂无凭证</h3>
        <p class="text-gray-500 mb-6">添加您的 Cloudflare API 凭证以开始管理隧道</p>
        <button
          onclick={openCreate}
          class="inline-flex items-center px-5 py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-medium rounded-xl shadow-lg hover:shadow-xl transition duration-200"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          添加第一个凭证
        </button>
      </div>
    {:else}
      <div class="space-y-4">
        {#each credentials as cred (cred.id)}
          <div class="bg-white rounded-2xl shadow-lg border border-gray-100 p-6 hover:shadow-xl transition duration-300">
            <div class="flex items-start justify-between">
              <div class="flex-1 min-w-0">
                <div class="flex items-center space-x-3 mb-3">
                  <div class="w-10 h-10 bg-gradient-to-r from-indigo-500 to-purple-500 rounded-xl flex items-center justify-center flex-shrink-0">
                    <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                    </svg>
                  </div>
                  <div class="min-w-0">
                    <div class="flex items-center space-x-2">
                      <h3 class="text-lg font-semibold text-gray-900 truncate">{cred.name}</h3>
                      {#if cred.is_default}
                        <span class="px-2.5 py-0.5 bg-indigo-100 text-indigo-700 text-xs font-medium rounded-full">默认</span>
                      {/if}
                    </div>
                    <p class="text-sm text-gray-500">Account ID: {cred.account_id}</p>
                  </div>
                </div>

                <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 ml-13">
                  <div class="flex items-center space-x-2 text-sm">
                    <span class="text-gray-400">API Token:</span>
                    <code class="px-2 py-0.5 bg-gray-100 rounded text-gray-700 font-mono text-xs">{maskToken(cred.api_token)}</code>
                  </div>
                  <div class="flex items-center space-x-2 text-sm text-gray-500">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span>创建于 {new Date(cred.created_at).toLocaleDateString('zh-CN')}</span>
                  </div>
                </div>

                {#if testResult && testResult.id === cred.id}
                  <div class="mt-3 ml-13 {testResult.success ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'} border rounded-lg p-3 flex items-center">
                    {#if testResult.success}
                      <svg class="w-4 h-4 text-green-500 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <span class="text-green-700 text-sm">{testResult.message}</span>
                    {:else}
                      <svg class="w-4 h-4 text-red-500 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <span class="text-red-700 text-sm">{testResult.message}</span>
                    {/if}
                  </div>
                {/if}
              </div>

              <div class="flex items-center space-x-2 ml-4 flex-shrink-0">
                {#if !cred.is_default}
                  <button
                    onclick={() => handleSetDefault(cred.id)}
                    class="px-3 py-1.5 text-xs font-medium text-indigo-600 bg-indigo-50 rounded-lg hover:bg-indigo-100 transition duration-200"
                  >
                    设为默认
                  </button>
                {/if}
                <button
                  onclick={() => handleTestCredential(cred)}
                  disabled={testingId === cred.id}
                  class="px-3 py-1.5 text-xs font-medium text-green-600 bg-green-50 rounded-lg hover:bg-green-100 transition duration-200 disabled:opacity-50"
                  title="测试凭证"
                >
                  {#if testingId === cred.id}
                    <span class="inline-flex items-center">
                      <svg class="animate-spin -ml-0.5 mr-1 h-3 w-3 text-green-600" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      测试中
                    </span>
                  {:else}
                    测试
                  {/if}
                </button>
                <button
                  onclick={() => openEdit(cred)}
                  class="p-2 text-gray-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition duration-200"
                  title="编辑"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
                <button
                  onclick={() => handleDelete(cred.id)}
                  disabled={deletingId === cred.id}
                  class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition duration-200 disabled:opacity-50"
                  title="删除"
                >
                  {#if deletingId === cred.id}
                    <div class="w-5 h-5 animate-spin rounded-full border-2 border-red-400 border-t-transparent"></div>
                  {:else}
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  {/if}
                </button>
              </div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </main>
</div>

{#if showModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
    <div class="fixed inset-0 bg-black/50 backdrop-blur-sm" onclick={closeModal} role="presentation"></div>
    <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md p-6">
      <div class="flex items-center justify-between mb-6">
        <h3 class="text-lg font-semibold text-gray-900">{editingId !== null ? '编辑凭证' : '添加凭证'}</h3>
        <button onclick={closeModal} aria-label="关闭" class="p-1 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition duration-200">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form onsubmit={e => { e.preventDefault(); handleSubmit() }} class="space-y-4">
        <div>
          <label for="name" class="block text-sm font-medium text-gray-700 mb-1.5">名称</label>
          <input
            id="name"
            type="text"
            bind:value={formData.name}
            required
            placeholder="例如：生产环境"
            class="w-full px-4 py-2.5 rounded-xl border border-gray-300 shadow-sm focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition duration-200"
          />
        </div>

        <div>
          <label for="account_id" class="block text-sm font-medium text-gray-700 mb-1.5">Account ID</label>
          <input
            id="account_id"
            type="text"
            bind:value={formData.account_id}
            required
            placeholder="Cloudflare Account ID"
            class="w-full px-4 py-2.5 rounded-xl border border-gray-300 shadow-sm focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition duration-200 font-mono text-sm"
          />
        </div>

        <div>
          <label for="api_token" class="block text-sm font-medium text-gray-700 mb-1.5">API Token</label>
          <input
            id="api_token"
            type="password"
            bind:value={formData.api_token}
            required
            placeholder="Cloudflare API Token"
            class="w-full px-4 py-2.5 rounded-xl border border-gray-300 shadow-sm focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition duration-200 font-mono text-sm"
          />
        </div>

        <label class="flex items-center space-x-3 cursor-pointer">
          <input
            type="checkbox"
            bind:checked={formData.is_default}
            class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
          />
          <span class="text-sm text-gray-700">设为默认凭证</span>
        </label>

        <div class="flex justify-end space-x-3 pt-4">
          <button
            type="button"
            onclick={closeModal}
            class="px-4 py-2.5 text-sm font-medium text-gray-700 bg-gray-100 rounded-xl hover:bg-gray-200 transition duration-200"
          >
            取消
          </button>
          <button
            type="submit"
            disabled={saving}
            class="px-5 py-2.5 text-sm font-medium text-white bg-gradient-to-r from-indigo-600 to-purple-600 rounded-xl shadow-lg hover:shadow-xl transition duration-200 disabled:opacity-50"
          >
            {#if saving}
              <span class="inline-flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                保存中...
              </span>
            {:else}
              {editingId !== null ? '保存' : '添加'}
            {/if}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
