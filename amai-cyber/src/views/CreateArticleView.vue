<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, useTemplateRef, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { marked, initCopyButtons, initCharts, processFootnotes } from '@/utils/marked-render'
import { useFilesStore } from '@/stores/files'
import { usePostsStore, usePostStore } from '../stores/articles'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const router = useRouter()
const filesStore = useFilesStore()
const postsStore = usePostsStore()
const postStore = usePostStore()
const { t } = useI18n()

marked.setOptions({
  breaks: true,
})

const editId = computed(() => route.query.editId as string | undefined)
const isEditMode = computed(() => Boolean(editId.value))

const existingPosterId = ref<string | null>(null)

const heroFile = ref<File | null>(null)
const heroPreviewUrl = ref<string | null>(null)
const heroInput = ref<HTMLInputElement | null>(null)

function openHeroPicker() {
  heroInput.value?.click()
}

function onHeroSelected(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  if (heroPreviewUrl.value && heroFile.value) {
    URL.revokeObjectURL(heroPreviewUrl.value)
  }

  heroFile.value = file
  heroPreviewUrl.value = URL.createObjectURL(file)
}

function removeHero() {
  if (heroPreviewUrl.value && heroFile.value) {
    URL.revokeObjectURL(heroPreviewUrl.value)
  }
  heroFile.value = null
  heroPreviewUrl.value = null
  existingPosterId.value = null
  if (heroInput.value) {
    heroInput.value.value = ''
  }
}

onBeforeUnmount(() => {
  if (heroPreviewUrl.value && heroFile.value) {
    URL.revokeObjectURL(heroPreviewUrl.value)
  }
  if (chartRenderTimeout) {
    clearTimeout(chartRenderTimeout)
  }
})

const title = ref('')
const markdownContent = ref('')

const parsedContent = computed(() => {
  return marked.parse(processFootnotes(markdownContent.value))
})

const previewRef = useTemplateRef<HTMLElement>('previewRef')

const errorKey = ref<string | null>(null)
const rawErrorMessage = ref<string | null>(null)

const errorMessage = computed(() => {
  if (errorKey.value) return t(errorKey.value)
  return rawErrorMessage.value
})

const showPreview = ref(false)

let chartRenderTimeout: ReturnType<typeof setTimeout> | null = null

function scheduleChartRender() {
  if (chartRenderTimeout) {
    clearTimeout(chartRenderTimeout)
  }
  chartRenderTimeout = setTimeout(() => {
    nextTick(() => {
      if (previewRef.value) {
        initCharts(previewRef.value)
      }
    })
  }, 300)
}

watch(parsedContent, () => {
  if (showPreview.value) {
    scheduleChartRender()
  }
})

watch(showPreview, (visible) => {
  if (visible) {
    scheduleChartRender()
  }
})

onMounted(async () => {
  if (isEditMode.value && editId.value) {
    await postStore.fetchPost(editId.value)
    const post = postStore.post
    if (post) {
      title.value = post.title
      markdownContent.value = post.body
      if (post.poster_id) {
        existingPosterId.value = post.poster_id
        heroPreviewUrl.value = filesStore.getFileUrl(post.poster_id)
      }
    } else if (postStore.error) {
      rawErrorMessage.value = postStore.error
    }
  }

  if (previewRef.value) {
    initCopyButtons(previewRef.value)
  }
})

const isSaving = ref(false)

async function onPublish() {
  errorKey.value = null
  rawErrorMessage.value = null

  if (!title.value.trim()) {
    errorKey.value = 'create_article.err_no_title'
    return
  }

  if (!markdownContent.value.trim()) {
    errorKey.value = 'create_article.err_no_text'
    return
  }

  if (!heroFile.value && !existingPosterId.value) {
    errorKey.value = 'create_article.err_no_img'
    return
  }

  isSaving.value = true

  try {
    let posterId = existingPosterId.value
    if (heroFile.value) {
      const fileId = await filesStore.upload(heroFile.value)
      if (!fileId) {
        if (filesStore.error) {
          rawErrorMessage.value = filesStore.error
        } else {
          errorKey.value = 'create_article.err_img_upload'
        }
        return
      }
      posterId = fileId
    }

    const success = isEditMode.value && editId.value
      ? await postsStore.update({
          id: editId.value,
          title: title.value.trim(),
          poster_id: posterId!,
          body: markdownContent.value,
        })
      : await postsStore.create({
          title: title.value.trim(),
          poster_id: posterId!,
          body: markdownContent.value,
        })

    if (!success) {
      if (postsStore.error) {
        rawErrorMessage.value = postsStore.error
      } else {
        errorKey.value = 'create_article.err_save'
      }
      return
    }

    router.push({ name: 'admin_articles' })
  } finally {
    isSaving.value = false
  }
}

