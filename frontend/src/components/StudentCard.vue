<template>
  <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 flex flex-col md:flex-row items-center gap-6">
    <div class="w-24 h-24 rounded-full overflow-hidden bg-gray-200 flex-shrink-0">
      <img v-if="student.photoUrl" :src="student.photoUrl" alt="Student Photo" class="w-full h-full object-cover">
      <div v-else class="w-full h-full flex items-center justify-center text-gray-400">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
        </svg>
      </div>
    </div>
    
    <div class="flex-1 text-center md:text-left">
      <h2 class="text-2xl font-bold text-gray-800 dark:text-white">{{ student.firstName }} {{ student.lastName }}</h2>
      <p class="text-gray-500 dark:text-gray-400 mb-2">{{ student.className }}</p>
      
      <div class="grid grid-cols-2 gap-4 mt-4">
        <div class="bg-blue-50 dark:bg-blue-900/20 p-3 rounded-lg">
          <span class="block text-sm text-gray-500 dark:text-gray-400">Average</span>
          <span class="text-xl font-bold text-blue-600 dark:text-blue-400">{{ student.averageGrade?.toFixed(1) || '-' }}</span>
        </div>
        <div class="bg-red-50 dark:bg-red-900/20 p-3 rounded-lg">
          <span class="block text-sm text-gray-500 dark:text-gray-400">Absences</span>
          <span class="text-xl font-bold text-red-600 dark:text-red-400">{{ student.totalAbsences || 0 }}</span>
        </div>
      </div>
    </div>

    <div v-if="selectable" class="mt-4 md:mt-0">
      <button @click="$emit('select')" class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors">
        View Details
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
// import { defineProps } from 'vue'; // defineProps is a macro

defineProps<{
  student: {
    firstName: string;
    lastName: string;
    photoUrl?: string;
    className: string;
    averageGrade?: number;
    totalAbsences?: number;
  };
  selectable?: boolean;
}>();

defineEmits(['select']);
</script>
