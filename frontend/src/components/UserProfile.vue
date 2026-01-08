<script setup lang="ts">

import { useAuthStore } from '@/store/AuthStore';
import { storeToRefs } from 'pinia';
import { ref, computed } from 'vue';

const authStore = useAuthStore();
const { user } = storeToRefs(authStore);

const activeSection = ref('public-data');
const { avatarColor } = storeToRefs(authStore)

const availableColors = [
  { name: 'Bleu', class: 'bg-blue-600' },
  { name: 'Rouge', class: 'bg-red-600' },
  { name: 'Vert', class: 'bg-green-600' },
  { name: 'Noir', class: 'bg-black' }
];

const userInitials = computed(() => {
  if (!user.value?.first_name || !user.value?.last_name) return '';
  return `${user.value.first_name.charAt(0)}${user.value.last_name.charAt(0)}`.toUpperCase();
});

</script>


<template>
  <div v-if="user" class="min-h-screen bg-white text-black">
    <!-- Header -->
    <div class="border-b-2 border-black p-8">
      <div class="max-w-6xl mx-auto">
        <h1 class="text-4xl font-black uppercase tracking-tight">MES PRÉFÉRENCES</h1>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-6xl mx-auto p-8">
      <div class="grid grid-cols-[250px_1fr] gap-8">
        
        <!-- Sidebar Menu -->
        <div class="border-2 border-black p-0">
          <button 
            @click="activeSection = 'public-data'"
            class="w-full p-4 border-b-2 border-black font-bold uppercase text-sm text-left hover:bg-black hover:text-white transition-none"
            :class="{ 'bg-black text-white': activeSection === 'public-data' }"
          >
            DONNÉES PUBLIQUES
          </button>
          <button 
            @click="activeSection = 'connection-data'"
            class="w-full p-4 font-bold uppercase text-sm text-left hover:bg-black hover:text-white transition-none"
            :class="{ 'bg-black text-white': activeSection === 'connection-data' }"
          >
            DONNÉES CONNEXION
          </button>
        </div>

        <!-- Main Content Area -->
        <div class="border-2 border-black p-8">
          
          <!-- Section: Données publiques -->
          <div v-if="activeSection === 'public-data'" class="space-y-8">
            
            <!-- Avatar Section -->
            <div class="border-2 border-black p-8">
              <div class="flex items-center gap-8">
                <!-- Avatar Carré Brutal -->
                <div class="border-4 border-black w-32 h-32 flex items-center justify-center flex-shrink-0" :class="avatarColor">
                  <span class="text-5xl font-black text-white">{{ userInitials }}</span>
                </div>
                
                <!-- Color Picker -->
                <div class="space-y-4">
                  <p class="font-bold uppercase text-sm">CHANGER LA COULEUR</p>
                  <div class="grid grid-cols-2 gap-2">
                    <button 
                      v-for="color in availableColors" 
                      :key="color.class"
                      @click="avatarColor = color.class"
                      class="border-2 border-black p-3 font-bold uppercase text-xs hover:bg-black hover:text-white transition-none"
                      :class="{ 'bg-black text-white': avatarColor === color.class }"
                    >
                      {{ color.name }}
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- Informations personnelles -->
            <div class="space-y-4">
              <h3 class="font-black uppercase tracking-wider text-lg">INFORMATIONS</h3>
              
              <div class="space-y-4">
                <!-- Prénom -->
                <div class="grid grid-cols-[150px_1fr] gap-4 items-center border-b-2 border-black pb-4">
                  <label class="font-bold uppercase text-xs tracking-widest">Prénom</label>
                  <input 
                    type="text" 
                    :value="user.first_name" 
                    class="border-2 border-black p-3 font-bold uppercase text-sm bg-white"
                    readonly 
                  />
                </div>

                <!-- Nom -->
                <div class="grid grid-cols-[150px_1fr] gap-4 items-center border-b-2 border-black pb-4">
                  <label class="font-bold uppercase text-xs tracking-widest">Nom</label>
                  <input 
                    type="text" 
                    :value="user.last_name" 
                    class="border-2 border-black p-3 font-bold uppercase text-sm bg-white"
                    readonly 
                  />
                </div>

                <!-- Email -->
                <div class="grid grid-cols-[150px_1fr] gap-4 items-center border-b-2 border-black pb-4">
                  <label class="font-bold uppercase text-xs tracking-widest">Email</label>
                  <input 
                    type="email" 
                    :value="user.email" 
                    class="border-2 border-black p-3 font-bold uppercase text-sm bg-white"
                    readonly 
                  />
                </div>

                <!-- Téléphone -->
                <div class="grid grid-cols-[150px_1fr] gap-4 items-center border-b-2 border-black pb-4">
                  <label class="font-bold uppercase text-xs tracking-widest">Téléphone</label>
                  <input 
                    type="tel" 
                    :value="user.phone_number" 
                    class="border-2 border-black p-3 font-bold uppercase text-sm bg-white"
                    readonly 
                  />
                </div>

                <!-- Rôles -->
                <div class="grid grid-cols-[150px_1fr] gap-4 items-start pt-4">
                  <label class="font-bold uppercase text-xs tracking-widest">Rôles</label>
                  <div class="flex flex-wrap gap-2">
                    <span v-for="role in user.roles" :key="role" class="border-2 border-black px-3 py-2 font-bold uppercase text-xs">
                      {{ role }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Section: Données de connexion -->
          <div v-if="activeSection === 'connection-data'" class="space-y-8">
            <h2 class="text-2xl font-black uppercase tracking-wider">DONNÉES DE CONNEXION</h2>
            
            <div class="space-y-4">
              <!-- Username -->
              <div class="grid grid-cols-[150px_1fr] gap-4 items-center border-b-2 border-black pb-4">
                <label class="font-bold uppercase text-xs tracking-widest">Username</label>
                <input 
                  type="text" 
                  :value="user.username" 
                  class="border-2 border-black p-3 font-bold uppercase text-sm bg-white"
                  readonly 
                />
              </div>
            </div>

            <!-- Change Password Button -->
            <button class="brutal-btn brutal-btn-warning">
              CHANGER LE MOT DE PASSE
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Si pas d'utilisateur connecté -->
  <div v-else class="min-h-screen bg-white text-black flex items-center justify-center">
    <div class="border-2 border-red-700 bg-red-50 p-8 max-w-md text-center">
      <p class="font-bold uppercase text-red-700">ERREUR</p>
      <p class="mt-2 font-bold">Vous devez être connecté pour voir cette page.</p>
    </div>
  </div>
</template>

<style scoped>
button {
  transition: none;
}
</style>