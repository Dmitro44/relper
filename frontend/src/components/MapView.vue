<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '../api/client'
import { errorMessage } from '../api/error'
import { auth, isAdmin } from '../auth'
import { label, placeCategories, placeCategoryLabels } from '../labels'
import type { Comment, Location, Place, Tip, TipInput } from '../types'
import MapCanvas from './MapCanvas.vue'
import PlacePopup from './PlacePopup.vue'
import ThemeToggle from './ThemeToggle.vue'
import TipComposer from './TipComposer.vue'

const emit = defineEmits<{
  login: []
  logout: []
  audit: []
}>()

const locations = ref<Location[]>([])
const selectedLocationId = ref('')
const places = ref<Place[]>([])
const selectedPlaceId = ref<string | null>(null)
const place = ref<Place | null>(null)
const comments = ref<Comment[]>([])
const tips = ref<Tip[]>([])
const category = ref('')
const sort = ref<'wilson' | 'avg_stars' | 'votes_count' | 'name'>('wilson')
const search = ref('')
const anchor = ref<{ x: number; y: number } | null>(null)
const tipComposerOpen = ref(false)
const loadingLocations = ref(false)
const loadingPlaces = ref(false)
const error = ref('')

const categories = placeCategories
const selectedLocation = computed(
  () => locations.value.find((item) => item.id === selectedLocationId.value) ?? null,
)
const center = computed(
  () => selectedLocation.value?.coordinates ?? { lat: 53.9006, lng: 27.559 },
)
const filteredPlaces = computed(() => {
  const query = search.value.trim().toLocaleLowerCase('ru')
  if (!query) {
    return places.value
  }

  return places.value.filter((item) =>
    [item.name, item.address, item.description].some((value) =>
      value?.toLocaleLowerCase('ru').includes(query),
    ),
  )
})
const popupStyle = computed<Record<string, string> | undefined>(() => {
  if (!anchor.value) {
    return undefined
  }
  return {
    '--anchor-x': `${anchor.value.x}px`,
    '--anchor-y': `${anchor.value.y}px`,
  }
})

function resetSelection(): void {
  selectedPlaceId.value = null
  place.value = null
  comments.value = []
  tips.value = []
  anchor.value = null
  tipComposerOpen.value = false
}

function closePlace(): void {
  resetSelection()
}

function setAnchor(point: { x: number; y: number } | null): void {
  anchor.value = point
}

async function loadLocations(): Promise<void> {
  loadingLocations.value = true
  error.value = ''

  const { data, error: requestError } = await api.GET('/locations')
  loadingLocations.value = false

  if (requestError || !data) {
    error.value = errorMessage(requestError)
    return
  }

  locations.value = data
  if (data.length && !data.some((item) => item.id === selectedLocationId.value)) {
    selectedLocationId.value = data[0].id
  }
}

async function loadPlaces(): Promise<void> {
  if (!selectedLocationId.value) {
    places.value = []
    return
  }

  loadingPlaces.value = true
  error.value = ''

  const { data, error: requestError } = await api.GET('/locations/{id}/places', {
    params: {
      path: { id: selectedLocationId.value },
      query: { category: category.value || undefined, sort: sort.value },
    },
  })
  loadingPlaces.value = false

  if (requestError || !data) {
    error.value = errorMessage(requestError)
    places.value = []
    return
  }

  places.value = data
  if (selectedPlaceId.value && !data.some((item) => item.id === selectedPlaceId.value)) {
    resetSelection()
  }
}

async function loadPlace(id: string): Promise<void> {
  error.value = ''

  const [placeResult, commentsResult, tipsResult] = await Promise.all([
    api.GET('/places/{id}', { params: { path: { id } } }),
    api.GET('/places/{id}/comments', { params: { path: { id } } }),
    api.GET('/places/{id}/tips', { params: { path: { id } } }),
  ])

  if (selectedPlaceId.value !== id) {
    return
  }

  if (placeResult.error || !placeResult.data) {
    error.value = errorMessage(placeResult.error)
    place.value = null
    return
  }

  place.value = placeResult.data
  comments.value = commentsResult.data ?? []
  tips.value = tipsResult.data ?? []
}

async function selectPlace(id: string): Promise<void> {
  if (selectedPlaceId.value === id && place.value) {
    return
  }

  selectedPlaceId.value = id
  place.value = null
  comments.value = []
  tips.value = []
  tipComposerOpen.value = false
  await loadPlace(id)
}

async function rate(stars: number): Promise<void> {
  if (!place.value) {
    return
  }

  const id = place.value.id
  const { error: requestError } = await api.PUT('/places/{id}/rating', {
    params: { path: { id } },
    body: { stars },
  })
  if (requestError) {
    throw new Error(errorMessage(requestError))
  }

  await Promise.all([loadPlace(id), loadPlaces()])
}

async function addComment(content: string): Promise<void> {
  if (!place.value) {
    return
  }

  const id = place.value.id
  const { error: requestError } = await api.POST('/comments', {
    body: { place_id: id, content },
  })
  if (requestError) {
    throw new Error(errorMessage(requestError))
  }

  await loadPlace(id)
}

