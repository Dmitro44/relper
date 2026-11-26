<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import {
  formatDate,
  label,
  placeCategoryLabels,
  tipCategoryLabels,
} from '../labels'
import type { Comment, LikeSummary, Place, Tip } from '../types'

const props = defineProps<{
  place: Place
  comments: Comment[]
  tips: Tip[]
  authed: boolean
  myRating: number | null
  rateAction: (stars: number) => Promise<void>
  commentAction: (content: string) => Promise<void>
  tipLikes: Record<string, LikeSummary>
  commentLikes: Record<string, LikeSummary>
  viewerId: string | null
  likeTipAction: (id: string, liked: boolean) => Promise<void>
  likeCommentAction: (id: string, liked: boolean) => Promise<void>
}>()

type Tab = 'main' | 'tips' | 'reviews'

const emit = defineEmits<{
  close: []
  login: []
  'add-tip': []
  'tab-change': [tab: Tab]
}>()

const activeTab = ref<Tab>('main')
const hoveredStar = ref(0)
const ratingBusy = ref(false)
const ratingError = ref('')
const commentText = ref('')
const commentBusy = ref(false)
const commentError = ref('')
const likeBusy = ref<Record<string, boolean>>({})
const likeError = ref('')

const filledStars = computed(() => {
  if (hoveredStar.value) {
    return hoveredStar.value
  }
  return props.myRating ?? 0
})

watch(
  () => props.place.id,
  () => {
    activeTab.value = 'main'
  },
)

function selectTab(tab: Tab): void {
  if (activeTab.value === tab) {
    return
  }
  activeTab.value = tab
  likeError.value = ''
  emit('tab-change', tab)
}

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

async function toggleTipLike(tip: Tip): Promise<void> {
  if (!props.authed) {
    emit('login')
    return
  }
  if (props.viewerId === tip.user_id || likeBusy.value[tip.id]) {
    return
  }

  likeBusy.value = { ...likeBusy.value, [tip.id]: true }
  likeError.value = ''
  try {
    await props.likeTipAction(tip.id, !!props.tipLikes[tip.id]?.my_like)
  } catch (error) {
    likeError.value = error instanceof Error ? error.message : 'Не удалось поставить лайк'
  } finally {
    likeBusy.value = { ...likeBusy.value, [tip.id]: false }
  }
}

async function toggleCommentLike(comment: Comment): Promise<void> {
  if (!props.authed) {
    emit('login')
    return
  }
  if (props.viewerId === comment.user_id || likeBusy.value[comment.id]) {
    return
  }

  likeBusy.value = { ...likeBusy.value, [comment.id]: true }
  likeError.value = ''
  try {
    await props.likeCommentAction(comment.id, !!props.commentLikes[comment.id]?.my_like)
  } catch (error) {
    likeError.value = error instanceof Error ? error.message : 'Не удалось поставить лайк'
  } finally {
    likeBusy.value = { ...likeBusy.value, [comment.id]: false }
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
        @click="selectTab('main')"
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
        @click="selectTab('tips')"
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
        @click="selectTab('reviews')"
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
            <div class="like-row">
              <button
                type="button"
                class="like-btn"
                :class="{ active: tipLikes[tip.id]?.my_like }"
                :disabled="viewerId === tip.user_id || likeBusy[tip.id]"
                :aria-pressed="!!tipLikes[tip.id]?.my_like"
                :title="viewerId === tip.user_id ? 'Нельзя лайкнуть свой совет' : authed ? 'Поставить лайк' : 'Войдите, чтобы поставить лайк'"
                @click="toggleTipLike(tip)"
              >
                <span aria-hidden="true">♥</span>
                {{ tipLikes[tip.id]?.count ?? 0 }}
              </button>
            </div>
          </li>
        </ul>
        <p v-if="likeError" class="form-error" role="alert">{{ likeError }}</p>

        <button v-if="authed" type="button" class="btn primary popup-add" @click="emit('add-tip')">
          Добавить совет
        </button>
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
            <div class="like-row">
              <button
                type="button"
                class="like-btn"
                :class="{ active: commentLikes[comment.id]?.my_like }"
                :disabled="viewerId === comment.user_id || likeBusy[comment.id]"
                :aria-pressed="!!commentLikes[comment.id]?.my_like"
                :title="viewerId === comment.user_id ? 'Нельзя лайкнуть свой отзыв' : authed ? 'Поставить лайк' : 'Войдите, чтобы поставить лайк'"
                @click="toggleCommentLike(comment)"
              >
                <span aria-hidden="true">♥</span>
                {{ commentLikes[comment.id]?.count ?? 0 }}
              </button>
            </div>
          </li>
        </ul>
        <p v-if="likeError" class="form-error" role="alert">{{ likeError }}</p>

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
