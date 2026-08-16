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
  min-height: calc(100vh - 400px);
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
}

.error-code {
  font-size: 6rem;
  font-weight: 500;
  color: var(--color-text);
  margin: 0;
  line-height: 1;
  letter-spacing: -2px;
  transition: color 0.2s ease;
}

.error-title {
  font-size: 1.5rem;
  font-weight: 400;
  color: var(--md-text);
  margin: 16px 0 12px;
  line-height: 1.3;
  transition: color 0.2s ease;
}

.error-message {
  font-size: 1rem;
  color: var(--md-text-secondary);
  line-height: 1.6;
  margin: 0 0 32px 0;
  transition: color 0.2s ease;
}

@media (max-width: 768px) {
  .error-code {
    font-size: 4.5rem;
  }

  .error-title {
    font-size: 1.25rem;
  }
}
</style>
