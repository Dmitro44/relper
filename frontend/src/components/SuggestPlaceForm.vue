<script setup lang="ts">
import { ref } from 'vue'
import { label, placeCategories, placeCategoryLabels } from '../labels'
import type { Location, SuggestionInput } from '../types'

const props = defineProps<{
  location: Location
  coordinates: { lat: number; lng: number }
  action: (payload: SuggestionInput) => Promise<void>
}>()

const emit = defineEmits<{
  close: []
}>()

const name = ref('')
const category = ref(placeCategories[0])
const address = ref('')
const description = ref('')
const busy = ref(false)
const error = ref('')
const done = ref(false)

async function submit(): Promise<void> {
  const value = name.value.trim()
  if (!value) {
    error.value = 'Укажите название места'
    return
  }

  busy.value = true
  error.value = ''
  done.value = false
  try {
    await props.action({
      name: value,
      category: category.value,
      address: address.value.trim(),
      description: description.value.trim(),
      coordinates: { lat: props.coordinates.lat, lng: props.coordinates.lng },
    })
    done.value = true
    name.value = ''
    address.value = ''
    description.value = ''
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : 'Не удалось отправить предложение'
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <div class="suggestion-backdrop" @click.self="emit('close')">
    <section class="panel suggestion-modal" aria-label="Предложить место">
      <header class="suggestion-head">
        <div>
          <p class="eyebrow">Новое место</p>
          <h2>Предложить место</h2>
          <p class="muted">{{ location.city }}, {{ location.country }}</p>
        </div>
        <button
          type="button"
          class="popup-close"
          aria-label="Закрыть форму предложения"
          @click="emit('close')"
        >
          ×
        </button>
      </header>

      <form class="inline-form" @submit.prevent="submit">
        <label class="field">
          <span>Название</span>
          <input v-model="name" maxlength="200" placeholder="Например, кофейня у метро" />
        </label>

        <label class="field">
          <span>Категория</span>
          <select v-model="category">
            <option v-for="item in placeCategories" :key="item" :value="item">
              {{ label(placeCategoryLabels, item) }}
            </option>
          </select>
        </label>

        <label class="field">
          <span>Адрес</span>
          <input v-model="address" maxlength="300" placeholder="Улица и дом" />
        </label>

        <label class="field">
          <span>Описание</span>
          <textarea
            v-model="description"
            rows="3"
            maxlength="2000"
            placeholder="Чем это место полезно?"
          ></textarea>
        </label>

        <div class="suggestion-point">
          <span class="muted">Точка на карте</span>
          <strong>{{ coordinates.lat.toFixed(5) }}, {{ coordinates.lng.toFixed(5) }}</strong>
        </div>

        <div class="suggestion-actions">
          <button class="btn primary" type="submit" :disabled="busy">
            {{ busy ? 'Отправка…' : 'Отправить' }}
          </button>
          <button class="btn ghost" type="button" @click="emit('close')">Закрыть</button>
        </div>

        <p v-if="error" class="form-error" role="alert">{{ error }}</p>
        <p v-if="done" class="form-success" role="status">
          Предложение отправлено. Администратор проверит его и добавит место на карту.
        </p>
      </form>
    </section>
  </div>
</template>
