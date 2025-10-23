<script setup lang="ts">

import { useAuthStore } from '@/store/AuthStore';

import { storeToRefs } from 'pinia';
import { ref, computed } from 'vue';

const authStore = useAuthStore();
const { user } = storeToRefs(authStore);

// Section active (pour gérer la navigation du menu)
const activeSection = ref('public-data');

// Couleur de fond pour l'avatar
const { avatarColor } = storeToRefs(authStore)


// Couleurs disponibles
const availableColors = [
  { name: 'Primary', class: 'bg-primary' },
  { name: 'Secondary', class: 'bg-secondary' },
  { name: 'Accent', class: 'bg-accent' },
  { name: 'Info', class: 'bg-info' }
];

// Calcul des initiales
const userInitials = computed(() => {
  if (!user.value?.first_name || !user.value?.last_name) return '';
  return `${user.value.first_name.charAt(0)}${user.value.last_name.charAt(0)}`.toUpperCase();
});

</script>


<template>
  <div v-if="user" class="container mx-auto p-6">
    <h1 class="text-3xl font-bold mb-6 text-primary">Mes préférences</h1>
    
    <div class="flex gap-6">
      <!-- Sidebar Menu -->
      <div class="w-64">
        <ul class="menu bg-base-200 rounded-box">
          <li>
            <a 
              :class="{ 'active': activeSection === 'public-data' }"
              @click="activeSection = 'public-data'"
            >
              Données publiques
            </a>
          </li>
          <li>
            <a 
              :class="{ 'active': activeSection === 'connection-data' }"
              @click="activeSection = 'connection-data'"
            >
              Données de connexion
            </a>
          </li>
        </ul>
      </div>

      <!-- Main Content -->
      <div class="flex-1 bg-base-100 rounded-box shadow-sm p-8">
        
        <!-- Section: Données publiques -->
        <div v-if="activeSection === 'public-data'" class="space-y-6">
          <!-- Photo de profil avec initiales -->
          <div class="flex items-center gap-6">
            <div class="avatar">
              <div class="w-32 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
                <div :class="[avatarColor, 'w-full h-full flex items-center justify-center']">
                  <span class="text-5xl font-bold text-white">{{ userInitials }}</span>
                </div>
              </div>
            </div>
            <div class="dropdown">
              <div tabindex="0" role="button" class="btn btn-sm btn-outline">Changer la couleur</div>
              <ul tabindex="0" class="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow">
                <li v-for="color in availableColors" :key="color.class">
                  <a @click="avatarColor = color.class">
                    <div :class="[color.class, 'w-4 h-4 rounded-full']"></div>
                    {{ color.name }}
                  </a>
                </li>
              </ul>
            </div>
          </div>

          <!-- Informations personnelles -->
          <div class="divider"></div>
          
          <div class="grid grid-cols-[180px_1fr] gap-x-4 gap-y-4 items-center">
            <label class="font-semibold text-right">Prénom</label>
            <input 
              type="text" 
              :value="user.first_name" 
              class="input input-bordered" 
              readonly 
            />

            <label class="font-semibold text-right">Nom</label>
            <input 
              type="text" 
              :value="user.last_name" 
              class="input input-bordered" 
              readonly 
            />

            <label class="font-semibold text-right">Email</label>
            <input 
              type="email" 
              :value="user.email" 
              class="input input-bordered" 
              readonly 
            />

            <label class="font-semibold text-right">Téléphone</label>
            <input 
              type="tel" 
              :value="user.phone_number" 
              class="input input-bordered" 
              readonly 
            />

            <label class="font-semibold text-right">Rôles</label>
            <div class="flex gap-2">
              <div v-for="role in user.roles" :key="role" class="badge badge-primary">
                {{ role }}
              </div>
            </div>
          </div>
        </div>

        <!-- Section: Données de connexion -->
        <div v-if="activeSection === 'connection-data'" class="space-y-6">
          <h2 class="text-xl font-semibold">Données de connexion</h2>
          
          <div class="grid grid-cols-[180px_1fr] gap-x-4 gap-y-4 items-center">
            <label class="font-semibold text-right">Nom d'utilisateur</label>
            <input 
              type="text" 
              :value="user.username" 
              class="input input-bordered" 
              readonly 
            />
          </div>

          <button class="btn btn-outline">Changer le mot de passe</button>
        </div>
      </div>
    </div>
  </div>

  <!-- Si pas d'utilisateur connecté -->
  <div v-else class="container mx-auto p-6">
    <div class="alert alert-warning">
      <span>Vous devez être connecté pour voir cette page.</span>
    </div>
  </div>
</template>