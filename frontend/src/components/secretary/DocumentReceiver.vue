<script setup lang="ts">
import { useSecretaryStore } from '@/stores/secretary';
import { onMounted, ref } from 'vue';
import Modal from '@/components/shared/Modal.vue';

const store = useSecretaryStore();
const selectedDoc = ref<any>(null);
const isReviewOpen = ref(false);
const rejectReason = ref('');
const isRejectMode = ref(false);

onMounted(() => {
    store.fetchInbox();
});

function openReview(doc: any) {
    selectedDoc.value = doc;
    isReviewOpen.value = true;
    isRejectMode.value = false;
    rejectReason.value = '';
}

async function handleApprove() {
    if (selectedDoc.value) {
        await store.approveDocument(selectedDoc.value.id);
        isReviewOpen.value = false;
        selectedDoc.value = null;
    }
}

async function handleReject() {
    if (selectedDoc.value) {
        await store.rejectDocument(selectedDoc.value.id, rejectReason.value);
        isReviewOpen.value = false;
        selectedDoc.value = null;
    }
}
</script>

<template>
  <div class="card bg-base-100 shadow">
      <div class="card-body">
          <div class="flex justify-between items-center mb-4">
            <h2 class="card-title">Document Inbox</h2>
            <button class="btn btn-sm btn-ghost" @click="store.fetchInbox()">Refresh</button>
          </div>
          
          <div class="overflow-x-auto">
              <table class="table w-full">
                  <thead>
                      <tr>
                          <th>Date</th>
                          <th>Type</th>
                          <th>Student</th>
                          <th>Workflow</th>
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
                          <td>
                              <!-- Mock workflow status -->
                              <ul class="steps steps-vertical lg:steps-horizontal text-xs">
                                  <li class="step step-primary">Generated</li>
                                  <li class="step step-primary">Received</li>
                                  <li class="step">Approval</li>
                              </ul>
                          </td>
                          <td class="flex gap-2">
                              <button class="btn btn-sm btn-primary" @click="openReview(doc)">
                                  Review
                              </button>
                          </td>
                      </tr>
                      <tr v-if="store.inbox.length === 0">
                          <td colspan="5" class="text-center text-gray-500 py-8">
                              No new documents pending approval
                          </td>
                      </tr>
                  </tbody>
              </table>
          </div>
      </div>

      <!-- Review Modal -->
      <Modal :isOpen="isReviewOpen" title="Review Document" @close="isReviewOpen = false">
          <div v-if="selectedDoc" class="space-y-4">
              <div class="alert alert-info shadow-sm">
                  <div>
                      <h3 class="font-bold">{{ selectedDoc.type }}</h3>
                      <div class="text-xs">Student: {{ selectedDoc.student }} ({{ selectedDoc.class }})</div>
                  </div>
              </div>
              
              <div class="mockup-window border border-base-300 bg-base-300 min-h-[200px] flex items-center justify-center">
                  <div class="text-center opacity-50">
                      [Document Preview Placeholder]<br>
                      {{ selectedDoc.type }}_Preview.pdf
                  </div>
              </div>

               <div v-if="isRejectMode" class="form-control w-full">
                  <label class="label">
                      <span class="label-text">Reason for rejection</span>
                  </label>
                  <textarea class="textarea textarea-bordered h-24" v-model="rejectReason" placeholder="e.g. Missing signature, incorrect grade..."></textarea>
              </div>
          </div>
          
          <template #actions>
               <button v-if="!isRejectMode" class="btn btn-error btn-outline" @click="isRejectMode = true">Reject</button>
               <button v-if="isRejectMode" class="btn btn-ghost" @click="isRejectMode = false">Cancel Reject</button>
               
               <button v-if="isRejectMode" class="btn btn-error" @click="handleReject" :disabled="!rejectReason">Confirm Rejection</button>
               <button v-if="!isRejectMode" class="btn btn-success" @click="handleApprove">Approve & Sign</button>
          </template>
      </Modal>
  </div>
</template>