function onCancel() {
  router.back()
}
</script>

<template>
  <div class="post-page container">
    <header class="post-header">
      <input
        ref="heroInput"
        type="file"
        accept="image/*"
        class="hidden-input"
        @change="onHeroSelected"
      />

      <input
        v-model="title"
        type="text"
        class="title-input"
        :placeholder="t('create_article.hint_title')"
      />

      <div class="post-hero" :class="{ 'post-hero-empty': !heroPreviewUrl }">
        <img v-if="heroPreviewUrl" :src="heroPreviewUrl" :alt="t('create_article.img')" class="post-hero-image" />

        <button v-if="!heroPreviewUrl" type="button" class="hero-upload-btn" @click="openHeroPicker">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="18" height="18" rx="2" />
            <circle cx="9" cy="9" r="2" />
            <path d="m21 15-5-5L5 21" />
          </svg>
          <span>[{{ t('create_article.upload_img') }}]</span>
        </button>

        <button v-else type="button" class="hero-remove-btn" :title="t('create_article.delete_img')" @click="removeHero">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6" />
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
          </svg>
        </button>
      </div>
    </header>

    <div class="editor-toolbar">
      <span class="pane-label">// {{ showPreview ? t('create_article.preview') : 'MARKDOWN_BUFFER' }}</span>
      <button type="button" class="preview-toggle-btn" @click="showPreview = !showPreview">
        <svg
          v-if="!showPreview"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7Z" />
          <circle cx="12" cy="12" r="3" />
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
        </svg>
        [{{ showPreview ? t('create_article.back_to_edit') : t('create_article.show_preview') }}]
      </button>
    </div>

    <div class="editor-layout">
      <div v-show="!showPreview" class="editor-pane">
        <textarea
          v-model="markdownContent"
          class="markdown-textarea"
          :placeholder="t('create_article.hint_text')"></textarea>
      </div>

      <div v-show="showPreview" class="preview-pane">
        <main ref="previewRef" class="markdown-body preview-content" v-html="parsedContent"></main>
      </div>
    </div>

    <p v-if="errorMessage" class="error-message">[SYS_ERROR] :: {{ errorMessage }}</p>

    <div class="publish-bar">
      <button type="button" class="cancel-btn" :disabled="isSaving" @click="onCancel">[{{ t('cancel') }}]</button>
      <button type="button" class="publish-btn" :disabled="isSaving" @click="onPublish">
        [{{ isSaving ? t('create_article.publishing') : t('create_article.publish') }}]
      </button>
    </div>
  </div>
</template>

<style scoped>
.post-page {
  padding-top: 30px;
  padding-bottom: 60px;
  width: 100%;
  max-width: 960px;
  margin: 0 auto;
  box-sizing: border-box;
}

.post-header {
  margin-bottom: 24px;
}

.hidden-input {
  display: none;
}

.post-hero {
  position: relative;
  width: 100%;
  margin-top: 20px;
  border: 1px solid var(--md-border);
  background-color: var(--color-surface);
  clip-path: polygon(0 0, calc(100% - 15px) 0, 100% 15px, 100% 100%, 15px 100%, 0 calc(100% - 15px));
  box-shadow: 0 0 15px rgba(0, 240, 255, 0.15);
}

.post-hero-empty {
  aspect-ratio: 2 / 1;
  border: 2px dashed var(--md-accent);
  background-color: var(--md-bg-blockquote);
  transition: all 0.25s ease;
}

.post-hero-empty:hover {
  border-color: var(--md-link);
  background-color: var(--md-bg-table-row-hover);
  box-shadow: 0 0 15px rgba(255, 0, 85, 0.2);
}

.post-hero-image {
  width: 100%;
  max-height: 450px;
  object-fit: cover;
  display: block;
}

.hero-upload-btn {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border: none;
  background: transparent;
  color: var(--md-accent);
  cursor: pointer;
  font-family: var(--main-font);
  font-size: 0.95rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  transition: all 0.2s ease;
}

.hero-upload-btn:hover {
  color: var(--md-link);
  text-shadow: 0 0 8px var(--md-link);
}

