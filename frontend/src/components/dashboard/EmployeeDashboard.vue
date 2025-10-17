<script setup lang="ts">
import EmployeeLayout from '../layout/EmployeeLayout.vue'
import ClockWidget from '@/components/widget/ClockWidget.vue'
import CalendarWidget from '@/components/widget/CalendarWidget.vue'
import ClockButton from '@/components/ClockButton.vue'
import { useAuthStore } from '@/store/AuthStore'
import { storeToRefs } from 'pinia'

const authStore = useAuthStore()
const { clockInTime, sessionStatus } = storeToRefs(authStore)
</script>

<template>
  <EmployeeLayout>
    <!-- Position 1 : Haut gauche -->
    <template #widget-1>
      <div class="bg-red-100 p-6 rounded h-full">Widget 1</div>
    </template>

    <!-- Position 3 : Clock (bas gauche) -->
    <template #clock>
      <ClockWidget :clockInTime="clockInTime" :status="sessionStatus" />
      <ClockButton />
    </template>

    <!-- Position 4 : Calendar (haut centre - 2 colonnes) -->
    <template #calendar>
      <div class="bg-blue-100 p-6 rounded h-full">
        <CalendarWidget />
      </div>
    </template>

    <!-- Position 5 : Widget 2 (bas centre gauche) -->
    <template #widget-2>
      <div class="bg-green-100 p-6 rounded h-full">Widget 2</div>
    </template>

    <!-- Position 6 : Widget 3 (bas centre droite) -->
    <template #widget-3>
      <div class="bg-yellow-100 p-6 rounded h-full">Widget 3</div>
    </template>

    <!-- Position 7 : Widget 4 (haut droite) -->
    <template #widget-4>
      <div class="bg-purple-100 p-6 rounded h-full">Widget 4</div>
    </template>

    <!-- Position 8 : Widget 5 (bas droite) -->
    <template #widget-5>
      <div class="bg-orange-100 p-6 rounded h-full">Widget 5</div>
    </template>
  </EmployeeLayout>
</template>