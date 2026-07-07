import { writable, derived } from 'svelte/store'
import type { AuthState, LoginRequest, RegisterRequest } from './types'
import * as authApi from './api'

function createAuthStore() {
  const storedToken = localStorage.getItem('token')
  const initialState: AuthState = {
    user: null,
    token: storedToken,
    isAuthenticated: !!storedToken,
    isLoading: false,
    error: null
  }

  const { subscribe, set, update } = writable<AuthState>(initialState)

  async function login(data: LoginRequest) {
    update(state => ({ ...state, isLoading: true, error: null }))
    try {
      const response = await authApi.login(data)
      if (response.code === 0 && response.data) {
        const { token, user } = response.data
        localStorage.setItem('token', token)
        set({
          user,
          token,
          isAuthenticated: true,
          isLoading: false,
          error: null
        })
        return true
      } else {
        update(state => ({
          ...state,
          isLoading: false,
          error: response.message
        }))
        return false
      }
    } catch (error: any) {
      update(state => ({
        ...state,
        isLoading: false,
        error: error.response?.data?.message || '登录失败'
      }))
      return false
    }
  }

  async function register(data: RegisterRequest) {
    update(state => ({ ...state, isLoading: true, error: null }))
    try {
      const response = await authApi.register(data)
      if (response.code === 0) {
        update(state => ({ ...state, isLoading: false }))
        return true
      } else {
        update(state => ({
          ...state,
          isLoading: false,
          error: response.message
        }))
        return false
      }
    } catch (error: any) {
      update(state => ({
        ...state,
        isLoading: false,
        error: error.response?.data?.message || '注册失败'
      }))
      return false
    }
  }

  async function fetchCurrentUser() {
    update(state => ({ ...state, isLoading: true, error: null }))
    try {
      const response = await authApi.getCurrentUser()
      if (response.code === 0 && response.data) {
        update(state => ({
          ...state,
          user: response.data!,
          isAuthenticated: true,
          isLoading: false
        }))
        return true
      } else {
        logout()
        return false
      }
    } catch (error) {
      logout()
      return false
    }
  }

  function logout() {
    localStorage.removeItem('token')
    set({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      error: null
    })
  }

  function clearError() {
    update(state => ({ ...state, error: null }))
  }

  return {
    subscribe,
    login,
    register,
    fetchCurrentUser,
    logout,
    clearError
  }
}

export const authStore = createAuthStore()

export const isAuthenticated = derived(authStore, $auth => $auth.isAuthenticated)
export const currentUser = derived(authStore, $auth => $auth.user)
export const authLoading = derived(authStore, $auth => $auth.isLoading)
export const authError = derived(authStore, $auth => $auth.error)