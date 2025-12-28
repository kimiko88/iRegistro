<template>
  <div class="p-4 space-y-6">
     <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100">Documents</h1>

    <div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
       <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
         <thead class="bg-gray-50 dark:bg-gray-700">
           <tr>
             <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Document Name</th>
             <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
             <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
             <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Action</th>
           </tr>
         </thead>
         <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
            <tr v-for="doc in documents" :key="doc.id">
              <td class="px-6 py-4 flex items-center">
                 <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                   <path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" clip-rule="evenodd" />
                 </svg>
                 {{ doc.name }}
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ doc.type }}</td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ new Date(doc.date).toLocaleDateString() }}</td>
              <td class="px-6 py-4">
                 <button class="text-indigo-600 hover:text-indigo-900 underline">Download</button>
              </td>
            </tr>
            <tr v-if="documents.length === 0">
              <td colspan="4" class="px-6 py-4 text-center text-gray-500">No documents found.</td>
            </tr>
         </tbody>
       </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { useParentStore } from '@/stores/parent';
import { useStudentStore } from '@/stores/student';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const parentStore = useParentStore();
const studentStore = useStudentStore();

const isParent = computed(() => auth.user?.role === 'parent');

const documents = computed(() => {
   return isParent.value ? parentStore.documents : studentStore.documents;
});

onMounted(() => {
   if (isParent.value) parentStore.fetchDocuments();
   else studentStore.fetchDocuments();
});
</script>
