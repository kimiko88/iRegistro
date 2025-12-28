<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useAdminStore } from '@/stores/admin';
import DataTable from '@/components/shared/DataTable.vue';
import Modal from '@/components/shared/Modal.vue';
import { useUIStore } from '@/stores/ui';

const adminStore = useAdminStore();
const uiStore = useUIStore();
const isCreateModalOpen = ref(false);

onMounted(() => {
  adminStore.fetchUsers();
});

const columns = [
    { key: 'name', label: 'Full Name' },
    { key: 'email', label: 'Email' },
    { key: 'role', label: 'Role' },
    { key: 'status', label: 'Status' }
];

const openCreateModal = () => isCreateModalOpen.value = true;
const closeCreateModal = () => isCreateModalOpen.value = false;

const createUser = () => {
    // Implement create logic
    uiStore.addNotification({ type: 'success', message: 'User created' });
    closeCreateModal();
};
</script>

<template>
  <div class="p-4 space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold">User Management</h1>
      <div class="space-x-2">
         <button class="btn btn-outline btn-sm">Import CSV</button>
         <button class="btn btn-primary btn-sm" @click="openCreateModal">+ Create User</button>
      </div>
    </div>

    <!-- Filters placeholder -->
    <div class="flex gap-2 mb-4 bg-base-100 p-2 rounded shadow-sm">
      <input type="text" placeholder="Search..." class="input input-bordered input-sm w-full max-w-xs" />
      <select class="select select-bordered select-sm">
         <option disabled selected>Filter by Role</option>
         <option>Student</option>
         <option>Teacher</option>
      </select>
    </div>

    <div class="card bg-base-100 shadow">
       <div class="card-body p-0">
          <DataTable :columns="columns" :items="adminStore.users" :loading="uiStore.isLoading" actions>
             <template #cell-status="{ item }">
                <div class="badge" :class="item.status === 'Active' ? 'badge-success' : 'badge-error'">
                    {{ item.status }}
                </div>
             </template>
             <template #actions>
                <button class="btn btn-ghost btn-xs">Edit</button>
                <button class="btn btn-ghost btn-xs text-error">Delete</button>
             </template>
          </DataTable>
       </div>
    </div>

    <Modal title="Create New User" :isOpen="isCreateModalOpen" @close="closeCreateModal">
         <div class="form-control w-full">
            <label class="label"><span class="label-text">Full Name</span></label>
            <input type="text" class="input input-bordered w-full" placeholder="John Doe" />
         </div>
         <div class="form-control w-full mt-2">
            <label class="label"><span class="label-text">Email</span></label>
            <input type="email" class="input input-bordered w-full" placeholder="john@example.com" />
         </div>
         <div class="form-control w-full mt-2">
            <label class="label"><span class="label-text">Role</span></label>
            <select class="select select-bordered">
               <option>Student</option>
               <option>Teacher</option>
               <option>Parent</option>
            </select>
         </div>
         
         <template #actions>
             <button class="btn" @click="closeCreateModal">Cancel</button>
             <button class="btn btn-primary" @click="createUser">Create</button>
         </template>
    </Modal>
  </div>
</template>
