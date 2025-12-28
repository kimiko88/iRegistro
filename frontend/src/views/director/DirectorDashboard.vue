<template>
    <div class="p-6">
        <h1 class="text-3xl font-bold mb-6">Director Dashboard</h1>

        <div v-if="store.loading && !store.kpis" class="flex justify-center p-10">
            <span class="loading loading-spinner loading-lg"></span>
        </div>

        <div v-else-if="store.kpis" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
            <div class="stat bg-base-100 shadow rounded-box">
                <div class="stat-figure text-primary">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>
                </div>
                <div class="stat-title">Total Students</div>
                <div class="stat-value text-primary">{{ store.kpis.totalStudents }}</div>
                <div class="stat-desc">Enrollment stable</div>
            </div>

            <div class="stat bg-base-100 shadow rounded-box">
                <div class="stat-figure text-secondary">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path></svg>
                </div>
                <div class="stat-title">Teachers</div>
                <div class="stat-value text-secondary">{{ store.kpis.totalTeachers }}</div>
                <div class="stat-desc">All positions filled</div>
            </div>

            <div class="stat bg-base-100 shadow rounded-box">
                <div class="stat-figure text-accent">
                   <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-8 h-8"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                </div>
                <div class="stat-title">Attendance Rate</div>
                <div class="stat-value text-accent">{{ store.kpis.attendanceRate }}%</div>
                <div class="stat-desc">Last 30 days</div>
            </div>

             <div class="stat bg-base-100 shadow rounded-box">
                <div class="stat-title">Documents to Sign</div>
                <div class="stat-value">{{ store.documentsToSign.length }}</div>
                <div class="stat-desc text-warning cursor-pointer hover:underline" @click="$router.push('/director/signing')">Action required</div>
            </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Recent Activity / Monitoring -->
             <div class="card bg-base-100 shadow-xl">
                <div class="card-body">
                    <h2 class="card-title">Real-time Monitoring</h2>
                    <ul class="steps steps-vertical">
                        <li class="step step-primary">Vote Uploads Complete (Class 5A)</li>
                        <li class="step step-primary">Attendance Registered (All Classes)</li>
                        <li class="step">Afternoon Labs ( In Progress )</li>
                        <li class="step">Digital Register Sync (Pending)</li>
                    </ul>
                </div>
            </div>

            <!-- Quick Actions -->
             <div class="card bg-base-100 shadow-xl">
                <div class="card-body">
                    <h2 class="card-title">Quick Actions</h2>
                     <div class="flex flex-col gap-3 mt-4">
                        <button class="btn btn-outline" @click="$router.push('/director/approvals')">
                            Review Pending Requests ({{ store.pendingRequests.length }})
                        </button>
                        <button class="btn btn-outline" @click="$router.push('/director/reports')">
                            Generate School Reports
                        </button>
                         <button class="btn btn-outline btn-error">
                            Emergency Communication
                        </button>
                     </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useDirectorStore } from '@/stores/director';

const store = useDirectorStore();

onMounted(() => {
    store.fetchDashboardData();
});
</script>
