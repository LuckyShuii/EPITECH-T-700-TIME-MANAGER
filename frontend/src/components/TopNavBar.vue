<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router';
import ThemeToggle from './ThemeToggle.vue';
import { useAuthStore } from '@/store/AuthStore';
import { storeToRefs } from 'pinia';


defineProps<{
    currentTheme: string
}>();

defineEmits<{
    'toggle-theme': []
}>();

const authStore = useAuthStore()

const { isAuthenticated } = storeToRefs(authStore);

const router = useRouter();

const handleLogout = async () => {
    await authStore.logout();
    router.push('/');
}

</script>
<template>
    <div class="navbar shadow-sm sticky top-0 z-50 base-content">
        <div class="flex-1">
            <RouterLink to="/" class="btn btn-ghost text-xl ">TML </RouterLink>
        </div>
        <div class="flex-none gap-2">
            <!-- Passe le thème actuel au bouton -->
            <ThemeToggle :currentTheme="currentTheme" @toggle-theme="$emit('toggle-theme')" />

            <div class="dropdown dropdown-end">
                <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
                    <div class="w-10 rounded-full">
                        <img alt="Tailwind CSS Navbar component"
                            src="https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.webp" />
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
                        <a @click="handleLogout" class="justify-between">
                            Logout
                        </a>
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>