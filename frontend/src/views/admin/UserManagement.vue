<template>
  <div class="p-6 space-y-6">
    <div class="flex justify-between items-center">
      <div class="flex items-center gap-3">
      <button v-if="currentSchoolId && authStore.user?.role !== 'Secretary'" @click="goBack" class="btn btn-circle btn-ghost btn-sm">
            <ArrowLeft class="w-5 h-5"/>
          </button>
          <h1 class="text-3xl font-bold">
            {{ currentSchoolId ? (route.query.schoolName || 'School Users') : 'User Management' }}
          </h1>
      </div>
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
      :showSchoolSelect="!currentSchoolId"
      :preselectedSchoolId="currentSchoolId"
      @close="showEditModal = false"
      @submit="handleUserSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { useAdminStore } from '@/stores/admin';
import { storeToRefs } from 'pinia';
import ActionButton from '@/components/shared/ActionButton.vue';
import DataTable from '@/components/shared/DataTable.vue';
import UserImportModal from '@/components/admin/UserImportModal.vue';
import UserEditModal from '@/components/admin/UserEditModal.vue';
import { Plus, Upload, MoreHorizontal, ArrowLeft } from 'lucide-vue-next';
import adminService from '@/services/admin';
import { useNotificationStore } from '@/stores/notification';

const props = defineProps({
    schoolId: {
        type: String,
        default: ''
    }
});

const authStore = useAuthStore();
const adminStore = useAdminStore();
const notificationStore = useNotificationStore();
const router = useRouter();
const route = useRoute();
const { users, totalUsers, loading } = storeToRefs(adminStore);

const currentSchoolId = computed(() => {
    if (props.schoolId) return props.schoolId;
    if (authStore.user?.role === 'Secretary' && authStore.user?.schoolId) {
        return authStore.user.schoolId.toString();
    }
    return ''; // Global admin view or no school context
});

const itemsPerPage = 10;
const filters = reactive({
  role: '',
  status: '',
  query: '',
  page: 1,
  limit: itemsPerPage,
  schoolId: currentSchoolId.value
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
   // Ensure schoolId is passed if present
   if (currentSchoolId.value) {
       filters.schoolId = currentSchoolId.value;
   }
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
    const payload = { ...formData };
    if (props.schoolId) {
        payload.schoolId = parseInt(props.schoolId as string); // Inject school context as number
    }

    if (selectedUser.value) {
       await adminService.updateUser((selectedUser.value as any).id, payload);
       notificationStore.success('User updated successfully');
    } else {
       await adminService.createUser(payload);
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
  if (!confirm(`Are you sure you want to ${user.status === 'active' ? 'deactivate' : 'activate'} ${user.firstName}?`)) return;
  
  try {
    const newStatus = user.status === 'active' ? 'inactive' : 'active';
    await adminService.updateUser(user.id, { status: newStatus, schoolId: props.schoolId ? parseInt(props.schoolId as string) : undefined });
    notificationStore.success(`User ${user.status === 'active' ? 'deactivated' : 'activated'} successfully`);
    fetchUsers();
  } catch (err) {
    notificationStore.error('Failed to update user status');
  }
};

const resetPassword = async (user: any) => {
   const newPassword = prompt(`Enter new password for ${user.firstName}:`);
   if (!newPassword) return;

   try {
     await adminService.updateUser(user.id, { password: newPassword, schoolId: props.schoolId ? parseInt(props.schoolId as string) : undefined });
     notificationStore.success('Password updated successfully');
   } catch (err) {
     notificationStore.error('Failed to reset password');
   }
};

const viewHistory = (user: any) => {
  router.push({ name: 'AuditLogs', query: { userId: user.id } });
};

const onImportSuccess = () => {
  fetchUsers();
};

const goBack = () => {
    router.back();
}

onMounted(() => {
  // Ensure filters has schoolId if prop is present
  if (props.schoolId) {
      filters.schoolId = props.schoolId;
  }
  fetchUsers();
});
</script>
