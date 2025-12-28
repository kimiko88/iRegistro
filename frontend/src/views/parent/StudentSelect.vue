<template>
  <div v-if="children.length > 1" class="flex items-center gap-4 overflow-x-auto w-full">
    <span class="mr-2 font-semibold text-gray-700 dark:text-gray-200 whitespace-nowrap">Select Child:</span>
    <div class="flex gap-3">
      <div 
        v-for="child in children" 
        :key="child.id"
        @click="$emit('select', child.id)"
        class="cursor-pointer flex items-center gap-2 px-4 py-2 rounded-full border transition-all"
        :class="selectedId === child.id ? 'bg-indigo-100 border-indigo-500 text-indigo-700 font-medium' : 'bg-white dark:bg-gray-700 border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-600 dark:text-gray-300'"
      >
        <div class="w-8 h-8 rounded-full bg-gray-200 dark:bg-gray-600 overflow-hidden">
             <img v-if="child.photoUrl" :src="child.photoUrl" class="w-full h-full object-cover">
        </div>
        <span class="whitespace-nowrap">{{ child.firstName }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  children: any[];
  selectedId: number | null;
}>();

defineEmits(['select']);
</script>
