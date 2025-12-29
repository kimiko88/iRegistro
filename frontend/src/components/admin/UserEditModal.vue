<template>
  <FormModal
    :isOpen="isOpen"
    :title="isEdit ? 'Edit User' : 'Create User'"
    :submitLabel="isEdit ? 'Update' : 'Create'"
    :loading="loading"
    @close="$emit('close')"
    @submit="onSubmit"
  >
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="form-control">
        <label class="label"><span class="label-text">First Name</span></label>
        <input type="text" v-model="form.firstName" class="input input-bordered" required />
      </div>

      <div class="form-control">
        <label class="label"><span class="label-text">Last Name</span></label>
        <input type="text" v-model="form.lastName" class="input input-bordered" required />
      </div>

      <div class="form-control">
        <label class="label"><span class="label-text">Email</span></label>
        <input type="email" v-model="form.email" class="input input-bordered" required />
      </div>

      <div class="form-control">
        <label class="label"><span class="label-text">Role</span></label>
        <select v-model="form.role" class="select select-bordered" @change="handleRoleChange">
          <option value="Teacher">Teacher</option>
          <option value="Student">Student</option>
          <option value="Parent">Parent</option>
          <option value="Secretary">Secretary</option>
          <option value="Principal">Principal</option>
          <option value="Admin">Admin</option>
          <option value="SuperAdmin">SuperAdmin</option>
        </select>
      </div>

      <!-- School Selection -->
      <div class="form-control md:col-span-2" v-if="shouldShowSchoolSelect">
         <label class="label"><span class="label-text">School</span></label>
         <div class="relative">
             <input 
                type="text" 
                class="input input-bordered w-full" 
                v-model="schoolSearch"
                @input="searchSchools"
                placeholder="Search school..." 
                v-if="!form.schoolId"
             />
             <div v-else class="flex items-center justify-between p-3 border rounded-lg bg-gray-50 dark:bg-gray-800">
                <span class="font-medium">{{ getSchoolName(form.schoolId) }}</span>
                <button @click="clearSchool" type="button" class="btn btn-xs btn-ghost btn-circle">✕</button>
             </div>
             
             <ul v-if="filteredSchools.length > 0 && !form.schoolId" class="absolute z-10 w-full p-2 mt-1 overflow-y-auto shadow menu bg-base-100 rounded-box max-h-48">
                 <li v-for="school in filteredSchools" :key="school.id">
                     <a @click="selectSchool(school)">{{ school.name }} ({{ school.code }})</a>
                 </li>
             </ul>
         </div>
         <label class="label" v-if="!form.schoolId && form.role !== 'SuperAdmin'">
             <span class="label-text-alt text-error">School is required for this role</span>
         </label>
      </div>

      <div class="form-control" v-if="!isEdit">
        <label class="label"><span class="label-text">Password</span></label>
        <div class="flex gap-2">
           <input :type="showPassword ? 'text' : 'password'" v-model="form.password" class="input input-bordered w-full" :disabled="autoGeneratePassword" />
           <button type="button" class="btn btn-square btn-ghost" @click="showPassword = !showPassword">
             👁️
           </button>
        </div>
        <label class="cursor-pointer label justify-start gap-2">
          <input type="checkbox" class="toggle toggle-sm" v-model="autoGeneratePassword" />
          <span class="label-text">Auto-generate</span>
        </label>
      </div>
      
      <div class="form-control md:col-span-2">
         <label class="cursor-pointer label justify-start gap-2">
          <input type="checkbox" class="toggle toggle-sm" v-model="form.sendCredentials" />
          <span class="label-text">Send email with credentials</span>
        </label>
      </div>
    </div>
  </FormModal>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import FormModal from '@/components/shared/FormModal.vue';
import { useAdminStore } from '@/stores/admin';
import { storeToRefs } from 'pinia';

const props = defineProps<{
  isOpen: boolean;
  user?: any;
  loading?: boolean;
  showSchoolSelect?: boolean;
  preselectedSchoolId?: string | number;
}>();

const emit = defineEmits(['close', 'submit']);

const adminStore = useAdminStore();
const { schools } = storeToRefs(adminStore);
// Ensure schools are loaded if needed
const schoolsLoaded = ref(false);

const isEdit = computed(() => !!props.user);
const showPassword = ref(false);
const autoGeneratePassword = ref(true);

const form = ref({
  firstName: '',
  lastName: '',
  email: '',
  role: 'Student',
  schoolId: null as number | null,
  password: '',
  sendCredentials: true
});

const schoolSearch = ref('');
const filteredSchools = ref<any[]>([]);

const shouldShowSchoolSelect = computed(() => {
    return props.showSchoolSelect && form.value.role !== 'SuperAdmin';
});

watch(() => props.user, (newUser) => {
  if (newUser) {
    form.value = { ...newUser, password: '', sendCredentials: false };
    if (newUser.schoolId) {
        form.value.schoolId = newUser.schoolId;
    }
  } else {
    resetForm();
  }
}, { immediate: true });

watch(() => props.isOpen, async (isOpen) => {
    if (isOpen && props.showSchoolSelect && !schoolsLoaded.value) {
        await adminStore.fetchSchools();
        schoolsLoaded.value = true;
    }
    // Set initial school if provided
    if (isOpen && !props.user && props.preselectedSchoolId) {
        form.value.schoolId = Number(props.preselectedSchoolId);
    }
});

function resetForm() {
  form.value = {
    firstName: '',
    lastName: '',
    email: '',
    role: 'Student',
    schoolId: props.preselectedSchoolId ? Number(props.preselectedSchoolId) : null,
    password: '',
    sendCredentials: true
  };
  autoGeneratePassword.value = true;
  schoolSearch.value = '';
}

const searchSchools = () => {
    if (!schoolSearch.value || schoolSearch.value.length < 2) {
        filteredSchools.value = [];
        return;
    }
    const q = schoolSearch.value.toLowerCase();
    filteredSchools.value = schools.value.filter((s:any) => 
        s.name.toLowerCase().includes(q) || s.code.toLowerCase().includes(q)
    ).slice(0, 5);
};

const selectSchool = (school: any) => {
    form.value.schoolId = school.id;
    filteredSchools.value = [];
    schoolSearch.value = '';
};

const clearSchool = () => {
    form.value.schoolId = null;
};

const getSchoolName = (id: number) => {
    const s = schools.value.find((s:any) => s.id === id);
    return s ? `${s.name} (${s.code})` : 'Unknown School';
};

const handleRoleChange = () => {
    if (form.value.role === 'SuperAdmin') {
        form.value.schoolId = null;
    } else if (props.preselectedSchoolId) {
        form.value.schoolId = Number(props.preselectedSchoolId);
    }
};

const onSubmit = () => {
    // Validation
    if (shouldShowSchoolSelect.value && !form.value.schoolId) {
        alert("Please select a school for this user.");
        return;
    }

    // If auto-generate, maybe clear password so backend handles it
    if (autoGeneratePassword.value && !isEdit.value) {
        form.value.password = '';
    }
    
    // Clean up payload
    const payload = { ...form.value };
    if (payload.role === 'SuperAdmin') {
        payload.schoolId = 0; // Or null, depending on backend. Backend uses 0 for SuperAdmin.
    }

    emit('submit', payload);
};
</script>
