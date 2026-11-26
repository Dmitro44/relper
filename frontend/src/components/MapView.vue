<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { api } from '../api/client'
import { errorMessage } from '../api/error'
import { auth, isAdmin } from '../auth'
import { label, placeCategories, placeCategoryLabels } from '../labels'
import type { Comment, LikeSummary, Location, Place, SuggestionInput, Tip, TipInput } from '../types'
import AccountMenu from './AccountMenu.vue'
import MapCanvas from './MapCanvas.vue'
import PlacePopup from './PlacePopup.vue'
import SuggestPlaceForm from './SuggestPlaceForm.vue'
import TipComposer from './TipComposer.vue'

const emit = defineEmits<{
  login: []
  logout: []
  audit: []
  suggestions: []
  profile: []
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
const myRating = ref<number | null>(null)
const tipLikes = ref<Record<string, LikeSummary>>({})
const commentLikes = ref<Record<string, LikeSummary>>({})
const tipComposerOpen = ref(false)
const suggestionOpen = ref(false)
const suggestionPicking = ref(false)
const suggestionPoint = ref<{ lat: number; lng: number } | null>(null)
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
  myRating.value = null
  tipLikes.value = {}
  commentLikes.value = {}
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

  const [placeResult, commentsResult, tipsResult, ratingResult] = await Promise.all([
    api.GET('/places/{id}', { params: { path: { id } } }),
    api.GET('/places/{id}/comments', { params: { path: { id } } }),
    api.GET('/places/{id}/tips', { params: { path: { id } } }),
    auth.user
      ? api.GET('/places/{id}/rating', { params: { path: { id } } })
      : Promise.resolve(null),
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
  myRating.value = ratingResult?.data?.stars ?? null

  const [nextTipLikes, nextCommentLikes] = await Promise.all([
    loadTipLikes(tips.value.map((item) => item.id)),
    loadCommentLikes(comments.value.map((item) => item.id)),
  ])
  if (selectedPlaceId.value !== id) {
    return
  }
  tipLikes.value = nextTipLikes
  commentLikes.value = nextCommentLikes
}

async function selectPlace(id: string): Promise<void> {
  if (selectedPlaceId.value === id && place.value) {
    return
  }

  selectedPlaceId.value = id
  place.value = null
  comments.value = []
  tips.value = []
  myRating.value = null
  tipLikes.value = {}
  commentLikes.value = {}
  tipComposerOpen.value = false
  await loadPlace(id)
}

async function loadTipLikes(ids: string[]): Promise<Record<string, LikeSummary>> {
  const entries = await Promise.all(
    ids.map(async (id) => {
      const result = auth.user
        ? await api.GET('/tips/{id}/likes/me', { params: { path: { id } } })
        : await api.GET('/tips/{id}/likes', { params: { path: { id } } })
      return result.data ? ([id, result.data] as const) : null
    }),
  )

  const result: Record<string, LikeSummary> = {}
  for (const entry of entries) {
    if (entry) {
      result[entry[0]] = entry[1]
    }
  }
  return result
}

async function loadCommentLikes(ids: string[]): Promise<Record<string, LikeSummary>> {
  const entries = await Promise.all(
    ids.map(async (id) => {
      const result = auth.user
        ? await api.GET('/comments/{id}/likes/me', { params: { path: { id } } })
        : await api.GET('/comments/{id}/likes', { params: { path: { id } } })
      return result.data ? ([id, result.data] as const) : null
    }),
  )

  const result: Record<string, LikeSummary> = {}
  for (const entry of entries) {
    if (entry) {
      result[entry[0]] = entry[1]
    }
  }
  return result
}

async function likeTip(id: string, liked: boolean): Promise<void> {
  const { data, error: requestError } = liked
    ? await api.DELETE('/tips/{id}/like', { params: { path: { id } } })
    : await api.PUT('/tips/{id}/like', { params: { path: { id } } })
  if (requestError) {
    throw new Error(errorMessage(requestError))
  }
  if (data) {
    tipLikes.value = { ...tipLikes.value, [id]: data }
  }
}

async function likeComment(id: string, liked: boolean): Promise<void> {
  const { data, error: requestError } = liked
    ? await api.DELETE('/comments/{id}/like', { params: { path: { id } } })
    : await api.PUT('/comments/{id}/like', { params: { path: { id } } })
  if (requestError) {
    throw new Error(errorMessage(requestError))
  }
  if (data) {
    commentLikes.value = { ...commentLikes.value, [id]: data }
  }
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

  myRating.value = stars
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

function openSuggestion(): void {
  if (!auth.user) {
    emit('login')
    return
  }
  resetSelection()
  suggestionOpen.value = false
  suggestionPicking.value = true
  suggestionPoint.value = null
}

function cancelSuggestion(): void {
  suggestionOpen.value = false
  suggestionPicking.value = false
  suggestionPoint.value = null
}

function handleMapClick(point: { lat: number; lng: number }): void {
  if (suggestionPicking.value) {
    suggestionPoint.value = point
    suggestionPicking.value = false
    suggestionOpen.value = true
    return
  }
  closePlace()
}

async function submitSuggestion(payload: SuggestionInput): Promise<void> {
  if (!selectedLocation.value) {
    throw new Error('Сначала выберите локацию')
  }

  const { error: requestError } = await api.POST('/place-suggestions', {
    body: {
      location_id: selectedLocation.value.id,
      name: payload.name,
      category: payload.category,
      address: payload.address,
      description: payload.description,
      coordinates: payload.coordinates,
    },
  })
  if (requestError) {
    throw new Error(errorMessage(requestError))
  }
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
  cancelSuggestion()
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
  <section class="map-page" :class="{ 'is-picking': suggestionPicking }">
    <MapCanvas
      :places="filteredPlaces"
      :selected-id="selectedPlaceId"
      :center="center"
      :picking="suggestionPicking"
      @select="selectPlace"
      @anchor="setAnchor"
      @map-click="handleMapClick"
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

      <button type="button" class="btn ghost small sidebar-suggest" @click="openSuggestion">
        Предложить место
      </button>

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
        <span class="sidebar-links">
          <button v-if="isAdmin()" type="button" class="link" @click="emit('suggestions')">
            Заявки
          </button>
          <button v-if="isAdmin()" type="button" class="link" @click="emit('audit')">
            Аудит
          </button>
        </span>
      </div>
    </aside>

    <AccountMenu
      :authed="!!auth.user"
      :nickname="auth.user?.nickname ?? ''"
      @login="emit('login')"
      @profile="emit('profile')"
      @logout="emit('logout')"
    />

    <div v-if="suggestionPicking" class="banner map-pick-hint" role="status">
      <span>Нажмите на карту, чтобы выбрать место.</span>
      <button type="button" class="link" @click="cancelSuggestion">Отмена</button>
    </div>

    <SuggestPlaceForm
      v-if="suggestionOpen && selectedLocation && auth.user && suggestionPoint"
      :location="selectedLocation"
      :coordinates="suggestionPoint"
      :action="submitSuggestion"
      @close="cancelSuggestion"
    />

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
        :my-rating="myRating"
        :rate-action="rate"
        :comment-action="addComment"
        :tip-likes="tipLikes"
        :comment-likes="commentLikes"
        :viewer-id="auth.user?.id ?? null"
        :like-tip-action="likeTip"
        :like-comment-action="likeComment"
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
