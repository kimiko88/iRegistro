<template>
  <div class="p-6 space-y-6">
    <h1 class="text-3xl font-bold">Audit Logs</h1>
    
    <div class="bg-base-100 p-4 rounded-lg shadow">
      <!-- Search & Filters -->
      <div class="flex gap-4 mb-4 items-end">
        <div class="form-control w-full max-w-xs">
          <label class="label"><span class="label-text">Action Type</span></label>
          <select class="select select-bordered" v-model="filters.action">
            <option value="">All Actions</option>
            <option value="LOGIN">Login</option>
            <option value="CREATE_USER">Create User</option>
            <option value="DELETE_USER">Delete User</option>
            <option value="UPDATE_SCHOOL">Update School</option>
          </select>
        </div>
         <div class="form-control w-full max-w-xs">
          <label class="label"><span class="label-text">User</span></label>
          <input type="text" class="input input-bordered" placeholder="Search user..." v-model="filters.user" />
        </div>
      </div>

      <DataTable
        :columns="columns"
        :data="logs"
        :loading="loading"
        @search="onSearch"
      >
        <template #cell-status="{ value }">
           <span class="badge" :class="value === 'SUCCESS' ? 'badge-success' : 'badge-error'">{{ value }}</span>
        </template>
        <template #item-actions="{ item }">
          <button class="btn btn-xs btn-ghost">Details</button>
        </template>
      </DataTable>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';
import DataTable from '@/components/shared/DataTable.vue';
import { useAdminStore } from '@/stores/admin';

const adminStore = useAdminStore();
// Reuse logs logic from store if available or local state
// store has loadAuditLogs

const filters = reactive({ action: '', user: '' });
const logs = ref([
  { id: 1, timestamp: '2023-10-27 10:00:00', user: 'admin@school.com', action: 'LOGIN', resource: 'Auth', status: 'SUCCESS' },
  { id: 2, timestamp: '2023-10-27 10:05:00', user: 'admin@school.com', action: 'CREATE_USER', resource: 'User: 123', status: 'SUCCESS' },
  { id: 3, timestamp: '2023-10-27 11:00:00', user: 'teacher@school.com', action: 'UPDATE_GRADE', resource: 'Grade: 456', status: 'SUCCESS' },
  { id: 4, timestamp: '2023-10-27 11:30:00', user: 'hacker@ip', action: 'LOGIN', resource: 'Auth', status: 'FAILURE' },
]);
const loading = ref(false);

const columns = [
  { key: 'timestamp', label: 'Timestamp', sortable: true },
  { key: 'user', label: 'User', sortable: true },
  { key: 'action', label: 'Action', sortable: true },
  { key: 'resource', label: 'Resource' },
  { key: 'status', label: 'Status' }
];

const onSearch = (q: string) => {
  // Client side search or api call
};

onMounted(() => {
  // adminStore.loadAuditLogs();
});
</script>
