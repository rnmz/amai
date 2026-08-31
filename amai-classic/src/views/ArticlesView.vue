<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ArticleBoxComponent from '@/components/ArticleBoxComponent.vue'
import { useArticlesStore } from '@/stores/articles'
import { useFilesStore } from '@/stores/files'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const router = useRouter()
const articlesStore = useArticlesStore()
const filesStore = useFilesStore()

const { t, locale } = useI18n()

const isAdmin = computed(() => Boolean(route.meta.isAdminPage))

onMounted(() => articlesStore.fetchArticles())

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
  await articlesStore.remove(id)
}

function goToPage(page: number) {
  articlesStore.fetchArticles(page)
}
</script>

<template>
  <div class="articles_list">
    <p v-if="articlesStore.isLoading">{{ t('common.loading') }}</p>
    <p v-else-if="articlesStore.error" class="error-message">{{ articlesStore.error }}</p>
    <p v-else-if="articlesStore.articles.length === 0">{{ t('article.no_articles') }}</p>

    <template v-else>
      <div class="articles-grid">
        <ArticleBoxComponent
          v-for="article in articlesStore.articles"
          :key="article.id"
          :id="article.id"
          :title="article.title"
          :published-at="formatDate(article.created_at)"
          :edited-at="article.updated_at ? formatDate(article.updated_at) : undefined"
          :image-url="filesStore.getFileUrl(article.poster_id)"
          :is-admin="isAdmin"
          @click="onCardClick"
          @edit="onEdit(article.id)"
          @delete="onDelete(article.id)"
        />
      </div>

      <div v-if="articlesStore.pages > 1" class="pagination">
        <button
          v-for="page in articlesStore.pages"
          :key="page"
          type="button"
          :disabled="page === articlesStore.currentPage"
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
  word-break: break-word;
  overflow: hidden;
  text-overflow: ellipsis;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 12px;
  font-size: 0.8rem;
  color: #6e7781;
}

.article-date {
  white-space: nowrap;
}

.error-message {
  color: #e54d58;
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
  min-width: 36px;
  height: 36px;
  padding: 0 10px;
  border: 1px solid #d0d7de;
  border-radius: 8px;
  background-color: #fff;
  color: #1f2328;
  font-size: 0.85rem;
  cursor: pointer;
  transition:
    background-color 0.15s ease,
    border-color 0.15s ease,
    color 0.15s ease;
  -webkit-tap-highlight-color: transparent;
}

.pagination button:hover:not(:disabled) {
  border-color: #0b7a30;
  color: #0b7a30;
}

.pagination button:disabled {
  background-color: #0b7a30;
  border-color: #0b7a30;
  color: #fff;
  cursor: default;
}
</style>
