<script lang="ts">
  import { onMount } from 'svelte'
  import { navigate } from 'svelte-routing'
  import { authStore, isAuthenticated } from './store'

  let { children }: { children: any } = $props()

  let isLoading = $state(true)

  onMount(async () => {
    const token = localStorage.getItem('token')
    if (!token) {
      navigate('/login')
      return
    }

    const success = await authStore.fetchCurrentUser()
    if (!success) {
      navigate('/login')
    }
    isLoading = false
  })
</script>

{#if isLoading}
  <div class="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center">
    <div class="text-center">
      <div class="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-r from-indigo-600 to-purple-600 rounded-2xl mb-4 shadow-lg animate-pulse">
        <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
        </svg>
      </div>
      <h2 class="text-xl font-semibold text-gray-900 mb-2">正在加载</h2>
      <p class="text-gray-600">请稍候，正在验证您的身份...</p>
      <div class="mt-4">
        <div class="w-12 h-1 bg-gradient-to-r from-indigo-600 to-purple-600 rounded-full mx-auto animate-pulse"></div>
      </div>
    </div>
  </div>
{:else if $isAuthenticated}
  {@render children()}
{/if}