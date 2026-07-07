<script lang="ts">
  import { authStore, authLoading, authError } from '../lib/auth/store'
  import { PrimaryButton, AuthLayout, AuthInput, AuthError } from '../lib/components'
  import type { RegisterRequest } from '../lib/auth/types'
  import { navigate } from 'svelte-routing'

  let formData: RegisterRequest = $state({
    nickname: '',
    username: '',
    email: '',
    password: ''
  })

  let confirmPassword = $state('')
  let passwordError = $state('')

  async function handleSubmit() {
    if (formData.password !== confirmPassword) {
      passwordError = '两次输入的密码不一致'
      return
    }
    passwordError = ''
    const success = await authStore.register(formData)
    if (success) {
      navigate('/login')
    }
  }

  function switchToLogin() {
    authStore.clearError()
    navigate('/login')
  }

  function checkPasswordMatch() {
    if (confirmPassword && formData.password !== confirmPassword) {
      passwordError = '两次输入的密码不一致'
    } else {
      passwordError = ''
    }
  }
</script>

<AuthLayout
  title="创建账户"
  subtitle="注册您的 Cloudflare Tunnel 管理账户"
>
  {#snippet icon()}
    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
  {/snippet}

  <AuthError error={$authError} />

  <form onsubmit={(e) => { e.preventDefault(); handleSubmit() }} class="space-y-5">
    <AuthInput
      id="nickname"
      name="nickname"
      label="昵称"
      placeholder="请输入昵称"
      required
      bind:value={formData.nickname}
    >
      {#snippet icon()}
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.121 17.804A13.937 13.937 0 0112 16c2.5 0 4.847.655 6.879 1.804M15 10a3 3 0 11-6 0 3 3 0 016 0zm6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
      {/snippet}
    </AuthInput>

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
      id="email"
      name="email"
      type="email"
      label="邮箱"
      placeholder="请输入邮箱"
      required
      bind:value={formData.email}
    >
      {#snippet icon()}
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
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
      oninput={checkPasswordMatch}
    >
      {#snippet icon()}
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
      {/snippet}
    </AuthInput>

    <AuthInput
      id="confirm-password"
      name="confirm-password"
      type="password"
      label="确认密码"
      placeholder="请再次输入密码"
      required
      bind:value={confirmPassword}
      oninput={checkPasswordMatch}
      error={passwordError}
    >
      {#snippet icon()}
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
      {/snippet}
    </AuthInput>

    <PrimaryButton
      type="submit"
      disabled={$authLoading}
      loading={$authLoading}
      loadingText="注册中..."
      fullWidth
    >
      注册
    </PrimaryButton>
  </form>

  <div class="mt-6 text-center">
    <button
      type="button"
      onclick={switchToLogin}
      class="text-sm text-indigo-600 hover:text-indigo-500 font-medium transition duration-200"
    >
      已有账户？<span class="underline">立即登录</span>
    </button>
  </div>
</AuthLayout>
