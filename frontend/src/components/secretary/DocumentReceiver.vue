<script setup lang="ts">
import { useSecretaryStore } from '@/stores/secretary';
import { onMounted } from 'vue';

const store = useSecretaryStore();

onMounted(() => {
    store.fetchInbox();
});

function onApprove(id: number) {
    store.approveDocument(id);
}

function onReject(id: number) {
    store.rejectDocument(id);
}
</script>

<template>
  <div class="card bg-base-100 shadow">
      <div class="card-body">
          <h2 class="card-title">Document Inbox</h2>
          <div class="overflow-x-auto">
              <table class="table w-full">
                  <thead>
                      <tr>
                          <th>Date</th>
                          <th>Type</th>
                          <th>Student</th>
                          <th>Action</th>
                      </tr>
                  </thead>
                  <tbody>
                      <tr v-for="doc in store.inbox" :key="doc.id">
                          <td>{{ doc.date }}</td>
                          <td>
                              <div class="badge badge-outline">{{ doc.type }}</div>
                          </td>
                          <td>
                              <div class="font-bold">{{ doc.student }}</div>
                              <div class="text-sm opacity-50">{{ doc.class }}</div>
                          </td>
                          <td class="flex gap-2">
                              <button class="btn btn-sm btn-success btn-circle" @click="onApprove(doc.id)" title="Approve">
                                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg>
                              </button>
                               <button class="btn btn-sm btn-error btn-circle" @click="onReject(doc.id)" title="Reject">
                                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
                              </button>
                          </td>
                      </tr>
                      <tr v-if="store.inbox.length === 0">
                          <td colspan="4" class="text-center text-gray-500">No new documents</td>
                      </tr>
                  </tbody>
              </table>
          </div>
      </div>
  </div>
</template>
