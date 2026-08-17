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

    <p v-if="filesStore.error" class="error-message">{{ filesStore.error }}</p>

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
          {{ filesStore.isUploading ? `${t('files.upload')} ${filesStore.uploadProgress}%` : t('files.add_file') }}
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
        {{ page }}
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
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  justify-content: start;
  gap: 16px;
  margin-top: 24px;
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
  gap: 8px;
  border: 2px dashed var(--md-border);
  border-radius: 6px;
  background-color: var(--md-bg-code);
  color: var(--md-accent);
  cursor: pointer;
  transition: all 0.15s ease;
}

.upload-tile:hover:not(:disabled) {
  border-color: var(--md-accent);
  background-color: var(--md-bg-table-row-hover);
  color: var(--md-accent);
}

.upload-tile:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.upload-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background-color: var(--color-surface);
  border: 1px solid var(--md-border);
  border-radius: 4px;
  color: inherit;
  transition: all 0.15s ease;
}

.upload-label {
  font-family: var(--main-font);
  font-size: 0.85rem;
  font-weight: 600;
}

.file-tile {
  display: flex;
  flex-direction: column;
  background-color: var(--color-surface);
  border: 1px solid var(--md-border);
  border-radius: 6px;
  color: var(--color-text);
  text-decoration: none;
  box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.08);
  transition: all 0.15s ease;
}

.file-tile:hover {
  border-color: var(--md-accent);
  box-shadow: 2px 2px 6px rgba(0, 0, 0, 0.12);
}

.file-tile-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: opacity 0.15s ease;
}

.file-tile:hover .file-tile-img {
  opacity: 0.9;
}

.file-tile-doc {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  background-color: var(--md-bg-pre);
}

.file-tile-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  background-color: var(--color-surface);
  border: 1px solid var(--md-border);
  border-radius: 4px;
  color: var(--md-accent);
}

.file-tile-ext {
  font-family: var(--main-font);
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--md-accent);
  text-transform: uppercase;
  background-color: var(--md-bg-code);
  padding: 2px 6px;
  border: 1px solid var(--md-border);
  border-radius: 3px;
}

.file-tile-footer {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 10px 8px 6px;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.75) 0%, transparent 100%);
}

.file-tile-doc ~ .file-tile-footer {
  position: static;
  background: transparent;
  padding: 8px;
  border-top: 1px solid var(--md-border-light);
}

.file-tile-name {
  display: block;
  font-family: var(--main-font);
  font-size: 0.75rem;
  font-weight: 600;
  color: #ffffff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-tile-doc ~ .file-tile-footer .file-tile-name {
  color: var(--md-text);
  text-align: center;
}

.tile-actions {
  position: absolute;
  top: 6px;
  left: 6px;
  display: flex;
  gap: 4px;
  z-index: 2;
  padding: 3px;
  background: var(--color-surface);
  border: 1px solid var(--md-border);
  border-radius: 4px;
  box-shadow: 1px 1px 2px rgba(0, 0, 0, 0.15);
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border: 1px solid var(--md-border-light);
  border-radius: 3px;
  background-color: var(--md-bg-code);
  color: var(--md-text);
  cursor: pointer;
  transition: all 0.15s ease;
}

.action-btn:active {
  transform: scale(0.92);
}

.copy-btn:hover {
  background-color: var(--md-bg-table-row-hover);
  color: var(--md-accent);
  border-color: var(--md-accent);
}

.copy-btn.active {
  background-color: var(--md-accent);
  border-color: var(--md-accent);
  color: #ffffff;
}

.delete-btn {
  color: var(--md-num-negative);
}

.delete-btn:hover {
  background-color: var(--md-num-negative);
  border-color: var(--md-num-negative);
  color: #ffffff;
}

.error-message {
  color: var(--md-num-negative);
  padding: 10px 14px;
  background-color: var(--md-bg-disclaimer);
  border: 1px solid var(--md-num-negative);
  border-radius: 6px;
  font-family: var(--main-font);
  font-size: 0.9rem;
  font-weight: 600;
  text-align: center;
  margin-bottom: 16px;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  margin-top: 24px;
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
</style>