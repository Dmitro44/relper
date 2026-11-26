<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { auth, isAdmin, loadMe, logout } from './auth'
import AuditView from './components/AuditView.vue'
import LoginView from './components/LoginView.vue'
import MapView from './components/MapView.vue'
import ProfileView from './components/ProfileView.vue'
import SuggestionsView from './components/SuggestionsView.vue'
import ThemeToggle from './components/ThemeToggle.vue'

type View = 'map' | 'login' | 'audit' | 'profile' | 'suggestions'

function viewFromPath(path: string): View {
  if (path === '/login' || path === '/login/') {
    return 'login'
  }
  if (path === '/audit' || path === '/audit/') {
    return 'audit'
  }
  if (path === '/profile' || path === '/profile/') {
    return 'profile'
  }
  if (path === '/suggestions' || path === '/suggestions/') {
    return 'suggestions'
  }
  return 'map'
}

function pathFor(view: View): string {
  return view === 'map' ? '/' : `/${view}`
}

const view = ref<View>(viewFromPath(window.location.pathname))

function go(next: View): void {
  const path = pathFor(next)
  if (window.location.pathname !== path) {
    window.history.pushState({}, '', path)
  }
  view.value = next
}

function onPopState(): void {
  view.value = viewFromPath(window.location.pathname)
}

function signOut(): void {
  logout()
  go('map')
}

function onAuthExpired(): void {
  auth.user = null
  auth.ready = true
  auth.error = 'Сессия истекла. Войдите снова.'
  go('login')
}

onMounted(() => {
  window.addEventListener('popstate', onPopState)
  window.addEventListener('relper:auth-expired', onAuthExpired)
  void loadMe()
})

onBeforeUnmount(() => {
  window.removeEventListener('popstate', onPopState)
  window.removeEventListener('relper:auth-expired', onAuthExpired)
})

watch(
  () => [auth.ready, auth.user, isAdmin()] as const,
  () => {
    if (view.value === 'login' && auth.user) {
      go('map')
    }
    if (view.value === 'audit' && auth.ready && !isAdmin()) {
      go('map')
    }
    if (view.value === 'suggestions' && auth.ready && !isAdmin()) {
      go('map')
    }
    if (view.value === 'profile' && auth.ready && !auth.user) {
      go('login')
    }
  },
)
</script>

<template>
  <LoginView v-if="view === 'login'" @back="go('map')" />

  <section v-else-if="view === 'audit' && isAdmin()" class="audit-page">
    <header class="page-bar">
      <button type="button" class="brand brand-button" @click="go('map')">
        <span class="brand-mark" aria-hidden="true">R</span>
        <span>
          <strong>Relper</strong>
          <small>Карта переезда</small>
        </span>
      </button>
      <div class="page-actions">
        <ThemeToggle />
        <button type="button" class="btn ghost small" @click="go('map')">На карту</button>
      </div>
    </header>
    <AuditView />
  </section>

  <ProfileView v-else-if="view === 'profile' && auth.user" @back="go('map')" />

  <section v-else-if="view === 'suggestions' && isAdmin()" class="audit-page">
    <header class="page-bar">
      <button type="button" class="brand brand-button" @click="go('map')">
        <span class="brand-mark" aria-hidden="true">R</span>
        <span>
          <strong>Relper</strong>
          <small>Карта переезда</small>
        </span>
      </button>
      <div class="page-actions">
        <ThemeToggle />
        <button type="button" class="btn ghost small" @click="go('map')">На карту</button>
      </div>
    </header>
    <SuggestionsView />
  </section>

  <MapView
    v-else
    @login="go('login')"
    @audit="go('audit')"
    @suggestions="go('suggestions')"
    @profile="go('profile')"
    @logout="signOut"
  />
</template>
