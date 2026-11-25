<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../api/client'
import { errorMessage } from '../api/error'
import { formatDate } from '../labels'
import type { AuditEvent } from '../types'

const events = ref<AuditEvent[]>([])
const source = ref('')
const action = ref('')
const limit = ref(50)
const loading = ref(false)
const error = ref('')

async function load(): Promise<void> {
  loading.value = true
  error.value = ''

  const { data, error: requestError } = await api.GET('/audit', {
    params: {
      query: {
        source: source.value || undefined,
        action: action.value.trim() || undefined,
        limit: limit.value,
      },
    },
  })
  loading.value = false

  if (requestError || !data) {
    error.value = errorMessage(requestError)
    events.value = []
    return
  }

  events.value = data
}

function sourceLabel(value: string): string {
  return value === 'trigger' ? 'Триггер' : 'Приложение'
}

function shortId(value?: string): string {
  return value ? value.slice(0, 8) : '—'
}

onMounted(load)
</script>

<template>
  <section class="audit-view">
    <header class="audit-head">
      <div>
        <h2>Аудит действий</h2>
        <p class="muted">Триггеры показывают изменения строк, приложение — вход, регистрацию и роли.</p>
      </div>
    </header>

    <form class="panel audit-filters" @submit.prevent="load">
      <label class="field">
        <span>Источник</span>
        <select v-model="source">
          <option value="">Все</option>
          <option value="app">Приложение</option>
          <option value="trigger">Триггеры</option>
        </select>
      </label>

      <label class="field">
        <span>Действие</span>
        <input v-model="action" placeholder="comments.update" />
      </label>

      <label class="field">
        <span>Лимит</span>
        <select v-model.number="limit">
          <option :value="25">25</option>
          <option :value="50">50</option>
          <option :value="100">100</option>
          <option :value="200">200</option>
        </select>
      </label>

      <button class="btn primary" type="submit" :disabled="loading">
        {{ loading ? 'Загрузка…' : 'Показать' }}
      </button>
    </form>

    <div v-if="error" class="banner error" role="alert">{{ error }}</div>
    <div v-if="loading" class="state">Загрузка событий…</div>
    <div v-else-if="!events.length" class="panel state">Событий пока нет.</div>

    <div v-else class="panel table-wrap">
      <table class="audit-table">
        <thead>
          <tr>
            <th>Время</th>
            <th>Источник</th>
            <th>Действие</th>
            <th>Автор</th>
            <th>Объект</th>
            <th>Детали</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="event in events" :key="event.id">
            <td class="nowrap">{{ formatDate(event.created_at) }}</td>
            <td>
              <span class="badge" :class="`source-${event.source}`">{{ sourceLabel(event.source) }}</span>
            </td>
            <td><code>{{ event.action }}</code></td>
            <td>{{ event.actor_nickname || '—' }}</td>
            <td>
              {{ event.entity_type }}
              <span class="muted">· {{ shortId(event.entity_id) }}</span>
            </td>
            <td>
              <details class="details">
                <summary>JSON</summary>
                <pre>{{ JSON.stringify(event.details, null, 2) }}</pre>
              </details>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
