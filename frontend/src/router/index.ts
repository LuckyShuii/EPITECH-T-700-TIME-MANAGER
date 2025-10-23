import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue';
import UserProfileview from '@/views/UserProfileview.vue';
import LoginView from '@/views/LoginView.vue';
import NotFoundView from '@/views/NotFoundView.vue';
import DashBoardView from '@/views/DashboardView.vue'

import { useAuthStore } from '@/store/AuthStore';



const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { hideTopBar: true }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { hideTopBar: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: UserProfileview
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: NotFoundView,
      meta: { hideTopBar: true }
    },
    {
      path: '/dashboard',
      name: 'dashboard-employee',
      component: DashBoardView
    }
  ]
})
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  // Pages publiques (accessibles sans auth)
  const publicPages = ['/', '/login']
  const isPublicRoute = publicPages.includes(to.path) || to.name === 'not-found'
  const authRequired = !isPublicRoute
  if (authRequired && !authStore.isAuthenticated) {
    return next('/login') // Redirige vers login
  }
  next()
})


export default router