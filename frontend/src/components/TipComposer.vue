<script setup lang="ts">
import { ref } from 'vue'
import { label, tipCategories, tipCategoryLabels } from '../labels'
import type { TipInput } from '../types'

const props = defineProps<{
  tipAction: (payload: TipInput) => Promise<void>
}>()

const emit = defineEmits<{
  close: []
  created: []
}>()

const title = ref('')
const category = ref('daily_life')
const content = ref('')
const busy = ref(false)
const error = ref('')

async function submit(): Promise<void> {
  const value = title.value.trim()
  const text = content.value.trim()
  if (!value || !text) {
    error.value = 'Заполните заголовок и текст совета'
    return
  }

  busy.value = true
  error.value = ''
  try {
    await props.tipAction({ title: value, content: text, category: category.value })
    title.value = ''
    content.value = ''
    emit('created')
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : 'Не удалось отправить совет'
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <section class="panel tip-composer" aria-label="Новый совет">
    <header class="tip-composer-head">
      <div>
        <p class="eyebrow">Новый совет</p>
        <h3>Чем поделиться?</h3>
      </div>
      <button
        type="button"
        class="popup-close"
        aria-label="Закрыть форму совета"
        @click="emit('close')"
      >
        ×
      </button>
    </header>

    <form class="inline-form tip-composer-form" @submit.prevent="submit">
      <label class="field">
        <span>Заголовок</span>
        <input v-model="title" maxlength="200" placeholder="Коротко о совете" />
      </label>
      <label class="field">
        <span>Категория</span>
        <select v-model="category">
          <option v-for="item in tipCategories" :key="item" :value="item">
            {{ label(tipCategoryLabels, item) }}
          </option>
        </select>
      </label>
      <label class="field">
        <span>Текст совета</span>
        <textarea
          v-model="content"
          rows="3"
          maxlength="4000"
          placeholder="Что поможет другому человеку?"
        ></textarea>
      </label>
      <button class="btn primary" type="submit" :disabled="busy">
        {{ busy ? 'Отправка…' : 'Добавить совет' }}
      </button>
      <p v-if="error" class="form-error" role="alert">{{ error }}</p>
    </form>
  </section>
</template>
