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
    <main ref="contentRef" class="markdown-body" v-html="parsedContent"></main>
  </div>
</template>

<style scoped>
.content-wrapper {
  padding-top: 60px;
  padding-bottom: 60px;
}
</style>