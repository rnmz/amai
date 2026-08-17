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
    <div v-if="postStore.isLoading" class="loading-state">
      <span>{{ t('common.loading') }}...</span>
    </div>
    <div v-else-if="postStore.error" class="error-state">
      {{ postStore.error }}
    </div>

    <article v-else-if="postStore.post" class="post-article">
      <header class="post-header">
        <div class="header-top">
          <div v-if="isAdmin" class="header-spacer"></div>

          <h1 class="post-title">{{ postStore.post.title }}</h1>

          <div v-if="isAdmin" class="admin-actions">
            <button class="action-btn edit-btn" :title="t('common.edit')" @click="onEdit">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z" />
              </svg>
            </button>
            <button class="action-btn delete-btn" :title="t('common.delete')" @click="onDelete">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
  padding-top: 24px;
  padding-bottom: 48px;
  width: 100%;
  max-width: 960px;
  margin: 0 auto;
  box-sizing: border-box;
}

.loading-state,
.error-state {
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

.error-state {
  color: var(--md-num-negative);
  border-color: var(--md-num-negative);
  background-color: var(--md-bg-disclaimer);
}

.post-article {
  background-color: var(--color-surface);
  border: 1px solid var(--md-border);
  border-radius: 8px;
  padding: 32px;
  box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.08);
}

.post-header {
  margin-bottom: 28px;
}

.header-top {
  display: grid;
  grid-template-columns: 80px 1fr 80px;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.header-spacer {
  width: 80px;
}

.post-title {
  grid-column: 2;
  font-family: var(--main-font);
  font-size: 2rem;
  font-weight: 800;
  color: var(--md-text);
  text-align: center;
  margin: 0;
  line-height: 1.3;
  word-break: keep-all;
  overflow-wrap: anywhere;
}

.admin-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border: 1px solid var(--md-border);
  border-radius: 4px;
  background-color: var(--md-bg-code);
  color: var(--md-text);
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.15s ease;
}

.action-btn:active {
  transform: scale(0.95);
}

.edit-btn:hover {
  background-color: var(--md-accent);
  border-color: var(--md-accent);
  color: #ffffff;
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

.post-hero {
  position: relative;
  width: 100%;
  overflow: hidden;
  border: 1px solid var(--md-border);
  border-radius: 6px;
  box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.08);
  background-color: var(--md-bg-pre);
}

.post-hero-image {
  width: 100%;
  max-height: 420px;
  object-fit: cover;
  display: block;
}

@media (max-width: 768px) {
  .post-article {
    padding: 20px 16px;
    border-radius: 6px;
  }

  .header-top {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
  }

  .header-spacer {
    display: none;
  }

  .admin-actions {
    justify-content: center;
  }

  .post-title {
    font-size: 1.5rem;
  }

  .post-hero-image {
    max-height: 260px;
  }
}
</style>