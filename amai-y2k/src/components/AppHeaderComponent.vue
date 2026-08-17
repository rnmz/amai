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
        <h1 class="logo">
          {{ authStore.isAuthenticated ? 'AMAI :: ADMIN' : 'AMAI' }}
        </h1>
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
            <option value="en-US">EN-US</option>
            <option value="ru-RU">RU-RU</option>
            <option value="ja-JP">JA-JP</option>
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
  font-family: var(--main-font);
  margin-top: 20px;
  padding: 16px 20px 12px;
  background: var(--color-surface);
  border: 1px solid var(--md-border-light);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.logo-link {
  text-decoration: none;
  color: inherit;
}

.logo {
  font-size: 1.8rem;
  font-weight: 800;
  margin: 0;
  line-height: 1;
  color: var(--md-accent);
  letter-spacing: -0.02em;
}

.settings {
  display: flex;
  align-items: center;
  gap: 12px;
}

.settings-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.85rem;
  font-weight: 600;
}

.settings-item label {
  color: var(--md-text-secondary);
}

.settings-item select {
  font-family: var(--main-font);
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--md-text);
  background-color: var(--md-bg-code);
  border: 1px solid var(--md-border);
  border-radius: 4px;
  padding: 4px 8px;
  outline: none;
  cursor: pointer;
}

.settings-item select option {
  background-color: var(--color-surface);
  color: var(--md-text);
}

.navigation {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  gap: 8px;
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px solid var(--md-border-light);
  flex-wrap: wrap;
}

.navigation a {
  text-decoration: none;
  color: var(--md-text);
  font-size: 0.9rem;
  font-weight: 600;
  padding: 5px 12px;
  background-color: var(--md-bg-code);
  border: 1px solid var(--md-border);
  border-radius: 4px;
  transition: all 0.15s ease;
}

.navigation a:hover {
  background-color: var(--md-bg-table-row-hover);
  color: var(--md-link);
  border-color: var(--md-link);
}

.navigation a.router-link-active {
  background-color: var(--md-accent);
  color: #ffffff;
  border-color: var(--md-accent);
}

.logout-btn {
  font-family: var(--main-font);
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--md-num-negative);
  background-color: transparent;
  border: 1px solid var(--md-num-negative);
  border-radius: 4px;
  cursor: pointer;
  padding: 5px 12px;
  transition: all 0.15s ease;
}

.logout-btn:hover {
  background-color: var(--md-num-negative);
  color: #ffffff;
}

@media (max-width: 768px) {
  .header-top {
    flex-direction: column;
    align-items: flex-start;
  }

  .settings {
    width: 100%;
    justify-content: flex-start;
  }

  .navigation {
    gap: 6px;
  }

  .navigation a,
  .logout-btn {
    font-size: 0.8rem;
    padding: 4px 8px;
  }
}
</style>