<script lang="ts">
  import { currentUser } from '../lib/auth/store'
  import { navigate } from 'svelte-routing'
  import { Layout } from '../lib/layout'
  import { QuickActionButton, StatCard } from '../lib/components'
  import { getDashboardStats } from '../lib/tunnel/api'
  import type { DashboardStatsVO } from '../lib/tunnel/types'

  let stats = $state<DashboardStatsVO | null>(null)
  let loading = $state(false)

  function formatBytes(bytes: number): string {
    if (bytes === 0) return '0 B'
    const units = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return `${(bytes / Math.pow(1024, i)).toFixed(2)} ${units[i]}`
  }

  async function loadStats() {
    loading = true
    try {
      const response = await getDashboardStats()
      if (response.code === 0 && response.data) {
        stats = response.data
      }
    } catch (err) {
      console.error('Failed to load dashboard stats:', err)
    } finally {
      loading = false
    }
  }

  $effect(() => {
    loadStats()
  })
</script>

<Layout title="欢迎回来，{$currentUser?.nickname || '用户'}！" subtitle="管理您的 Cloudflare Tunnel 配置和状态">
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
    <StatCard
      iconGradientFrom="from-green-500"
      iconGradientTo="to-emerald-500"
      title="隧道状态"
      description="所有隧道正常运行"
    >
      {#snippet badge()}
        <span class="px-3 py-1 bg-green-100 text-green-800 text-xs font-medium rounded-full">
          {stats?.running_count ?? 0} 运行中
        </span>
      {/snippet}
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
    </StatCard>

    <StatCard
      iconGradientFrom="from-blue-500"
      iconGradientTo="to-cyan-500"
      title="隧道数量"
      description="已配置的隧道总数"
    >
      {#snippet badge()}
        <span class="text-2xl font-bold text-gray-900">{stats?.total_count ?? 0}</span>
      {/snippet}
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
    </StatCard>

    <StatCard
      iconGradientFrom="from-purple-500"
      iconGradientTo="to-pink-500"
      title="入站流量"
      description="接收的数据量"
    >
      {#snippet badge()}
        <span class="text-2xl font-bold text-gray-900">{formatBytes(stats?.bytes_in ?? 0)}</span>
      {/snippet}
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3" />
    </StatCard>

    <StatCard
      iconGradientFrom="from-orange-500"
      iconGradientTo="to-red-500"
      title="出站流量"
      description="发送的数据量"
    >
      {#snippet badge()}
        <span class="text-2xl font-bold text-gray-900">{formatBytes(stats?.bytes_out ?? 0)}</span>
      {/snippet}
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18" />
    </StatCard>

    <StatCard
      iconGradientFrom="from-indigo-500"
      iconGradientTo="to-violet-500"
      title="总请求数"
      description="处理的请求总数"
    >
      {#snippet badge()}
        <span class="text-2xl font-bold text-gray-900">{stats?.total_requests ?? 0}</span>
      {/snippet}
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
    </StatCard>
  </div>

  <div class="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden">
    <div class="px-6 py-4 border-b border-gray-200 bg-gray-50/50">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-semibold text-gray-900">快速操作</h3>
        <button class="text-sm text-indigo-600 hover:text-indigo-500 font-medium transition duration-200">
          查看全部
        </button>
      </div>
    </div>
    <div class="p-6">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <QuickActionButton
          onclick={() => navigate('/tunnels')}
          gradientFrom="from-indigo-50"
          gradientTo="to-purple-50"
          borderFrom="border-indigo-100"
          hoverBorderFrom="hover:border-indigo-300"
          iconGradientFrom="from-indigo-500"
          iconGradientTo="to-purple-500"
          label="隧道列表"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
        </QuickActionButton>

        <QuickActionButton
          onclick={() => navigate('/config')}
          gradientFrom="from-green-50"
          gradientTo="to-emerald-50"
          borderFrom="border-green-100"
          hoverBorderFrom="hover:border-green-300"
          iconGradientFrom="from-green-500"
          iconGradientTo="to-emerald-500"
          label="配置管理"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </QuickActionButton>

        <QuickActionButton
          onclick={() => navigate('/dns')}
          gradientFrom="from-blue-50"
          gradientTo="to-cyan-50"
          borderFrom="border-blue-100"
          hoverBorderFrom="hover:border-blue-300"
          iconGradientFrom="from-blue-500"
          iconGradientTo="to-cyan-500"
          label="DNS管理"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
        </QuickActionButton>

        <QuickActionButton
          gradientFrom="from-orange-50"
          gradientTo="to-red-50"
          borderFrom="border-orange-100"
          hoverBorderFrom="hover:border-orange-300"
          iconGradientFrom="from-orange-500"
          iconGradientTo="to-red-500"
          label="日志查看"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
        </QuickActionButton>
      </div>
    </div>
  </div>


</Layout>
