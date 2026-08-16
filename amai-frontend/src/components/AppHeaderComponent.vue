<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useTheme } from '@/composables/useTheme'

const { t, locale } = useI18n()
const authStore = useAuthStore()
const { theme: currentTheme } = useTheme()

const lang = ref(localStorage.getItem('user-lang') || locale.value)

function changeLanguage() {
  locale.value = lang.value
  localStorage.setItem('user-lang', lang.value)
  document.documentElement.lang = lang.value.split('-')[0] ?? 'en'
}

onMounted(() => {
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
      <RouterLink :to="authStore.isAuthenticated ? '/admin/articles' : '/'" class="logo-link">
        <h1 class="logo">{{ authStore.isAuthenticated ? 'Amai: Admin' : 'Amai' }}</h1>
      </RouterLink>

      <div class="settings">
        <div class="settings-item">
          <label for="theme-select">{{ t('nav.theme') }}:</label>
          <select id="theme-select" v-model="currentTheme">
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

    <nav v-if="authStore.isAuthenticated" class="navigation">
      <RouterLink to="/">{{ t('nav.main') }}</RouterLink>
      <RouterLink :to="{ name: 'admin_articles' }">{{ t('nav.articles') }}</RouterLink>
      <RouterLink :to="{ name: 'admin_upload' }">{{ t('nav.files') }}</RouterLink>
      <RouterLink :to="{ name: 'admin_create_article' }">{{ t('nav.create') }}</RouterLink>
      <RouterLink to="/about">{{ t('nav.about') }}</RouterLink>
      <button type="button" class="logout-btn" @click="authStore.logout()">{{ t('nav.logout') }}</button>
    </nav>

    <nav v-else class="navigation">
      <RouterLink to="/">{{ t('nav.main') }}</RouterLink>
      <RouterLink to="/articles">{{ t('nav.articles') }}</RouterLink>
      <RouterLink to="/about">{{ t('nav.about') }}</RouterLink>
    </nav>
  </header>
</template>

<style scoped>
.header {
  font-family: 'Nunito', serif;
  border-bottom: 5px solid var(--color-border);
  margin-top: 50px;
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.logo-link {
  text-decoration: none;
  color: inherit;
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
  color: var(--color-text);
}

.settings-item select {
  font-family: inherit;
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--color-text);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 4px 10px;
  outline: none;
  cursor: pointer;
  transition: border-color 0.2s ease, background-color 0.2s ease;
}

.settings-item select:hover {
  border-color: var(--md-link);
}

.settings-item select:focus {
  border-color: var(--md-link);
  box-shadow: 0 0 0 2px var(--md-selection-bg);
}

.settings-item select option {
  background-color: var(--color-surface);
  color: var(--color-text);
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
  color: var(--color-text);
  font-size: 1.3rem;
  transition: opacity 0.2s;
}

.navigation a:hover {
  opacity: 0.6;
}

.logout-btn {
  font-family: inherit;
  font-size: 1.3rem;
  color: var(--color-text);
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  transition: opacity 0.2s;
}

.logout-btn:hover {
  opacity: 0.6;
}
</style>