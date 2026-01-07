<script setup lang="ts">
import DraggableLayout from './DraggableLayout.vue'
</script>

<template>
  <DraggableLayout dashboard-name="admin">
    <template #add-employee>
      <slot name="add-employee"></slot>
    </template>

    <template #staff-settings>
      <slot name="staff-settings"></slot>
    </template>

    <template #kpi-monthly>
      <slot name="team-gestion"></slot>
    </template>

    <template #kpi-history>
      <slot name="kpi-history"></slot>
    </template>

    <template #calendar>
      <slot name="calendar"></slot>
    </template>

    <template #presence-rate>
  <slot name="presence-rate"></slot>
    </template>
    <template #export-button>
  <slot name="export-button"></slot>
  </template>

  </DraggableLayout>
</template>