<script lang="ts">
  import type { Snippet } from 'svelte'

  let {
    onclick,
    gradientFrom,
    gradientTo,
    borderFrom,
    hoverBorderFrom,
    iconGradientFrom,
    iconGradientTo,
    label,
    children
  }: {
    onclick?: () => void
    gradientFrom: string
    gradientTo: string
    borderFrom: string
    hoverBorderFrom: string
    iconGradientFrom: string
    iconGradientTo: string
    label: string
    children: Snippet
  } = $props()
</script>

<button
  {onclick}
  class="flex flex-col items-center p-4 bg-linear-to-br {gradientFrom} {gradientTo} rounded-xl border {borderFrom} {hoverBorderFrom} hover:shadow-md transition duration-200"
>
  <div class="w-10 h-10 bg-linear-to-r {iconGradientFrom} {iconGradientTo} rounded-lg flex items-center justify-center mb-3">
    <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      {@render children()}
    </svg>
  </div>
  <span class="text-sm font-medium text-gray-900">{label}</span>
</button>
