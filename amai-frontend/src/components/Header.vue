<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n';

const { t, locale } = useI18n()

const currentTheme = ref('light')

function changeTheme() {
  document.documentElement.className = currentTheme.value
  localStorage.setItem('user-theme', currentTheme.value)
}

const lang = ref(localStorage.getItem('user-lang') || locale.value)

function changeLanguage() {
  locale.value = lang.value
  localStorage.setItem('user-lang', lang.value)
}

onMounted(() => {
  const savedTheme = localStorage.getItem('user-theme')
  if (savedTheme) {
    currentTheme.value = savedTheme
  }
  changeTheme()

  const savedLanguage = localStorage.getItem('user-lang')
  if (savedLanguage) {
    lang.value = savedLanguage
    locale.value = savedLanguage
  }
})

</script>

<template>
  <header class="header">
    <div class="header-top">
      <h1 class="logo" v-if="$route.meta.isAdminPage">Amai: Admin</h1>
      <h1 class="logo">Amai</h1>
      <div class="settings">
        <div class="settings-item">
          <label for="theme-select">theme: </label>
          <select id="theme-select" v-model="currentTheme" @change="changeTheme()">
            <option value="light">light</option>
            <option value="dark">dark</option>
          </select>
        </div>
        <div class="settings-item">
          <label for="lang-select">language: </label>
          <select id="lang-select" v-model="lang" @change="changeLanguage()">
            <option value="en-US">english</option> 
            <option value="ru-RU">russian</option>
          </select>
        </div>
      </div>
      <nav>
        <RouterLink to="/"><h3>{{ t('nav.main') }}</h3></RouterLink>
        <RouterLink to="/about"><h3>{{ t('nav.about') }}</h3></RouterLink>
        <RouterLink to="/articles"><h3>{{ t('nav.articles') }}</h3></RouterLink>
      </nav>
    </div>
  </header>
</template>