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
      <span>{{ t('common.loading') }}...</span>
    </div>
    <div v-else-if="postsStore.error" class="error-message">
      {{ postsStore.error }}
    </div>
    <div v-else-if="postsStore.posts.length === 0" class="system-status">
      {{ t('article.no_articles') }}
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
  padding-top: 24px;
  gap: 24px;
}

.articles-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  width: 100%;
}

.system-status {
  text-align: center;
  font-family: var(--main-font);
  font-size: 1rem;
  font-weight: 600;
  padding: 24px 20px;
  color: var(--md-accent);
  background-color: var(--color-surface);
  border: 1px solid var(--md-border);
  border-radius: 8px;
  box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.08);
}

.error-message {
  color: var(--md-num-negative);
  padding: 16px;
  background-color: var(--md-bg-disclaimer);
  border: 1px solid var(--md-num-negative);
  border-radius: 8px;
  font-family: var(--main-font);
  font-size: 0.95rem;
  font-weight: 600;
  text-align: center;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  margin-top: 12px;
  flex-wrap: wrap;
}

.pagination button {
  min-width: 36px;
  height: 36px;
  padding: 0 10px;
  border: 1px solid var(--md-border);
  border-radius: 4px;
  background-color: var(--color-surface);
  color: var(--md-text);
  font-family: var(--main-font);
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  box-shadow: 1px 1px 2px rgba(0, 0, 0, 0.05);
  -webkit-tap-highlight-color: transparent;
}

.pagination button:hover:not(:disabled) {
  border-color: var(--md-accent);
  color: var(--md-accent);
  background-color: var(--md-bg-table-row-hover);
}

.pagination button:active:not(:disabled) {
  transform: scale(0.96);
}

.pagination button:disabled {
  background-color: var(--md-accent);
  border-color: var(--md-accent);
  color: #ffffff;
  cursor: default;
}

@media (max-width: 768px) {
  .articles-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}
</style>