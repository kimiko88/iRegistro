<template>
  <div class="p-6 space-y-6">
    <h1 class="text-3xl font-bold">School Settings</h1>

    <div role="tablist" class="tabs tabs-lifted">
      
      <!-- General Info Tab -->
      <input type="radio" name="school_tabs" role="tab" class="tab" aria-label="General Info" checked />
      <div role="tabpanel" class="tab-content bg-base-100 border-base-300 rounded-box p-6">
        <h2 class="text-xl font-bold mb-4">School Information</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="form-control">
            <label class="label"><span class="label-text">School Name</span></label>
            <input type="text" v-model="school.name" class="input input-bordered" />
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text">Meccanographic Code</span></label>
            <input type="text" v-model="school.code" class="input input-bordered" />
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text">Address</span></label>
            <input type="text" v-model="school.address" class="input input-bordered" />
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text">City</span></label>
            <input type="text" v-model="school.city" class="input input-bordered" />
          </div>
        </div>
        <div class="mt-4 flex justify-end">
          <ActionButton label="Save Changes" @click="saveGeneralInfo" :loading="loading" />
        </div>
      </div>

      <!-- Plessi & Indirizzi Tab -->
      <input type="radio" name="school_tabs" role="tab" class="tab" aria-label="Branches & Addresses" />
      <div role="tabpanel" class="tab-content bg-base-100 border-base-300 rounded-box p-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
             <div class="flex justify-between items-center mb-2">
               <h3 class="font-bold">Branches (Plessi)</h3>
               <button class="btn btn-xs btn-primary" @click="addBranch">Add</button>
             </div>
             <ul class="list-disc pl-4 bg-base-200 p-4 rounded-lg">
               <li v-for="(branch, index) in branches" :key="index" class="flex justify-between items-center mb-2">
                 <span>{{ branch.name }}</span>
                 <button class="btn btn-circle btn-xs btn-ghost text-error" @click="removeBranch(index)">✕</button>
               </li>
             </ul>
          </div>
          <div>
             <div class="flex justify-between items-center mb-2">
               <h3 class="font-bold">Education Paths (Indirizzi)</h3>
               <button class="btn btn-xs btn-primary" @click="addPath">Add</button>
             </div>
             <ul class="list-disc pl-4 bg-base-200 p-4 rounded-lg">
               <li v-for="(path, index) in paths" :key="index" class="flex justify-between items-center mb-2">
                 <span>{{ path.name }}</span>
                 <button class="btn btn-circle btn-xs btn-ghost text-error" @click="removePath(index)">✕</button>
               </li>
             </ul>
          </div>
        </div>
      </div>

      <!-- Academic Years Tab -->
      <input type="radio" name="school_tabs" role="tab" class="tab" aria-label="Academic Years" />
      <div role="tabpanel" class="tab-content bg-base-100 border-base-300 rounded-box p-6">
         <h2 class="text-xl font-bold mb-4">Academic Years</h2>
         <div class="flex justify-end mb-4">
            <button class="btn btn-sm btn-primary">Add Year</button>
         </div>
         <table class="table">
           <thead>
             <tr>
               <th>Year</th>
               <th>Start Date</th>
               <th>End Date</th>
               <th>Status</th>
               <th>Actions</th>
             </tr>
           </thead>
           <tbody>
             <tr v-for="year in years" :key="year.id">
               <td>{{ year.name }}</td>
               <td>{{ year.startDate }}</td>
               <td>{{ year.endDate }}</td>
               <td>
                 <input type="checkbox" class="toggle toggle-success toggle-sm" :checked="year.active" />
               </td>
               <td>
                 <button class="btn btn-xs">Edit</button>
               </td>
             </tr>
           </tbody>
         </table>
      </div>

      <!-- Templates & Reports Tab -->
      <input type="radio" name="school_tabs" role="tab" class="tab" aria-label="Branding & Reports" />
      <div role="tabpanel" class="tab-content bg-base-100 border-base-300 rounded-box p-6">
        <div class="space-y-6">
           <div>
             <h3 class="font-bold mb-2">School Logo</h3>
             <div class="flex items-center gap-4">
               <div class="w-24 h-24 bg-base-200 rounded flex items-center justify-center">
                 <img v-if="school.logoUrl" :src="school.logoUrl" class="max-w-full max-h-full" />
                 <span v-else class="text-xs text-gray-500">No Logo</span>
               </div>
               <input type="file" class="file-input file-input-bordered file-input-sm" />
             </div>
           </div>
           
           <div>
             <h3 class="font-bold mb-2">Theme Colors</h3>
             <div class="flex gap-4">
               <div class="form-control">
                 <label class="label"><span class="label-text">Primary</span></label>
                 <input type="color" v-model="school.primaryColor" class="input h-10 w-20 p-1" />
               </div>
               <div class="form-control">
                 <label class="label"><span class="label-text">Secondary</span></label>
                 <input type="color" v-model="school.secondaryColor" class="input h-10 w-20 p-1" />
               </div>
             </div>
           </div>

           <div>
             <h3 class="font-bold mb-2">Report Card Template</h3>
             <textarea class="textarea textarea-bordered w-full" placeholder="Customize header/footer text for PDF reports..."></textarea>
           </div>
           
           <div class="flex justify-end">
             <ActionButton label="Update Branding" variant="secondary" />
           </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import ActionButton from '@/components/shared/ActionButton.vue';

const loading = ref(false);

const school = reactive({
  name: 'Istituto Comprensivo G. Verdi',
  code: 'MIIC8A0002',
  address: 'Via Roma 1',
  city: 'Milano',
  logoUrl: '',
  primaryColor: '#570df8',
  secondaryColor: '#f000b8'
});

const branches = ref([
  { name: 'Plesso Centrale via Roma' },
  { name: 'Succursale via Dante' }
]);

const paths = ref([
  { name: 'Scuola Primaria' },
  { name: 'Scuola Secondaria di I Grado' }
]);

const years = ref([
  { id: 1, name: '2023/2024', startDate: '2023-09-12', endDate: '2024-06-08', active: true },
  { id: 2, name: '2022/2023', startDate: '2022-09-10', endDate: '2023-06-09', active: false }
]);

const saveGeneralInfo = async () => {
  loading.value = true;
  // TODO: Implement API call
  setTimeout(() => { loading.value = false; alert('Saved'); }, 1000);
};

const addBranch = () => {
  const name = prompt('Branch Name:');
  if (name) branches.value.push({ name });
};

const removeBranch = (index: number) => {
  if (confirm('Remove branch?')) branches.value.splice(index, 1);
};

const addPath = () => {
    const name = prompt('Path Name:');
   if (name) paths.value.push({ name });
};

const removePath = (index: number) => {
    if (confirm('Remove path?')) paths.value.splice(index, 1);
};
</script>
