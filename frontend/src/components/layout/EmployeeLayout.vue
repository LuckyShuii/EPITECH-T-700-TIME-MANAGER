<script setup lang="ts">
</script>

<template>
  <div class="grid grid-cols-4 grid-rows-4 gap-3">
    <!-- Position 1 : Haut gauche (clock) -->
    <div class="col-start-1 row-start-1 row-span-2">
      <slot name="clock"></slot>
    </div>
    
    <!-- Position 2 : team-presence (haut droite - 3 rows) -->
    <div class="row-span-3 col-start-4 row-start-1">
      <slot name="team-presence"></slot>
    </div>
    
    <!-- Position 3 : widget-5 (bas droite) -->
    <div class="col-start-4 row-start-4">
      <slot name="widget-4"></slot>
    </div>
    
    <!-- Position 4 : widget-4 (milieu gauche) -->
    <div class="col-start-1 row-start-3">
      <slot name="widget-6"></slot>
    </div>
    
    <!-- Position 5 : (bas gauche) - vide pour l'instant ou réutilisable -->
    <div class="col-start-1 row-start-4">
      <slot name="extra-widget"></slot>
    </div>
    
    <!-- Position 6 : calendar (haut centre - 2 colonnes, 2 rows) -->
    <div class="col-span-2 row-span-2 col-start-2 row-start-1">
      <slot name="calendar"></slot>
    </div>
    
    <!-- Position 7 : widget-2 (milieu centre - 2 colonnes) -->
    <div class="col-span-2 col-start-2 row-start-3">
      <slot name="widget-2"></slot>
    </div>
    
    <!-- Position 8 : widget-3 (bas centre - 2 colonnes) -->
    <div class="col-span-2 col-start-2 row-start-4">
      <slot name="widget-3"></slot>
    </div>
  </div>
</template>