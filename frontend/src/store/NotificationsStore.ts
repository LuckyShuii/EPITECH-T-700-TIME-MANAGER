import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Notification {
  id: string
  status: 'success' | 'error' | 'warning' | 'info'
  title: string
  description?: string
  duration?: number
}

export const useNotificationsStore = defineStore('notifications', () => {
  const notifications = ref<Notification[]>([])

  function addNotification(notification: Omit<Notification, 'id'>) {
    const id = Date.now().toString() + Math.random().toString(36)
    const duration = notification.duration || 4000

    const newNotification: Notification = {
      ...notification,
      id
    }

    notifications.value.push(newNotification)

    // Auto-suppression après la durée
    setTimeout(() => {
      removeNotification(id)
    }, duration)

    return id
  }

  function removeNotification(id: string) {
    notifications.value = notifications.value.filter(n => n.id !== id)
  }

  return {
    notifications,
    addNotification,
    removeNotification
  }
})