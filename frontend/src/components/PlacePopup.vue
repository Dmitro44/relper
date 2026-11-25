<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import {
  formatDate,
  label,
  placeCategoryLabels,
  tipCategories,
  tipCategoryLabels,
} from '../labels'
import type { Comment, Place, Tip, TipInput } from '../types'

const props = defineProps<{
  place: Place
  comments: Comment[]
  tips: Tip[]
  authed: boolean
  rateAction: (stars: number) => Promise<void>
  commentAction: (content: string) => Promise<void>
  tipAction: (payload: TipInput) => Promise<void>
}>()

const emit = defineEmits<{
  close: []
  login: []
}>()

type Tab = 'main' | 'tips' | 'reviews'

const activeTab = ref<Tab>('main')
const hoveredStar = ref(0)
const ratingBusy = ref(false)
const ratingError = ref('')
const commentText = ref('')
const commentBusy = ref(false)
const commentError = ref('')
const tipTitle = ref('')
const tipCategory = ref('daily_life')
const tipContent = ref('')
const tipBusy = ref(false)
const tipError = ref('')

const filledStars = computed(() => {
  if (hoveredStar.value) {
    return hoveredStar.value
  }
  return Math.round(props.place.rating.avg_stars ?? 0)
})

watch(
  () => props.place.id,
  () => {
    activeTab.value = 'main'
  },
)

function votesLabel(count: number): string {
  const mod10 = count % 10
  const mod100 = count % 100
  if (mod10 === 1 && mod100 !== 11) {
    return `${count} оценка`
  }
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14)) {
    return `${count} оценки`
  }
  return `${count} оценок`
}

async function rate(stars: number): Promise<void> {
  if (ratingBusy.value) {
    return
  }

  ratingBusy.value = true
  ratingError.value = ''
  try {
    await props.rateAction(stars)
  } catch (error) {
    ratingError.value = error instanceof Error ? error.message : 'Не удалось поставить оценку'
  } finally {
    ratingBusy.value = false
  }
}

async function submitComment(): Promise<void> {
  const content = commentText.value.trim()
  if (!content) {
    commentError.value = 'Напишите текст отзыва'
    return
  }

  commentBusy.value = true
  commentError.value = ''
  try {
    await props.commentAction(content)
    commentText.value = ''
  } catch (error) {
    commentError.value = error instanceof Error ? error.message : 'Не удалось отправить отзыв'
  } finally {
    commentBusy.value = false
  }
}

async function submitTip(): Promise<void> {
  const title = tipTitle.value.trim()
  const content = tipContent.value.trim()
  if (!title || !content) {
    tipError.value = 'Заполните заголовок и текст совета'
    return
  }

  tipBusy.value = true
  tipError.value = ''
  try {
    await props.tipAction({ title, content, category: tipCategory.value })
    tipTitle.value = ''
    tipContent.value = ''
  } catch (error) {
    tipError.value = error instanceof Error ? error.message : 'Не удалось отправить совет'
  } finally {
    tipBusy.value = false
  }
}
</script>

