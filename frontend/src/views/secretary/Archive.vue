<script setup lang="ts">
import { useSecretaryStore } from '@/stores/secretary';
import { onMounted, ref, watch } from 'vue';
import DocumentPrinter from '@/components/secretary/DocumentPrinter.vue';

const store = useSecretaryStore();
const selected = ref<number[]>([]);
const showPrinter = ref(false);

const filters = ref({
    search: '',
    type: 'all',
    dateFrom: '',
    dateTo: ''
});

onMounted(() => {
    store.fetchArchive();
});

watch(filters, (newVal) => {
    store.fetchArchive(newVal);
}, { deep: true });

function toggleSelection(id: number) {
    if (selected.value.includes(id)) {
        selected.value = selected.value.filter(i => i !== id);
    } else {
        selected.value.push(id);
    }
}

function selectAll(e: any) {
    if (e.target.checked) {
        selected.value = store.archive.map(d => d.id);
    } else {
        selected.value = [];
    }
}

function handlePrint(templateId: string) {
    alert(`Printing ${selected.value.length} documents with template: ${templateId}`);
    // store.printDocuments(selected.value);
    console.log('Printing', selected.value);
    showPrinter.value = false;
    selected.value = [];
}

async function handleBulkExport() {
    await store.bulkExport(selected.value);
    selected.value = [];
}
</script>

<template>
  <div class="p-4 space-y-4">
      <div class="flex flex-col gap-4">
          <div class="flex justify-between items-center">
              <h1 class="text-2xl font-bold">Document Archive</h1>
                 <div class="flex gap-2">
                  <button class="btn btn-secondary" :disabled="selected.length === 0" @click="showPrinter = true">
                      Print Selected ({{ selected.length }})
                  </button>
                   <button class="btn btn-accent" :disabled="selected.length === 0" @click="handleBulkExport">
                      Export Archive ({{ selected.length }})
                  </button>
              </div>
          </div>
          
          <div class="card bg-base-100 p-4 shadow-sm grid grid-cols-1 md:grid-cols-4 gap-4">
             <div class="form-control">
                  <label class="label"><span class="label-text">Search Student</span></label>
                  <input type="text" v-model="filters.search" placeholder="Name..." class="input input-bordered input-sm" />
             </div>
             <div class="form-control">
                  <label class="label"><span class="label-text">Type</span></label>
                  <select v-model="filters.type" class="select select-bordered select-sm">
                      <option value="all">All Types</option>
                      <option value="REPORT_CARD">Report Card</option>
                      <option value="CERTIFICATE">Certificate</option>
                      <option value="PDP">PDP</option>
                  </select>
             </div>
             <div class="form-control">
                  <label class="label"><span class="label-text">From</span></label>
                  <input type="date" v-model="filters.dateFrom" class="input input-bordered input-sm" />
             </div>
               <div class="form-control">
                  <label class="label"><span class="label-text">To</span></label>
                  <input type="date" v-model="filters.dateTo" class="input input-bordered input-sm" />
             </div>
          </div>
      </div>

      <div class="overflow-x-auto bg-base-100 rounded-box shadow">
          <table class="table w-full">
              <thead>
                  <tr>
                      <th>
                          <label><input type="checkbox" class="checkbox" @change="selectAll" /></label>
                      </th>
                      <th>Type</th>
                      <th>Student</th>
                      <th>Class</th>
                      <th>Date</th>
                      <th>Status</th>
                  </tr>
              </thead>
              <tbody>
                  <tr v-for="doc in store.archive" :key="doc.id" class="hover">
                      <th>
                          <label>
                              <input type="checkbox" class="checkbox" 
                                     :checked="selected.includes(doc.id)" 
                                     @change="toggleSelection(doc.id)" />
                          </label>
                      </th>
                      <td>{{ doc.type }}</td>
                      <td>{{ doc.student }}</td>
                      <td>{{ doc.class }}</td>
                      <td>{{ doc.date }}</td>
                      <td>
                          <div class="flex items-center gap-2">
                              <div class="badge badge-ghost">{{ doc.status }}</div>
                              <a v-if="doc.status === 'SIGNED' && doc.data?.file_path" 
                                 :href="`/api/files/download?path=${doc.data.file_path}`" 
                                 target="_blank"
                                 class="btn btn-xs btn-outline btn-primary">
                                  Download
                              </a>
                          </div>
                      </td>
                  </tr>
                   <tr v-if="store.archive.length === 0">
                      <td colspan="6" class="text-center py-8 text-gray-400">No documents found matching filters</td>
                  </tr>
              </tbody>
          </table>
      </div>

      <DocumentPrinter :selectedIds="selected" v-if="showPrinter" @cancel="showPrinter = false" @print="handlePrint" />
  </div>
</template>
