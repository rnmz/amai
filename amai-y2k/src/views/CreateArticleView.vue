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
      <span class="pane-label">{{ showPreview ? t('create_article.preview') : 'Editor' }}</span>
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
  padding-top: 24px;
  padding-bottom: 48px;
  width: 100%;
  max-width: 960px;
  margin: 0 auto;
  box-sizing: border-box;
}

.post-header {
  margin-bottom: 20px;
}

.hidden-input {
  display: none;
}

.post-hero {
  position: relative;
  width: 100%;
  margin-top: 16px;
  border: 1px solid var(--md-border);
  border-radius: 6px;
  background-color: var(--color-surface);
  box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.post-hero-empty {
  aspect-ratio: 2.2 / 1;
  border: 2px dashed var(--md-border);
  background-color: var(--md-bg-code);
  transition: all 0.15s ease;
}

.post-hero-empty:hover {
  border-color: var(--md-accent);
  background-color: var(--md-bg-table-row-hover);
}

.post-hero-image {
  width: 100%;
  max-height: 420px;
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
  color: var(--md-accent);
  cursor: pointer;
  font-family: var(--main-font);
  font-size: 0.9rem;
  font-weight: 600;
  transition: all 0.15s ease;
}

.hero-upload-btn:hover {
  color: var(--md-link-hover);
}

.hero-remove-btn {
  position: absolute;
  top: 10px;
  left: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid var(--md-border);
  border-radius: 4px;
  background-color: var(--color-surface);
  color: var(--md-num-negative);
  cursor: pointer;
  box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.15);
  transition: all 0.15s ease;
}

.hero-remove-btn:hover {
  background-color: var(--md-num-negative);
  border-color: var(--md-num-negative);
  color: #ffffff;
}

.title-input {
  width: 100%;
  font-size: 1.8rem;
  font-weight: 800;
  color: var(--md-text);
  text-align: center;
  border: 1px solid var(--md-border);
  border-radius: 6px;
  padding: 12px 16px;
  line-height: 1.25;
  font-family: var(--main-font);
  background: var(--color-surface);
  transition: all 0.15s ease;
  box-sizing: border-box;
}

.title-input::placeholder {
  color: var(--md-text-secondary);
  opacity: 0.6;
}

.title-input:focus {
  outline: none;
  border-color: var(--md-accent);
  box-shadow: 0 0 0 2px var(--md-bg-table-row-hover);
}

.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
  font-family: var(--main-font);
}

.pane-label {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--md-accent);
  text-transform: uppercase;
}

.preview-toggle-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border: 1px solid var(--md-border);
  border-radius: 4px;
  background-color: var(--color-surface);
  color: var(--md-text);
  font-family: var(--main-font);
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  box-shadow: 1px 1px 2px rgba(0, 0, 0, 0.05);
}

.preview-toggle-btn:hover {
  background-color: var(--md-bg-table-row-hover);
  border-color: var(--md-accent);
  color: var(--md-accent);
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
  min-height: 480px;
  padding: 16px;
  border: 1px solid var(--md-border);
  border-radius: 6px;
  background-color: var(--color-surface);
  color: var(--md-text);
  font-family: var(--main-font);
  font-size: 0.95rem;
  line-height: 1.6;
  resize: vertical;
  box-sizing: border-box;
  transition: all 0.15s ease;
}

.markdown-textarea:focus {
  outline: none;
  border-color: var(--md-accent);
  box-shadow: 0 0 0 2px var(--md-bg-table-row-hover);
}

.preview-content {
  min-height: 480px;
  padding: 16px;
  border: 1px solid var(--md-border);
  border-radius: 6px;
  background-color: var(--color-surface);
  color: var(--md-text);
  overflow-y: auto;
  box-sizing: border-box;
}

.error-message {
  margin: 16px 0 0 0;
  padding: 10px 14px;
  background-color: var(--md-bg-disclaimer);
  color: var(--md-num-negative);
  border: 1px solid var(--md-num-negative);
  border-radius: 6px;
  font-family: var(--main-font);
  font-size: 0.9rem;
  font-weight: 600;
  text-align: center;
}

.publish-bar {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 24px;
}

.cancel-btn,
.publish-btn {
  padding: 8px 24px;
  font-family: var(--main-font);
  font-size: 0.9rem;
  font-weight: 600;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s ease;
  box-shadow: 1px 1px 2px rgba(0, 0, 0, 0.08);
}

.cancel-btn {
  border: 1px solid var(--md-border);
  background-color: var(--color-surface);
  color: var(--md-text);
}

.cancel-btn:hover:not(:disabled) {
  background-color: var(--md-bg-code);
  border-color: var(--md-border);
}

.publish-btn {
  border: 1px solid var(--md-accent);
  background-color: var(--md-accent);
  color: #ffffff;
}

.publish-btn:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}

.publish-btn:active:not(:disabled) {
  transform: scale(0.97);
}

.cancel-btn:disabled,
.publish-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  box-shadow: none;
}

@media (max-width: 768px) {
  .title-input {
    font-size: 1.4rem;
    padding: 10px 12px;
  }

  .markdown-textarea,
  .preview-content {
    min-height: 320px;
  }
}
</style>