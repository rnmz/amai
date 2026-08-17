<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

interface PostProps {
  id?: string | number
  title?: string
  publishedAt?: string
  editedAt?: string | null
  imageUrl?: string
  isAdmin?: boolean
}

const props = withDefaults(defineProps<PostProps>(), {
  id: undefined,
  title: '...',
  publishedAt: '2025-07-17',
  editedAt: null,
  imageUrl: '',
  isAdmin: false,
})

const { t } = useI18n()

const emit = defineEmits<{
  (e: 'edit'): void
  (e: 'delete'): void
  (e: 'click', id: PostProps['id']): void
}>()

const formattedTitle = computed(() => {
  const text = props.title ?? ''
  const limit = 25

  if (text.length <= limit) {
    return text
  }
  return text.slice(0, limit).trim() + '...'
})

function handleCardClick() {
  emit('click', props.id)
}
</script>

<template>
  <div class="card" role="button" tabindex="0" @click="handleCardClick" @keydown.enter="handleCardClick">
    <div class="card-image-wrapper">
      <img :src="imageUrl" :alt="title" class="card-image" />
      <div class="card-scanline"></div>

      <div v-if="isAdmin" class="admin-actions">
        <button class="action-btn delete-btn" title="Удалить" @click.stop="emit('delete')">
          <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
          </svg>
        </button>
        <button class="action-btn edit-btn" title="Редактировать" @click.stop="emit('edit')">
          <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
          </svg>
        </button>
      </div>
    </div>

    <div class="card-content">
      <h3 class="card-title">{{ formattedTitle }}</h3>
      <div class="card-meta">
        <p>// {{ t('articles.published') }}: <span class="meta-val">{{ publishedAt }}</span></p>
        <p v-if="editedAt != null">// {{ t('articles.edited') }}: <span class="meta-val">{{ editedAt }}</span></p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  display: flex;
  align-items: center;
  border: 1px solid var(--md-border-light);
  background-color: var(--color-surface);
  color: var(--color-text);
  width: 100%;
  max-width: 520px;
  cursor: pointer;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 15px), calc(100% - 15px) 100%, 0 100%);
  transition: all 0.25s ease;
  position: relative;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
}

.card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  background-color: var(--md-accent);
  transition: background-color 0.25s ease, box-shadow 0.25s ease;
}

.card:hover {
  transform: translateY(-2px);
  border-color: var(--md-accent);
  box-shadow: 0 0 15px rgba(0, 240, 255, 0.25);
}

.card:hover::before {
  background-color: var(--md-link);
  box-shadow: 0 0 10px var(--md-link);
}

.card.is-loading {
  opacity: 0.4;
  pointer-events: none;
  filter: grayscale(0.8);
}

.card-image-wrapper {
  position: relative;
  width: 42%;
  height: 100%;
  flex-shrink: 0;
  overflow: hidden;
  aspect-ratio: 4 / 3;
  background-color: var(--md-bg-pre);
  border-right: 1px solid var(--md-border-light);
}

.card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.4s ease, filter 0.4s ease;
}

.card-scanline {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to bottom,
    rgba(255,255,255,0),
    rgba(255,255,255,0) 50%,
    rgba(0, 0, 0, 0.3) 50%,
    rgba(0, 0, 0, 0.3)
  );
  background-size: 100% 4px;
  pointer-events: none;
  opacity: 0.6;
}

.card:hover .card-image {
  transform: scale(1.08);
  filter: contrast(1.1);
}

.admin-actions {
  position: absolute;
  top: 8px;
  left: 8px;
  display: flex;
  gap: 6px;
  z-index: 2;
  padding: 3px;
  background: rgba(5, 8, 17, 0.85);
  border: 1px solid var(--md-accent);
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%);
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: 1px solid transparent;
  color: #ffffff;
  cursor: pointer;
  transition: all 0.2s ease;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 4px), calc(100% - 4px) 100%, 0 100%);
}

.action-btn:active {
  transform: scale(0.9);
}

.delete-btn {
  background-color: var(--md-num-negative);
}

.delete-btn:hover {
  background-color: #ff0055;
  box-shadow: 0 0 8px #ff0055;
}

.edit-btn {
  background-color: var(--md-accent);
  color: #050811;
}

.edit-btn:hover {
  background-color: #00f0ff;
  box-shadow: 0 0 8px #00f0ff;
}

.card-content {
  flex-grow: 1;
  padding: 16px 20px;
  text-align: left;
  min-width: 0;
}

.card-title {
  font-size: 1.05rem;
  font-weight: 700;
  margin: 0 0 10px 0;
  color: #ffffff;
  line-height: 1.3;
  font-family: var(--main-font);
  text-transform: uppercase;
  letter-spacing: 0.03em;
  overflow: hidden;
  word-break: break-word;
  overflow-wrap: anywhere;
  transition: color 0.2s ease, text-shadow 0.2s ease;
}

.card:hover .card-title {
  color: var(--md-accent);
  text-shadow: 0 0 8px var(--md-accent);
}

.card-meta p {
  margin: 4px 0;
  font-size: 0.8rem;
  color: var(--md-text-secondary);
  font-family: var(--main-font);
  font-weight: 700;
  letter-spacing: 0.03em;
}

.meta-val {
  color: var(--md-text);
}

@media (max-width: 768px) {
  .card {
    flex-direction: column;
    align-items: stretch;
  }

  .card-image-wrapper {
    width: 100%;
    height: auto;
    aspect-ratio: 16 / 9;
    border-right: none;
    border-bottom: 1px solid var(--md-border-light);
  }

  .card-content {
    padding: 14px 16px;
    width: 100%;
    box-sizing: border-box;
  }

  .card-title {
    font-size: 0.95rem;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
  }
}
</style>