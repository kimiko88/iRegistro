<template>
  <div class="p-6 space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-3xl font-bold">Backup Management</h1>
      <ActionButton
        label="Create Manual Backup"
        variant="primary"
        @click="createBackup"
        :loading="creatingBackup"
        :icon="Save"
      />
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-base-100 p-6 rounded-lg shadow md:col-span-2">
        <h2 class="text-xl font-bold mb-4">Recent Backups</h2>
        <table class="table w-full">
          <thead>
            <tr>
              <th>Date</th>
              <th>Size</th>
              <th>Type</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="backup in backups" :key="backup.id">
              <td>{{ backup.date }}</td>
              <td>{{ backup.size }}</td>
              <td>{{ backup.type }}</td>
              <td><span class="badge badge-success">{{ backup.status }}</span></td>
              <td>
                <button class="btn btn-xs btn-ghost text-info">Restore</button>
                <button class="btn btn-xs btn-ghost text-primary">Download</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="bg-base-100 p-6 rounded-lg shadow">
        <h2 class="text-xl font-bold mb-4">Configuration</h2>
        <div class="form-control">
           <label class="cursor-pointer label">
            <span class="label-text font-bold">Automatic Backups</span> 
            <input type="checkbox" class="toggle toggle-primary" v-model="config.enabled" />
          </label>
        </div>
        
        <div v-if="config.enabled" class="space-y-4 mt-4">
           <div class="form-control">
             <label class="label"><span class="label-text">Frequency</span></label>
             <select class="select select-bordered" v-model="config.frequency">
               <option value="daily">Daily</option>
               <option value="weekly">Weekly</option>
               <option value="monthly">Monthly</option>
             </select>
           </div>
           
           <div class="form-control">
             <label class="label"><span class="label-text">Time (UTC)</span></label>
             <input type="time" class="input input-bordered" v-model="config.time" />
           </div>

           <div class="form-control">
             <label class="label"><span class="label-text">Retention (Days)</span></label>
             <input type="number" class="input input-bordered" v-model="config.retention" />
           </div>
           
           <button class="btn btn-secondary w-full">Save Config</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import ActionButton from '@/components/shared/ActionButton.vue';
import { Save } from 'lucide-vue-next';
import { useNotificationStore } from '@/stores/notification';

const notificationStore = useNotificationStore();
const creatingBackup = ref(false);

const backups = ref([
  { id: 1, date: '2023-10-27 02:00', size: '1.2 GB', type: 'Auto', status: 'Completed' },
  { id: 2, date: '2023-10-26 02:00', size: '1.1 GB', type: 'Auto', status: 'Completed' },
  { id: 3, date: '2023-10-25 15:30', size: '1.1 GB', type: 'Manual', status: 'Completed' },
]);

const config = ref({
  enabled: true,
  frequency: 'daily',
  time: '02:00',
  retention: 30
});

const createBackup = async () => {
  creatingBackup.value = true;
  // Simulate API
  setTimeout(() => {
    creatingBackup.value = false;
    notificationStore.success('Backup started successfully');
    backups.value.unshift({
       id: Date.now(),
       date: new Date().toLocaleString(),
       size: '0 MB',
       type: 'Manual',
       status: 'In Progress'
    });
  }, 1000);
};
</script>
