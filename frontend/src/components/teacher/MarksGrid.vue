<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
    students: any[];
    marks: any[]; // Flat array of marks
}>();

const emit = defineEmits<{
    (e: 'add-mark', studentId: number): void;
    (e: 'edit-mark', mark: any): void;
}>();

// Transform marks into easy lookup format
const getStudentMarks = (studentId: number) => {
    return props.marks.filter(m => m.studentId === studentId);
};

const getBackgroundColor = (value: number) => {
    if (value < 6) return 'bg-error/20 text-error-content';
    if (value >= 9) return 'bg-success/20 text-success-content';
    return 'bg-base-200';
};
</script>

<template>
  <div class="overflow-x-auto">
      <table class="table table-zebra table-sm w-full">
          <thead>
              <tr>
                  <th class="w-1/4">Student</th>
                  <th>Marks</th>
                  <th class="w-16">Avg</th>
                  <th class="w-16">Actions</th>
              </tr>
          </thead>
          <tbody>
              <tr v-for="student in students" :key="student.id">
                  <td class="font-bold">{{ student.name }}</td>
                  <td>
                      <div class="flex flex-wrap gap-2">
                          <div v-for="mark in getStudentMarks(student.id)" :key="mark.id"
                               class="badge badge-lg cursor-pointer hover:scale-110 transition-transform"
                               :class="getBackgroundColor(mark.value)"
                               @click="emit('edit-mark', mark)"
                               :title="`${mark.type} - ${mark.date}`">
                              {{ mark.value }}
                          </div>
                          <button class="btn btn-ghost btn-xs btn-circle border-dashed border-2 border-base-300" 
                                  @click="emit('add-mark', student.id)">
                              +
                          </button>
                      </div>
                  </td>
                  <td>
                      <!-- Calculate Avg -->
                      {{ (getStudentMarks(student.id).reduce((acc, m) => acc + m.value, 0) / (getStudentMarks(student.id).length || 1)).toFixed(1) }}
                  </td>
                  <td>
                      <button class="btn btn-ghost btn-xs">History</button>
                  </td>
              </tr>
          </tbody>
      </table>
  </div>
</template>
