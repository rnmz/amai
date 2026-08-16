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

      <div v-if="isAdmin" class="admin-actions">
        <button class="action-btn delete-btn" title="Удалить" @click.stop="emit('delete')">
          <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
          </svg>
        </button>
        <button class="action-btn edit-btn" title="Редактировать" @click.stop="emit('edit')">
          <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
          </svg>
        </button>
      </div>
    </div>

    <div class="card-content">
      <h3 class="card-title">{{ formattedTitle }}</h3>
      <div class="card-meta">
        <p>{{ t('articles.published') }}: {{ publishedAt }}</p>
        <p v-if="editedAt != null">{{ t('articles.edited') }}: {{ editedAt }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  display: flex;
  align-items: center;
  border: 1px solid var(--color-border);
  border-radius: 20px;
  overflow: hidden;
  background-color: var(--color-surface);
  color: var(--color-text);
  width: 100%;
  max-width: 520px;
  cursor: pointer;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.04);
  transition: transform 0.25s cubic-bezier(0.2, 0, 0, 1), 
              box-shadow 0.25s cubic-bezier(0.2, 0, 0, 1),
              opacity 0.15s ease;
}

.card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.1);
}

.card.is-loading {
  opacity: 0.6;
  pointer-events: none;
}

.card-image-wrapper {
  position: relative;
  width: 42%;
  height: 100%;
  flex-shrink: 0;
  overflow: hidden;
  aspect-ratio: 4 / 3;
}

.card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.4s ease;
}

.card:hover .card-image {
  transform: scale(1.05);
}

.admin-actions {
  position: absolute;
  top: 10px;
  left: 10px;
  display: flex;
  gap: 6px;
  z-index: 2;
  padding: 4px;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(8px);
  border-radius: 10px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  color: #ffffff;
  cursor: pointer;
  transition: transform 0.1s ease, opacity 0.2s;
}

.action-btn:hover {
  opacity: 0.9;
}

.action-btn:active {
  transform: scale(0.92);
}

.delete-btn {
  background-color: var(--md-num-negative);
}

.edit-btn {
  background-color: var(--md-link);
}

.card-content {
  flex-grow: 1;
  padding: 18px 22px;
  text-align: left;
  min-width: 0;
}

.card-title {
  font-size: 1.15rem;
  font-weight: 700;
  margin: 0 0 10px 0;
  color: var(--color-text);
  line-height: 1.4;
  font-family: var(--main-font);
  overflow: hidden;
  word-break: break-word;
  overflow-wrap: anywhere;
}

.card-meta p {
  margin: 4px 0;
  font-size: 0.85rem;
  color: var(--md-text-secondary);
  display: flex;
  align-items: center;
  gap: 6px;
}

@media (max-width: 768px) {
  .card {
    flex-direction: column;
    align-items: stretch;
    border-radius: 16px;
  }

  .card-image-wrapper {
    width: 100%;
    aspect-ratio: 16 / 9;
  }

  .card-content {
    padding: 16px;
  }
}
</style>