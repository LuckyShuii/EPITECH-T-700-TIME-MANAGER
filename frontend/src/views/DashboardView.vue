<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/store/AuthStore'
import EmployeeDashboard from '@/components/dashboard/EmployeeDashboard.vue'
import AdminDashboard from '@/components/dashboard/AdminDashboard.vue'

const authStore = useAuthStore()
const { user } = storeToRefs(authStore)

// Computed qui détermine quel dashboard afficher
const currentDashboard = computed(() => {
  if (authStore.isAdmin) {
    return AdminDashboard
  }
  
  // Pour l'instant, manager et employee utilisent le même dashboard
  return EmployeeDashboard
})
</script>

<template>
  <div class="p-6">
    <!-- Dynamic component : affiche le dashboard selon le rôle -->
    <component :is="currentDashboard" />
  </div>
</template>