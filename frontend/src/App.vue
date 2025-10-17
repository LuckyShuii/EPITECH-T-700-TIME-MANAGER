<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import { onMounted, ref } from 'vue'
import TopNavBar from './components/TopNavBar.vue'
import { useAuthStore } from './store/AuthStore'

const theme = ref('dark')
const authStore = useAuthStore()
const route = useRoute()

const toggleTheme = () => {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
}

onMounted(() => {
  authStore.fetchUserProfile()
})

</script>
<template>
  <main>
    <div :data-theme="theme" class="min-h-screen">
      <TopNavBar v-if="!route.meta.hideTopBar" :currentTheme="theme" @toggle-theme="toggleTheme" />
      <RouterView />
    </div>
  </main>
</template>