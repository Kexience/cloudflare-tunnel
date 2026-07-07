<script lang="ts">
  let {
    error,
    ondismiss
  }: {
    error: string | null
    ondismiss?: () => void
  } = $props()
</script>

{#if error}
  <div class="mb-6 p-4 bg-red-50 border border-red-200 rounded-xl">
    <div class="flex items-center">
      <svg class="w-5 h-5 text-red-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span class="text-red-700 text-sm">{error}</span>
    </div>
  </div>
{/if}
