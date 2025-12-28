<template>
  <div class="calendar-container">
    <!-- Simplified Calendar for brevity, assume a library or custom grid -->
    <div class="grid grid-cols-7 gap-1 text-center">
      <div v-for="day in ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']" :key="day" class="font-bold text-gray-500">
        {{ day }}
      </div>
      <div 
        v-for="date in calendarDays" 
        :key="date.toISOString()"
        class="aspect-square border rounded p-1 relative"
        :class="getDayClass(date)"
      >
        <span class="text-sm">{{ date.getDate() }}</span>
        <div v-if="hasAbsence(date)" class="absolute bottom-1 right-1 w-2 h-2 rounded-full bg-red-500"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  absences: any[]; // { date: string, type: 'absence' | 'late', justified: boolean }
  month: number; // 0-11
  year: number;
}>();

const calendarDays = computed(() => {
  const days = [];
  const start = new Date(props.year, props.month, 1);
  const end = new Date(props.year, props.month + 1, 0);
  
  // Fill text days before start
  for (let i = 0; i < start.getDay(); i++) {
    days.push(new Date(props.year, props.month, -i)); // simple placeholder logic
  }
  // This logic is simplified. Better to use a full calendar logic loop.
  // Implementing a basic 30-day view for now.
  for(let i=1; i<=end.getDate(); i++) {
    days.push(new Date(props.year, props.month, i));
  }
  return days;
});

const hasAbsence = (date: Date) => {
  return props.absences.some(a => new Date(a.date).toDateString() === date.toDateString());
};

const getDayClass = (date: Date) => {
  if (hasAbsence(date)) return 'bg-red-50 dark:bg-red-900/10 border-red-200';
  return 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700';
};
</script>
