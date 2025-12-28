<script setup lang="ts">
import DocumentReceiver from '@/components/secretary/DocumentReceiver.vue';
</script>

<template>
  <div class="p-4 space-y-6">
      <h1 class="text-2xl font-bold">Document Management</h1>
      
      <!-- Notifications / Stats Area -->
      <div class="stats shadow bg-base-100 w-full mb-4">
          <div class="stat">
              <div class="stat-title">New Documents</div>
              <div class="stat-value text-primary">2</div>
              <div class="stat-desc">Pending approval</div>
          </div>
           <div class="stat">
              <div class="stat-title">Processed Today</div>
              <div class="stat-value">5</div>
          </div>
      </div>

      <DocumentReceiver />
      
      <div class="flex gap-4 mt-6">
          <router-link to="/secretary/archive" class="btn btn-outline">Go to Archive</router-link>
          <button class="btn btn-outline btn-secondary">Bulk Export</button>
      </div>
  </div>
</template>
