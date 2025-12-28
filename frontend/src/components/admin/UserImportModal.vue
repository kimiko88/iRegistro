<template>
  <FormModal
    :isOpen="isOpen"
    title="Import Users"
    submitLabel="Import"
    :loading="loading"
    :disabled="!parsedData.length || loading"
    @close="$emit('close')"
    @submit="importUsers"
  >
    <div class="space-y-4">
      <div v-if="!parsedData.length">
        <label class="block mb-2 text-sm font-medium text-gray-900">Upload CSV File</label>
        <input 
          type="file" 
          accept=".csv" 
          class="file-input file-input-bordered w-full" 
          @change="onFileChange"
        />
        <p class="mt-1 text-sm text-gray-500">Supported columns: email, name, surname, role</p>
      </div>

      <div v-else>
        <!-- Mapping Section -->
        <div class="bg-base-200 p-4 rounded-lg mb-4">
          <h4 class="font-bold mb-2">Column Mapping</h4>
          <div class="grid grid-cols-2 gap-4">
             <div v-for="(target, key) in targets" :key="key" class="form-control">
               <label class="label"><span class="label-text">{{ target.label }}</span></label>
               <select class="select select-bordered select-sm" v-model="mapping[key]">
                 <option value="">-- Select Column --</option>
                 <option v-for="header in headers" :key="header" :value="header">
                   {{ header }}
                 </option>
               </select>
             </div>
          </div>
        </div>

        <!-- Preview Table -->
        <h4 class="font-bold mb-2">Preview (First 5 rows)</h4>
        <div class="overflow-x-auto">
          <table class="table table-xs">
            <thead>
              <tr>
                <th v-for="header in headers" :key="header">{{ header }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(row, i) in parsedData.slice(0, 5)" :key="i">
                <td v-for="header in headers" :key="header">{{ row[header] }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Progress Bar if importing -->
        <div v-if="loading" class="mt-4">
          <progress class="progress progress-primary w-full" :value="progress" max="100"></progress>
          <p class="text-center text-xs mt-1">{{ progress }}%</p>
        </div>
      </div>
      
      <div v-if="error" class="alert alert-error text-sm mt-4">
        {{ error }}
      </div>
    </div>
  </FormModal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import FormModal from '@/components/shared/FormModal.vue';
import { useUserImport } from '@/composables/useUserImport';

const props = defineProps<{ isOpen: boolean }>();
const emit = defineEmits(['close', 'import-success']);

const { 
  parseCSV, 
  headers, 
  parsedData, 
  mapping, 
  loading, 
  progress, 
  error, 
  importUsers: runImport,
  importStats
} = useUserImport();

const targets = {
  email: { label: 'Email' },
  firstName: { label: 'First Name' },
  lastName: { label: 'Last Name' },
  role: { label: 'Role' }
};

const onFileChange = (e: Event) => {
  const file = (e.target as HTMLInputElement).files?.[0];
  if (file) {
    parseCSV(file);
  }
};

const importUsers = async () => {
  await runImport();
  if (!error.value) {
    emit('import-success', importStats.value);
    emit('close');
  }
};

// Reset state when modal opens/closes if needed
watch(() => props.isOpen, (newVal) => {
  if (!newVal) {
    parsedData.value = [];
    headers.value = [];
    error.value = null;
  }
});
</script>
