<script setup>
const emit = defineEmits(['toggle-theme'])
</script>

<template>
  <button @click="emit('toggle-theme')" class="btn">
    Change theme
  </button>
</template>