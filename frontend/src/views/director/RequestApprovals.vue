<template>
    <div class="p-6">
        <h1 class="text-3xl font-bold mb-6">Pending Approvals</h1>

        <div class="grid gap-4">
             <div v-for="req in store.pendingRequests" :key="req.id" class="card bg-base-100 shadow-lg border-l-4 border-warning">
                <div class="card-body">
                    <div class="flex justify-between items-start">
                        <div>
                            <h2 class="card-title">{{ req.type.replace('_', ' ') }}</h2>
                            <p class="text-gray-500 text-sm">Requested by: {{ req.requester }} on {{ req.date }}</p>
                            <p class="mt-2 text-lg">{{ req.details }}</p>
                        </div>
                        <div class="flex flex-col gap-2">
                            <button class="btn btn-success btn-sm text-white" @click="store.approveRequest(req.id)">Approve</button>
                            <button class="btn btn-error btn-sm text-white" @click="store.rejectRequest(req.id)">Reject</button>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="store.pendingRequests.length === 0" class="text-center py-10 text-gray-400">
                No pending requests.
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useDirectorStore } from '@/stores/director';

const store = useDirectorStore();

onMounted(() => {
    store.fetchRequests();
});
</script>
