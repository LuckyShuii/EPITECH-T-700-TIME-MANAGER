<script setup lang="ts">
import { RouterView } from 'vue-router'
import { ref } from 'vue'
import OuiView from '@/components/OuiView.vue'

const theme = ref('aqua')

</script>

<template>
    <main>
      <div :data-theme="theme" class="min-h-screen">
        <RouterView />
        <OuiView @toggle-theme="theme = theme === 'dark' ? 'aqua' : 'dark'" />
      </div>
    </main>
</template>


