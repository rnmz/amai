<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, useTemplateRef, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { marked, initCopyButtons, initCharts, processFootnotes } from '@/utils/marked-render'
import { useFilesStore } from '@/stores/files'
import { usePostsStore, usePostStore } from '@/stores/posts'
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
          <span>{{ t('create_article.upload_img') }}</span>
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
      <span class="pane-label">{{ showPreview ? t('create_article.preview') : 'Markdown' }}</span>
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
        {{ showPreview ? t('create_article.back_to_edit') : t('create_article.show_preview') }}
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

    <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>

    <div class="publish-bar">
      <button type="button" class="cancel-btn" :disabled="isSaving" @click="onCancel">{{ t('cancel') }}</button>
      <button type="button" class="publish-btn" :disabled="isSaving" @click="onPublish">
        {{ isSaving ? t('create_article.publishing') : t('create_article.publish') }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.post-page {
  padding-top: 40px;
  padding-bottom: 60px;
}

.post-header {
  margin-bottom: 32px;
}

.hidden-input {
  display: none;
}

.post-hero {
  position: relative;
  width: 100%;
  border-radius: 16px;
  overflow: hidden;
  margin-top: 20px;
}

.post-hero-empty {
  aspect-ratio: 2 / 1;
  border: 2px dashed var(--md-border);
  background-color: var(--md-bg-disclaimer);
  transition: border-color 0.2s ease, background-color 0.2s ease;
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
  gap: 8px;
  border: none;
  background: transparent;
  color: var(--md-text-secondary);
  cursor: pointer;
  font-size: 0.9rem;
  transition: color 0.15s ease;
}

.hero-upload-btn:hover {
  color: var(--md-link);
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
  border-radius: 8px;
  background-color: var(--md-num-negative);
  color: #ffffff;
  cursor: pointer;
  transition:
    transform 0.1s ease,
    opacity 0.2s;
}

.hero-remove-btn:hover {
  opacity: 0.9;
}

.hero-remove-btn:active {
  transform: scale(0.95);
}

.title-input {
  width: 100%;
  font-size: 2.2rem;
  font-weight: 400;
  color: var(--color-text);
  text-align: center;
  border: none;
  border-bottom: 2px solid transparent;
  padding: 8px 8px 16px;
  line-height: 1.25;
  font-family: inherit;
  background: transparent;
  transition: border-color 0.15s ease, color 0.2s ease;
}

.title-input::placeholder {
  color: var(--md-text-secondary);
  opacity: 0.7;
}

.title-input:focus {
  outline: none;
  border-bottom-color: var(--md-link);
}

.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.pane-label {
  font-size: 0.8rem;
  font-weight: 600;
  letter-spacing: 0.03em;
  text-transform: uppercase;
  color: var(--md-text-secondary);
}

.preview-toggle-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  background-color: var(--color-surface);
  color: var(--md-text-secondary);
  font-size: 0.82rem;
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    border-color 0.2s ease,
    color 0.2s ease;
}

.preview-toggle-btn:hover {
  background-color: var(--md-bg-pre);
  border-color: var(--md-border);
  color: var(--color-text);
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
  min-height: 560px;
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: 10px;
  background-color: var(--color-surface);
  color: var(--color-text);
  font-family: 'SF Mono', 'Roboto Mono', Consolas, monospace;
  font-size: 0.9rem;
  line-height: 1.6;
  resize: vertical;
  box-sizing: border-box;
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease,
    background-color 0.2s ease,
    color 0.2s ease;
}

.markdown-textarea:focus {
  outline: none;
  border-color: var(--md-link);
  box-shadow: 0 0 0 3px var(--md-selection-bg);
}

.preview-content {
  min-height: 560px;
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: 10px;
  background-color: var(--color-surface);
  color: var(--color-text);
  overflow-y: auto;
  box-sizing: border-box;
  transition: background-color 0.2s ease, border-color 0.2s ease;
}

.error-message {
  margin: 20px 0 0 0;
  padding: 10px 14px;
  border-radius: 8px;
  background-color: var(--md-mark-red-bg);
  color: var(--md-num-negative);
  font-size: 0.9rem;
  text-align: center;
}

.publish-bar {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 32px;
}

.cancel-btn,
.publish-btn {
  padding: 11px 28px;
  border-radius: 8px;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    border-color 0.2s ease,
    opacity 0.2s ease;
}

.cancel-btn {
  border: 1px solid var(--color-border);
  background-color: var(--color-surface);
  color: var(--md-text-secondary);
}

.cancel-btn:hover:not(:disabled) {
  background-color: var(--md-bg-pre);
  color: var(--color-text);
}

.publish-btn {
  border: none;
  background-color: var(--md-link);
  color: #ffffff;
}

.publish-btn:hover:not(:disabled) {
  background-color: var(--md-link-hover);
}

.cancel-btn:disabled,
.publish-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .title-input {
    font-size: 1.6rem;
  }

  .markdown-textarea,
  .preview-content {
    min-height: 360px;
  }
}
</style>