<template>
  <div class="p-4 space-y-6 h-screen flex flex-col">
    <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100 flex-shrink-0">Messages</h1>

    <div class="flex-1 flex gap-4 overflow-hidden min-h-0">
      <!-- Sidebar / Conversations List -->
      <div class="w-1/3 bg-white dark:bg-gray-800 rounded-lg shadow flex flex-col">
        <div class="p-4 border-b dark:border-gray-700">
           <input type="text" placeholder="Search..." class="w-full rounded-md border-gray-300 dark:bg-gray-700 dark:border-gray-600">
        </div>
        <div class="flex-1 overflow-y-auto">
           <div 
             v-for="msg in messages" 
             :key="msg.id" 
             @click="selectedMessage = msg"
             class="p-4 border-b hover:bg-gray-50 dark:hover:bg-gray-700 cursor-pointer"
             :class="{'bg-indigo-50 dark:bg-indigo-900/20': selectedMessage?.id === msg.id, 'font-bold': !msg.read}"
           >
             <div class="flex justify-between mb-1">
               <span class="text-sm font-semibold">{{ msg.sender }}</span>
               <span class="text-xs text-gray-500">{{ new Date(msg.date).toLocaleDateString() }}</span>
             </div>
             <p class="text-sm text-gray-600 dark:text-gray-400 truncate">{{ msg.subject }}</p>
           </div>
        </div>
      </div>

      <!-- Detail View -->
      <div class="flex-1 bg-white dark:bg-gray-800 rounded-lg shadow flex flex-col">
         <div v-if="selectedMessage" class="flex flex-col h-full">
            <div class="p-6 border-b dark:border-gray-700">
               <h2 class="text-2xl font-bold mb-2">{{ selectedMessage.subject }}</h2>
               <div class="flex justify-between items-center text-sm text-gray-500">
                 <span>From: {{ selectedMessage.sender }}</span>
                 <span>{{ new Date(selectedMessage.date).toLocaleString() }}</span>
               </div>
            </div>
            
            <div class="p-6 flex-1 overflow-y-auto whitespace-pre-wrap">
               {{ selectedMessage.body }}
            </div>
            
            <!-- Reply Area -->
            <div class="p-4 border-t dark:border-gray-700 bg-gray-50 dark:bg-gray-900 text-right">
               <button class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700">Reply</button>
            </div>
         </div>
         <div v-else class="flex-1 flex items-center justify-center text-gray-400">
            Select a message to view
         </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useParentStore } from '@/stores/parent';
import { useStudentStore } from '@/stores/student';
import { useAuthStore } from '@/stores/auth';
import parentApi from '@/services/parent';
// import studentApi from '@/services/student';

const auth = useAuthStore();
const parentStore = useParentStore();
const studentStore = useStudentStore();

const isParent = computed(() => auth.user?.role === 'parent');

const messages = computed(() => {
   return isParent.value ? parentStore.messages : studentStore.messages;
});

const selectedMessage = ref<any>(null);

onMounted(async () => {
   if (isParent.value) {
     if (parentStore.selectedChildId) {
        // Direct API call if store action missing, or implement store action
        // Using direct API for now to fix build, assuming logic is here
        try {
            const res = await parentApi.getMessages(parentStore.selectedChildId);
            parentStore.messages = res.data;
        } catch (e) { console.error(e); }
     }
   } else {
     await studentStore.fetchMessages();
   }
});
</script>
