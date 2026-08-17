<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useFilesStore } from '@/stores/files'
import { useI18n } from 'vue-i18n'

const filesStore = useFilesStore()
const { t } = useI18n()

const IMAGE_EXTENSIONS = ['png', 'jpg', 'jpeg', 'gif', 'webp', 'svg', 'avif']

function isImage(ext: string): boolean {
  return IMAGE_EXTENSIONS.includes(ext.toLowerCase().replace('.', ''))
}

const fileInput = ref<HTMLInputElement | null>(null)
const copiedId = ref<string | null>(null)

onMounted(() => filesStore.fetchFiles())

function getAbsoluteUrl(url: string): string {
  return new URL(url, window.location.origin).href
}

async function copyLink(id: string) {
  const fullUrl = getAbsoluteUrl(filesStore.getFileUrl(id))

  try {
    await navigator.clipboard.writeText(fullUrl)
    copiedId.value = id
    setTimeout(() => {
      if (copiedId.value === id) copiedId.value = null
    }, 2000)
  } catch (err) {
    console.error(t('files.err_copy'), err)
  }
}

function openFilePicker() {
  fileInput.value?.click()
}

async function onFilesSelected(event: Event) {
  const input = event.target as HTMLInputElement
  const selected = input.files
  if (!selected || selected.length === 0) return

  for (const file of Array.from(selected)) {
    await filesStore.upload(file)
  }

  await filesStore.fetchFiles(filesStore.currentPage)
  input.value = ''
}

async function onDelete(id: string) {
  await filesStore.remove(id)
}

function goToPage(page: number) {
  filesStore.fetchFiles(page)
}
</script>

<template>
  <div class="gallery container">
    <input
      ref="fileInput"
      type="file"
      multiple
      class="hidden-input"
      @change="onFilesSelected"
    />

    <p v-if="filesStore.error" class="error-message">[SYS_ERROR] :: {{ filesStore.error }}</p>

    <div class="gallery-grid">
      <button
        type="button"
        class="upload-tile"
        :disabled="filesStore.isUploading"
        @click="openFilePicker"
      >
        <span class="upload-icon">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 5v14" />
            <path d="M5 12h14" />
          </svg>
        </span>
        <span class="upload-label">
          {{ filesStore.isUploading ? `${t('files.upload')} ${filesStore.uploadProgress}%` : `[ ${t('files.add_file')} ]` }}
        </span>
      </button>

      <a
        v-for="item in filesStore.files"
        :key="item.id"
        class="file-tile"
        :href="filesStore.getFileUrl(item.id)"
        target="_blank"
        rel="noopener"
      >
        <img
          v-if="isImage(item.ext)"
          :src="filesStore.getFileUrl(item.id)"
          :alt="item.id"
          class="file-tile-img"
        />

        <div v-else class="file-tile-doc">
          <span class="file-tile-icon">
            <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
              <polyline points="14 2 14 8 20 8" />
            </svg>
          </span>
          <span class="file-tile-ext">{{ item.ext || 'file' }}</span>
        </div>

        <div class="file-tile-footer">
          <span class="file-tile-name">{{ item.id }}{{ item.ext }}</span>
        </div>

        <div class="tile-actions">
          <button
            type="button"
            class="action-btn copy-btn"
            :class="{ active: copiedId === item.id }"
            :title="copiedId === item.id ? t('files.link_copied') : t('files.copy_link')"
            @click.prevent.stop="copyLink(item.id)"
          >
            <svg v-if="copiedId === item.id" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12" />
            </svg>
            <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
              <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
            </svg>
          </button>

          <button
            type="button"
            class="action-btn delete-btn"
            :title="t('common.delete')"
            @click.prevent.stop="onDelete(item.id)"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3 6 5 6 21 6" />
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
            </svg>
          </button>
        </div>
      </a>
    </div>

    <div v-if="filesStore.pages > 1" class="pagination">
      <button
        v-for="page in filesStore.pages"
        :key="page"
        type="button"
        :disabled="page === filesStore.currentPage"
        @click="goToPage(page)"
      >
        [{{ page < 10 ? `0${page}` : page }}]
      </button>
    </div>
  </div>
</template>

<style scoped>
.gallery {
  width: 100%;
}

.hidden-input {
  display: none;
}

.gallery-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(170px, 1fr));
  justify-content: start;
  gap: 20px;
  margin-top: 40px;
}

