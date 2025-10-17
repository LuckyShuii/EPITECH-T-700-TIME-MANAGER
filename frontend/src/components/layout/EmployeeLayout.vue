<script setup lang="ts">
</script>
<template>
<div class="grid grid-cols-4 grid-rows-4 gap-3">
  <!-- Position 1 : Haut gauche -->
  <div class="row-span-2">
    <slot name="clock"></slot>
  </div>
  
  <!-- Position 3 : clock (bas gauche) -->
  <div class="row-span-2 col-start-1 row-start-3">
    <slot name="widget-4"></slot>
  </div>
  
  <!-- Position 4 : calendar (haut centre - 2 colonnes) -->
  <div class="col-span-2 row-span-2 col-start-2 row-start-1">
    <slot name="calendar"></slot>
  </div>
  
  <!-- Position 5 : widget-2 (bas centre gauche) -->
  <div class="row-span-2 col-start-2 row-start-3">
    <slot name="widget-2"></slot>
  </div>
  
  <!-- Position 6 : widget-3 (bas centre droite) -->
  <div class="row-span-2 col-start-3 row-start-3">
    <slot name="widget-3"></slot>
  </div>
  
  <!-- Position 7 : widget-4 (haut droite) -->
  <div class="row-span-2 col-start-4 row-start-1">
    <slot name="widget-1"></slot>
  </div>
  
  <!-- Position 8 : widget-5 (bas droite) -->
  <div class="row-span-2 row-start-3">
    <slot name="widget-5"></slot>
  </div>
</div>
</template>