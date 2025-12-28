<template>
  <div class="h-full flex flex-col bg-base-100 rounded-box shadow p-4">
      <h2 class="text-xl font-bold mb-4">Weekly Schedule</h2>
      
      <div class="overflow-x-auto">
          <table class="table table-fixed w-full border border-base-200">
              <thead>
                  <tr>
                      <th class="w-20 bg-base-200">Hour</th>
                      <th v-for="day in weekDays" :key="day" class="text-center bg-base-200">{{ day }}</th>
                  </tr>
              </thead>
              <tbody>
                  <tr v-for="hour in hours" :key="hour">
                      <th class="text-center bg-base-100 text-sm">{{ hour }}:00</th>
                      <td v-for="day in weekDays" :key="day" class="border border-base-200 p-1 h-24 align-top">
                          <div v-if="getLesson(day, hour)" class="bg-primary/10 p-2 rounded h-full text-xs hover:bg-primary/20 cursor-pointer transition-colors border-l-4 border-primary">
                              <div class="font-bold text-primary">{{ getLesson(day, hour)?.class }}</div>
                              <div>{{ getLesson(day, hour)?.subject }}</div>
                              <div class="opacity-60 mt-1">{{ getLesson(day, hour)?.room }}</div>
                          </div>
                      </td>
                  </tr>
              </tbody>
          </table>
      </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const weekDays = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
const hours = [8, 9, 10, 11, 12, 13, 14, 15];

// Mock Schedule Data - In real app, fetch from store/API
const lessons = ref([
    { day: 'Monday', hour: 8, class: '1A', subject: 'Matematica', room: 'Aula 101' },
    { day: 'Monday', hour: 9, class: '1A', subject: 'Matematica', room: 'Aula 101' },
    { day: 'Tuesday', hour: 10, class: '2B', subject: 'Fisica', room: 'Lab 2' },
    { day: 'Wednesday', hour: 11, class: '1A', subject: 'Matematica', room: 'Aula 101' },
    { day: 'Friday', hour: 9, class: '5C', subject: 'Matematica', room: 'Aula 304' },
]);

const getLesson = (day: string, hour: number) => {
    return lessons.value.find(l => l.day === day && l.hour === hour);
};
</script>
