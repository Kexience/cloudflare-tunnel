<script lang="ts">
  import type { TunnelVO } from '../tunnel/types'

  let {
    tunnel,
    onView,
    onDelete
  }: {
    tunnel: TunnelVO
    onView: (tunnelId: string) => void
    onDelete: (tunnelId: string) => void
  } = $props()
</script>

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
        onclick={() => onView(tunnel.id)}
        class="p-2 text-blue-600 hover:text-blue-900 hover:bg-blue-50 rounded-md transition duration-200"
        title="查看详情"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
      </button>
      <button
        onclick={() => onDelete(tunnel.id)}
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
