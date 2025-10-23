<script setup lang="ts">
import FormProfile from '@/components/FormProfile.vue'
</script>

<template>
    <div class="min-h-screen flex items-center justify-center bg-gradient-primary p-8">
        <!-- Card de connexion - S'ADAPTE AU THÈME -->
        <div class="bg-base-100 rounded-3xl p-12 shadow-card-hover w-full max-w-md">
            <div class="text-center mb-8">
                <h2 class="text-3xl font-extrabold text-base-content mb-2">Connexion</h2>
                <p class="text-base-content/70">Bienvenue ! Connectez-vous à votre compte</p>
            </div>

            <FormProfile />
        </div>
    </div>
</template>