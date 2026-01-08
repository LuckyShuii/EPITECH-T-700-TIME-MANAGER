<script setup lang="ts">
import { computed } from 'vue'
import type { Notification } from '@/store/NotificationsStore';
import { CheckCircleIcon, XCircleIcon, ExclamationTriangleIcon, InformationCircleIcon } from '@heroicons/vue/24/solid'

const props = defineProps<{
  notification: Notification
}>()

const emit = defineEmits<{
  close: []
}>()

const notificationColors = computed(() => {
  const colors: Record<string, { border: string; bg: string; icon: string }> = {
    success: { border: 'border-green-600', bg: 'bg-green-50', icon: 'text-green-600' },
    error: { border: 'border-red-600', bg: 'bg-red-50', icon: 'text-red-600' },
    warning: { border: 'border-yellow-600', bg: 'bg-yellow-50', icon: 'text-yellow-600' },
    info: { border: 'border-blue-600', bg: 'bg-blue-50', icon: 'text-blue-600' }
  }
  return colors[props.notification.status] || colors.info
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
    :class="[notificationColors.border, notificationColors.bg, 'border-2 shadow-lg transition-all duration-300 ease-in-out flex items-center gap-4 p-4']"
  >
    <!-- Icon -->
    <component 
      :is="iconComponent" 
      :class="[notificationColors.icon, 'w-6 h-6 flex-shrink-0']" 
    />

    <!-- Content -->
    <div class="flex-1">
      <h3 class="font-black uppercase text-sm tracking-wider">{{ notification.title }}</h3>
      <p v-if="notification.description" class="text-xs font-bold mt-1">
        {{ notification.description }}
      </p>
    </div>

    <!-- Close Button -->
    <button 
      class="border-2 w-6 h-6 flex items-center justify-center font-bold text-xs hover:bg-black hover:text-white transition-none flex-shrink-0"
      :class="notificationColors.border"
      @click="emit('close')"
    >
      X
    </button>
  </div>
</template>

<style scoped>
button {
  transition: none;
}
</style>