<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goToLogin = () => {
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen bg-white text-black flex flex-col">
    <!-- Header/Navbar -->
    <header class="border-b-2 border-black p-6">
      <div class="max-w-6xl mx-auto flex justify-between items-center">
        <h1 class="text-2xl font-bold uppercase tracking-wider">TML</h1>
        <button @click="goToLogin" class="brutal-btn brutal-btn-primary text-sm">
          LOGIN
        </button>
      </div>
    </header>

    <!-- Hero Section -->
    <section class="flex-1 flex flex-col items-center justify-center px-6 py-20">
      <div class="max-w-4xl w-full space-y-12">
        <!-- Titre Principal -->
        <div class="border-2 border-black p-12">
          <h2 class="text-5xl md:text-6xl font-black uppercase tracking-tight leading-tight">
            Pointage<br/>Temps Réel
          </h2>
        </div>

        <!-- Sous-titre -->
        <div class="space-y-4">
          <p class="text-xl font-bold uppercase tracking-wide">SYSTÈME SIMPLE ET EFFICACE</p>
          <div class="border-l-4 border-black pl-6">
            <p class="text-base leading-relaxed">
              Clock in. Clock out. Consultez vos stats. C'est tout ce qu'il faut.
            </p>
          </div>
        </div>

        <!-- CTA Principal -->
        <div class="flex gap-4">
          <button 
            @click="goToLogin"
            class="brutal-btn brutal-btn-success flex-1"
          >
            Commencer
          </button>
        </div>
      </div>
    </section>

    <!-- Features Section - Brutal -->
    <section class="border-t-2 border-black bg-black text-white">
      <div class="max-w-6xl mx-auto p-12">
        <h3 class="text-3xl font-black uppercase tracking-wider mb-12">FONCTIONNALITÉS</h3>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- Feature 1 -->
          <div class="border-2 border-white p-8">
            <p class="text-sm font-bold uppercase tracking-widest text-gray-300 mb-4">01</p>
            <h4 class="text-2xl font-black uppercase mb-4">POINTAGE</h4>
            <p class="text-sm leading-relaxed">
              Un seul bouton. Entrez. Sortez. Pas d'excuses, pas de complications.
            </p>
          </div>

          <!-- Feature 2 -->
          <div class="border-2 border-white p-8">
            <p class="text-sm font-bold uppercase tracking-widest text-gray-300 mb-4">02</p>
            <h4 class="text-2xl font-black uppercase mb-4">HISTORIQUE</h4>
            <p class="text-sm leading-relaxed">
              Visualisez chaque jour. Chaque heure. Chaque pause. Totale transparence.
            </p>
          </div>

          <!-- Feature 3 -->
          <div class="border-2 border-white p-8">
            <p class="text-sm font-bold uppercase tracking-widest text-gray-300 mb-4">03</p>
            <h4 class="text-2xl font-black uppercase mb-4">STATISTIQUES</h4>
            <p class="text-sm leading-relaxed">
              Vos données en temps réel. Moyennes. Totaux. Sans filtre.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section - Roles -->
    <section class="border-t-2 border-black">
      <div class="max-w-6xl mx-auto p-12">
        <h3 class="text-3xl font-black uppercase tracking-wider mb-12">POUR CHACUN</h3>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- Employees -->
          <div class="border-2 border-black p-8">
            <h4 class="text-xl font-black uppercase mb-6 tracking-wider">EMPLOYÉ</h4>
            <ul class="space-y-3">
              <li class="text-sm font-bold uppercase">- Pointage rapide</li>
              <li class="text-sm font-bold uppercase">- Historique personnel</li>
              <li class="text-sm font-bold uppercase">- Vue équipe</li>
            </ul>
          </div>

          <!-- Managers -->
          <div class="border-2 border-black p-8">
            <h4 class="text-xl font-black uppercase mb-6 tracking-wider">MANAGER</h4>
            <ul class="space-y-3">
              <li class="text-sm font-bold uppercase">- Dashboard équipe</li>
              <li class="text-sm font-bold uppercase">- Présence en temps réel</li>
              <li class="text-sm font-bold uppercase">- Rapports détaillés</li>
            </ul>
          </div>

          <!-- Admins -->
          <div class="border-2 border-black p-8">
            <h4 class="text-xl font-black uppercase mb-6 tracking-wider">ADMIN</h4>
            <ul class="space-y-3">
              <li class="text-sm font-bold uppercase">- Gestion complète</li>
              <li class="text-sm font-bold uppercase">- Configuration équipes</li>
              <li class="text-sm font-bold uppercase">- Vue d'ensemble</li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Final -->
    <section class="border-t-2 border-black bg-black text-white">
      <div class="max-w-6xl mx-auto p-12 flex flex-col items-center justify-center">
        <h2 class="text-4xl font-black uppercase tracking-wider mb-8 text-center">
          COMMENCEZ MAINTENANT
        </h2>
        <button 
          @click="goToLogin"
          class="brutal-btn brutal-btn-success"
        >
          Accéder à TimeTracker
        </button>
      </div>
    </section>

    <!-- Footer -->
    <footer class="border-t-2 border-black bg-black text-white p-8">
      <div class="max-w-6xl mx-auto">
        <div class="flex justify-between items-center">
          <div>
            <p class="font-black text-lg uppercase">TML</p>
            <p class="text-xs text-gray-400 mt-2">TIMETRACKER</p>
          </div>
          <p class="text-xs text-gray-400">2025 - TOUS DROITS RÉSERVÉS</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
/* Pas d'animations fluides - juste du brutal */
button {
  transition: none;
}
</style>