.hero-remove-btn {
  position: absolute;
  top: 12px;
  left: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background-color: var(--md-num-negative);
  color: #ffffff;
  cursor: pointer;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%);
  box-shadow: 0 0 10px var(--md-num-negative);
  transition: all 0.2s ease;
}

.hero-remove-btn:hover {
  background-color: #ff0055;
  transform: scale(1.05);
}

.hero-remove-btn:active {
  transform: scale(0.92);
}

.title-input {
  width: 100%;
  font-size: 2rem;
  font-weight: 900;
  color: #ffffff;
  text-align: center;
  border: none;
  border-bottom: 2px solid var(--md-border);
  padding: 8px 8px 16px;
  line-height: 1.25;
  font-family: var(--main-font);
  letter-spacing: 0.03em;
  text-transform: uppercase;
  background: transparent;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.title-input::placeholder {
  color: var(--md-text-secondary);
  opacity: 0.5;
}

.title-input:focus {
  outline: none;
  border-bottom-color: var(--md-accent);
  text-shadow: 0 0 8px var(--md-accent);
}

.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  font-family: var(--main-font);
}

.pane-label {
  font-size: 0.85rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  color: var(--md-accent);
  text-shadow: 0 0 6px var(--md-accent);
}

.preview-toggle-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  border: 1px solid var(--md-border-light);
  background-color: var(--color-surface);
  color: var(--md-text-secondary);
  font-family: var(--main-font);
  font-size: 0.85rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  cursor: pointer;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%);
  transition: all 0.2s ease;
}

.preview-toggle-btn:hover {
  background-color: rgba(0, 240, 255, 0.1);
  border-color: var(--md-accent);
  color: var(--md-accent);
  box-shadow: 0 0 10px rgba(0, 240, 255, 0.2);
}

.editor-layout {
  display: block;
}

.editor-pane,
.preview-pane {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.markdown-textarea {
  width: 100%;
  min-width: 0;
  min-height: 520px;
  padding: 18px;
  border: 1px solid var(--md-border-light);
  background-color: var(--color-surface);
  color: var(--md-text);
  font-family: 'Orbitron', monospace;
  font-size: 0.9rem;
  line-height: 1.6;
  resize: vertical;
  box-sizing: border-box;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 12px), calc(100% - 12px) 100%, 0 100%);
  transition: all 0.2s ease;
}

.markdown-textarea:focus {
  outline: none;
  border-color: var(--md-accent);
  box-shadow: 0 0 12px rgba(0, 240, 255, 0.2);
}

.preview-content {
  min-height: 520px;
  padding: 18px;
  border: 1px solid var(--md-border-light);
  background-color: var(--color-surface);
  color: var(--md-text);
  overflow-y: auto;
  box-sizing: border-box;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 12px), calc(100% - 12px) 100%, 0 100%);
}

.error-message {
  margin: 20px 0 0 0;
  padding: 12px 16px;
  background-color: var(--md-mark-red-bg);
  color: var(--md-num-negative);
  border: 1px solid var(--md-num-negative);
  font-family: var(--main-font);
  font-size: 0.9rem;
  font-weight: 700;
  text-align: center;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 8px), calc(100% - 8px) 100%, 0 100%);
}

.publish-bar {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 32px;
}

.cancel-btn,
.publish-btn {
  padding: 10px 28px;
  font-family: var(--main-font);
  font-size: 0.95rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  cursor: pointer;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 8px), calc(100% - 8px) 100%, 0 100%);
  transition: all 0.2s ease;
}

.cancel-btn {
  border: 1px solid var(--md-border-light);
  background-color: var(--color-surface);
  color: var(--md-text-secondary);
}

.cancel-btn:hover:not(:disabled) {
  background-color: rgba(255, 255, 255, 0.05);
  border-color: var(--md-text);
  color: var(--md-text);
}

.publish-btn {
  border: 1px solid var(--md-accent);
  background-color: var(--md-accent);
  color: #050811;
  box-shadow: 0 0 10px var(--md-accent);
}

.publish-btn:hover:not(:disabled) {
  background-color: #00f0ff;
  box-shadow: 0 0 16px #00f0ff;
  transform: translateY(-1px);
}

.publish-btn:active:not(:disabled) {
  transform: scale(0.97);
}

.cancel-btn:disabled,
.publish-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  box-shadow: none;
}

@media (max-width: 768px) {
  .title-input {
    font-size: 1.5rem;
  }

  .markdown-textarea,
  .preview-content {
    min-height: 360px;
  }
}
</style>