<template>
  <article class="panel place-popup-card">
    <header class="popup-head">
      <div class="popup-title">
        <p class="eyebrow">{{ label(placeCategoryLabels, place.category) }}</p>
        <h2>{{ place.name }}</h2>
        <p class="muted">{{ place.address || 'Адрес не указан' }}</p>
      </div>

      <div class="popup-head-side">
        <div class="rating-summary">
          <strong>{{ place.rating.votes_count ? place.rating.avg_stars.toFixed(1) : '—' }}</strong>
          <span>
            {{ place.rating.votes_count ? `из 5 · ${votesLabel(place.rating.votes_count)}` : 'Пока нет оценок' }}
          </span>
        </div>
        <button type="button" class="popup-close" aria-label="Закрыть карточку" @click="emit('close')">
          ×
        </button>
      </div>
    </header>

    <nav class="popup-tabs" role="tablist" aria-label="Информация о месте">
      <button
        id="tab-main"
        type="button"
        role="tab"
        class="popup-tab"
        :class="{ active: activeTab === 'main' }"
        :aria-selected="activeTab === 'main'"
        aria-controls="panel-main"
        @click="activeTab = 'main'"
      >
        Основное
      </button>
      <button
        id="tab-tips"
        type="button"
        role="tab"
        class="popup-tab"
        :class="{ active: activeTab === 'tips' }"
        :aria-selected="activeTab === 'tips'"
        aria-controls="panel-tips"
        @click="activeTab = 'tips'"
      >
        Советы <span class="count">{{ tips.length }}</span>
      </button>
      <button
        id="tab-reviews"
        type="button"
        role="tab"
        class="popup-tab"
        :class="{ active: activeTab === 'reviews' }"
        :aria-selected="activeTab === 'reviews'"
        aria-controls="panel-reviews"
        @click="activeTab = 'reviews'"
      >
        Отзывы <span class="count">{{ comments.length }}</span>
      </button>
    </nav>

    <div class="popup-body">
      <section
        v-if="activeTab === 'main'"
        id="panel-main"
        role="tabpanel"
        aria-labelledby="tab-main"
        class="popup-section"
      >
        <p v-if="place.description" class="place-description">{{ place.description }}</p>
        <p v-else class="muted">Описание пока не добавлено.</p>

        <dl class="place-facts">
          <div>
            <dt>Категория</dt>
            <dd>{{ label(placeCategoryLabels, place.category) }}</dd>
          </div>
          <div>
            <dt>Оценок</dt>
            <dd>{{ place.rating.votes_count ? votesLabel(place.rating.votes_count) : 'Пока нет' }}</dd>
          </div>
        </dl>

        <div v-if="authed" class="rating-control">
          <span>Ваша оценка</span>
          <div class="stars" @mouseleave="hoveredStar = 0">
            <button
              v-for="star in 5"
              :key="star"
              type="button"
              class="star-btn"
              :class="{ filled: star <= filledStars }"
              :aria-label="`Оценка ${star}`"
              :disabled="ratingBusy"
              @click="rate(star)"
              @mouseenter="hoveredStar = star"
            >
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path
                  d="M12 2.5l2.9 5.9 6.5.9-4.7 4.6 1.1 6.5L12 17.3 6.2 20.4l1.1-6.5L2.6 9.3l6.5-.9L12 2.5z"
                />
              </svg>
            </button>
          </div>
          <span v-if="ratingError" class="form-error" role="alert">{{ ratingError }}</span>
        </div>
        <p v-else class="login-hint">
          Войдите, чтобы оценить место.
          <button type="button" class="link" @click="emit('login')">Войти</button>
        </p>
      </section>

      <section
        v-else-if="activeTab === 'tips'"
        id="panel-tips"
        role="tabpanel"
        aria-labelledby="tab-tips"
        class="popup-section"
      >
        <p v-if="!tips.length" class="muted">Пока нет советов.</p>
        <ul v-else class="tip-list">
          <li v-for="tip in tips" :key="tip.id">
            <div class="tip-head">
              <strong>{{ tip.title }}</strong>
              <span class="badge">{{ label(tipCategoryLabels, tip.category) }}</span>
            </div>
            <p>{{ tip.content }}</p>
          </li>
        </ul>

        <form v-if="authed" class="inline-form" @submit.prevent="submitTip">
          <label class="field">
            <span>Заголовок</span>
            <input v-model="tipTitle" maxlength="200" placeholder="Коротко о совете" />
          </label>
          <label class="field">
            <span>Категория</span>
            <select v-model="tipCategory">
              <option v-for="item in tipCategories" :key="item" :value="item">
                {{ label(tipCategoryLabels, item) }}
              </option>
            </select>
          </label>
          <label class="field">
            <span>Текст совета</span>
            <textarea
              v-model="tipContent"
              rows="3"
              maxlength="4000"
              placeholder="Что поможет другому человеку?"
            ></textarea>
          </label>
          <button class="btn primary" type="submit" :disabled="tipBusy">
            {{ tipBusy ? 'Отправка…' : 'Добавить совет' }}
          </button>
          <p v-if="tipError" class="form-error" role="alert">{{ tipError }}</p>
        </form>
        <p v-else class="login-hint">
          Войдите, чтобы добавить совет.
          <button type="button" class="link" @click="emit('login')">Войти</button>
        </p>
      </section>

      <section
        v-else
        id="panel-reviews"
        role="tabpanel"
        aria-labelledby="tab-reviews"
        class="popup-section"
      >
        <p v-if="!comments.length" class="muted">Пока нет отзывов.</p>
        <ul v-else class="comment-list">
          <li v-for="comment in comments" :key="comment.id">
            <div class="comment-head">
              <strong>{{ comment.nickname }}</strong>
              <time :datetime="comment.created_at">{{ formatDate(comment.created_at) }}</time>
            </div>
            <p>{{ comment.content }}</p>
          </li>
        </ul>

        <form v-if="authed" class="inline-form" @submit.prevent="submitComment">
          <label class="field">
            <span>Оставить отзыв</span>
            <textarea
              v-model="commentText"
              rows="3"
              maxlength="2000"
              placeholder="Что стоит знать об этом месте?"
            ></textarea>
          </label>
          <button class="btn primary" type="submit" :disabled="commentBusy">
            {{ commentBusy ? 'Отправка…' : 'Отправить отзыв' }}
          </button>
          <p v-if="commentError" class="form-error" role="alert">{{ commentError }}</p>
        </form>
        <p v-else class="login-hint">
          Войдите, чтобы оставить отзыв.
          <button type="button" class="link" @click="emit('login')">Войти</button>
        </p>
      </section>
    </div>
  </article>
</template>
