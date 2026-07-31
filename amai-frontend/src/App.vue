<script setup lang="ts">
import { RouterView } from 'vue-router'
import AppHeaderComponentVue from '@/components/AppHeaderComponent.vue'
</script>

<template>
  <div class="app-wrapper">
    <AppHeaderComponentVue />
    <main class="content-container">
      <RouterView />
    </main>
  </div>
</template>

<style>
@font-face {
  font-family: 'Noto Sans JP';
  src: url('@/assets/font/NotoSansJP-Regular.ttf') format('truetype');
  font-weight: normal;
  font-style: normal;
  font-display: swap;
}

@font-face {
  font-family: 'Nunito';
  src: url('@/assets/font/Nunito-Regular.ttf') format('truetype');
  font-weight: normal;
  font-style: normal;
  font-display: swap;
}

:root {
  --main-font: 'Nunito', 'Noto Sans JP', sans-serif;
}
html[lang='ja'] {
  --main-font: 'Noto Sans JP', 'Nunito', sans-serif;
}

body {
  margin: 0;
  padding: 0;
  font-family: var(--main-font), serif;
  background-color: #f9f9f9;
  color: #111;
}
.content-container {
  width: 100%;
  max-width: 900px;
  margin: 0 auto;
  padding: 0 20px 10px 20px;
  box-sizing: border-box;
}
</style>
