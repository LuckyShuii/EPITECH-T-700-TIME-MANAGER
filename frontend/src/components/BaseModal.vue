<script setup lang="ts">
import { watch } from 'vue'

interface Props {
  modelValue: boolean  // v-model pour ouvrir/fermer
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  size: 'lg'
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

// Fonction pour fermer le modal
const closeModal = () => {
  emit('update:modelValue', false)
}

// Gérer la touche ESC
const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && props.modelValue) {
    closeModal()
  }
}

// Ajouter/retirer l'event listener
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    document.addEventListener('keydown', handleKeydown)
  } else {
    document.removeEventListener('keydown', handleKeydown)
  }
})
</script>

<template>
  <Transition name="modal">
    <div 
      v-if="modelValue"
      class="modal modal-open"
      @click.self="closeModal"
    >
      <div 
        class="modal-box"
        :class="{
          'max-w-sm': size === 'sm',
          'max-w-2xl': size === 'md',
          'max-w-5xl': size === 'lg',
          'max-w-7xl': size === 'xl'
        }"
      >
        <!-- Header avec titre et bouton fermer -->
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-bold text-lg">{{ title }}</h3>
          <button 
            class="btn btn-sm btn-circle btn-ghost"
            @click="closeModal"
          >
            ✕
          </button>
        </div>

        <!-- Contenu du modal (slot) -->
        <div class="py-4">
          <slot></slot>
        </div>

        <!-- Footer optionnel (slot) -->
        <div v-if="$slots.footer" class="modal-action">
          <slot name="footer"></slot>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
/* Animations pour le modal */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-box,
.modal-leave-active .modal-box {
  transition: transform 0.3s ease;
}

.modal-enter-from .modal-box,
.modal-leave-to .modal-box {
  transform: scale(0.95) translateY(-20px);
}
</style>