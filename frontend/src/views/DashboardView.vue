<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/store/AuthStore'
import EmployeeDashboard from '@/components/dashboard/EmployeeDashboard.vue'
import ManagerDashboard from '@/components/dashboard/ManagerDashboard.vue'
import AdminDashboard from '@/components/dashboard/AdminDashboard.vue'

const authStore = useAuthStore()

// Computed qui détermine quel dashboard afficher
const currentDashboard = computed(() => {
  if (authStore.isAdmin) {
    return AdminDashboard
  }
  if (authStore.isManager) {
    return ManagerDashboard
  }
  return EmployeeDashboard
})
</script>

<template>
  <div class="p-6">
    <component :is="currentDashboard" />
  </div>
</template>