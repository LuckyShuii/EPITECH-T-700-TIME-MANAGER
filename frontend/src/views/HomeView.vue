<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goToLogin = () => {
  router.push('/login')
}

const features = [
  {
    icon: '⚡',
    title: 'Pointage en un clic',
    description: 'Clock in/out instantané, où que vous soyez'
  },
  {
    icon: '📅',
    title: 'Suivez votre temps',
    description: 'Visualisez votre historique et calendrier de présence'
  },
  {
    icon: '📊',
    title: 'Rapports en temps réel',
    description: 'Consultez vos statistiques et totaux automatiquement'
  }
]
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-primary/10 via-base-100 to-secondary/10">
    <!-- Hero Section -->
    <div class="hero min-h-[70vh] px-4">
      <div class="hero-content text-center max-w-4xl">
        <div class="space-y-8">
          <!-- Titre principal -->
          <h1 class="text-5xl md:text-7xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent animate-fade-in">
            Gérez votre temps de travail en toute simplicité
          </h1>
          
          <!-- Sous-titre -->
          <p class="text-xl md:text-2xl text-base-content/70 max-w-2xl mx-auto">
            TimeTracker vous aide à pointer, suivre et analyser votre temps de travail efficacement
          </p>
          
          <!-- CTA Principal -->
          <button 
            @click="goToLogin"
            class="btn btn-primary btn-lg text-lg px-8 shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-105"
          >
            Commencer
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Features Section -->
    <div class="container mx-auto px-4 py-16 max-w-6xl">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div 
          v-for="(feature, index) in features" 
          :key="index"
          class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 hover:-translate-y-2"
        >
          <div class="card-body items-center text-center space-y-4">
            <!-- Icône -->
            <div class="text-6xl mb-2">
              {{ feature.icon }}
            </div>
            
            <!-- Titre -->
            <h3 class="card-title text-2xl">
              {{ feature.title }}
            </h3>
            
            <!-- Description -->
            <p class="text-base-content/70">
              {{ feature.description }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- CTA Final -->
    <div class="bg-gradient-to-r from-primary to-secondary py-16">
      <div class="container mx-auto px-4 text-center">
        <h2 class="text-3xl md:text-4xl font-bold text-white mb-6">
          Prêt à gagner du temps ?
        </h2>
        <button 
          @click="goToLogin"
          class="btn btn-lg bg-white text-primary hover:bg-base-200 border-none shadow-xl hover:shadow-2xl transition-all duration-300 hover:scale-105"
        >
          Se connecter
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Footer Simple -->
    <footer class="footer footer-center p-10 bg-base-200 text-base-content">
      <div>
        <p class="font-bold text-lg">TimeTracker</p>
        <p class="text-sm opacity-70">Gestion du temps de travail simplifiée</p>
        <p class="text-xs opacity-50 mt-4">© 2025 TimeTracker. Tous droits réservés.</p>
      </div>
    </footer>
  </div>
</template>

<style scoped>
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fade-in 0.8s ease-out;
}
</style>