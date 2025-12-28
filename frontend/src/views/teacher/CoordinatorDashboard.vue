<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-4">Cruscotto Coordinatore</h1>
    <div class="alert alert-info shadow-sm">
      <span class="i-heroicons-information-circle w-6 h-6"></span>
      <div>
        <h3 class="font-bold">Work in Progress</h3>
        <div class="text-xs">Questa dashboard mostrerà le statistiche aggregate della classe coordinata.</div>
      </div>
    </div>
    
    <div class="stats shadow mt-6 w-full">
      <div class="stat">
        <div class="stat-title">Media Classe</div>
        <div class="stat-value">7.2</div>
        <div class="stat-desc">↗︎ 0.4 (Ultimo mese)</div>
      </div>
      <div class="stat">
        <div class="stat-title">Assenze Totali</div>
        <div class="stat-value text-warning">45</div>
        <div class="stat-desc">Da inizio anno</div>
      </div>
      <div class="stat">
        <div class="stat-title">Note Disciplinari</div>
        <div class="stat-value text-error">3</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Logic for coordinator would go here, fetching specific coordinator stats
</script>
