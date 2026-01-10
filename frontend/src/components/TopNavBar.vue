<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router';
import CustomizeLayoutButton from '@/components/CustomizeLayoutButton.vue'
import { useAuthStore } from '@/store/AuthStore';
import { storeToRefs } from 'pinia';
import { computed } from 'vue'; 
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
const editModeStore = useEditModeStore()

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
  <nav class="sticky top-0 z-50">
    <div class="flex justify-between items-center px-6 py-4 max-w-full">
      <!-- Logo avec couleur d'avatar -->
      <RouterLink :to="isAuthenticated ? '/dashboard' : ''" class="font-black text-2xl uppercase tracking-wider hover:opacity-70">
        <span v-if="isAuthenticated" :class="avatarColor" class="text-white px-2 py-1">
          TML
        </span>
        <span v-else>
          TML
        </span>
      </RouterLink>

      <!-- Right side: Theme + Customize + Avatar -->
      <div class="flex items-center gap-4">
        <!-- Customize Layout Button (si on est sur un dashboard) -->
        <CustomizeLayoutButton v-if="editModeStore.currentDashboard" />

        

        <!-- Avatar Dropdown -->
<div class="dropdown dropdown-end">
  <!-- Avatar Button -->
  <button 
tabindex="0"
:class="[isAuthenticated ? avatarColor : 'bg-gray-300', 'border-2 border-black w-10 h-10 flex items-center justify-center font-bold uppercase text-xs text-white hover:opacity-70 transition-none rounded-none']"
>
    <!-- Si connecté : Initiales -->
    <span v-if="isAuthenticated">
      {{ userInitials }}
    </span>
    <!-- Si NON connecté : U pour User -->
    <span v-else>
      U
    </span>
  </button>
          <!-- Dropdown Menu -->
          <ul 
            tabindex="0" 
            class="dropdown-content menu border-2 border-black bg-white z-40 w-48 p-0 shadow-lg"
          >
            <!-- Profile -->
            <li class="border-b-2 border-black">
              <RouterLink to="/profile" class="p-4 font-bold uppercase text-sm hover:bg-black hover:text-white transition-none rounded-none">
                PROFIL
              </RouterLink>
            </li>

            <!-- Login (si non authentifié) -->
            <li v-if="!isAuthenticated" class="border-b-2 border-black">
              <RouterLink to="/login" class="p-4 font-bold uppercase text-sm hover:bg-black hover:text-white transition-none rounded-none">
                LOGIN
              </RouterLink>
            </li>

            <!-- Dashboard (si authentifié) -->
            <li v-if="isAuthenticated" class="border-b-2 border-black">
              <RouterLink to="/dashboard" class="p-4 font-bold uppercase text-sm hover:bg-black hover:text-white transition-none rounded-none">
                DASHBOARD
              </RouterLink>
            </li>

            <!-- Logout (si authentifié) -->
            <li v-if="isAuthenticated">
              <a 
                @click="handleLogout" 
                class="p-4 font-bold uppercase text-sm hover:bg-black hover:text-white transition-none rounded-none cursor-pointer"
              >
                LOGOUT
              </a>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </nav>
</template>

<style scoped>
/* Pas d'animations fluides */
* {
  transition: none;
}

button:hover,
a:hover {
  transition: none;
}
</style>