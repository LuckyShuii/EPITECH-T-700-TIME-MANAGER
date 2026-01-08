<script setup lang="ts">
import DraggableLayout from './DraggableLayout.vue'
</script>

<template>
  <DraggableLayout dashboard-name="manager">
    <template #clock>
      <slot name="clock"></slot>
    </template>

    <template #kpi-stats>
      <slot name="kpi-stats"></slot>
    </template>

    <template #team-view>
      <slot name="team-view"></slot>
    </template>

    <template #calendar>
      <slot name="calendar"></slot>
    </template>

    <template #report-button>
      <slot name="report-button"></slot>
    </template>

    <template #team-presence>
      <slot name="team-presence"></slot>
    </template>

    <template #modal-team>
      <slot name="modal-team"></slot>
    </template>

    <template #kpi-history>
      <slot name="kpi-history"></slot>
    </template>

  </DraggableLayout>
</template>