<script setup lang="ts">
import { computed } from 'vue'
import { marked } from 'marked'

import rawMarkdown from '@/assets/post/test.md?raw'

interface Props {
  isAdmin?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isAdmin: true,
})

marked.setOptions({
  breaks: true,
})

// 2. Превращаем сырой текст Markdown в HTML
const parsedContent = computed(() => {
  return marked.parse(rawMarkdown)
})

function onEdit() {
  console.log('Редактировать')
}

function onDelete() {
  console.log('Удалить')
}
</script>

<template>
  <div class="post-page container">
    <header class="post-header">
      <div v-if="props.isAdmin" class="header-spacer"></div>

      <h1 class="post-title">Тестовый пост из MD</h1>

      <div v-if="props.isAdmin" class="admin-actions">
        <button class="action-btn edit-btn" title="Редактировать" @click="onEdit">
          <svg
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z" />
          </svg>
        </button>
        <button class="action-btn delete-btn" title="Удалить" @click="onDelete">
          <svg
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <polyline points="3 6 5 6 21 6" />
            <path
              d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
            />
          </svg>
        </button>
      </div>
    </header>

    <main class="post-body" v-html="parsedContent"></main>
  </div>
</template>

<style scoped>
.post-page {
  padding-top: 40px;
  padding-bottom: 60px;
}

.post-header {
  display: grid;
  grid-template-columns: 94px 1fr 94px;
  align-items: start;
  gap: 16px;
  margin-bottom: 32px;
}

.header-spacer {
  width: 94px;
}

.post-title {
  font-size: 2.2rem;
  font-weight: 400;
  color: #111;
  text-align: center;
  margin: 0;
  line-height: 1.25;
  word-break: keep-all;
  overflow-wrap: anywhere;
}

.admin-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  border: none;
  border-radius: 10px;
  color: #fff;
  cursor: pointer;
  flex-shrink: 0;
  transition:
    transform 0.1s ease,
    opacity 0.2s;
}

.action-btn:hover {
  opacity: 0.88;
}
.action-btn:active {
  transform: scale(0.95);
}
.edit-btn {
  background-color: #55f081;
}
.delete-btn {
  background-color: #ff7878;
}

.post-body {
  font-size: 1.1rem;
  line-height: 1.6;
  color: #222;
  text-align: justify;
  overflow-wrap: break-word;
}

.post-body :deep(p) {
  margin: 0 0 1.2em 0;
}

.post-body :deep(blockquote p) {
  margin: 0;
  line-height: 1.6;
  color: #2b2b2b;
  position: relative;
  z-index: 1;
}

.post-body :deep(h1),
.post-body :deep(h2),
.post-body :deep(h3) {
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  font-weight: 600;
  line-height: 1.3;
}

.post-body :deep(ul),
.post-body :deep(ol) {
  padding-left: 1.5em;
  margin-bottom: 1.2em;
}

.post-body :deep(li) {
  margin-bottom: 0.4em;
}

.post-body :deep(blockquote) {
  position: relative;
  margin: 1.5em 0;
  padding: 14px 18px;
  border-left: 4px solid #14eb51;
  border-radius: 10px;
  background-color: rgba(14, 164, 57, 0.05);
  overflow: hidden;
}

.post-body :deep(blockquote::before) {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: url('@/assets/blockquote.svg');
  background-repeat: repeat;
  opacity: 0.5;
  z-index: 0;
}

.post-body :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  display: block;
  margin: 1.5em auto;
}

@media (max-width: 768px) {
  .post-header {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .header-spacer {
    display: none;
  }
  .admin-actions {
    justify-content: center;
  }
  .post-title {
    font-size: 1.75rem;
  }
}
</style>
