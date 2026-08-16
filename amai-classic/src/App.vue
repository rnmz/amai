<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import AppHeaderComponentVue from '@/components/AppHeaderComponent.vue'
import { useTheme } from '@/composables/useTheme'
import { initMarkdownTheming } from '@/utils/marked-render'
import '@/assets/style.css'

const route = useRoute()
useTheme()
initMarkdownTheming()
</script>

<template>
  <div class="app-wrapper">
    <AppHeaderComponentVue v-if="!route.meta.hideHeader" />
    <main class="content-container">
      <RouterView />
    </main>
  </div>
</template>
