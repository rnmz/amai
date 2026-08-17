<script setup lang="ts">
import { computed, onMounted, ref, useTemplateRef } from 'vue'
import { marked, initCopyButtons, processFootnotes } from '@/utils/marked-render'
import rawRuMarkdown from '@/assets/md/ru/main.md?raw'
import rawEnMarkdown from '@/assets/md/en/main.md?raw'
import rawJaMarkdown from '@/assets/md/jp/main.md?raw'

marked.setOptions({ breaks: true })

const CONTENT_BY_LANG: Record<string, string> = {
  'en-US': rawEnMarkdown,
  'ru-RU': rawRuMarkdown,
  'ja-JP': rawJaMarkdown,
}

const DEFAULT_LANG = 'en-US'

const lang = ref(localStorage.getItem('user-lang') ?? DEFAULT_LANG)

const parsedContent = computed(() => {
  const raw = CONTENT_BY_LANG[lang.value] ?? CONTENT_BY_LANG[DEFAULT_LANG] ?? ''
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
  <div class="content-wrapper">
    <div class="sys-header">
      <span class="sys-title">Main</span>
      <span class="sys-status">Welcome</span>
    </div>
    <main ref="contentRef" class="markdown-body" v-html="parsedContent"></main>
  </div>
</template>

<style scoped>
.content-wrapper {
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
  .content-wrapper {
    padding: 20px 16px;
    margin-top: 16px;
    border-radius: 6px;
  }
}
</style>