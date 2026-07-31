<script setup lang="ts">
import ArticleBoxComponent from '@/components/ArticleBoxComponent.vue'
</script>

<template>
  <div class="articles_list">
    <ArticleBoxComponent />
    <ArticleBoxComponent />
  </div>
</template>
<style scoped>
.articles_list {
  width: 100%;
  display: flex;
  gap: 30px;
}

@media (max-width: 768px) {
  .articles_list {
    flex-direction: column;
    display: flex;
    gap: 30px;
  }
}
</style>
