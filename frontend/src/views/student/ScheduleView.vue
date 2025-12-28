<template>
  <div class="p-4 space-y-6">
    <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100">Weekly Schedule</h1>

    <div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
      <div class="grid grid-cols-6 border-b dark:border-gray-700 bg-gray-50 dark:bg-gray-700">
        <div class="p-4 text-center font-bold text-gray-500 dark:text-gray-300">Time</div>
        <div v-for="day in days" :key="day" class="p-4 text-center font-bold text-gray-500 dark:text-gray-300">
          {{ day }}
        </div>
      </div>

      <div class="divide-y divide-gray-200 dark:divide-gray-700">
        <div v-for="slot in timeSlots" :key="slot" class="grid grid-cols-6">
          <div class="p-4 text-center text-sm text-gray-500 border-r dark:border-gray-700 flex items-center justify-center">
            {{ slot }}
          </div>
          <div v-for="day in days" :key="day" class="p-2 border-r dark:border-gray-700 min-h-[80px]">
             <div v-if="getLesson(day, slot)" class="bg-indigo-50 dark:bg-indigo-900/30 p-2 rounded border border-indigo-100 dark:border-indigo-800 h-full">
               <p class="font-bold text-indigo-700 dark:text-indigo-300 text-sm">{{ getLesson(day, slot)?.subject }}</p>
               <p class="text-xs text-gray-600 dark:text-gray-400">{{ getLesson(day, slot)?.teacher }}</p>
               <p class="text-xs text-gray-500 dark:text-gray-500">{{ getLesson(day, slot)?.room }}</p>
             </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import studentApi from '@/services/student';

const days = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday'];
const timeSlots = ['08:00', '09:00', '10:00', '11:00', '12:00', '13:00'];

// Mock data structure: { day: 'Monday', time: '08:00', subject: 'Math', teacher: 'Rossi', room: '1A' }
const schedule = ref<any[]>([]);

const getLesson = (day: string, time: string) => {
  return schedule.value.find(s => s.day === day && s.startTime.startsWith(time));
};

onMounted(async () => {
  // In a real app, fetch from API
  // schedule.value = (await studentApi.getSchedule()).data;
  
  // Mock data for demo
  schedule.value = [
    { day: 'Monday', startTime: '08:00', subject: 'Math', teacher: 'Rossi', room: '1A' },
    { day: 'Monday', startTime: '09:00', subject: 'History', teacher: 'Bianchi', room: '1A' },
    { day: 'Tuesday', startTime: '10:00', subject: 'Physics', teacher: 'Verdi', room: 'Lab 2' },
  ];
});
</script>
