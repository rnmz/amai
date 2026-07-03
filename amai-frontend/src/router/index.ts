import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/user/MainView.vue'),
      meta: { isAdminPage: false }
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/user/AboutView.vue'),
      meta: { isAdminPage: false }
    },
    {
      path: '/articles',
      name: 'articles',
      component: () => import('../views/user/ArticlesView.vue'),
      meta: { isAdminPage: false }
    },
    {
      path: '/article/:id',
      name: 'article',
      component: () => import("../views/user/ReadArticleView.vue")
    },
    {
      path: '/admin/articles',
      name: 'admin_articles',
      component: () => import('../views/admin/ArticlesView.vue'),
      meta: { isAdminPage: true }
    }
  ],
})

export default router
