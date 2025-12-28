<template>
  <div class="p-4 space-y-6">
    <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100">Marks & Grades</h1>

    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4">
      <div class="flex flex-col sm:flex-row gap-4 mb-4">
        <!-- Filters -->
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Subject</label>
          <select v-model="filterSubject" class="mt-1 block w-full rounded-md border-gray-300 dark:bg-gray-700 dark:border-gray-600">
            <option value="">All Subjects</option>
            <option v-for="s in uniqueSubjects" :key="s" :value="s">{{ s }}</option>
          </select>
        </div>
        <div class="flex-1">
           <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Period</label>
           <select v-model="filterPeriod" class="mt-1 block w-full rounded-md border-gray-300 dark:bg-gray-700 dark:border-gray-600">
            <option value="">All Year</option>
            <option value="1">Trimester 1</option>
            <option value="2">Trimester 2</option>
            <option value="3">Trimester 3</option>
           </select>
        </div>
      </div>
      
      <!-- Chart -->
      <div class="mb-8">
        <MarksChart :marks="filteredMarks" />
      </div>

       <!-- Table -->
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead class="bg-gray-50 dark:bg-gray-700">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Date</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Subject</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Grade</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Teacher</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Notes</th>
            </tr>
          </thead>
          <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
            <tr v-for="mark in filteredMarks" :key="mark.id">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">{{ new Date(mark.date).toLocaleDateString() }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">{{ mark.subject }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-bold" :class="getGradeColor(mark.grade)">
                {{ mark.grade }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">{{ mark.teacher }}</td>
              <td class="px-6 py-4 text-sm text-gray-500 dark:text-gray-400">
                {{ mark.notes }}
                <span v-if="mark.appeal" class="ml-2 text-yellow-600 bg-yellow-100 px-2 py-0.5 rounded text-xs">Appealed</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useParentStore } from '@/stores/parent';
import { useStudentStore } from '@/stores/student';
import { useAuthStore } from '@/stores/auth';
import MarksChart from '@/components/MarksChart.vue';

const auth = useAuthStore();
const parentStore = useParentStore();
const studentStore = useStudentStore();

const filterSubject = ref("");
const filterPeriod = ref("");

const isParent = computed(() => auth.user?.role === 'parent');

const marks = computed(() => {
  if (isParent.value) {
    return parentStore.marks;
  } else {
    return studentStore.marks;
  }
});

const uniqueSubjects = computed(() => {
  const subjects = new Set(marks.value.map(m => m.subject));
  return Array.from(subjects);
});

const filteredMarks = computed(() => {
  return marks.value.filter(m => {
    if (filterSubject.value && m.subject !== filterSubject.value) return false;
    // Period logic optional - assume date based or m.period check
    if (filterPeriod.value && m.period != filterPeriod.value) return false; 
    return true;
  });
});

const getGradeColor = (g: number) => {
  if (g < 6) return 'text-red-600';
  if (g >= 8) return 'text-green-600';
  return 'text-gray-900 dark:text-gray-100';
};

onMounted(() => {
  if (isParent.value) {
    parentStore.fetchMarks();
  } else {
    studentStore.fetchMarks();
  }
});
</script>
