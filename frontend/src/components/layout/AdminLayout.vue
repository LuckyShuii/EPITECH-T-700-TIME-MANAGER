<script setup lang="ts">
</script>

<template>
  <div class="grid grid-cols-4 grid-rows-4 gap-3">
    <!-- Colonne 1 : Boutons/Actions -->
    <div class="col-start-1 row-start-1">
      <slot name="add-employee"></slot>
    </div>
    
    <div class="col-start-1 row-start-2">
      <slot name="staff-settings"></slot>
    </div>
    
    <div class="col-start-1 row-start-3">
      <slot name="kpi-monthly"></slot>
    </div>
    
    <div class="col-start-1 row-start-4">
      <slot name="kpi-history"></slot>
    </div>
    
    <!-- Centre : Calendrier (2 colonnes, 2 lignes) -->
    <div class="col-span-2 row-span-2 col-start-2 row-start-1">
      <slot name="calendar"></slot>
    </div>
    
    <!-- Centre bas : Widget 6 (à définir) -->
    <div class="col-span-2 row-span-2 col-start-2 row-start-3">
      <slot name="widget-6"></slot>
    </div>
    
    <!-- Droite haut : TT/Absent -->
    <div class="row-span-2 col-start-4 row-start-1">
      <slot name="remote-absence"></slot>
    </div>
    
    <!-- Droite bas : Rapport manager -->
    <div class="row-span-2 col-start-4 row-start-3">
      <slot name="manager-report"></slot>
    </div>
  </div>
</template>