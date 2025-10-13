<script setup lang="ts">
import { RouterView } from 'vue-router'
import { ref } from 'vue'
import TopNavBar from './components/TopNavBar.vue'

const theme = ref('dark')

const toggleTheme = () => {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
}
</script>
<template>
  <main>
    <div :data-theme="theme" class="min-h-screen">
      <TopNavBar :currentTheme="theme" @toggle-theme="toggleTheme" />
      <RouterView />
    </div>
  </main>
</template>