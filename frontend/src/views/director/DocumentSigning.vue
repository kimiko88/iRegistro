<template>
    <div class="p-6">
        <h1 class="text-3xl font-bold mb-6">Document Signing</h1>

        <div v-if="store.documentsToSign.length === 0" class="alert alert-success">
            <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
            <span>No pending documents to sign. Great job!</span>
        </div>

        <div v-else class="overflow-x-auto">
            <table class="table w-full bg-base-100 shadow-xl rounded-box">
                <thead>
                    <tr>
                        <th>Title</th>
                        <th>Type</th>
                        <th>Student / Target</th>
                        <th>Date Created</th>
                        <th>Action</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="doc in store.documentsToSign" :key="doc.id">
                        <td class="font-bold">{{ doc.title }}</td>
                        <td><div class="badge badge-ghost">{{ doc.type }}</div></td>
                        <td>{{ doc.studentName }}</td>
                        <td>{{ doc.date }}</td>
                        <td>
                            <button class="btn btn-sm btn-primary" @click="openSignModal(doc)">
                                Sign Digitally
                            </button>
                             <button class="btn btn-sm btn-ghost ml-2">View</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <DigitalSignatureModal 
            :isOpen="isModalOpen" 
            :loading="store.loading"
            @confirm="handleSign" 
            @cancel="isModalOpen = false" 
        />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useDirectorStore } from '@/stores/director';
// @ts-ignore
import DigitalSignatureModal from '@/components/director/DigitalSignatureModal.vue';

const store = useDirectorStore();
const isModalOpen = ref(false);
const selectedDocId = ref<number | null>(null);

onMounted(() => {
    store.fetchDocumentsToSign();
});

const openSignModal = (doc: any) => {
    selectedDocId.value = doc.id;
    isModalOpen.value = true;
};

const handleSign = async (pin: string) => {
    if (selectedDocId.value) {
        await store.signDocument(selectedDocId.value, pin);
        isModalOpen.value = false;
        selectedDocId.value = null;
        // Show success toast? (Assuming global toast or simple alert for now)
    }
};
</script>
