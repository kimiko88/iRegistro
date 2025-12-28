<template>
  <div class="p-4 space-y-6">
    <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100">Absences & Attendance</h1>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Calendar View -->
      <div class="lg:col-span-1 bg-white dark:bg-gray-800 rounded-lg shadow p-4">
         <h3 class="text-lg font-bold mb-4">Calendar</h3>
         <AbsenceCalendar :absences="absences" :month="currentMonth" :year="currentYear" />
         <!-- Simple month nav -->
      </div>

      <!-- List & Justification -->
      <div class="lg:col-span-2 bg-white dark:bg-gray-800 rounded-lg shadow p-4">
        <h3 class="text-lg font-bold mb-4">Details</h3>
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-700">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
                 <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                 <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                 <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Action</th>
              </tr>
            </thead>
             <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="absence in absences" :key="absence.id">
                <td class="px-6 py-4 text-sm">{{ new Date(absence.date).toLocaleDateString() }}</td>
                <td class="px-6 py-4 text-sm">{{ absence.type }}</td>
                <td class="px-6 py-4 text-sm">
                  <span v-if="absence.justified" class="text-green-600 font-semibold">Justified</span>
                  <span v-else class="text-red-600 font-semibold">Unjustified</span>
                </td>
                <td class="px-6 py-4 text-sm">
                  <button 
                    v-if="!absence.justified && isParent" 
                    @click="justify(absence)"
                    class="text-indigo-600 hover:text-indigo-900 underline"
                  >
                    Justify
                  </button>
                   <span v-if="!isParent && !absence.justified" class="text-gray-400 text-xs">Parent action needed</span>
                </td>
              </tr>
             </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useParentStore } from '@/stores/parent';
import { useStudentStore } from '@/stores/student';
import { useAuthStore } from '@/stores/auth';
import AbsenceCalendar from '@/components/AbsenceCalendar.vue';

const auth = useAuthStore();
const parentStore = useParentStore();
const studentStore = useStudentStore();

const currentMonth = ref(new Date().getMonth());
const currentYear = ref(new Date().getFullYear());

const isParent = computed(() => auth.user?.role === 'parent');

const absences = computed(() => {
  if (isParent.value) {
    return parentStore.absences;
  } else {
    return studentStore.absences;
  }
});

const justify = async (absence: any) => {
  const reason = prompt("Enter justification reason:");
  if (reason) {
     await parentStore.justifyAbsence(absence.id, reason);
  }
};

onMounted(() => {
   if (isParent.value) {
    parentStore.fetchAbsences();
  } else {
    studentStore.fetchAbsences();
  }
});
</script>
