<script setup lang="ts">

import { useAuthStore } from '@/store/AuthStore';

import { storeToRefs } from 'pinia';
import { ref } from 'vue';

const authStore = useAuthStore();
const { user } = storeToRefs(authStore);

// Section active (pour gérer la navigation du menu)
const activeSection = ref('public-data');

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
          <li>
            <a 
              :class="{ 'active': activeSection === 'notifications' }"
              @click="activeSection = 'notifications'"
            >
              Notifications
            </a>
          </li>
        </ul>
      </div>

      <!-- Main Content -->
      <div class="flex-1 bg-base-100 rounded-box shadow-sm p-8">
        
        <!-- Section: Données publiques -->
        <div v-if="activeSection === 'public-data'" class="space-y-6">
          <!-- Photo de profil -->
          <div class="flex items-center gap-6">
            <div class="avatar">
              <div class="w-32 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
                <img src="https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.webp" alt="Photo de profil" />
              </div>
            </div>
            <div>
              <button class="btn btn-sm btn-outline">Modifier la photo</button>
            </div>
          </div>

          <!-- Informations personnelles -->
          <div class="divider"></div>
          
          <div class="grid grid-cols-2 gap-4">
            <div class="form-control">
              <label class="label">
                <span class="label-text font-semibold">Prénom</span>
              </label>
              <input 
                type="text" 
                :value="user.first_name" 
                class="input input-bordered" 
                readonly 
              />
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text font-semibold">Nom</span>
              </label>
              <input 
                type="text" 
                :value="user.last_name" 
                class="input input-bordered" 
                readonly 
              />
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text font-semibold">Email</span>
              </label>
              <input 
                type="email" 
                :value="user.email" 
                class="input input-bordered" 
                readonly 
              />
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text font-semibold">Téléphone</span>
              </label>
              <input 
                type="tel" 
                :value="user.phone_number" 
                class="input input-bordered" 
                readonly 
              />
            </div>
          </div>

          <!-- Rôles -->
          <div class="form-control">
            <label class="label">
              <span class="label-text font-semibold">Rôles</span>
            </label>
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
          
          <div class="form-control">
            <label class="label">
              <span class="label-text font-semibold">Nom d'utilisateur</span>
            </label>
            <input 
              type="text" 
              :value="user.username" 
              class="input input-bordered" 
              readonly 
            />
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text font-semibold">UUID</span>
            </label>
            <input 
              type="text" 
              :value="user.user_uuid" 
              class="input input-bordered font-mono text-sm" 
              readonly 
            />
          </div>

          <button class="btn btn-outline">Changer le mot de passe</button>
        </div>

        <!-- Section: Notifications -->
        <div v-if="activeSection === 'notifications'" class="space-y-6">
          <h2 class="text-xl font-semibold">Notifications</h2>
          <p class="text-base-content/70">Section en construction...</p>
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