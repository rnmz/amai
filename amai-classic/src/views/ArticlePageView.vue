<script setup lang="ts">
import { computed, onMounted, ref, useTemplateRef, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { marked, initCopyButtons, processFootnotes, initCharts } from '@/utils/marked-render'
import { usePostStore, usePostsStore } from '@/stores/posts'
import { useFilesStore } from '@/stores/files'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const router = useRouter()
const postStore = usePostStore()
const postsStore = usePostsStore()
const filesStore = useFilesStore()

const { t } = useI18n()

const isAdmin = computed(() => Boolean(route.meta.isAdminPage))

marked.setOptions({
  breaks: true,
})

const parsedContent = computed(() => {
  if (!postStore.post?.body) return ''
  return marked.parse(processFootnotes(postStore.post.body))
})

const contentRef = useTemplateRef<HTMLElement>('contentRef')

const deleteErrorMessage = ref<string | null>(null)

async function loadPost(id: string) {
  await postStore.fetchPost(id)
  if (contentRef.value) {
    initCopyButtons(contentRef.value)
    initCharts(contentRef.value)
  }
}

onMounted(() => {
  const id = route.params.id as string
  if (id) loadPost(id)
})

watch(
  () => route.params.id,
  (newId) => {
    if (newId) loadPost(newId as string)
  }
)

function onEdit() {
  if (!postStore.post) return
  router.push({ name: 'admin_create_article', query: { editId: postStore.post.id } })
}

async function onDelete() {
  if (!postStore.post) return

  const isConfirmed = confirm(t('article.delete'))
  if (!isConfirmed) return

  deleteErrorMessage.value = null

  const success = await postsStore.remove(postStore.post.id)
  if (success) {
    router.push({ name: 'admin_articles' })
  } else {
    deleteErrorMessage.value = postsStore.error ?? t('article.err_delete')
  }
}
</script>

<template>
  <div class="post-page container">
    <div v-if="postStore.isLoading" class="loading-state">{{ t('common.loading') }}</div>
    <div v-else-if="postStore.error" class="error-state">{{ postStore.error }}</div>

    <article v-else-if="postStore.post">
      <header class="post-header">
        <div class="header-top">
          <div v-if="isAdmin" class="header-spacer"></div>

          <h1 class="post-title">{{ postStore.post.title }}</h1>

          <div v-if="isAdmin" class="admin-actions">
            <button class="action-btn edit-btn" :title="t('common.edit')" @click="onEdit">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z" />
              </svg>
            </button>
            <button class="action-btn delete-btn" :title="t('common.delete')" @click="onDelete">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6" />
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
              </svg>
            </button>
          </div>
        </div>

        <div v-if="postStore.post.poster_id" class="post-hero">
          <img
            class="post-hero-image"
            :src="filesStore.getFileUrl(postStore.post.poster_id)"
            :alt="postStore.post.title"
          />
        </div>
      </header>

      <main ref="contentRef" class="markdown-body" v-html="parsedContent"></main>
    </article>
  </div>
</template>

<style scoped>
.post-page {
  padding-top: 40px;
  padding-bottom: 60px;
}

.loading-state,
.error-state {
  text-align: center;
  font-size: 1.2rem;
  padding: 40px 0;
  color: var(--color-text);
}

.error-state {
  color: #ff7878;
}

.post-header {
  margin-bottom: 32px;
}

.header-top {
  display: grid;
  grid-template-columns: 94px 1fr 94px;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.header-spacer {
  width: 94px;
}

.post-title {
  grid-column: 2;
  font-size: 2.2rem;
  font-weight: 400;
  color: var(--color-text);
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

.post-hero {
  width: 100%;
  border-radius: 16px;
  overflow: hidden;
}

.post-hero-image {
  width: 100%;
  max-height: 450px;
  object-fit: cover;
  display: block;
}

@media (max-width: 768px) {
  .header-top {
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

  .post-hero-image {
    max-height: 280px;
  }
}
</style>