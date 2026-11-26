<script setup lang="ts">
import { ref } from 'vue'
import { api } from '../api/client'
import { errorMessage } from '../api/error'
import { auth } from '../auth'
import ThemeToggle from './ThemeToggle.vue'

const emit = defineEmits<{
  back: []
}>()

const nickname = ref(auth.user?.nickname ?? '')
const email = ref(auth.user?.email ?? '')
const busy = ref(false)
const error = ref('')
const saved = ref(false)

async function save(): Promise<void> {
  if (!auth.user) {
    return
  }

  busy.value = true
  error.value = ''
  saved.value = false
  try {
    const { data, error: requestError } = await api.PATCH('/users/{id}', {
      params: { path: { id: auth.user.id } },
      body: {
        email: email.value.trim(),
        nickname: nickname.value.trim(),
      },
    })
    if (requestError || !data) {
      error.value = errorMessage(requestError)
      return
    }

    auth.user = { ...auth.user, ...data }
    nickname.value = data.nickname
    email.value = data.email
    saved.value = true
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : 'Не удалось сохранить профиль'
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <section class="audit-page profile-page">
    <header class="page-bar">
      <button type="button" class="brand brand-button" @click="emit('back')">
        <span class="brand-mark" aria-hidden="true">R</span>
        <span>
          <strong>Relper</strong>
          <small>Карта переезда</small>
        </span>
      </button>
      <div class="page-actions">
        <ThemeToggle />
        <button type="button" class="btn ghost small" @click="emit('back')">На карту</button>
      </div>
    </header>

    <section class="panel auth-card profile-card">
      <header class="auth-head">
        <h2>Профиль</h2>
        <p class="muted">Никнейм и почта, которые видны в сессии.</p>
      </header>

      <form class="auth-form" @submit.prevent="save">
        <label class="field">
          <span>Никнейм</span>
          <input
            v-model="nickname"
            type="text"
            maxlength="50"
            autocomplete="nickname"
            required
          />
        </label>

        <label class="field">
          <span>Почта</span>
          <input v-model="email" type="email" autocomplete="email" required />
        </label>

        <div class="profile-actions">
          <button class="btn primary" type="submit" :disabled="busy">
            {{ busy ? 'Сохранение…' : 'Сохранить' }}
          </button>
          <button class="btn ghost" type="button" @click="emit('back')">Отмена</button>
        </div>

        <p v-if="error" class="form-error" role="alert">{{ error }}</p>
        <p v-if="saved" class="form-success" role="status">Профиль сохранён.</p>
      </form>
    </section>
  </section>
</template>
