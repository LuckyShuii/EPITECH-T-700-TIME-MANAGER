<script setup lang="ts">
import { computed } from 'vue'
import type { Notification } from '@/store/notificationsStore'
import { CheckCircleIcon, XCircleIcon, ExclamationTriangleIcon, InformationCircleIcon } from '@heroicons/vue/24/solid'

const props = defineProps<{
  notification: Notification
}>()

const emit = defineEmits<{
  close: []
}>()

const alertClass = computed(() => {
  const classes: Record<string, string> = {
    success: 'alert-success',
    error: 'alert-error',
    warning: 'alert-warning',
    info: 'alert-info'
  }
  return classes[props.notification.status] || 'alert-info'
})

const iconComponent = computed(() => {
  const icons = {
    success: CheckCircleIcon,
    error: XCircleIcon,
    warning: ExclamationTriangleIcon,
    info: InformationCircleIcon
  }
  return icons[props.notification.status] || InformationCircleIcon
})
</script>

<template>
  <div 
    class="alert shadow-lg transition-all duration-300 ease-in-out flex items-center"
    :class="alertClass"
  >
    <component :is="iconComponent" class="w-6 h-6 flex-shrink-0" />
    <div class="flex-1">
      <h3 class="font-bold">{{ notification.title }}</h3>
      <div v-if="notification.description" class="text-xs">
        {{ notification.description }}
      </div>
    </div>
    <div class="flex-none">
      <button 
        class="btn btn-sm btn-ghost"
        @click="emit('close')"
      >
        ✕
      </button>
    </div>
  </div>
</template>