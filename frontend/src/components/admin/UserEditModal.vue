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
        <select v-model="form.role" class="select select-bordered">
          <option value="Teacher">Teacher</option>
          <option value="Student">Student</option>
          <option value="Parent">Parent</option>
          <option value="Secretary">Secretary</option>
          <option value="Principal">Principal</option>
           <option value="Admin">Admin</option>
        </select>
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
import { ref, computed, watch } from 'vue';
import FormModal from '@/components/shared/FormModal.vue';

const props = defineProps<{
  isOpen: boolean;
  user?: any;
  loading?: boolean;
}>();

const emit = defineEmits(['close', 'submit']);

const isEdit = computed(() => !!props.user);
const showPassword = ref(false);
const autoGeneratePassword = ref(true);

const form = ref({
  firstName: '',
  lastName: '',
  email: '',
  role: 'Student',
  password: '',
  sendCredentials: true
});

watch(() => props.user, (newUser) => {
  if (newUser) {
    form.value = { ...newUser, password: '', sendCredentials: false };
  } else {
    resetForm();
  }
}, { immediate: true });

function resetForm() {
  form.value = {
    firstName: '',
    lastName: '',
    email: '',
    role: 'Student',
    password: '',
    sendCredentials: true
  };
  autoGeneratePassword.value = true;
}

const onSubmit = () => {
    // If auto-generate, maybe clear password so backend handles it
    if (autoGeneratePassword.value && !isEdit.value) {
        form.value.password = '';
    }
    emit('submit', form.value);
};
</script>
