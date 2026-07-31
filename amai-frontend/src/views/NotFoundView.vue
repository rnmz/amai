<script setup lang="ts">
</script>

<template>
  <div class="not-found-page container">
    <div class="not-found-content">
      <h1 class="error-code">404</h1>
      <h2 class="error-title">Page not found / ページが見つかりません</h2>
      <p class="error-message">
     </p>
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
  color: #111;
  margin: 0;
  line-height: 1;
  letter-spacing: -2px;
}

.error-title {
  font-size: 1.5rem;
  font-weight: 400;
  color: #222;
  margin: 16px 0 12px;
  line-height: 1.3;
}

.error-message {
  font-size: 1rem;
  color: #666;
  line-height: 1.6;
  margin: 0 0 32px 0;
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
