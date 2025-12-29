<script setup lang="ts">
import DocumentReceiver from '@/components/secretary/DocumentReceiver.vue';
import { useSecretaryStore } from '@/stores/secretary';
import { onMounted, ref } from 'vue';
import { 
  FileText, 
  CheckCircle, 
  AlertTriangle, 
  Inbox, 
  Truck,
  Archive,
  Download
} from 'lucide-vue-next';

const store = useSecretaryStore();
const activeTab = ref('inbox');

onMounted(async () => {
   await store.fetchInbox();
   await store.fetchStats();
});
</script>

<template>
  <div class="p-6 space-y-8 bg-gray-50 min-h-screen dark:bg-gray-900">
      <div class="flex justify-between items-center">
        <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100 flex items-center gap-3">
          <FileText class="w-8 h-8 text-primary" />
          Document Management
        </h1>
        <div class="flex gap-3">
             <router-link to="/secretary/archive" class="btn btn-outline btn-sm gap-2">
                <Archive class="w-4 h-4"/>
                Archive
             </router-link>
             <button class="btn btn-secondary btn-sm gap-2" @click="store.bulkExport([])">
                <Download class="w-4 h-4"/>
                Export All
             </button>
        </div>
      </div>
      
      <!-- Stats / KPI Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="stat bg-white dark:bg-gray-800 shadow-md rounded-2xl p-6 border border-gray-100 dark:border-gray-700">
              <div class="stat-figure text-primary">
                  <div class="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center">
                    <Inbox class="w-6 h-6" />
                  </div>
              </div>
              <div class="stat-title font-medium text-gray-500">New Documents</div>
              <div class="stat-value text-primary text-4xl mt-1">{{ store.stats.new_documents }}</div>
              <div class="stat-desc text-gray-400 mt-1">Pending approval</div>
          </div>
          
           <div class="stat bg-white dark:bg-gray-800 shadow-md rounded-2xl p-6 border border-gray-100 dark:border-gray-700">
              <div class="stat-figure text-success">
                  <div class="w-12 h-12 rounded-full bg-success/10 flex items-center justify-center">
                    <CheckCircle class="w-6 h-6" />
                  </div>
              </div>
              <div class="stat-title font-medium text-gray-500">Processed Today</div>
              <div class="stat-value text-gray-800 dark:text-gray-100 text-4xl mt-1">{{ store.stats.processed_today }}</div>
              <div class="stat-desc text-green-600 mt-1 font-medium">Synced</div>
          </div>

          <div class="stat bg-white dark:bg-gray-800 shadow-md rounded-2xl p-6 border border-gray-100 dark:border-gray-700">
              <div class="stat-figure text-error">
                  <div class="w-12 h-12 rounded-full bg-error/10 flex items-center justify-center">
                    <AlertTriangle class="w-6 h-6" />
                  </div>
              </div>
              <div class="stat-title font-medium text-gray-500">Delivery Issues</div>
              <div class="stat-value text-error text-4xl mt-1">{{ store.stats.delivery_issues }}</div>
              <div class="stat-desc text-error/80 mt-1">Requires attention</div>
          </div>
      </div>

      <!-- Navigation Tabs -->
      <div role="tablist" class="tabs tabs-boxed bg-white dark:bg-gray-800 p-1 rounded-xl shadow-sm inline-flex">
          <a role="tab" class="tab rounded-lg transition-all duration-200" 
             :class="{ 'tab-active bg-primary text-white shadow': activeTab === 'inbox', 'text-gray-500 hover:text-gray-700': activeTab !== 'inbox' }" 
             @click="activeTab = 'inbox'">
             <Inbox class="w-4 h-4 mr-2" /> Inbox
          </a>
          <a role="tab" class="tab rounded-lg transition-all duration-200" 
             :class="{ 'tab-active bg-primary text-white shadow': activeTab === 'delivery', 'text-gray-500 hover:text-gray-700': activeTab !== 'delivery' }" 
             @click="activeTab = 'delivery'; store.fetchDeliveryReports()">
             <Truck class="w-4 h-4 mr-2" /> Delivery Tracking
          </a>
      </div>

      <!-- Tab Content: Inbox -->
      <div v-if="activeTab === 'inbox'" class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
          <DocumentReceiver />
      </div>

      <!-- Tab Content: Delivery -->
      <div v-else-if="activeTab === 'delivery'" class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6">
           <div class="flex items-center justify-between mb-6">
                <h2 class="card-title flex items-center gap-2">
                    <Truck class="w-5 h-5 text-gray-500"/>
                    Delivery Status
                </h2>
           </div>
           
           <div class="alert alert-info shadow-sm mb-6 bg-blue-50 text-blue-800 border-blue-100 rounded-xl">
               <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
               <span>This view tracks which parents/students have opened their documents, ensuring accountability.</span>
           </div>

           <div class="overflow-x-auto">
               <table class="table w-full">
                   <thead>
                       <tr class="bg-gray-50 dark:bg-gray-700/50 text-gray-600 dark:text-gray-300">
                           <th class="rounded-l-lg py-4">Document</th>
                           <th class="py-4">Recipient</th>
                           <th class="py-4">Status</th>
                           <th class="rounded-r-lg py-4">Last Action</th>
                       </tr>
                   </thead>
                    <tbody class="text-sm">
                       <tr v-if="store.deliveryReports.length === 0">
                           <td colspan="4" class="text-center py-8 text-gray-500">No delivery reports found</td>
                       </tr>
                       <tr v-for="report in store.deliveryReports" :key="report.id" class="hover:bg-gray-50 dark:hover:bg-gray-700/30 transition-colors">
                           <td class="font-medium">{{ report.documentName }}</td>
                           <td class="text-gray-500">{{ report.recipientName }} ({{ report.recipientRole }})</td>
                           <td>
                               <span v-if="report.status === 'DELIVERED'" class="badge badge-success gap-1 pl-1 pr-3 py-3">
                                   <CheckCircle class="w-3 h-3"/> Delivered
                               </span>
                               <span v-else-if="report.status === 'READ'" class="badge badge-info gap-1 pl-1 pr-3 py-3">
                                   <CheckCircle class="w-3 h-3"/> Read
                               </span>
                               <span v-else class="badge badge-warning gap-1 pl-1 pr-3 py-3">
                                   <AlertTriangle class="w-3 h-3"/> {{ report.status }}
                               </span>
                           </td>
                           <td class="text-gray-500 font-mono text-xs">{{ new Date(report.lastActionAt).toLocaleString() }}</td>
                       </tr>
                   </tbody>
               </table>
           </div>
      </div>
  </div>
</template>
