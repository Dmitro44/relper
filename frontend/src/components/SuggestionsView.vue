<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../api/client'
import { errorMessage } from '../api/error'
import { formatDate, label, placeCategoryLabels } from '../labels'
import type { Suggestion, SuggestionStatus } from '../types'

const statuses: SuggestionStatus[] = ['pending', 'approved', 'rejected']
const statusLabels: Record<string, string> = {
  pending: 'На проверке',
  approved: 'Одобрено',
  rejected: 'Отклонено',
}

const items = ref<Suggestion[]>([])
const status = ref<SuggestionStatus>('pending')
const loading = ref(false)
const error = ref('')
const notice = ref('')
const busy = ref<Record<string, boolean>>({})
const rejecting = ref<string | null>(null)
const reason = ref('')
const actionError = ref('')

async function load(): Promise<void> {
  loading.value = true
  error.value = ''
  actionError.value = ''

  const { data, error: requestError } = await api.GET('/place-suggestions', {
    params: { query: { status: status.value } },
  })
  loading.value = false

  if (requestError || !data) {
    error.value = errorMessage(requestError)
    items.value = []
    return
  }

  items.value = data as Suggestion[]
}

async function approve(item: Suggestion): Promise<void> {
  busy.value = { ...busy.value, [item.id]: true }
  actionError.value = ''
  notice.value = ''
  try {
    const { error: requestError } = await api.PATCH('/place-suggestions/{id}/approve', {
      params: { path: { id: item.id } },
    })
    if (requestError) {
      actionError.value = errorMessage(requestError)
      return
    }
    items.value = items.value.filter((row) => row.id !== item.id)
    notice.value = `Место «${item.name}» добавлено на карту.`
  } finally {
    busy.value = { ...busy.value, [item.id]: false }
  }
}

function startReject(item: Suggestion): void {
  rejecting.value = item.id
  reason.value = ''
  actionError.value = ''
  notice.value = ''
}

async function reject(item: Suggestion): Promise<void> {
  const text = reason.value.trim()
  if (!text) {
    actionError.value = 'Напишите причину отказа'
    return
  }

  busy.value = { ...busy.value, [item.id]: true }
  actionError.value = ''
  try {
    const { error: requestError } = await api.PATCH('/place-suggestions/{id}/reject', {
      params: { path: { id: item.id } },
      body: { reason: text },
    })
    if (requestError) {
      actionError.value = errorMessage(requestError)
      return
    }
    items.value = items.value.filter((row) => row.id !== item.id)
    rejecting.value = null
    reason.value = ''
    notice.value = `Предложение «${item.name}» отклонено.`
  } finally {
    busy.value = { ...busy.value, [item.id]: false }
  }
}

onMounted(load)
</script>

<template>
  <section class="audit-view suggestions-view">
    <header class="audit-head">
      <div>
        <h2>Предложенные места</h2>
        <p class="muted">Одобренное место сразу появляется на карте выбранной локации.</p>
      </div>
    </header>

    <form class="panel audit-filters" @submit.prevent="load">
      <label class="field">
        <span>Статус</span>
        <select v-model="status" @change="load">
          <option v-for="item in statuses" :key="item" :value="item">
            {{ statusLabels[item] }}
          </option>
        </select>
      </label>
      <button class="btn primary" type="submit" :disabled="loading">
        {{ loading ? 'Загрузка…' : 'Показать' }}
      </button>
    </form>

    <div v-if="error" class="banner error" role="alert">{{ error }}</div>
    <div v-if="notice" class="banner success" role="status">{{ notice }}</div>
    <div v-if="actionError" class="banner error" role="alert">{{ actionError }}</div>

    <div v-if="loading" class="state">Загрузка предложений…</div>
    <div v-else-if="!items.length" class="panel state">
      {{ status === 'pending' ? 'Нет предложений на проверке.' : 'В этом статусе пусто.' }}
    </div>

    <ul v-else class="suggestion-list">
      <li v-for="item in items" :key="item.id" class="panel suggestion-card">
        <header class="suggestion-card-head">
          <div>
            <p class="eyebrow">{{ label(placeCategoryLabels, item.category) }}</p>
            <h3>{{ item.name }}</h3>
            <p class="muted">{{ item.address || 'Адрес не указан' }}</p>
          </div>
          <span class="badge" :class="`status-${item.status}`">{{ statusLabels[item.status] }}</span>
        </header>

        <p v-if="item.description" class="suggestion-text">{{ item.description }}</p>
        <dl class="place-facts">
          <div>
            <dt>Автор</dt>
            <dd>{{ item.nickname || item.user_id.slice(0, 8) }}</dd>
          </div>
          <div>
            <dt>Отправлено</dt>
            <dd>{{ formatDate(item.created_at) }}</dd>
          </div>
          <div>
            <dt>Координаты</dt>
            <dd>{{ item.coordinates.lat.toFixed(5) }}, {{ item.coordinates.lng.toFixed(5) }}</dd>
          </div>
          <div>
            <dt>Контакты</dt>
            <dd>{{ item.phone || item.website_url || '—' }}</dd>
          </div>
        </dl>

        <p v-if="item.rejection_reason" class="form-error">
          Причина отказа: {{ item.rejection_reason }}
        </p>

        <template v-if="item.status === 'pending'">
          <div v-if="rejecting === item.id" class="suggestion-reject">
            <label class="field">
              <span>Причина отказа</span>
              <input v-model="reason" maxlength="500" placeholder="Коротко, почему не подходит" />
            </label>
            <div class="suggestion-actions">
              <button
                class="btn primary small"
                type="button"
                :disabled="busy[item.id]"
                @click="reject(item)"
              >
                Отклонить
              </button>
              <button class="btn ghost small" type="button" @click="rejecting = null">
                Отмена
              </button>
            </div>
          </div>
          <div v-else class="suggestion-actions">
            <button
              class="btn primary small"
              type="button"
              :disabled="busy[item.id]"
              @click="approve(item)"
            >
              Одобрить
            </button>
            <button class="btn ghost small" type="button" @click="startReject(item)">
              Отклонить
            </button>
          </div>
        </template>
      </li>
    </ul>
  </section>
</template>
