<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import ThemeToggle from './ThemeToggle.vue'

const props = defineProps<{
  authed: boolean
  nickname: string
}>()

const emit = defineEmits<{
  login: []
  profile: []
  logout: []
}>()

const root = ref<HTMLElement | null>(null)
const open = ref(false)

function toggle(): void {
  open.value = !open.value
}

function close(): void {
  open.value = false
}

function act(action: 'login' | 'profile' | 'logout'): void {
  close()
  if (action === 'login') {
    emit('login')
  } else if (action === 'profile') {
    emit('profile')
  } else {
    emit('logout')
  }
}

function onDocumentClick(event: MouseEvent): void {
  if (open.value && root.value && !event.composedPath().includes(root.value)) {
    close()
  }
}

function onKeydown(event: KeyboardEvent): void {
  if (event.key === 'Escape') {
    close()
  }
}

onMounted(() => {
  document.addEventListener('click', onDocumentClick)
  document.addEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onDocumentClick)
  document.removeEventListener('keydown', onKeydown)
})

const buttonLabel = () => (props.authed ? 'Меню аккаунта' : 'Войти')
</script>

<template>
  <div ref="root" class="map-account">
    <button
      type="button"
      class="account-toggle"
      :aria-label="buttonLabel()"
      :aria-expanded="open"
      :title="buttonLabel()"
      @click="toggle"
    >
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <circle cx="12" cy="8" r="3.5" />
        <path d="M4.5 20a7.5 7.5 0 0 1 15 0" />
      </svg>
    </button>

    <div v-if="open" class="panel account-popover">
      <template v-if="authed">
        <p class="account-name">{{ nickname }}</p>
        <div class="account-row">
          <button
            type="button"
            class="account-icon-btn"
            aria-label="Профиль"
            title="Профиль"
            @click="act('profile')"
          >
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <circle cx="12" cy="8" r="3.5" />
              <path d="M4.5 20a7.5 7.5 0 0 1 15 0" />
            </svg>
          </button>
          <ThemeToggle />
          <button type="button" class="link account-logout" @click="act('logout')">Выйти</button>
        </div>
      </template>
      <template v-else>
        <button type="button" class="btn primary small" @click="act('login')">Войти</button>
        <div class="account-row">
          <ThemeToggle />
        </div>
      </template>
    </div>
  </div>
</template>
