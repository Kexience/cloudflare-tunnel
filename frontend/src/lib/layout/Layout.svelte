<script lang="ts">
  import { navigate } from 'svelte-routing'
  import { currentUser } from '../auth/store'
  import { authStore } from '../auth/store'

  let { title, subtitle, children }: {
    title: string
    subtitle?: string
    children: import('svelte').Snippet
  } = $props()

  function handleLogout() {
    authStore.logout()
    navigate('/login')
  }
</script>

<div class="min-h-screen bg-linear-to-br from-blue-50 via-indigo-50 to-purple-50">
  <nav class="bg-white/80 backdrop-blur-md shadow-sm border-b border-gray-200/50 sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16 items-center">
        <div class="flex items-center space-x-3">
          <button onclick={() => navigate('/dashboard')} class="flex items-center space-x-3">
            <div class="w-10 h-10 bg-linear-to-r from-indigo-600 to-purple-600 rounded-xl flex items-center justify-center shadow-lg">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
              </svg>
            </div>
            <div>
              <h1 class="text-xl font-bold text-gray-900">Cloudflare Tunnel</h1>
              <p class="text-xs text-gray-500">{subtitle || '管理系统'}</p>
            </div>
          </button>
        </div>

        <div class="flex items-center space-x-4">
          <div class="hidden sm:flex items-center space-x-2 bg-gray-100 rounded-xl px-4 py-2">
            <div class="w-8 h-8 bg-linear-to-r from-indigo-500 to-purple-500 rounded-lg flex items-center justify-center">
              <span class="text-white text-sm font-medium">{$currentUser?.nickname?.[0] || 'U'}</span>
            </div>
            <div class="text-sm">
              <p class="font-medium text-gray-900">{$currentUser?.nickname || '用户'}</p>
            </div>
          </div>
          <button
            onclick={handleLogout}
            class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-xl text-white bg-linear-to-r from-red-500 to-pink-500 hover:from-red-600 hover:to-pink-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition duration-200 shadow-lg hover:shadow-xl transform hover:scale-[1.02]"
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
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-gray-900 mb-2">{title}</h2>
      {#if subtitle}
        <p class="text-gray-600">{subtitle}</p>
      {/if}
    </div>

    {@render children()}
  </main>
</div>
