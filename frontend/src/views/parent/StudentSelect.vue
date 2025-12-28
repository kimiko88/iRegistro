<script setup lang="ts">
import { useParentStore } from '@/stores/parent';
const parentStore = useParentStore();
</script>

<template>
  <div class="join">
      <template v-for="child in parentStore.children" :key="child.id">
          <input class="join-item btn" type="radio" 
                 :name="`child-${child.id}`" 
                 :aria-label="child.name"
                 :checked="parentStore.selectedChildId === child.id"
                 @change="parentStore.selectChild(child.id)" />
      </template>
  </div>
</template>
