<script setup lang="ts">
import DraggableLayout from './DraggableLayout.vue'
</script>

<template>
  <DraggableLayout dashboard-name="employee">
    <!-- On passe tous les slots au DraggableLayout -->
    <template #clock>
      <slot name="clock"></slot>
    </template>

    <template #team-presence>
      <slot name="team-presence"></slot>
    </template>

    <template #widget-4>
      <slot name="widget-4"></slot>
    </template>

    <template #widget-6>
      <slot name="widget-6"></slot>
    </template>

    <template #extra-widget>
      <slot name="extra-widget"></slot>
    </template>

    <template #calendar>
      <slot name="calendar"></slot>
    </template>

    <template #widget-2>
      <slot name="widget-2"></slot>
    </template>

    <template #widget-3>
      <slot name="widget-3"></slot>
    </template>
  </DraggableLayout>
</template>