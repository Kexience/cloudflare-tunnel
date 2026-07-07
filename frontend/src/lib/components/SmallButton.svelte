<script lang="ts">
  import type { Snippet } from 'svelte'

  let {
    onclick,
    disabled = false,
    variant = 'indigo',
    children
  }: {
    onclick?: () => void
    disabled?: boolean
    variant?: 'indigo' | 'green' | 'red' | 'gray'
    children: Snippet
  } = $props()

  const variantClasses: Record<string, string> = {
    indigo: 'text-indigo-600 bg-indigo-50 hover:bg-indigo-100',
    green: 'text-green-600 bg-green-50 hover:bg-green-100',
    red: 'text-red-600 bg-red-50 hover:bg-red-100',
    gray: 'text-gray-600 bg-gray-50 hover:bg-gray-100'
  }
</script>

<button
  type="button"
  {onclick}
  {disabled}
  class="px-3 py-1.5 text-xs font-medium rounded-lg transition duration-200 disabled:opacity-50 {variantClasses[variant]}"
>
  {@render children()}
</button>
