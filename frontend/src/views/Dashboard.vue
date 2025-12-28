<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
</script>

<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
    <div class="stats shadow">
      <div class="stat">
        <div class="stat-figure text-secondary">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
        </div>
        <div class="stat-title">Role</div>
        <div class="stat-value">{{ auth.user?.role || 'User' }}</div>
        <div class="stat-desc">Current Session</div>
      </div>
    </div>

    <!-- More stats placeholders -->
    <div class="card w-full bg-base-100 shadow-xl col-span-1 md:col-span-2">
      <div class="card-body">
        <h2 class="card-title">Welcome back!</h2>
        <p>This is your school management dashboard. Select an option from the sidebar to get started.</p>
        <div class="card-actions justify-end">
          <button class="btn btn-primary">View Details</button>
        </div>
      </div>
    </div>
  </div>
</template>
