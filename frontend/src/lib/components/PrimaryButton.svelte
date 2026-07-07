<script lang="ts">
  import type { Snippet } from 'svelte'

  let {
    onclick,
    disabled = false,
    type = 'button',
    loading = false,
    loadingText = '保存中...',
    fullWidth = false,
    children
  }: {
    onclick?: () => void
    disabled?: boolean
    type?: 'button' | 'submit' | 'reset'
    loading?: boolean
    loadingText?: string
    fullWidth?: boolean
    children: Snippet
  } = $props()
</script>

<button
  {type}
  {onclick}
  {disabled}
  class="{fullWidth ? 'w-full flex justify-center' : 'inline-flex'} items-center px-5 py-2.5 bg-linear-to-r from-indigo-600 to-purple-600 text-white font-medium rounded-xl shadow-lg hover:shadow-xl transition duration-200 transform hover:scale-[1.02] disabled:opacity-50"
>
  {#if loading}
    <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
    {loadingText}
  {:else}
    {@render children()}
  {/if}
</button>
