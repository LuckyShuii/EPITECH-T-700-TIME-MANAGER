import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue';
//import Profile from '@/view/ProfileView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    } 
    // {
    //   path: '/',
    //   name: 'profile',
    //   component: Profile
    // }
  ]
})
export default router
