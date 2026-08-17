<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
  <div class="not-found-page container">
    <div class="not-found-content">
      <div class="sys-badge">// ERR_404 :: SIGNAL_LOST</div>
      <h1 class="error-code">404</h1>
      <h2 class="error-title">{{ t('not_found.title') }}</h2>
      <p class="error-message">{{ t('not_found.text') }}</p>
      
      <RouterLink to="/" class="home-btn">
        [{{ t('nav.main') }}]
      </RouterLink>
    </div>
  </div>
</template>

<style scoped>
.not-found-page {
  min-height: calc(80vh - 120px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  box-sizing: border-box;
}

.not-found-content {
  text-align: center;
  max-width: 520px;
  width: 100%;
  margin: 0 auto;
  padding: 36px 28px;
  background-color: var(--color-surface);
  border: 1px solid var(--md-border);
  clip-path: polygon(
    0 0,
    calc(100% - 20px) 0,
    100% 20px,
    100% 100%,
    20px 100%,
    0 calc(100% - 20px)
  );
  box-shadow: 0 0 20px rgba(0, 240, 255, 0.15);
  position: relative;
}

.sys-badge {
  font-family: var(--main-font);
  font-size: 0.8rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: var(--md-link);
  margin-bottom: 12px;
  text-shadow: 0 0 8px var(--md-link);
}

.error-code {
  font-family: var(--main-font);
  font-size: 6rem;
  font-weight: 900;
  margin: 0;
  line-height: 1;
  letter-spacing: 0.02em;
  color: #ffffff;
  text-shadow: 0 0 15px var(--md-accent);
}

.error-title {
  font-family: var(--main-font);
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--md-accent);
  margin: 16px 0 12px;
  line-height: 1.3;
  letter-spacing: 0.03em;
  text-transform: uppercase;
  text-shadow: 0 0 8px var(--md-accent);
}

.error-message {
  font-family: var(--main-font);
  font-size: 0.95rem;
  color: var(--md-text-secondary);
  line-height: 1.6;
  margin: 0 0 28px 0;
}

.home-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 10px 28px;
  border: 1px solid var(--md-accent);
  background-color: var(--md-accent);
  color: #050811;
  font-family: var(--main-font);
  font-size: 0.95rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  text-decoration: none;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 8px), calc(100% - 8px) 100%, 0 100%);
  box-shadow: 0 0 12px var(--md-accent);
  transition: all 0.2s ease;
}

.home-btn:hover {
  background-color: #00f0ff;
  border-color: #00f0ff;
  box-shadow: 0 0 18px #00f0ff;
  transform: translateY(-1px);
}

.home-btn:active {
  transform: scale(0.97);
}

@media (max-width: 768px) {
  .not-found-content {
    padding: 28px 18px;
  }

  .error-code {
    font-size: 4.5rem;
  }

  .error-title {
    font-size: 1.25rem;
  }
}
</style>