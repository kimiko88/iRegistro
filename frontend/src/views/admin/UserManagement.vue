<template>
  <div class="p-6 space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-3xl font-bold">User Management</h1>
      <div class="flex gap-2">
        <ActionButton 
          label="Import CSV" 
          variant="secondary" 
          @click="showImportModal = true"
          :icon="Upload"
        />
        <ActionButton 
          label="Create User" 
          variant="primary" 
          @click="openCreateModal"
          :icon="Plus"
        />
      </div>
    </div>

    <!-- Stats Row (Optional, user requirement B didn't explicitly ask but good for UX) -->
    
    <!-- Filters -->
    <div class="bg-base-100 p-4 rounded-lg shadow-sm flex flex-wrap gap-4 items-end">
       <div class="form-control w-full max-w-xs">
          <label class="label"><span class="label-text">Role</span></label>
          <select class="select select-bordered" v-model="filters.role" @change="fetchUsers">
            <option value="">All Roles</option>
            <option value="Teacher">Teacher</option>
            <option value="Student">Student</option>
            <option value="Parent">Parent</option>
            <option value="Secretary">Secretary</option>
            <option value="Principal">Principal</option>
          </select>
       </div>
       <div class="form-control w-full max-w-xs">
          <label class="label"><span class="label-text">Status</span></label>
          <select class="select select-bordered" v-model="filters.status" @change="fetchUsers">
            <option value="">All Statuses</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
       </div>
    </div>

    <div class="bg-base-100 p-4 rounded-lg shadow">
      <DataTable
        :columns="columns"
        :data="users"
        :loading="loading"
        :totalItems="totalUsers"
        :itemsPerPage="itemsPerPage"
        @search="onSearch"
        @update:page="onPageChange"
      >
        <template #item-actions="{ item }">
          <div class="dropdown dropdown-left">
            <label tabindex="0" class="btn btn-ghost btn-xs">
               <MoreHorizontal class="w-4 h-4" />
            </label>
            <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
              <li><a @click="editUser(item)">Edit</a></li>
              <li><a @click="resetPassword(item)">Reset Password</a></li>
              <li><a @click="viewHistory(item)">View History</a></li>
              <li v-if="item.status === 'active'"><a class="text-error" @click="toggleStatus(item)">Deactivate</a></li>
              <li v-else><a class="text-success" @click="toggleStatus(item)">Activate</a></li>
            </ul>
          </div>
        </template>
        
        <template #cell-status="{ value }">
           <div class="badge" :class="value === 'active' ? 'badge-success' : 'badge-ghost'">
             {{ value }}
           </div>
        </template>
      </DataTable>
    </div>

    <!-- Modals -->
    <UserImportModal 
      :isOpen="showImportModal"
      @close="showImportModal = false"
      @import-success="onImportSuccess"
    />

    <UserEditModal
      :isOpen="showEditModal"
      :user="selectedUser"
      :loading="modifying"
      @close="showEditModal = false"
      @submit="handleUserSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useAdminStore } from '@/stores/admin';
import { storeToRefs } from 'pinia';
import ActionButton from '@/components/shared/ActionButton.vue';
import DataTable from '@/components/shared/DataTable.vue';
import UserImportModal from '@/components/admin/UserImportModal.vue';
import UserEditModal from '@/components/admin/UserEditModal.vue';
import { Plus, Upload, MoreHorizontal } from 'lucide-vue-next';
import adminService from '@/services/admin';
import { useNotificationStore } from '@/stores/notification';

const adminStore = useAdminStore();
const notificationStore = useNotificationStore();
const { users, totalUsers, loading } = storeToRefs(adminStore);

const itemsPerPage = 10;
const filters = reactive({
  role: '',
  status: '',
  query: '',
  page: 1,
  limit: itemsPerPage
});

const showImportModal = ref(false);
const showEditModal = ref(false);
const selectedUser = ref(null);
const modifying = ref(false);

const columns = [
  { key: 'firstName', label: 'First Name', sortable: true },
  { key: 'lastName', label: 'Last Name', sortable: true },
  { key: 'email', label: 'Email', sortable: true },
  { key: 'role', label: 'Role', sortable: true },
  { key: 'status', label: 'Status', sortable: true },
  { key: 'createdAt', label: 'Created At', sortable: true },
];

const fetchUsers = () => {
   adminStore.fetchUsers(filters);
};

const onSearch = (query: string) => {
  filters.query = query;
  filters.page = 1;
  fetchUsers();
};

const onPageChange = (page: number) => {
  filters.page = page;
  fetchUsers();
};

const openCreateModal = () => {
  selectedUser.value = null;
  showEditModal.value = true;
};

const editUser = (user: any) => {
  selectedUser.value = user;
  showEditModal.value = true;
};

const handleUserSubmit = async (formData: any) => {
  modifying.value = true;
  try {
    if (selectedUser.value) {
       await adminService.updateUser((selectedUser.value as any).id, formData);
       notificationStore.success('User updated successfully');
    } else {
       await adminService.createUser(formData);
       notificationStore.success('User created successfully');
    }
    showEditModal.value = false;
    fetchUsers();
  } catch (err) {
    notificationStore.error('Failed to save user');
  } finally {
    modifying.value = false;
  }
};

const toggleStatus = async (user: any) => {
  // Implementation for deactivate/activate
  notificationStore.info(`Simulated: ${user.status === 'active' ? 'Deactivated' : 'Activated'} user`);
};

const resetPassword = async (user: any) => {
  // Implementation
   notificationStore.info(`Simulated: Password reset email sent to ${user.email}`);
};

const viewHistory = (user: any) => {
  // View audit logs for user
  notificationStore.info('Not implemented: View History');
};

const onImportSuccess = () => {
  fetchUsers();
};

onMounted(() => {
  fetchUsers();
});
</script>