.upload-tile,
.file-tile {
  position: relative;
  width: 100%;
  aspect-ratio: 1 / 1;
  overflow: hidden;
  box-sizing: border-box;
}

.upload-tile {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border: 1px dashed var(--md-accent);
  background-color: var(--md-bg-blockquote);
  color: var(--md-accent);
  cursor: pointer;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 10px), calc(100% - 10px) 100%, 0 100%);
  transition: all 0.2s ease;
}

.upload-tile:hover:not(:disabled) {
  border-color: var(--md-link);
  background-color: var(--md-bg-table-row-hover);
  color: var(--md-link);
  box-shadow: 0 0 15px rgba(255, 0, 85, 0.2);
}

.upload-tile:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  box-shadow: none;
}

.upload-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  background-color: var(--color-surface);
  border: 1px solid var(--md-accent);
  color: inherit;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%);
  transition: all 0.2s ease;
}

.upload-tile:hover:not(:disabled) .upload-icon {
  border-color: var(--md-link);
  background-color: var(--md-link);
  color: #050811;
  box-shadow: 0 0 10px var(--md-link);
}

.upload-label {
  font-family: var(--main-font);
  font-size: 0.85rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.file-tile {
  display: flex;
  flex-direction: column;
  background-color: var(--color-surface);
  border: 1px solid var(--md-border-light);
  color: var(--color-text);
  text-decoration: none;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 10px), calc(100% - 10px) 100%, 0 100%);
  transition: all 0.2s ease;
}

.file-tile:hover {
  border-color: var(--md-accent);
  box-shadow: 0 0 15px rgba(0, 240, 255, 0.2);
}

.file-tile-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.3s ease;
}

.file-tile:hover .file-tile-img {
  transform: scale(1.06);
}

.file-tile-doc {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background-color: var(--md-bg-pre);
}

.file-tile-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background-color: var(--color-surface);
  border: 1px solid var(--md-border-light);
  color: var(--md-accent);
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%);
}

.file-tile-ext {
  font-family: var(--main-font);
  font-size: 0.7rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  color: var(--md-accent);
  text-transform: uppercase;
  background-color: var(--md-bg-blockquote);
  padding: 3px 8px;
  border: 1px solid var(--md-accent);
}

.file-tile-footer {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 12px 10px 8px;
  background: linear-gradient(to top, rgba(5, 8, 17, 0.9) 0%, transparent 100%);
}

.file-tile-doc ~ .file-tile-footer {
  position: static;
  background: transparent;
  padding: 10px;
  border-top: 1px solid var(--md-border-light);
}

.file-tile-name {
  display: block;
  font-family: var(--main-font);
  font-size: 0.75rem;
  font-weight: 700;
  color: #ffffff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  text-shadow: 0 0 6px rgba(0, 0, 0, 0.8);
}

.file-tile-doc ~ .file-tile-footer .file-tile-name {
  color: var(--md-text);
  text-shadow: none;
  text-align: center;
}

.tile-actions {
  position: absolute;
  top: 8px;
  left: 8px;
  display: flex;
  gap: 6px;
  z-index: 2;
  padding: 3px;
  background: rgba(5, 8, 17, 0.85);
  border: 1px solid var(--md-accent);
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%);
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 4px), calc(100% - 4px) 100%, 0 100%);
}

.action-btn:active {
  transform: scale(0.92);
}

.copy-btn {
  background-color: transparent;
  color: var(--md-text);
}

.copy-btn:hover {
  background-color: rgba(0, 240, 255, 0.2);
  color: var(--md-accent);
}

.copy-btn.active {
  background-color: var(--md-accent);
  color: #050811;
  box-shadow: 0 0 8px var(--md-accent);
}

.delete-btn {
  background-color: var(--md-num-negative);
  color: #ffffff;
}

.delete-btn:hover {
  background-color: #ff0055;
  box-shadow: 0 0 8px #ff0055;
}

.error-message {
  color: var(--md-num-negative);
  padding: 14px;
  background-color: var(--md-mark-red-bg);
  border: 1px solid var(--md-num-negative);
  font-family: var(--main-font);
  font-size: 0.9rem;
  font-weight: 700;
  text-align: center;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 8px), calc(100% - 8px) 100%, 0 100%);
  margin-bottom: 20px;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  margin-top: 32px;
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
</style>