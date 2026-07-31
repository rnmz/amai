<script setup lang="ts">
import { computed } from 'vue'
import { marked } from 'marked'
import rawRuAboutMarkdown from '@/assets/md/ru/about.md?raw'
import rawEnAboutMarkdown from '@/assets/md/en/about.md?raw'
import rawJpAboutMarkdown from '@/assets/md/jp/about.md?raw'

marked.setOptions({ breaks: true })

const lang = localStorage.getItem('user-lang')

const parsedContent = computed(() => {
  if (lang === 'en-US') {
    return marked.parse(rawEnAboutMarkdown)
  } else if (lang === 'ru-RU') {
    return marked.parse(rawRuAboutMarkdown)
  } else if (lang === 'jp-JP') {
    return marked.parse(rawJpAboutMarkdown)
  } else {
    return marked.parse(rawEnAboutMarkdown)
  }
})
</script>

<template>
  <div class="content-wrapper">
    <main class="markdown-body" v-html="parsedContent"></main>
  </div>
</template>

<style scoped>
.content-wrapper {
  width: 100%;
}
</style>
