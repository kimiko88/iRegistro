<script setup lang="ts">
import { ref } from 'vue';

const activeTab = ref('general');
</script>

<template>
  <div class="p-4 space-y-6">
    <h1 class="text-2xl font-bold">School Settings</h1>

    <div role="tablist" class="tabs tabs-boxed">
      <a role="tab" class="tab" :class="{ 'tab-active': activeTab === 'general' }" @click="activeTab = 'general'">General</a>
      <a role="tab" class="tab" :class="{ 'tab-active': activeTab === 'backup' }" @click="activeTab = 'backup'">Backups</a>
      <a role="tab" class="tab" :class="{ 'tab-active': activeTab === 'audit' }" @click="activeTab = 'audit'">Audit Logs</a>
    </div>

    <div v-if="activeTab === 'general'" class="card bg-base-100 shadow">
      <div class="card-body">
         <h2 class="card-title">School Information</h2>
         <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
             <div class="form-control">
                <label class="label">School Name</label>
                <input type="text" class="input input-bordered" value="Liceo Scientifico Galileo Galilei" />
             </div>
             <div class="form-control">
                <label class="label">Address</label>
                <input type="text" class="input input-bordered" value="Via Roma 1, Milano" />
             </div>
         </div>
         <div class="divider"></div>
         <h3 class="font-bold">Branding</h3>
         <div class="flex gap-4 items-center">
            <div class="avatar placeholder">
                <div class="bg-neutral text-neutral-content rounded-full w-16">
                    <span class="text-xl">LG</span>
                </div>
            </div>
            <button class="btn btn-sm">Upload Logo</button>
         </div>
      </div>
    </div>

    <div v-else-if="activeTab === 'backup'" class="card bg-base-100 shadow">
       <div class="card-body">
          <h2 class="card-title">Backup Management</h2>
          <div class="alert alert-info">
             <span>Next scheduled backup: Tonight at 02:00 AM</span>
          </div>
          <button class="btn btn-warning w-fit">Trigger Manual Backup</button>
          
          <h3 class="font-bold mt-4">History</h3>
          <ul class="menu bg-base-200 w-full rounded-box">
             <li><a>2023-12-28_backup.zip (240MB) <span class="badge badge-sm">Success</span></a></li>
             <li><a>2023-12-27_backup.zip (238MB) <span class="badge badge-sm">Success</span></a></li>
          </ul>
       </div>
    </div>
    
    <div v-else-if="activeTab === 'audit'" class="card bg-base-100 shadow">
       <div class="card-body">
           <h2 class="card-title">System Audit Logs</h2>
           <!-- Placeholder for AuditLogs component/table -->
           <div class="overflow-x-auto">
               <table class="table table-xs">
                   <thead><tr><th>Time</th><th>User</th><th>Action</th></tr></thead>
                   <tbody>
                       <tr><td>10:00 AM</td><td>Admin</td><td>Updated Settings</td></tr>
                       <tr><td>09:45 AM</td><td>SuperAdmin</td><td>Viewed School</td></tr>
                   </tbody>
               </table>
           </div>
       </div>
    </div>

  </div>
</template>
