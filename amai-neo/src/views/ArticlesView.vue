<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ArticleBoxComponent from '@/components/ArticleBoxComponent.vue'
import { usePostsStore } from '@/stores/posts'
import { useFilesStore } from '@/stores/files'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const router = useRouter()
const postsStore = usePostsStore()
const filesStore = useFilesStore()

const { t, locale } = useI18n()

const isAdmin = computed(() => Boolean(route.meta.isAdminPage))

onMounted(() => postsStore.fetchPosts())

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(locale.value)
}

function onCardClick(id: string | number | undefined) {
  if (!id) return
  const name = isAdmin.value ? 'admin_article' : 'article'
  router.push({ name, params: { id: String(id) } })
}

function onEdit(id: string) {
  router.push({ name: 'admin_create_article', query: { editId: id } })
}

async function onDelete(id: string) {
  const isConfirmed = confirm(t('article.delete'))
  if (!isConfirmed) return
  await postsStore.remove(id)
}

function goToPage(page: number) {
  postsStore.fetchPosts(page)
}
</script>

<template>
  <div class="articles_list">
    <p v-if="postsStore.isLoading">{{ t('common.loading') }}</p>
    <p v-else-if="postsStore.error" class="error-message">{{ postsStore.error }}</p>
    <p v-else-if="postsStore.posts.length === 0">{{ t('article.no_articles') }}</p>

    <template v-else>
      <div class="articles-grid">
        <ArticleBoxComponent
          v-for="post in postsStore.posts"
          :key="post.id"
          :id="post.id"
          :title="post.title"
          :published-at="formatDate(post.created_at)"
          :edited-at="post.updated_at ? formatDate(post.updated_at) : undefined"
          :image-url="filesStore.getFileUrl(post.poster_id)"
          :is-admin="isAdmin"
          @click="onCardClick"
          @edit="onEdit(post.id)"
          @delete="onDelete(post.id)"
        />
      </div>

      <div v-if="postsStore.pages > 1" class="pagination">
        <button
          v-for="page in postsStore.pages"
          :key="page"
          type="button"
          :disabled="page === postsStore.currentPage"
          @click="goToPage(page)"
        >
          {{ page }}
        </button>
      </div>
    </template>
  </div>
</template>

<style scoped>
.articles_list {
  width: 100%;
  display: flex;
  flex-direction: column;
  padding-top: 40px;
  gap: 30px;
}

.articles-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 30px;
  width: 100%;
}

.article-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
}

.article-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--color-text);
  word-break: break-word;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.2s ease;
}

.article-card:hover .article-title {
  color: var(--md-accent);
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 12px;
  font-size: 0.85rem;
  color: var(--md-text-secondary);
  font-weight: 500;
}

.article-date {
  white-space: nowrap;
}

.error-message {
  color: var(--md-num-negative);
  padding: 12px 16px;
  border-radius: 10px;
  background-color: var(--md-mark-red-bg);
  border: 1px solid rgba(239, 68, 68, 0.3);
  font-weight: 500;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  margin-top: 10px;
  flex-wrap: wrap;
}

.pagination button {
  min-width: 38px;
  height: 38px;
  padding: 0 12px;
  border: 1px solid var(--md-border);
  border-radius: 8px;
  background-color: var(--color-surface);
  color: var(--color-text);
  font-family: var(--main-font);
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  -webkit-tap-highlight-color: transparent;
}

.pagination button:hover:not(:disabled) {
  border-color: var(--md-accent);
  color: var(--md-accent);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.15);
}

.pagination button:active:not(:disabled) {
  transform: scale(0.95);
}

.pagination button:disabled {
  background-color: var(--md-accent);
  border-color: var(--md-accent);
  color: #ffffff;
  cursor: default;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.25);
}

@media (max-width: 768px) {
  .articles-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
}
</style>