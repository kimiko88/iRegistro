<script setup lang="ts">
import { useSecretaryStore } from '@/stores/secretary';
import { onMounted, ref } from 'vue';
import DocumentPrinter from '@/components/secretary/DocumentPrinter.vue';

const store = useSecretaryStore();
const selected = ref<number[]>([]);
const showPrinter = ref(false);

onMounted(() => {
    store.fetchArchive();
});

function toggleSelection(id: number) {
    if (selected.value.includes(id)) {
        selected.value = selected.value.filter(i => i !== id);
    } else {
        selected.value.push(id);
    }
}

function handlePrint() {
    alert(`Printing ${selected.value.length} documents...`);
    showPrinter.value = false;
    selected.value = [];
}
</script>

<template>
  <div class="p-4 space-y-4">
      <div class="flex justify-between items-center">
          <h1 class="text-2xl font-bold">Document Archive</h1>
          <div class="flex gap-2">
              <input type="text" placeholder="Search..." class="input input-bordered" />
              <button class="btn btn-secondary" :disabled="selected.length === 0" @click="showPrinter = true">
                  Print Selected ({{ selected.length }})
              </button>
          </div>
      </div>

      <div class="overflow-x-auto bg-base-100 rounded-box shadow">
          <table class="table w-full">
              <thead>
                  <tr>
                      <th>
                          <label><input type="checkbox" class="checkbox" disabled /></label>
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
              </tbody>
          </table>
      </div>

      <DocumentPrinter :selectedIds="selected" v-if="showPrinter" @cancel="showPrinter = false" @print="handlePrint" />
  </div>
</template>
