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
          <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
          </svg>
        </button>
        <button class="action-btn edit-btn" title="Редактировать" @click.stop="emit('edit')">
          <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
          </svg>
        </button>
      </div>
    </div>

    <div class="card-content">
      <h3 class="card-title">{{ formattedTitle }}</h3>
      <div class="card-meta">
        <p>{{ t('articles.published') }}: <span class="meta-val">{{ publishedAt }}</span></p>
        <p v-if="editedAt != null">{{ t('articles.edited') }}: <span class="meta-val">{{ editedAt }}</span></p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  display: flex;
  align-items: center;
  border: 1px solid var(--md-border);
  border-radius: 6px;
  background-color: var(--color-surface);
  color: var(--color-text);
  width: 100%;
  max-width: 520px;
  cursor: pointer;
  transition: all 0.15s ease;
  position: relative;
  box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.1);
}

.card:hover {
  border-color: var(--md-accent);
  background-color: var(--md-bg-table-row-hover);
  box-shadow: 2px 2px 6px rgba(0, 0, 0, 0.15);
}

.card.is-loading {
  opacity: 0.5;
  pointer-events: none;
}

.card-image-wrapper {
  position: relative;
  width: 40%;
  height: 100%;
  flex-shrink: 0;
  overflow: hidden;
  aspect-ratio: 4 / 3;
  background-color: var(--md-bg-pre);
  border-right: 1px solid var(--md-border-light);
  border-radius: 5px 0 0 5px;
}

.card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: opacity 0.2s ease;
}

.card:hover .card-image {
  opacity: 0.9;
}

.admin-actions {
  position: absolute;
  top: 6px;
  left: 6px;
  display: flex;
  gap: 4px;
  z-index: 2;
  padding: 3px;
  background: var(--color-surface);
  border: 1px solid var(--md-border);
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border: 1px solid var(--md-border);
  border-radius: 3px;
  background-color: var(--md-bg-code);
  color: var(--md-text);
  cursor: pointer;
  transition: all 0.15s ease;
}

.delete-btn {
  color: var(--md-num-negative);
  border-color: var(--md-border-light);
}

.delete-btn:hover {
  background-color: var(--md-num-negative);
  border-color: var(--md-num-negative);
  color: #ffffff;
}

.edit-btn {
  color: var(--md-text);
  border-color: var(--md-border-light);
}

.edit-btn:hover {
  background-color: var(--md-accent);
  border-color: var(--md-accent);
  color: #ffffff;
}

.card-content {
  flex-grow: 1;
  padding: 14px 18px;
  text-align: left;
  min-width: 0;
}

.card-title {
  font-size: 1rem;
  font-weight: 700;
  margin: 0 0 8px 0;
  color: var(--md-text);
  line-height: 1.35;
  font-family: var(--main-font);
  overflow: hidden;
  word-break: break-word;
  overflow-wrap: anywhere;
  transition: color 0.15s ease;
}

.card:hover .card-title {
  color: var(--md-link);
}

.card-meta p {
  margin: 3px 0;
  font-size: 0.8rem;
  color: var(--md-text-secondary);
  font-family: var(--main-font);
  font-weight: 600;
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
    border-radius: 5px 5px 0 0;
  }

  .card-content {
    padding: 12px 14px;
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