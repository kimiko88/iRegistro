<script setup lang="ts">
import { onMounted } from 'vue';
import { useAdminStore } from '@/stores/admin';
import DataTable from '@/components/shared/DataTable.vue';
import { useUIStore } from '@/stores/ui';

const adminStore = useAdminStore();
const uiStore = useUIStore();

onMounted(() => {
  adminStore.fetchSchools();
});

const schoolColumns = [
    { key: 'name', label: 'School Name' },
    { key: 'type', label: 'Type' },
    { key: 'users', label: 'Users' },
    { key: 'storage', label: 'Storage' }
];
</script>

<template>
  <div class="p-4 space-y-6">
    <h1 class="text-3xl font-bold">SuperAdmin Dashboard</h1>

    <!-- KPI Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="stats shadow bg-base-100">
        <div class="stat">
          <div class="stat-title">Total Schools</div>
          <div class="stat-value text-primary">{{ adminStore.schools.length }}</div>
          <div class="stat-desc">across 5 regions</div>
        </div>
      </div>
      <div class="stats shadow bg-base-100">
        <div class="stat">
          <div class="stat-title">Total Users</div>
          <div class="stat-value text-secondary">1,200</div>
          <div class="stat-desc">↗︎ 40 (24h)</div>
        </div>
      </div>
      <div class="stats shadow bg-base-100">
        <div class="stat">
          <div class="stat-title">System Storage</div>
          <div class="stat-value">25%</div>
          <div class="stat-desc">150GB / 600GB</div>
        </div>
      </div>
    </div>

    <!-- Map Placeholder -->
    <div class="card bg-base-100 shadow">
      <div class="card-body">
         <h2 class="card-title">Schools Distribution</h2>
         <div class="h-64 bg-base-200 rounded flex items-center justify-center">
            <span class="text-base-content/50">Interactive Italy Map Placeholder</span>
         </div>
      </div>
    </div>

    <!-- Schools Table -->
    <div class="card bg-base-100 shadow">
      <div class="card-body">
        <div class="flex justify-between items-center mb-4">
           <h2 class="card-title">Managed Schools</h2>
           <button class="btn btn-primary btn-sm">+ Add School</button>
        </div>
        <DataTable :columns="schoolColumns" :items="adminStore.schools" :loading="uiStore.isLoading">
             <template #cell-name="{ item }">
                 <span class="font-bold">{{ item.name }}</span>
             </template>
        </DataTable>
      </div>
    </div>
  </div>
</template>
