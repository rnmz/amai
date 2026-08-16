import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/MainView.vue'),
      meta: { isAdminPage: false }
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('@/views/AboutView.vue'),
      meta: { isAdminPage: false }
    },
    {
      path: '/articles',
      name: 'articles',
      component: () => import('@/views/ArticlesView.vue'),
      meta: { isAdminPage: false }
    },
    {
      path: '/article/:id',
      name: 'article',
      component: () => import("@/views/ArticlePageView.vue"),
      meta: { isAdminPage: false }
    },
    {
      path: '/admin/articles',
      name: 'admin_articles',
      component: () => import('@/views/ArticlesView.vue'),
      meta: { isAdminPage: true }
    },
    {
      path: '/admin/article/:id',
      name: 'admin_article',
      component: () => import('@/views/ArticlePageView.vue'),
      meta: { isAdminPage: true }
    },
    {
      path: '/admin/files',
      name: 'admin_upload',
      component: () => import('@/views/FilesView.vue'),
      meta: { isAdminPage: true }
    },
    {
      path: '/admin/create',
      name: 'admin_create_article',
      component: () => import('@/views/CreateArticleView.vue'),
      meta: { isAdminPage: true }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { hideHeader: true, isAdminPage: false },
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not_found',
      component: () => import('@/views/NotFoundView.vue'),
      meta: { isAdminPage: false }
    }
  ],
})

router.beforeEach((to) => {
  if (!to.meta.isAdminPage) return true

  const authStore = useAuthStore()
  if (!authStore.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  return true
})

export default router