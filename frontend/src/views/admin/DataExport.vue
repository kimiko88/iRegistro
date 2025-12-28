<template>
  <div class="p-6 space-y-6">
    <h1 class="text-3xl font-bold">Data Export</h1>
    
    <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
       <!-- New Export Request -->
       <div class="bg-base-100 p-6 rounded-lg shadow">
         <h2 class="text-xl font-bold mb-4">New Export</h2>
         
         <div class="form-control mb-4">
           <label class="label"><span class="label-text">Data Type</span></label>
           <select class="select select-bordered" v-model="exportType">
             <option value="users">Users</option>
             <option value="classes">Classes</option>
             <option value="grades">Grades / Marks</option>
             <option value="attendance">Attendance / Absences</option>
             <option value="audit_logs">Audit Logs</option>
           </select>
         </div>
         
         <div class="form-control mb-4">
           <label class="label"><span class="label-text">Format</span></label>
           <div class="flex gap-4">
             <label class="label cursor-pointer justify-start gap-2">
               <input type="radio" name="format" class="radio radio-primary" value="csv" v-model="format" />
               <span class="label-text">CSV</span>
             </label>
             <label class="label cursor-pointer justify-start gap-2">
               <input type="radio" name="format" class="radio radio-primary" value="json" v-model="format" />
               <span class="label-text">JSON</span>
             </label>
             <label class="label cursor-pointer justify-start gap-2">
               <input type="radio" name="format" class="radio radio-primary" value="pdf" v-model="format" />
               <span class="label-text">PDF (Report)</span>
             </label>
           </div>
         </div>
         
         <!-- Optional: Filters based on type -->
         <div v-if="exportType === 'grades'" class="form-control mb-4">
            <label class="label"><span class="label-text">Class Filter (Optional)</span></label>
            <input type="text" class="input input-bordered" placeholder="e.g. 1A" />
         </div>

         <div class="mt-6">
            <ActionButton 
              label="Generate Export" 
              variant="primary" 
              class="w-full"
              :loading="loading"
              @click="handleExport"
              :icon="Download"
            />
         </div>
       </div>
       
       <!-- History -->
       <div class="bg-base-100 p-6 rounded-lg shadow">
         <h2 class="text-xl font-bold mb-4">Export History</h2>
         <div class="overflow-y-auto max-h-[400px]">
           <table class="table w-full">
             <thead>
               <tr>
                 <th>Date</th>
                 <th>Type</th>
                 <th>Status</th>
                 <th>Action</th>
               </tr>
             </thead>
             <tbody>
               <tr v-for="item in history" :key="item.id">
                 <td>{{ item.date }}</td>
                 <td>{{ item.type }} ({{ item.format }})</td>
                 <td><span class="badge badge-success">Ready</span></td>
                 <td><button class="btn btn-xs btn-link">Download</button></td>
               </tr>
             </tbody>
           </table>
         </div>
       </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import ActionButton from '@/components/shared/ActionButton.vue';
import { useDataExport } from '@/composables/useDataExport';
import { Download } from 'lucide-vue-next';

const { exportData, loading } = useDataExport();

const exportType = ref('users');
const format = ref<any>('csv');

const history = ref([
  { id: 1, date: '2023-10-27', type: 'Users', format: 'CSv' },
  { id: 2, date: '2023-10-25', type: 'Grades', format: 'PDF' },
]);

const handleExport = async () => {
   await exportData(exportType.value, format.value);
   // Add to history mock
   history.value.unshift({
     id: Date.now(),
     date: new Date().toLocaleDateString(),
     type: exportType.value,
     format: format.value.toUpperCase()
   });
};
</script>
