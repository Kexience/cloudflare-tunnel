<script lang="ts">
  import { Router, Route, navigate } from 'svelte-routing'
  import { onMount } from 'svelte'
  import { authStore, isAuthenticated } from './lib/auth/store'
  import Login from './pages/Login.svelte'
  import Register from './pages/Register.svelte'
  import Dashboard from './pages/Dashboard.svelte'
  import Config from './pages/Config.svelte'
  import ProtectedRoute from './lib/auth/ProtectedRoute.svelte'

  export let url = ''

  onMount(() => {
    const token = localStorage.getItem('token')
    if (token && window.location.pathname === '/') {
      authStore.fetchCurrentUser().then(success => {
        if (success) {
          navigate('/dashboard')
        } else {
          navigate('/login')
        }
      })
    } else if (!token && window.location.pathname !== '/register') {
      navigate('/login')
    }
  })
</script>

<Router {url}>
  <div>
    <Route path="/login">
      <Login />
    </Route>
    <Route path="/register">
      <Register />
    </Route>
    <Route path="/dashboard">
      <ProtectedRoute>
        <Dashboard />
      </ProtectedRoute>
    </Route>
    <Route path="/config">
      <ProtectedRoute>
        <Config />
      </ProtectedRoute>
    </Route>
    <Route path="/">
      <div class="min-h-screen flex items-center justify-center">
        <div class="animate-spin rounded-full h-32 w-32 border-b-2 border-indigo-600"></div>
      </div>
    </Route>
  </div>
</Router>