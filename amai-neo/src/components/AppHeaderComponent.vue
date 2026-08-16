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
  border-bottom: 2px solid var(--md-accent);
  margin-top: 40px;
  padding-bottom: 8px;
  position: relative;
}

.header::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, var(--md-accent) 0%, var(--md-link) 100%);
  opacity: 0.8;
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo-link {
  text-decoration: none;
  color: inherit;
  transition: opacity 0.2s ease;
}

.logo-link:hover {
  opacity: 0.9;
}

.logo {
  font-size: 2.2rem;
  font-weight: 700;
  margin: 0;
  line-height: 1;
  font-family: 'Nunito', serif;
  letter-spacing: -0.02em;
  background: linear-gradient(135deg, var(--color-text) 0%, var(--md-text-secondary) 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.settings {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.settings-item {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  font-family: 'Nunito', sans-serif;
  font-size: 1rem;
  font-weight: 500;
}

.settings-item label {
  color: var(--md-text-secondary);
}

.settings-item select {
  font-family: inherit;
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--color-text);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 6px 12px;
  outline: none;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.settings-item select:hover {
  border-color: var(--md-accent);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

.settings-item select:focus {
  border-color: var(--md-accent);
  box-shadow: 0 0 0 3px var(--md-selection-bg);
}

.settings-item select option {
  background-color: var(--color-surface);
  color: var(--color-text);
  padding: 8px;
}

.navigation {
  font-family: 'Nunito', sans-serif;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 28px;
  margin-top: 20px;
  margin-bottom: 12px;
}

.navigation a {
  text-decoration: none;
  color: var(--color-text);
  font-size: 1.1rem;
  font-weight: 600;
  position: relative;
  padding: 4px 0;
  transition: color 0.2s ease;
}

.navigation a::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  width: 0;
  height: 2px;
  background-color: var(--md-accent);
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  transform: translateX(-50%);
  border-radius: 2px;
}

.navigation a:hover {
  color: var(--md-accent);
}

.navigation a:hover::after {
  width: 100%;
}

.logout-btn {
  font-family: inherit;
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--md-num-negative);
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.logout-btn:hover {
  background-color: var(--md-mark-red-bg);
  opacity: 1;
}
</style>