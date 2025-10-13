<script setup lang="ts">
import { MoonIcon, SunIcon } from '@heroicons/vue/24/outline';

defineProps<{
  currentTheme: string
}>();

defineEmits<{
  'toggle-theme': []
}>();

</script>
<template>
  <button @click="$emit('toggle-theme')" class="btn btn-ghost btn-circle">
    <MoonIcon v-if="currentTheme === 'light'" class="h-5 w-5" />
    <SunIcon v-else class="h-5 w-5" />
  </button>
</template>