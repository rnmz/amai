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
      <span class="sys-title">// SYSTEM_INFO :: ABOUT_MODULE</span>
      <span class="sys-status">[ONLINE]</span>
    </div>
    <main ref="contentRef" class="markdown-body" v-html="parsedContent"></main>
  </div>
</template>

<style scoped>
.about-page {
  margin: 40px auto 60px;
  padding: 32px;
  width: 100%;
  max-width: 960px;
  box-sizing: border-box;
  background-color: var(--color-surface);
  border: 1px solid var(--md-border);
  clip-path: polygon(
    0 0, 
    calc(100% - 20px) 0, 
    100% 20px, 
    100% 100%, 
    20px 100%, 
    0 calc(100% - 20px)
  );
  box-shadow: 0 0 20px rgba(0, 240, 255, 0.1);
  position: relative;
}

.sys-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 16px;
  margin-bottom: 24px;
  border-bottom: 1px solid var(--md-border-light);
  font-family: var(--main-font);
  font-size: 0.85rem;
  font-weight: 700;
  letter-spacing: 0.05em;
}

.sys-title {
  color: var(--md-accent);
  text-shadow: 0 0 8px var(--md-accent);
}

.sys-status {
  color: var(--md-num-positive);
  text-shadow: 0 0 8px var(--md-num-positive);
}

@media (max-width: 768px) {
  .about-page {
    padding: 20px 16px;
    margin-top: 20px;
  }
}
</style>