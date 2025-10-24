<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router';
import ThemeToggle from './ThemeToggle.vue';
import CustomizeLayoutButton from '@/components/CustomizeLayoutButton.vue'
import { useAuthStore } from '@/store/AuthStore';
import { storeToRefs } from 'pinia';
import { computed, inject } from 'vue'; 
import { useEditModeStore } from '@/store/EditModeStore'

const props = defineProps<{
  currentTheme: string
}>();

const emit = defineEmits<{
  'toggle-theme': []
}>();

const authStore = useAuthStore()
const { isAuthenticated, user, avatarColor } = storeToRefs(authStore);
const router = useRouter();

// AJOUTE CES LIGNES ↓
const editModeStore = useEditModeStore()

console.log('🔍 Dans NavBar - Dashboard actif:', editModeStore.currentDashboard)

const userInitials = computed(() => {
  if (!user.value?.first_name || !user.value?.last_name) return '';
  return `${user.value.first_name.charAt(0)}${user.value.last_name.charAt(0)}`.toUpperCase();
});

const handleLogout = async () => {
  await authStore.logout();
  router.push('/');
}
</script>

<template>
  <div class="navbar shadow-sm sticky top-0 z-50 base-content">
    <div class="flex-1">
      <RouterLink :to="isAuthenticated ? '/dashboard' : '/'" class="btn btn-ghost text-xl">
        TML
      </RouterLink>
    </div>
    
    <div class="flex-none gap-2">
      <CustomizeLayoutButton v-if="editModeStore.currentDashboard" />
      
      <ThemeToggle :currentTheme="props.currentTheme" @toggle-theme="emit('toggle-theme')" />
      
      <div class="dropdown dropdown-end">
        <div tabindex="0" class="btn btn-ghost btn-circle avatar">
          <div class="w-10 rounded-full">
            <!-- Si connecté : Initiales -->
            <div v-if="isAuthenticated" :class="[avatarColor, 'w-full h-full flex items-center justify-center']">
              <span class="text-primary-content font-semibold text-sm">{{ userInitials }}</span>
            </div>
            <!-- Si NON connecté : Icône utilisateur générique -->
            <div v-else class="w-full h-full bg-base-300 flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-base-content/50" fill="none"
                viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </div>
          </div>
        </div>
        <ul tabindex="0" class="menu menu-sm dropdown-content bg-base-100 rounded-box z-1 mt-3 w-52 p-2 shadow">
          <li>
            <RouterLink to="profile" class="justify-between">
              Profile
            </RouterLink>
          </li>
          <li v-if="!isAuthenticated">
            <RouterLink to="login" class="justify-between">
              Login
            </RouterLink>
          </li>
          <li v-if="isAuthenticated">
            <RouterLink to="dashboard" class="justify-between">
              Dashboard
            </RouterLink>
          </li>
          <li v-if="isAuthenticated">
            <a @click="handleLogout" class="justify-between">
              Logout
            </a>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>