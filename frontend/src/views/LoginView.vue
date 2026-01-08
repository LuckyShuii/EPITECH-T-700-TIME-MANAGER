<script setup lang="ts">
import FormProfile from '@/components/FormProfile.vue'
</script>

<template>
  <div class="min-h-screen bg-white text-black flex flex-col">
    <!-- Header simple -->
    <header class="border-b-2 border-black p-6">
      <div class="max-w-6xl mx-auto">
        <h1 class="text-2xl font-bold uppercase tracking-wider">TML</h1>
      </div>
    </header>

    <!-- Login Container -->
    <div class="flex-1 flex items-center justify-center px-6 py-12">
      <div class="w-full max-w-md">
        <!-- Title Section -->
        <div class="border-2 border-black p-8 mb-8">
          <h2 class="text-4xl font-black uppercase tracking-tight mb-2">
            CONNEXION
          </h2>
          <p class="text-sm font-bold uppercase tracking-widest text-gray-600">
            Accédez à votre compte TimeTracker
          </p>
        </div>

        <!-- Form Section -->
        <div class="border-2 border-black p-8">
          <FormProfile />
        </div>

        <!-- Info Footer -->
        <div class="border-2 border-black border-t-0 p-6 text-xs font-bold uppercase tracking-widest text-gray-600">
          Identifiants oubliés ? Contactez votre administrateur.
        </div>
      </div>
    </div>

    <!-- Footer -->
    <footer class="border-t-2 border-black bg-black text-white p-8">
      <div class="max-w-6xl mx-auto">
        <div class="flex justify-between items-center">
          <div>
            <p class="font-black text-lg uppercase">TML</p>
            <p class="text-xs text-gray-400 mt-2">TIMETRACKER</p>
          </div>
          <p class="text-xs text-gray-400">2025</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
/* Pas de fioritures */
button {
  transition: none;
}
</style>