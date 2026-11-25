<script setup lang="ts">
import { ref } from 'vue'
import { auth, login, register } from '../auth'

const mode = ref<'login' | 'register'>('login')
const email = ref('')
const password = ref('')
const nickname = ref('')

function toggle(): void {
  mode.value = mode.value === 'login' ? 'register' : 'login'
  auth.error = ''
}

async function submit(): Promise<void> {
  if (mode.value === 'login') {
    await login(email.value.trim(), password.value)
    return
  }
  await register(email.value.trim(), password.value, nickname.value.trim())
}
</script>

<template>
  <section class="panel auth-card">
    <div class="auth-head">
      <h2>{{ mode === 'login' ? 'Вход' : 'Регистрация' }}</h2>
      <p class="muted">Войдите, чтобы добавлять отзывы и советы.</p>
    </div>

    <form class="auth-form" @submit.prevent="submit">
      <label class="field">
        <span>Email</span>
        <input v-model="email" type="email" autocomplete="email" required />
      </label>

      <label v-if="mode === 'register'" class="field">
        <span>Имя</span>
        <input v-model="nickname" type="text" autocomplete="nickname" maxlength="50" required />
      </label>

      <label class="field">
        <span>Пароль</span>
        <input
          v-model="password"
          type="password"
          :autocomplete="mode === 'login' ? 'current-password' : 'new-password'"
          minlength="8"
          required
        />
      </label>

      <button class="btn primary" type="submit" :disabled="auth.busy">
        {{ auth.busy ? 'Подождите…' : mode === 'login' ? 'Войти' : 'Создать аккаунт' }}
      </button>

      <p v-if="auth.error" class="form-error" role="alert">{{ auth.error }}</p>
    </form>

    <button type="button" class="link" @click="toggle">
      {{ mode === 'login' ? 'Создать аккаунт' : 'У меня уже есть аккаунт' }}
    </button>
  </section>
</template>
