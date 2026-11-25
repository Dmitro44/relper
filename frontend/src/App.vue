<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { auth, isAdmin, loadMe, logout } from './auth'
import AuditView from './components/AuditView.vue'
import LoginView from './components/LoginView.vue'
import MapView from './components/MapView.vue'
import ThemeToggle from './components/ThemeToggle.vue'

type View = 'map' | 'login' | 'audit'

function viewFromPath(path: string): View {
  if (path === '/login' || path === '/login/') {
    return 'login'
  }
  if (path === '/audit' || path === '/audit/') {
    return 'audit'
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

onMounted(() => {
  window.addEventListener('popstate', onPopState)
  void loadMe()
})

onBeforeUnmount(() => {
  window.removeEventListener('popstate', onPopState)
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

  <MapView v-else @login="go('login')" @audit="go('audit')" @logout="signOut" />
</template>
