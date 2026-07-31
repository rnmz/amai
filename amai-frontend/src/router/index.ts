import { createRouter, createWebHistory } from 'vue-router'

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
      component: () => import("@/views/ArticlePageView.vue")
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
      path: '/admin/upload',
      name: 'admin_upload',
      component: () => import('@/views/UploadView.vue'),
      meta: { isAdminPage: true }
    },
    {
      path: '/admin/create',
      name: 'admin_create_article',
      component: () => import('@/views/ArticlePageView.vue'),
      meta: { isAdminPage: true }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not_found',
      component: () => import ('@/views/NotFoundView.vue'),
      meta: { isAdminPage: false }
    }
  ],
})

export default router
