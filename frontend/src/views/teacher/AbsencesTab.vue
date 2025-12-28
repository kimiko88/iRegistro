<script setup lang="ts">
import { useTeacherStore } from '@/stores/teacher';
import { ref } from 'vue';

const teacherStore = useTeacherStore();
const period = ref('2023-12-25'); // Week start

const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri'];
</script>

<template>
  <div class="space-y-4">
     <div class="flex justify-between items-center bg-base-100 p-2 rounded">
         <h3 class="font-bold">Attendance Register</h3>
         <div class="join">
             <button class="btn btn-sm join-item">«</button>
             <button class="btn btn-sm join-item">Current Week</button>
             <button class="btn btn-sm join-item">»</button>
         </div>
         <button class="btn btn-sm btn-primary">Mark All Present</button>
     </div>
     
     <div class="overflow-x-auto">
         <table class="table table-zebra table-sm">
             <thead>
                 <tr>
                     <th>Student</th>
                     <th v-for="d in days" :key="d" class="text-center">{{ d }}</th>
                 </tr>
             </thead>
             <tbody>
                 <tr v-for="student in teacherStore.students" :key="student.id">
                     <td class="font-medium">{{ student.name }}</td>
                     <td v-for="d in days" :key="d" class="text-center">
                         <input type="checkbox" class="checkbox checkbox-sm checkbox-primary" />
                         <!-- Helper for detailed status (Late, etc) would be a dropdown/popover -->
                     </td>
                 </tr>
             </tbody>
         </table>
     </div>
  </div>
</template>
