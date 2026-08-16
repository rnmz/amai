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
    <main ref="contentRef" class="markdown-body" v-html="parsedContent"></main>
  </div>
</template>

<style scoped>
.about-page {
  padding-top: 60px;
  padding-bottom: 60px;
}
</style>