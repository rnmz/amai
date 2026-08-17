<script setup lang="ts">
import { computed, onMounted, ref, useTemplateRef } from 'vue'
import { marked, initCopyButtons, processFootnotes } from '@/utils/marked-render'
import rawRuAboutMarkdown from '@/assets/md/ru/about.md?raw'
import rawEnAboutMarkdown from '@/assets/md/en/about.md?raw'
import rawJaAboutMarkdown from '@/assets/md/jp/about.md?raw'

marked.setOptions({ breaks: true })

const ABOUT_BY_LANG: Record<string, string> = {
  'en-US': rawEnAboutMarkdown,
  'ru-RU': rawRuAboutMarkdown,
  'ja-JP': rawJaAboutMarkdown,
}

const DEFAULT_LANG = 'en-US'
const lang = ref(localStorage.getItem('user-lang') ?? DEFAULT_LANG)

const parsedContent = computed(() => {
  const raw = ABOUT_BY_LANG[lang.value] ?? ABOUT_BY_LANG[DEFAULT_LANG] ?? ''
  return marked.parse(processFootnotes(raw))
})

const contentRef = useTemplateRef<HTMLElement>('contentRef')

onMounted(() => {
  if (contentRef.value) {
    initCopyButtons(contentRef.value)
  }
})
</script>

<template>
  <div class="about-page container">
    <div class="sys-header">
      <span class="sys-title">About</span>
      <span class="sys-status">Information</span>
    </div>
    <main ref="contentRef" class="markdown-body" v-html="parsedContent"></main>
  </div>
</template>

<style scoped>
.about-page {
  margin: 24px auto 40px;
  padding: 28px 32px;
  width: 100%;
  max-width: 960px;
  box-sizing: border-box;
  background-color: var(--color-surface);
  border: 1px solid var(--md-border);
  border-radius: 8px;
  box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.08);
  position: relative;
}

.sys-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 12px;
  margin-bottom: 20px;
  border-bottom: 1px solid var(--md-border-light);
  font-family: var(--main-font);
  font-size: 0.85rem;
  font-weight: 600;
}

.sys-title {
  color: var(--md-accent);
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.sys-status {
  color: var(--md-text-secondary);
  font-size: 0.8rem;
}

@media (max-width: 768px) {
  .about-page {
    padding: 20px 16px;
    margin-top: 16px;
    border-radius: 6px;
  }
}
</style>