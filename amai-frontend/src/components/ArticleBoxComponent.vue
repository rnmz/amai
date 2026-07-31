<script setup lang="ts">
import { computed } from 'vue'

interface PostProps {
  title?: string
  publishedAt?: string
  editedAt?: string
  imageUrl?: string
  isAdmin?: boolean
}

const props = withDefaults(defineProps<PostProps>(), {
  title: '異世界に転生したけれど言語の壁が高すぎるので文字の装飾から始めることにした件',
  publishedAt: '2025-07-17',
  editedAt: '2025-07-17',
  imageUrl: new URL('@/assets/img/bg.jpg', import.meta.url).href,
  isAdmin: true,
})

const formattedTitle = computed(() => {
  const text = props.title ?? ''
  const limit = 25

  if (text.length <= limit) {
    return text
  }
  return text.slice(0, limit).trim() + '...'
})

const emit = defineEmits<{
  (e: 'edit'): void
  (e: 'delete'): void
}>()
</script>

<template>
  <div class="card">
    <div class="card-image-wrapper">
      <img :src="imageUrl" :alt="title" class="card-image" />

      <div v-if="isAdmin" class="admin-actions">
        <button class="action-btn delete-btn" title="Удалить" @click="emit('delete')">
          <svg
            viewBox="0 0 24 24"
            width="20"
            height="20"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <polyline points="3 6 5 6 21 6"></polyline>
            <path
              d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
            ></path>
          </svg>
        </button>
        <button class="action-btn edit-btn" title="Редактировать" @click="emit('edit')">
          <svg
            viewBox="0 0 24 24"
            width="20"
            height="20"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
          </svg>
        </button>
      </div>
    </div>

    <div class="card-content">
      <h3 class="card-title">{{ formattedTitle }}</h3>
      <div class="card-meta">
        <p>Published: {{ publishedAt }}</p>
        <p v-if="editedAt != null">Edited: {{ editedAt }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  display: flex;
  align-items: center;
  border: 1px solid #333;
  border-radius: 28px;
  overflow: hidden;
  background-color: #fff;
  width: 100%;
  max-width: 500px;
}

.card-image-wrapper {
  position: relative;
  width: 45%;
  height: 100%;
  flex-shrink: 0;
  border-radius: 24px;
  overflow: hidden;
  aspect-ratio: 4 / 3;
}

.card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.admin-actions {
  position: absolute;
  top: 10px;
  left: 10px;
  display: flex;
  gap: 8px;
  z-index: 2;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 8px;
  color: white;
  cursor: pointer;
  transition:
    transform 0.1s ease,
    opacity 0.2s;
}

.action-btn:hover {
  opacity: 0.9;
}

.action-btn:active {
  transform: scale(0.95);
}

.delete-btn {
  background-color: #e54d58;
}

.edit-btn {
  background-color: #4cd964;
}

.card-content {
  flex-grow: 1;
  padding: 16px 20px;
  text-align: center;
  min-width: 0;
}

.card-title {
  font-size: 1.25rem;
  font-weight: 500;
  margin: 0 0 16px 0;
  color: #000;
  line-height: 1.3;
  font-family: 'Nunito', 'Noto Sans JP', sans-serif;
}

.card-meta p {
  margin: 4px 0;
  font-size: 0.85rem;
  color: #444;
}

@media (max-width: 768px) {
  .card {
    flex-direction: column;
    border-radius: 20px;
  }

  .card-image-wrapper {
    width: 100%;
    border-radius: 20px 20px 0 0;
  }

  .card-content {
    padding: 16px;
  }
}
</style>
