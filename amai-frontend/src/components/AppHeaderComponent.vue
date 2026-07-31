<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'

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
  document.documentElement.lang = lang.value.split('-')[0] ?? 'en'
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
  }

  changeLanguage()
})
</script>

<template>
  <header class="header content-container">
    <div class="header-top">
      <h1 class="logo" v-if="$route.meta.isAdminPage">Amai: Admin</h1>
      <h1 class="logo" v-else>Amai</h1>

      <div class="settings">
        <div class="settings-item">
          <label for="theme-select">{{ t('nav.theme') }}:</label>
          <select id="theme-select" v-model="currentTheme" @change="changeTheme()">
            <option value="light">{{ t('theme.light') }}</option>
            <option value="dark">{{ t('theme.dark') }}</option>
          </select>
        </div>
        <div class="settings-item">
          <label for="lang-select">{{ t('nav.lang') }}:</label>
          <select id="lang-select" v-model="lang" @change="changeLanguage()">
            <option value="en-US" class="lang_select_item">en</option>
            <option value="ru-RU" class="lang_select_item">ru</option>
            <option value="ja-JP" class="lang_select_item">jp</option>
          </select>
        </div>
      </div>
    </div>
    <nav class="navigation">
      <RouterLink to="/">{{ t('nav.main') }}</RouterLink>
      <RouterLink to="/articles">{{ t('nav.articles') }}</RouterLink>
      <RouterLink to="/about">{{ t('nav.about') }}</RouterLink>
    </nav>
  </header>
</template>

<style scoped>
.header {
  font-family: 'MPPlusS1p', serif;
  border-bottom: 5px solid #d8d8d8;
  margin-top: 50px;
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.logo {
  font-size: 2.2rem;
  font-weight: normal;
  margin-top: 10px;
  line-height: 1;
  font-family: 'Nunito', serif;
}
.settings {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.settings-item {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  font-family: 'Nunito', sans-serif;
  font-size: 1.1rem;
}

.settings-item label {
  color: #000;
}

.settings-item select {
  font-family: inherit;
  font-size: inherit;
  border: none;
  background: transparent;
  outline: none;
  cursor: pointer;
  padding: 0;
  margin: 0;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
}

.navigation {
  font-family: 'Nunito', sans-serif;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 25px;
  margin-top: 15px;
  margin-bottom: 10px;
}

.navigation a {
  text-decoration: none;
  color: #000;
  font-size: 1.3rem;
  transition: opacity 0.2s;
}

.navigation a:hover {
  opacity: 0.6;
}
</style>
