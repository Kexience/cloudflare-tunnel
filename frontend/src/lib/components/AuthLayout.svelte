<script lang="ts">
  import type { Snippet } from 'svelte'

  let {
    title,
    subtitle,
    icon,
    children
  }: {
    title: string
    subtitle: string
    icon: Snippet
    children: Snippet
  } = $props()
</script>

<div class="min-h-screen bg-linear-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center p-4">
  <div class="w-full max-w-md">
    <div class="text-center mb-8">
      <div class="inline-flex items-center justify-center w-16 h-16 bg-linear-to-r from-indigo-600 to-purple-600 rounded-2xl mb-4 shadow-lg">
        <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          {@render icon()}
        </svg>
      </div>
      <h1 class="text-3xl font-bold text-gray-900 mb-2">{title}</h1>
      <p class="text-gray-600">{subtitle}</p>
    </div>

    <div class="bg-white rounded-2xl shadow-xl p-8 border border-gray-100">
      {@render children()}
    </div>

    <div class="mt-8 text-center text-sm text-gray-500">
      <p>Cloudflare Tunnel 管理系统 &copy; 2024</p>
    </div>
  </div>
</div>
