<script setup lang="ts">
import DraggableLayout from './DraggableLayout.vue'
</script>

<template>
  <DraggableLayout dashboard-name="admin">
    <!-- On passe tous les slots au DraggableLayout -->
    <template #add-employee>
      <slot name="add-employee"></slot>
    </template>

    <template #staff-settings>
      <slot name="staff-settings"></slot>
    </template>

    <template #kpi-monthly>
      <slot name="kpi-monthly"></slot>
    </template>

    <template #kpi-history>
      <slot name="kpi-history"></slot>
    </template>

    <template #calendar>
      <slot name="calendar"></slot>
    </template>

    <template #widget-6>
      <slot name="widget-6"></slot>
    </template>

    <template #remote-absence>
      <slot name="remote-absence"></slot>
    </template>

    <template #manager-report>
      <slot name="manager-report"></slot>
    </template>
  </DraggableLayout>
</template>