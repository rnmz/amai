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
        <h1 class="logo" :data-text="authStore.isAuthenticated ? 'AMAI_SYS::ADMIN' : 'AMAI_SYS'">
          {{ authStore.isAuthenticated ? 'AMAI_SYS::ADMIN' : 'AMAI_SYS' }}
        </h1>
      </RouterLink>

      <div class="settings">
        <div class="settings-item">
          <label for="theme-select">//{{ t('nav.theme') }}:</label>
          <select id="theme-select" v-model="currentTheme">
            <option value="light">{{ t('theme.light') }}</option>
            <option value="dark">{{ t('theme.dark') }}</option>
          </select>
        </div>
        <div class="settings-item">
          <label for="lang-select">//{{ t('nav.lang') }}:</label>
          <select id="lang-select" v-model="lang" @change="changeLanguage()">
            <option value="en-US" class="lang_select_item">EN-US</option>
            <option value="ru-RU" class="lang_select_item">RU-RU</option>
            <option value="ja-JP" class="lang_select_item">JA-JP</option>
          </select>
        </div>
      </div>
    </div>

    <nav v-if="authStore.isAuthenticated" class="navigation">
      <RouterLink to="/">[{{ t('nav.main') }}]</RouterLink>
      <RouterLink :to="{ name: 'admin_articles' }">[{{ t('nav.articles') }}]</RouterLink>
      <RouterLink :to="{ name: 'admin_upload' }">[{{ t('nav.files') }}]</RouterLink>
      <RouterLink :to="{ name: 'admin_create_article' }">[{{ t('nav.create') }}]</RouterLink>
      <RouterLink to="/about">[{{ t('nav.about') }}]</RouterLink>
      <button type="button" class="logout-btn" @click="authStore.logout()">[ {{ t('nav.logout') }} ]</button>
    </nav>

    <nav v-else class="navigation">
      <RouterLink to="/">[{{ t('nav.main') }}]</RouterLink>
      <RouterLink to="/articles">[{{ t('nav.articles') }}]</RouterLink>
      <RouterLink to="/about">[{{ t('nav.about') }}]</RouterLink>
    </nav>
  </header>
</template>

<style scoped>
.header {
  font-family: var(--main-font);
  border-bottom: 2px solid var(--md-accent);
  margin-top: 30px;
  padding-bottom: 12px;
  position: relative;
  background-color: var(--color-surface);
  box-shadow: 0 0 15px rgba(0, 240, 255, 0.1);
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 15px), calc(100% - 15px) 100%, 0 100%);
  padding-top: 16px;
}

.header::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, var(--md-accent) 0%, var(--md-link) 50%, transparent 100%);
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 10px;
  gap: 12px;
}

.logo-link {
  text-decoration: none;
  color: inherit;
  display: inline-block;
}

.logo {
  font-size: 1.8rem;
  font-weight: 900;
  margin: 0;
  line-height: 1;
  font-family: var(--main-font);
  letter-spacing: 0.05em;
  color: #ffffff;
  text-shadow: 0 0 8px var(--md-accent);
  position: relative;
}

.logo::before {
  content: attr(data-text);
  position: absolute;
  left: -2px;
  text-shadow: 2px 0 var(--md-link);
  top: 0;
  color: #ffffff;
  overflow: hidden;
  clip: rect(0, 900px, 0, 0);
  animation: glitch-effect 3s infinite linear alternate-reverse;
}

@keyframes glitch-effect {
  0% { clip: rect(10px, 9999px, 30px, 0); }
  5% { clip: rect(30px, 9999px, 10px, 0); }
  10% { clip: rect(0, 0, 0, 0); }
  100% { clip: rect(0, 0, 0, 0); }
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
  gap: 10px;
  font-family: var(--main-font);
  font-size: 0.85rem;
  font-weight: 700;
  letter-spacing: 0.05em;
}

.settings-item label {
  color: var(--md-accent);
  text-transform: uppercase;
}

.settings-item select {
  font-family: var(--main-font);
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--md-accent);
  background-color: var(--md-bg-pre);
  border: 1px solid var(--md-accent);
  padding: 4px 10px;
  outline: none;
  cursor: pointer;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%);
  transition: all 0.2s ease;
}

.settings-item select:hover,
.settings-item select:focus {
  background-color: var(--md-accent);
  color: #050811;
  box-shadow: 0 0 10px var(--md-accent);
}

.settings-item select option {
  background-color: #050811;
  color: var(--md-accent);
  font-weight: 700;
}

.navigation {
  font-family: var(--main-font);
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 20px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.navigation a {
  text-decoration: none;
  color: var(--md-text-secondary);
  font-size: 0.9rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  transition: all 0.2s ease;
  padding: 4px 6px;
  white-space: nowrap;
}

.navigation a:hover,
.navigation a.router-link-active {
  color: var(--md-link);
  text-shadow: 0 0 8px var(--md-link);
  background-color: rgba(255, 0, 85, 0.1);
}

.logout-btn {
  font-family: var(--main-font);
  font-size: 0.9rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--md-num-negative);
  background: transparent;
  border: 1px solid var(--md-num-negative);
  cursor: pointer;
  padding: 4px 10px;
  transition: all 0.2s ease;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%);
  white-space: nowrap;
}

.logout-btn:hover {
  background-color: var(--md-num-negative);
  color: #ffffff;
  box-shadow: 0 0 10px var(--md-num-negative);
}

@media (max-width: 768px) {
  .header {
    margin-top: 15px;
    padding-top: 12px;
  }

  .header-top {
    flex-direction: column;
    align-items: center;
    gap: 12px;
  }

  .logo {
    font-size: 1.5rem;
  }

  .settings {
    flex-direction: row;
    justify-content: center;
    gap: 12px;
    width: 100%;
  }

  .settings-item {
    font-size: 0.75rem;
  }

  .navigation {
    gap: 8px 12px;
    margin-top: 16px;
  }

  .navigation a,
  .logout-btn {
    font-size: 0.8rem;
    padding: 3px 6px;
  }
}
</style>