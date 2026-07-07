<script lang="ts">
  import { authStore, authLoading, authError } from '../lib/auth/store'
  import { PrimaryButton, AuthLayout, AuthInput, AuthError } from '../lib/components'
  import type { LoginRequest } from '../lib/auth/types'
  import { navigate } from 'svelte-routing'

  let formData: LoginRequest = $state({
    username: '',
    password: ''
  })

  async function handleSubmit() {
    const success = await authStore.login(formData)
    if (success) {
      navigate('/dashboard')
    }
  }

  function switchToRegister() {
    authStore.clearError()
    navigate('/register')
  }
</script>

<AuthLayout
  title="欢迎回来"
  subtitle="登录您的 Cloudflare Tunnel 管理账户"
>
  {#snippet icon()}
    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
  {/snippet}

  <AuthError error={$authError} />

  <form onsubmit={(e) => { e.preventDefault(); handleSubmit() }} class="space-y-6">
    <AuthInput
      id="username"
      name="username"
      label="用户名"
      placeholder="请输入用户名"
      required
      bind:value={formData.username}
    >
      {#snippet icon()}
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
      {/snippet}
    </AuthInput>

    <AuthInput
      id="password"
      name="password"
      type="password"
      label="密码"
      placeholder="请输入密码"
      required
      bind:value={formData.password}
    >
      {#snippet icon()}
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
      {/snippet}
    </AuthInput>

    <PrimaryButton
      type="submit"
      disabled={$authLoading}
      loading={$authLoading}
      loadingText="登录中..."
    >
      登录
    </PrimaryButton>
  </form>

  <div class="mt-6 text-center">
    <button
      type="button"
      onclick={switchToRegister}
      class="text-sm text-indigo-600 hover:text-indigo-500 font-medium transition duration-200"
    >
      没有账户？<span class="underline">立即注册</span>
    </button>
  </div>
</AuthLayout>