async function addTip(payload: TipInput): Promise<void> {
  if (!place.value) {
    return
  }

  const id = place.value.id
  const { error: requestError } = await api.POST('/tips', {
    body: { place_id: id, ...payload },
  })
  if (requestError) {
    throw new Error(errorMessage(requestError))
  }

  await loadPlace(id)
}

async function retry(): Promise<void> {
  if (!locations.value.length) {
    await loadLocations()
  }
  await loadPlaces()
  if (selectedPlaceId.value) {
    await loadPlace(selectedPlaceId.value)
  }
}

watch(selectedLocationId, async () => {
  resetSelection()
  await loadPlaces()
})

watch([category, sort], async () => {
  if (selectedLocationId.value) {
    await loadPlaces()
  }
})

watch(filteredPlaces, (current) => {
  if (selectedPlaceId.value && !current.some((item) => item.id === selectedPlaceId.value)) {
    resetSelection()
  }
})

onMounted(loadLocations)
</script>

<template>
  <section class="map-page">
    <MapCanvas
      :places="filteredPlaces"
      :selected-id="selectedPlaceId"
      :center="center"
      @select="selectPlace"
      @anchor="setAnchor"
      @close="closePlace"
    />

    <aside class="map-sidebar" aria-label="Поиск мест">
      <div class="sidebar-head">
        <div class="brand">
          <span class="brand-mark" aria-hidden="true">R</span>
          <div>
            <h1>Relper</h1>
            <span>Карта переезда</span>
          </div>
        </div>

        <div class="sidebar-actions">
          <ThemeToggle />
          <div class="sidebar-session">
            <template v-if="auth.user">
              <span class="session-name">{{ auth.user.nickname }}</span>
              <button type="button" class="link" @click="emit('logout')">Выйти</button>
            </template>
            <button v-else type="button" class="btn small" @click="emit('login')">Войти</button>
          </div>
        </div>
      </div>

      <label class="field sidebar-search">
        <span>Поиск</span>
        <input
          v-model="search"
          type="search"
          placeholder="Название, адрес или описание"
          aria-label="Поиск мест"
        />
      </label>

      <div class="sidebar-filters">
        <label class="field">
          <span>Категория</span>
          <select v-model="category">
            <option value="">Все</option>
            <option v-for="item in categories" :key="item" :value="item">
              {{ label(placeCategoryLabels, item) }}
            </option>
          </select>
        </label>

        <label class="field">
          <span>Сортировка</span>
          <select v-model="sort">
            <option value="wilson">Рекомендации</option>
            <option value="avg_stars">Средняя оценка</option>
            <option value="votes_count">Число оценок</option>
            <option value="name">Название</option>
          </select>
        </label>
      </div>

      <label v-if="locations.length > 1" class="field">
        <span>Локация</span>
        <select v-model="selectedLocationId">
          <option v-for="location in locations" :key="location.id" :value="location.id">
            {{ location.city }}, {{ location.country }}
          </option>
        </select>
      </label>

      <div class="sidebar-meta">
        <span>{{ filteredPlaces.length }} мест</span>
        <span v-if="selectedLocation" class="muted">
          {{ selectedLocation.city }}, {{ selectedLocation.country }}
        </span>
      </div>

      <div v-if="loadingLocations || loadingPlaces" class="state">Загрузка мест…</div>
      <div v-else-if="!locations.length" class="state">Пока нет локаций.</div>
      <div v-else-if="!filteredPlaces.length" class="state">
        {{ search ? 'Ничего не найдено.' : 'В этой локации пока нет мест.' }}
      </div>
      <ul v-else class="place-list sidebar-list">
        <li v-for="item in filteredPlaces" :key="item.id">
          <button
            type="button"
            class="place-row"
            :class="{ selected: item.id === selectedPlaceId }"
            :aria-pressed="item.id === selectedPlaceId"
            @click="selectPlace(item.id)"
          >
            <span class="place-row-main">
              <strong>{{ item.name }}</strong>
              <span class="muted">{{ label(placeCategoryLabels, item.category) }}</span>
            </span>
            <span class="rating-inline">
              {{ item.rating.votes_count ? item.rating.avg_stars.toFixed(1) : '—' }}
              <span aria-hidden="true">★</span>
            </span>
          </button>
        </li>
      </ul>

      <div class="sidebar-foot">
        <span class="muted">{{ selectedLocation?.city ?? 'Минск' }}</span>
        <button v-if="isAdmin()" type="button" class="link" @click="emit('audit')">
          Аудит
        </button>
      </div>
    </aside>

    <div v-if="error" class="banner error map-error" role="alert">
      <span>{{ error }}</span>
      <button type="button" class="btn ghost small" @click="retry">Повторить</button>
    </div>

    <div
      v-if="place && anchor"
      class="place-popup"
      :class="{ 'has-composer': tipComposerOpen }"
      :style="popupStyle"
    >
      <PlacePopup
        :place="place"
        :comments="comments"
        :tips="tips"
        :authed="!!auth.user"
        :rate-action="rate"
        :comment-action="addComment"
        @add-tip="tipComposerOpen = true"
        @tab-change="tipComposerOpen = false"
        @close="closePlace"
        @login="emit('login')"
      />
      <TipComposer
        v-if="tipComposerOpen"
        :tip-action="addTip"
        @close="tipComposerOpen = false"
        @created="tipComposerOpen = false"
      />
    </div>
  </section>
</template>
