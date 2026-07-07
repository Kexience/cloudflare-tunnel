<script lang="ts">
  import type { Snippet } from 'svelte'

  let {
    id,
    name,
    type = 'text',
    label,
    placeholder,
    required = false,
    value = $bindable(''),
    icon,
    error,
    oninput,
    children
  }: {
    id: string
    name: string
    type?: string
    label: string
    placeholder: string
    required?: boolean
    value?: string
    icon: Snippet
    error?: string
    oninput?: (e: Event) => void
    children?: Snippet
  } = $props()

  let showPassword = $state(false)
</script>

<div>
  <label for={id} class="block text-sm font-medium text-gray-700 mb-2">{label}</label>
  <div class="relative">
    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
      <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        {@render icon()}
      </svg>
    </div>
    <input
      {id}
      {name}
      type={type === 'password' && showPassword ? 'text' : type}
      {required}
      bind:value
      {oninput}
      class="block w-full pl-10 {type === 'password' ? 'pr-12' : 'pr-3'} py-3 border {error ? 'border-red-500' : 'border-gray-300'} rounded-xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition duration-200"
      {placeholder}
    />
    {#if type === 'password'}
      <button
        type="button"
        onclick={() => showPassword = !showPassword}
        class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 transition duration-200"
      >
        {#if showPassword}
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
          </svg>
        {:else}
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
          </svg>
        {/if}
      </button>
    {/if}
  </div>
  {#if error}
    <p class="mt-2 text-sm text-red-600">{error}</p>
  {/if}
  {#if children}
    {@render children()}
  {/if}
</div>
