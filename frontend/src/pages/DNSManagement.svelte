<script lang="ts">
  import { onMount } from 'svelte'
  import { Layout } from '../lib/layout'
  import { listDNSRecords, createDNSRecord, updateDNSRecord, deleteDNSRecord } from '../lib/tunnel/api'
  import { getCredentials } from '../lib/config/api'
  import type { DNSRecordVO, CreateDNSRecordRequest, UpdateDNSRecordRequest } from '../lib/tunnel/types'
  import type { Credential } from '../lib/config/types'

  let records = $state<DNSRecordVO[]>([])
  let credentials = $state<Credential[]>([])
  let loading = $state(false)
  let error = $state<string | null>(null)
  let selectedCredentialId = $state<number | null>(null)
  let zoneId = $state('')
  let showCreateForm = $state(false)
  let editingRecord = $state<DNSRecordVO | null>(null)
  let saving = $state(false)
  let deletingId = $state<string | null>(null)

  let formData = $state<CreateDNSRecordRequest>({
    credential_id: 0,
    zone_id: '',
    name: '',
    content: '',
    proxied: false,
    ttl: 1
  })

  async function loadCredentials() {
    try {
      const response = await getCredentials()
      if (response.code === 0 && response.data) {
        credentials = response.data
        if (credentials.length > 0) {
          selectedCredentialId = credentials[0].id
        }
      }
    } catch (err) {
      console.error('Failed to load credentials:', err)
    }
  }

  async function loadDNSRecords() {
    if (!selectedCredentialId || !zoneId.trim()) return
    loading = true
    error = null
    try {
      const response = await listDNSRecords({
        credential_id: selectedCredentialId,
        zone_id: zoneId.trim()
      })
      if (response.code === 0 && response.data) {
        records = response.data
      } else {
        error = response.message || '加载DNS记录失败'
      }
    } catch (err) {
      error = '加载DNS记录失败'
      console.error('Failed to load DNS records:', err)
    } finally {
      loading = false
    }
  }

  async function handleCreateRecord() {
    if (!selectedCredentialId || !zoneId.trim()) return
    saving = true
    error = null
    try {
      const request: CreateDNSRecordRequest = {
        ...formData,
        credential_id: selectedCredentialId,
        zone_id: zoneId.trim()
      }
      const response = await createDNSRecord(request)
      if (response.code === 0) {
        showCreateForm = false
        resetForm()
        await loadDNSRecords()
      } else {
        error = response.message || '创建DNS记录失败'
      }
    } catch (err) {
      error = '创建DNS记录失败'
      console.error('Failed to create DNS record:', err)
    } finally {
      saving = false
    }
  }

  async function handleUpdateRecord() {
    if (!editingRecord || !selectedCredentialId || !zoneId.trim()) return
    saving = true
    error = null
    try {
      const request: UpdateDNSRecordRequest = {
        credential_id: selectedCredentialId,
        zone_id: zoneId.trim(),
        name: formData.name,
        content: formData.content,
        proxied: formData.proxied,
        ttl: formData.ttl
      }
      const response = await updateDNSRecord(editingRecord.id, request)
      if (response.code === 0) {
        editingRecord = null
        resetForm()
        await loadDNSRecords()
      } else {
        error = response.message || '更新DNS记录失败'
      }
    } catch (err) {
      error = '更新DNS记录失败'
      console.error('Failed to update DNS record:', err)
    } finally {
      saving = false
    }
  }

  async function handleDeleteRecord(recordId: string) {
    if (!selectedCredentialId || !zoneId.trim()) return
    if (!confirm('确定要删除这条DNS记录吗？')) return
    deletingId = recordId
    error = null
    try {
      const response = await deleteDNSRecord(recordId, {
        credential_id: selectedCredentialId,
        zone_id: zoneId.trim()
      })
      if (response.code === 0) {
        await loadDNSRecords()
      } else {
        error = response.message || '删除DNS记录失败'
      }
    } catch (err) {
      error = '删除DNS记录失败'
      console.error('Failed to delete DNS record:', err)
    } finally {
      deletingId = null
    }
  }

  function openEditForm(record: DNSRecordVO) {
    editingRecord = record
    formData = {
      credential_id: selectedCredentialId || 0,
      zone_id: zoneId,
      name: record.name,
      content: record.content,
      proxied: record.proxied || false,
      ttl: record.ttl
    }
  }

  function resetForm() {
    formData = {
      credential_id: selectedCredentialId || 0,
      zone_id: zoneId,
      name: '',
      content: '',
      proxied: false,
      ttl: 1
    }
  }

  onMount(loadCredentials)
