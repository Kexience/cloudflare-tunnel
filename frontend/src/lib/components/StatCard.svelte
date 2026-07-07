<script lang="ts">
  import type { Snippet } from 'svelte'

  let {
    iconGradientFrom,
    iconGradientTo,
    title,
    description,
    badge,
    children
  }: {
    iconGradientFrom: string
    iconGradientTo: string
    title: string
    description: string
    badge?: Snippet
    children: Snippet
  } = $props()
</script>

<div class="bg-white rounded-2xl shadow-lg p-6 border border-gray-100 hover:shadow-xl transition duration-300">
  <div class="flex items-center justify-between mb-4">
    <div class="w-12 h-12 bg-linear-to-r {iconGradientFrom} {iconGradientTo} rounded-xl flex items-center justify-center">
      <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        {@render children()}
      </svg>
    </div>
    {#if badge}
      {@render badge()}
    {/if}
  </div>
  <h3 class="text-lg font-semibold text-gray-900 mb-1">{title}</h3>
  <p class="text-gray-600 text-sm">{description}</p>
</div>
