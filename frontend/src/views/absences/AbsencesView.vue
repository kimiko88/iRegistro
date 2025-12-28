<script setup lang="ts">
// Simple generic absence view
defineProps<{
    absences: any[];
    percentage?: number;
}>();
</script>

<template>
   <div class="space-y-4">
       <div class="alert alert-info shadow-lg" v-if="percentage && percentage < 75">
          <div>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current flex-shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            <span>Attendance is below 75% threshold ({{ percentage }}%). Please verify justifications.</span>
          </div>
       </div>

       <div class="card bg-base-100 shadow">
           <div class="card-body">
               <h3 class="card-title">Absences History</h3>
               <table class="table w-full">
                   <thead>
                       <tr>
                           <th>Date</th>
                           <th>Type</th>
                           <th>Note</th>
                           <th>Status</th>
                       </tr>
                   </thead>
                   <tbody>
                       <tr v-for="a in absences" :key="a.id">
                           <td>{{ a.date }}</td>
                           <td>{{ a.type || 'Absence' }}</td>
                           <td>{{ a.note || '-' }}</td>
                           <td>
                               <button v-if="!a.justified" class="btn btn-xs btn-outline btn-warning">Justify</button>
                               <span v-else class="text-success">Justified</span>
                           </td>
                       </tr>
                   </tbody>
               </table>
           </div>
       </div>
   </div>
</template>