</script>

<Layout title="DNS记录管理" subtitle="管理您的Cloudflare DNS记录">
  <div class="flex items-center justify-between mb-8">
    <div></div>
    <button
      onclick={() => showCreateForm = true}
      class="inline-flex items-center px-5 py-2.5 bg-linear-to-r from-indigo-600 to-purple-600 text-white font-medium rounded-xl shadow-lg hover:shadow-xl transition duration-200 transform hover:scale-[1.02]"
    >
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
      </svg>
      添加记录
    </button>
  </div>

  <!-- 筛选条件 -->
  <div class="bg-white rounded-2xl shadow-lg border border-gray-100 p-6 mb-6">
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div>
        <label for="credential" class="block text-sm font-medium text-gray-700 mb-2">选择凭证</label>
        <select
          id="credential"
          bind:value={selectedCredentialId}
          class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
        >
          {#each credentials as credential}
            <option value={credential.id}>{credential.name}</option>
          {/each}
        </select>
      </div>
      <div>
        <label for="zoneId" class="block text-sm font-medium text-gray-700 mb-2">Zone ID</label>
        <input
          type="text"
          id="zoneId"
          bind:value={zoneId}
          placeholder="输入Zone ID"
          class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
        />
      </div>
      <div class="flex items-end">
        <button
          onclick={loadDNSRecords}
          disabled={loading || !selectedCredentialId || !zoneId.trim()}
          class="w-full px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
        >
          {loading ? '加载中...' : '查询记录'}
        </button>
      </div>
    </div>
  </div>

  {#if error}
    <div class="mb-6 bg-red-50 border border-red-200 rounded-xl p-4 flex items-center">
      <svg class="w-5 h-5 text-red-500 mr-3 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span class="text-red-700 text-sm">{error}</span>
      <button onclick={() => error = null} aria-label="关闭" class="ml-auto text-red-400 hover:text-red-600">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  {/if}

  <!-- DNS记录列表 -->
  <div class="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden">
    <div class="px-6 py-4 border-b border-gray-200 bg-gray-50/50">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-semibold text-gray-900">DNS记录列表</h3>
        <button
          onclick={loadDNSRecords}
          class="text-sm text-indigo-600 hover:text-indigo-500 font-medium transition duration-200"
        >
          刷新
        </button>
      </div>
    </div>
    <div class="p-6">
      {#if loading}
        <div class="flex justify-center items-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
        </div>
      {:else if records.length === 0}
        <div class="text-center py-8">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">暂无DNS记录</h3>
          <p class="mt-1 text-sm text-gray-500">点击"添加记录"开始创建</p>
        </div>
      {:else}
        <div class="space-y-4">
          {#each records as record}
            <div class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition duration-200">
              <div class="flex items-center justify-between">
                <div class="flex items-center space-x-4">
                  <div class="w-10 h-10 bg-linear-to-r from-blue-500 to-cyan-500 rounded-lg flex items-center justify-center">
                    <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
                    </svg>
                  </div>
                  <div>
                    <div class="flex items-center space-x-2">
                      <h4 class="text-sm font-medium text-gray-900">{record.name}</h4>
                      <span class="px-2 py-1 text-xs font-medium rounded-full bg-blue-100 text-blue-800">
                        {record.type}
                      </span>
                      {#if record.proxied}
                        <span class="px-2 py-1 text-xs font-medium rounded-full bg-orange-100 text-orange-800">
                          已代理
                        </span>
                      {/if}
                    </div>
                    <p class="text-sm text-gray-500">{record.content}</p>
                  </div>
                </div>
                <div class="flex items-center space-x-2">
                  <div class="text-right mr-4">
                    <p class="text-xs text-gray-500">TTL: {record.ttl === 1 ? '自动' : record.ttl}</p>
                    <p class="text-xs text-gray-500">ID: {record.id}</p>
                  </div>
                  <button
                    onclick={() => openEditForm(record)}
                    class="p-2 text-gray-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition duration-200"
                    title="编辑"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button
                    onclick={() => handleDeleteRecord(record.id)}
                    disabled={deletingId === record.id}
                    class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition duration-200 disabled:opacity-50"
                    title="删除"
                  >
                    {#if deletingId === record.id}
                      <div class="w-4 h-4 animate-spin rounded-full border-2 border-red-400 border-t-transparent"></div>
                    {:else}
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    {/if}
                  </button>
                </div>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>

  <!-- 创建/编辑DNS记录模态框 -->
  {#if showCreateForm || editingRecord}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
      <div class="fixed inset-0 bg-black/50 backdrop-blur-sm" onclick={() => { showCreateForm = false; editingRecord = null }} role="presentation"></div>
      <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md p-6">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold text-gray-900">{editingRecord ? '编辑DNS记录' : '添加DNS记录'}</h3>
          <button onclick={() => { showCreateForm = false; editingRecord = null }} aria-label="关闭" class="p-1 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition duration-200">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form onsubmit={e => { e.preventDefault(); editingRecord ? handleUpdateRecord() : handleCreateRecord() }} class="space-y-4">
          <div>
            <label for="recordName" class="block text-sm font-medium text-gray-700 mb-1.5">记录名称</label>
            <input
              id="recordName"
              type="text"
              bind:value={formData.name}
              required
              placeholder="例如：www.example.com"
              class="w-full px-4 py-2.5 rounded-xl border border-gray-300 shadow-sm focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition duration-200"
            />
          </div>

          <div>
            <label for="recordType" class="block text-sm font-medium text-gray-700 mb-1.5">记录类型</label>
            <select
              id="recordType"
              class="w-full px-4 py-2.5 rounded-xl border border-gray-300 shadow-sm focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition duration-200"
            >
              <option value="A">A</option>
              <option value="AAAA">AAAA</option>
              <option value="CNAME">CNAME</option>
              <option value="MX">MX</option>
              <option value="TXT">TXT</option>
              <option value="NS">NS</option>
            </select>
          </div>

          <div>
            <label for="recordContent" class="block text-sm font-medium text-gray-700 mb-1.5">记录内容</label>
            <input
              id="recordContent"
              type="text"
              bind:value={formData.content}
              required
              placeholder="例如：192.168.1.1"
              class="w-full px-4 py-2.5 rounded-xl border border-gray-300 shadow-sm focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition duration-200"
            />
          </div>

          <div>
            <label for="recordTTL" class="block text-sm font-medium text-gray-700 mb-1.5">TTL (秒)</label>
            <input
              id="recordTTL"
              type="number"
              bind:value={formData.ttl}
              min="1"
              placeholder="1 表示自动"
              class="w-full px-4 py-2.5 rounded-xl border border-gray-300 shadow-sm focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition duration-200"
            />
          </div>

          <label class="flex items-center space-x-3 cursor-pointer">
            <input
              type="checkbox"
              bind:checked={formData.proxied}
              class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
            />
            <span class="text-sm text-gray-700">启用Cloudflare代理</span>
          </label>

          <div class="flex justify-end space-x-3 pt-4">
            <button
              type="button"
              onclick={() => { showCreateForm = false; editingRecord = null }}
              class="px-4 py-2.5 text-sm font-medium text-gray-700 bg-gray-100 rounded-xl hover:bg-gray-200 transition duration-200"
            >
              取消
            </button>
            <button
              type="submit"
              disabled={saving}
              class="px-5 py-2.5 text-sm font-medium text-white bg-linear-to-r from-indigo-600 to-purple-600 rounded-xl shadow-lg hover:shadow-xl transition duration-200 disabled:opacity-50"
            >
              {#if saving}
                <span class="inline-flex items-center">
                  <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  保存中...
                </span>
              {:else}
                {editingRecord ? '保存' : '添加'}
              {/if}
            </button>
          </div>
        </form>
      </div>
    </div>
  {/if}
</Layout>
