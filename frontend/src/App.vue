<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import { onMounted, ref } from 'vue'
import TopNavBar from './components/TopNavBar.vue'
import ToastNotification from './components/ToastNotification.vue'
import { useAuthStore } from './store/AuthStore'
import { useNotificationsStore } from './store/NotificationsStore'


const theme = ref('light')
const authStore = useAuthStore()
const notificationsStore = useNotificationsStore()
const route = useRoute()


onMounted(() => {
  authStore.fetchUserProfile()
})
</script>

<template>
  <main>
    <div :data-theme="theme" class="min-h-screen">
      <RouterView />
      
      <!-- Conteneur des notifications -->
      <div class="toast toast-top toast-end z-50">
        <TransitionGroup name="slide-fade">
          <ToastNotification
            v-for="notification in notificationsStore.notifications"
            :key="notification.id"
            :notification="notification"
            @close="notificationsStore.removeNotification(notification.id)"
          />
        </TransitionGroup>
      </div>
    </div>
  </main>
</template>

<style scoped>
.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.2s ease-in;
}

.slide-fade-enter-from {
  transform: translateX(20px);
  opacity: 0;
}

.slide-fade-leave-to {
  transform: translateX(20px);
  opacity: 0;
}
</style>