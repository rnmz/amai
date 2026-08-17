<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
  <div class="not-found-page container">
    <div class="not-found-content">
      <div class="sys-badge">Error 404</div>
      <h1 class="error-code">404</h1>
      <h2 class="error-title">{{ t('not_found.title') }}</h2>
      <p class="error-message">{{ t('not_found.text') }}</p>
      
      <RouterLink to="/" class="home-btn">
        {{ t('nav.main') }}
      </RouterLink>
    </div>
  </div>
</template>

<style scoped>
.not-found-page {
  min-height: calc(75vh - 120px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px 20px;
  box-sizing: border-box;
}

.not-found-content {
  text-align: center;
  max-width: 480px;
  width: 100%;
  margin: 0 auto;
  padding: 32px 24px;
  background-color: var(--color-surface);
  border: 1px solid var(--md-border);
  border-radius: 8px;
  box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.08);
  position: relative;
}

.sys-badge {
  font-family: var(--main-font);
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--md-text-secondary);
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.error-code {
  font-family: var(--main-font);
  font-size: 5rem;
  font-weight: 800;
  margin: 0;
  line-height: 1;
  color: var(--md-accent);
}

.error-title {
  font-family: var(--main-font);
  font-size: 1.35rem;
  font-weight: 700;
  color: var(--md-text);
  margin: 12px 0 8px;
  line-height: 1.3;
}

.error-message {
  font-family: var(--main-font);
  font-size: 0.9rem;
  color: var(--md-text-secondary);
  line-height: 1.5;
  margin: 0 0 24px 0;
}

.home-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 8px 22px;
  border: 1px solid var(--md-accent);
  border-radius: 4px;
  background-color: var(--md-accent);
  color: #ffffff;
  font-family: var(--main-font);
  font-size: 0.9rem;
  font-weight: 600;
  text-decoration: none;
  box-shadow: 1px 1px 2px rgba(0, 0, 0, 0.08);
  transition: all 0.15s ease;
}

.home-btn:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.home-btn:active {
  transform: scale(0.97);
}

@media (max-width: 768px) {
  .not-found-content {
    padding: 24px 16px;
    border-radius: 6px;
  }

  .error-code {
    font-size: 4rem;
  }

  .error-title {
    font-size: 1.15rem;
  }
}
</style>