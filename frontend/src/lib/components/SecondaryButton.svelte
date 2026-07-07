<script lang="ts">
  let {
    onclick,
    disabled = false,
    text = '取消'
  }: {
    onclick?: () => void
    disabled?: boolean
    text?: string
  } = $props()
</script>

<button
  type="button"
  {onclick}
  {disabled}
  class="px-4 py-2.5 text-sm font-medium text-gray-700 bg-gray-100 rounded-xl hover:bg-gray-200 transition duration-200 disabled:opacity-50"
>
  {text}
</button>
