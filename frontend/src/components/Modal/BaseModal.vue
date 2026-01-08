<script setup lang="ts">
import { watch } from 'vue'

interface Props {
  modelValue?: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
  title: '',
  size: 'lg'
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const closeModal = () => {
  emit('update:modelValue', false)
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && props.modelValue) {
    closeModal()
  }
}

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
        class="modal-box border-2 border-black rounded-none"
        :class="{
          'max-w-sm': size === 'sm',
          'max-w-2xl': size === 'md',
          'max-w-5xl': size === 'lg',
          'max-w-7xl': size === 'xl'
        }"
      >
        <!-- Header avec titre et bouton fermer -->
        <div class="flex justify-between items-center mb-4 pb-2 border-b-2 border-black">
          <h3 class="font-bold text-lg uppercase">{{ title }}</h3>
          <button 
            class="border-2 border-black w-6 h-6 flex items-center justify-center text-sm font-bold hover:bg-black hover:text-white"
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
        <div v-if="$slots.footer" class="modal-action border-t-2 border-black pt-4 mt-4 flex gap-2 justify-end">
          <slot name="footer"></slot>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
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