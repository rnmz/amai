<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n();
</script>

<template>
  <div class="not-found-page container">
    <div class="not-found-content">
      <h1 class="error-code">404</h1>
      <h2 class="error-title">{{ t('not_found.title')}}</h2>
      <p class="error-message">{{ t('not_found.text') }}</p>
    </div>
  </div>
</template>

<style scoped>
.not-found-page {
  min-height: calc(100vh - 300px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding-top: 40px;
  padding-bottom: 60px;
}

.not-found-content {
  text-align: center;
  max-width: 540px;
  margin: 0 auto;
  padding: 40px 24px;
  background-color: var(--color-surface);
  border: none;
  border-radius: 20px;
  box-shadow: none;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

.error-code {
  font-family: var(--main-font);
  font-size: 6.5rem;
  font-weight: 800;
  margin: 0;
  line-height: 1;
  letter-spacing: -0.04em;
  background: linear-gradient(135deg, var(--md-accent) 0%, var(--md-link) 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  transition: all 0.25s ease;
}

.error-title {
  font-family: var(--main-font);
  font-size: 1.6rem;
  font-weight: 700;
  color: var(--color-text);
  margin: 16px 0 12px;
  line-height: 1.3;
  transition: color 0.25s ease;
}

.error-message {
  font-family: var(--main-font);
  font-size: 1rem;
  color: var(--md-text-secondary);
  line-height: 1.6;
  margin: 0 0 32px 0;
  transition: color 0.25s ease;
}

.home-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 12px 28px;
  border-radius: 10px;
  background-color: var(--md-accent);
  color: #ffffff;
  font-family: var(--main-font);
  font-size: 0.95rem;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.home-btn:hover {
  background-color: var(--md-accent-hover);
  transform: translateY(-1px);
}

.home-btn:active {
  transform: scale(0.97);
}

@media (max-width: 768px) {
  .not-found-content {
    padding: 32px 18px;
  }

  .error-code {
    font-size: 4.5rem;
  }

  .error-title {
    font-size: 1.3rem;
  }
}
</style>
