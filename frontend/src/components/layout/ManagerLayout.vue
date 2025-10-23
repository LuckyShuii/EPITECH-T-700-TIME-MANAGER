<script setup lang="ts">
</script>

<template>
  <div class="grid grid-cols-4 grid-rows-4 gap-3">
    <!-- Position 1 : Pointeuse (col 1, rows 1-2) -->
    <div className="row-span-2">
      <slot name="clock"></slot>
    </div>

    <!-- Position 2 : KPI Stats (col 1, rows 3-4) -->
    <div className="col-start-1 row-start-3">
      <slot name="kpi-stats"></slot>
    </div>

    <!-- Position 3 : Calendrier (col 2-3, rows 1-2) -->
   <div className="col-span-2 row-span-2 col-start-2 row-start-1">
      <slot name="calendar"></slot>
    </div>

    <div className="row-start-4">
      <slot name="team-view"> test </slot>
    </div>

    <!-- Position 4 : Présence équipe (col 2-3, rows 3-4) -->
   <div className="row-span-3 col-start-4 row-start-1">
      <slot name="team-presence"></slot>
    </div>

    <!-- Position 5 : KPI Carousel (col 4, rows 1-3) -->
    <div className="col-start-4 row-start-4">
      <slot name="kpi-carousel"></slot>
    </div>

    <!-- Position 6 : Bouton rapport (col 4, row 4) -->
    <div className="col-span-2 row-span-2 col-start-2 row-start-3">
      <slot name="report-button"></slot>
    </div>
  </div>
</template>
