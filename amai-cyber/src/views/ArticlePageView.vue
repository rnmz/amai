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
      <span class="glitch-text">// {{ t('common.loading') }}...</span>
    </div>
    <div v-else-if="postStore.error" class="error-state">
      [SYS_ERROR] :: {{ postStore.error }}
    </div>

    <article v-else-if="postStore.post" class="post-article">
      <header class="post-header">
        <div class="header-top">
          <div v-if="isAdmin" class="header-spacer"></div>

          <h1 class="post-title" :data-text="postStore.post.title">{{ postStore.post.title }}</h1>

          <div v-if="isAdmin" class="admin-actions">
            <button class="action-btn edit-btn" :title="t('common.edit')" @click="onEdit">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z" />
              </svg>
            </button>
            <button class="action-btn delete-btn" :title="t('common.delete')" @click="onDelete">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
          <div class="hero-scanline"></div>
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
  width: 100%;
  max-width: 960px;
  margin: 0 auto;
  box-sizing: border-box;
}

.loading-state,
.error-state {
  text-align: center;
  font-family: var(--main-font);
  font-size: 1.1rem;
  font-weight: 700;
  padding: 30px 20px;
  color: var(--md-accent);
  background-color: var(--color-surface);
  border: 1px solid var(--md-border);
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 12px), calc(100% - 12px) 100%, 0 100%);
  box-shadow: 0 0 15px rgba(0, 240, 255, 0.15);
  letter-spacing: 0.05em;
}

.error-state {
  color: var(--md-num-negative);
  border-color: var(--md-num-negative);
  background-color: var(--md-mark-red-bg);
  box-shadow: 0 0 15px rgba(255, 0, 85, 0.2);
}

.post-article {
  background-color: var(--color-surface);
  border: 1px solid var(--md-border-light);
  padding: 32px;
  clip-path: polygon(
    0 0, 
    calc(100% - 20px) 0, 
    100% 20px, 
    100% 100%, 
    20px 100%, 
    0 calc(100% - 20px)
  );
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.5);
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
  font-family: var(--main-font);
  font-size: 2.2rem;
  font-weight: 900;
  color: #ffffff;
  text-align: center;
  margin: 0;
  line-height: 1.25;
  letter-spacing: 0.03em;
  text-transform: uppercase;
  text-shadow: 0 0 10px var(--md-accent);
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
  width: 38px;
  height: 38px;
  border: 1px solid transparent;
  color: #050811;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.2s ease;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%);
}

.action-btn:active {
  transform: scale(0.92);
}

.edit-btn {
  background-color: var(--md-accent);
}

.edit-btn:hover {
  background-color: #00f0ff;
  box-shadow: 0 0 12px var(--md-accent);
}

.delete-btn {
  background-color: var(--md-num-negative);
  color: #ffffff;
}

.delete-btn:hover {
  background-color: #ff0055;
  box-shadow: 0 0 12px var(--md-num-negative);
}

.post-hero {
  position: relative;
  width: 100%;
  overflow: hidden;
  border: 1px solid var(--md-accent);
  clip-path: polygon(0 0, calc(100% - 15px) 0, 100% 15px, 100% 100%, 15px 100%, 0 calc(100% - 15px));
  box-shadow: 0 0 15px rgba(0, 240, 255, 0.2);
  background-color: var(--md-bg-pre);
}

.post-hero-image {
  width: 100%;
  max-height: 450px;
  object-fit: cover;
  display: block;
}

.hero-scanline {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to bottom,
    rgba(255,255,255,0),
    rgba(255,255,255,0) 50%,
    rgba(0, 0, 0, 0.3) 50%,
    rgba(0, 0, 0, 0.3)
  );
  background-size: 100% 4px;
  pointer-events: none;
  opacity: 0.6;
}

@media (max-width: 768px) {
  .post-article {
    padding: 20px 16px;
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
    font-size: 1.6rem;
  }

  .post-hero-image {
    max-height: 280px;
  }
}
</style>