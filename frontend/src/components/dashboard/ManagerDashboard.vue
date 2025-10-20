<script setup lang="ts">
import ManagerLayout from '../layout/ManagerLayout.vue'
import ClockWidget from '@/components/widget/ClockWidget.vue'
import CalendarWidget from '@/components/widget/CalendarWidget.vue'
import ClockButton from '@/components/ClockButton.vue'
import { useAuthStore } from '@/store/AuthStore'
import { storeToRefs } from 'pinia'

const authStore = useAuthStore()
const { clockInTime, sessionStatus } = storeToRefs(authStore)
</script>

<template>
  <ManagerLayout>
    <!-- Clock + Boutons -->
    <template #clock>
      <ClockWidget
        :clockInTime="clockInTime"
        :status="sessionStatus"
      />
      <ClockButton />
    </template>

    <!-- KPI Stats -->
    <template #kpi-stats>
      <div class="bg-purple-100 p-6 rounded h-full">KPI Stats</div>
    </template>

    <!-- Calendrier -->
    <template #calendar>
      <div class="bg-blue-100 p-6 rounded h-full">
        <CalendarWidget />
      </div>
    </template>

    <!-- Présence équipe -->
    <template #team-presence>
      <div class="bg-green-100 p-6 rounded h-full">Présence de l'équipe</div>
    </template>

    <!-- KPI Carousel -->
    <template #kpi-carousel>
      <div class="bg-yellow-100 p-6 rounded h-full">KPI Carousel</div>
    </template>

    <!-- Bouton rapport -->
    <template #report-button>
      <div class="bg-orange-100 p-6 rounded h-full">
        <button class="btn btn-primary w-full">📊 Rapport</button>
      </div>
    </template>
  </ManagerLayout>
</template>