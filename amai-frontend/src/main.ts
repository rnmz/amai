import { createApp } from 'vue'
import { createPinia } from 'pinia'
import i18n from '@/i18n/index.ts'

import App from './App.vue'
import router from './router'

import './assets/markdown.css'
import 'highlight.js/styles/github.css'
import 'katex/dist/katex.min.css'
import { useAuthStore } from './stores/auth.ts'

async function bootstrap() {
  const app = createApp(App);
  const pinia = createPinia();
  app.use(pinia);

  const authStore = useAuthStore();
  await authStore.checkAuth();

  app.use(router);
  app.use(i18n);
  app.mount("#app");
}

bootstrap();
