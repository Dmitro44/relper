<script setup lang="ts">
import AuthPanel from './AuthPanel.vue'
import ThemeToggle from './ThemeToggle.vue'

const emit = defineEmits<{
  back: []
}>()
</script>

<template>
  <main class="login-page">
    <div class="login-shell">
      <section class="login-intro">
        <div class="login-top">
          <div class="brand login-brand">
            <span class="brand-mark" aria-hidden="true">R</span>
            <div>
              <h1>Relper</h1>
              <span>Карта переезда</span>
            </div>
          </div>
          <ThemeToggle />
        </div>

        <div class="login-copy">
          <p class="eyebrow">Минск · демо</p>
          <h2>Места, отзывы и советы — на одной карте.</h2>
          <p class="muted">
            Войдите, чтобы оценивать места, оставлять отзывы и делиться советами.
          </p>
        </div>

        <button type="button" class="link" @click="emit('back')">Вернуться на карту</button>
      </section>

      <section class="login-box">
        <AuthPanel />
        <p class="muted login-note">Демо-админ: demo@relper.local / demo12345</p>
      </section>
    </div>
  </main>
</template>
