<script setup lang="ts">
interface Props {
  modelValue: boolean
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  variant?: 'error' | 'warning' | 'info'
}

const props = withDefaults(defineProps<Props>(), {
  confirmText: 'Confirmer',
  cancelText: 'Annuler',
  variant: 'error'
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  confirm: []
  cancel: []
}>()

const handleConfirm = () => {
  emit('confirm')
  emit('update:modelValue', false)
}

const handleCancel = () => {
  emit('cancel')
  emit('update:modelValue', false)
}

const variantClass = {
  error: 'btn-error',
  warning: 'btn-warning',
  info: 'btn-info'
}
</script>

<template>
  <dialog :open="modelValue" class="modal">
    <div class="modal-box">
      <h3 class="font-bold text-lg">{{ title }}</h3>
      <p class="py-4" v-html="message"></p>
      
      <div class="modal-action">
        <button class="btn btn-ghost" @click="handleCancel">
          {{ cancelText }}
        </button>
        <button 
          class="btn" 
          :class="variantClass[variant]"
          @click="handleConfirm"
        >
          {{ confirmText }}
        </button>
      </div>
    </div>
    
    <!-- Backdrop cliquable pour fermer -->
    <div class="modal-backdrop" @click="handleCancel"></div>
  </dialog>
</template>