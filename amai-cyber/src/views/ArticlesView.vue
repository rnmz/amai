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
    <div v-if="postsStore.isLoading" class="system-status">
      <span class="glitch-text">// {{ t('common.loading') }}...</span>
    </div>
    <div v-else-if="postsStore.error" class="error-message">
      [SYS_ERROR] :: {{ postsStore.error }}
    </div>
    <div v-else-if="postsStore.posts.length === 0" class="system-status">
      // {{ t('article.no_articles') }}
    </div>

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
          [{{ page < 10 ? `0${page}` : page }}]
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
  padding-top: 30px;
  gap: 30px;
}

.articles-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  width: 100%;
}

.system-status {
  text-align: center;
  font-family: var(--main-font);
  font-size: 1.1rem;
  font-weight: 700;
  padding: 30px 20px;
  color: var(--md-accent);
  background-color: var(--color-surface);
  border: 1px solid var(--md-border-light);
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 10px), calc(100% - 10px) 100%, 0 100%);
  letter-spacing: 0.05em;
}

.error-message {
  color: var(--md-num-negative);
  padding: 16px;
  background-color: var(--md-mark-red-bg);
  border: 1px solid var(--md-num-negative);
  font-family: var(--main-font);
  font-size: 0.95rem;
  font-weight: 700;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 10px), calc(100% - 10px) 100%, 0 100%);
  text-align: center;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  margin-top: 10px;
  flex-wrap: wrap;
}

.pagination button {
  min-width: 42px;
  height: 42px;
  padding: 0 12px;
  border: 1px solid var(--md-border-light);
  background-color: var(--color-surface);
  color: var(--md-text-secondary);
  font-family: var(--main-font);
  font-size: 0.9rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%);
  -webkit-tap-highlight-color: transparent;
}

.pagination button:hover:not(:disabled) {
  border-color: var(--md-accent);
  color: var(--md-accent);
  background-color: rgba(0, 240, 255, 0.1);
  box-shadow: 0 0 10px rgba(0, 240, 255, 0.2);
}

.pagination button:active:not(:disabled) {
  transform: scale(0.95);
}

.pagination button:disabled {
  background-color: var(--md-accent);
  border-color: var(--md-accent);
  color: #050811;
  cursor: default;
  box-shadow: 0 0 12px var(--md-accent);
}

@media (max-width: 768px) {
  .articles-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}
</style>