<script setup lang="ts">
import { onMounted } from 'vue';
import { useParentStore } from '@/stores/parent';
import StudentSelect from './StudentSelect.vue';
import MarksView from '@/views/marks/MarksView.vue'; // Reuse as component or link to route
// Using generic views as components inside dashboard for quick overview, 
// or providing links to full views. For dashboard, better to show summary.

const parentStore = useParentStore();

onMounted(() => {
    parentStore.fetchChildren();
});
</script>

<template>
  <div class="p-4 space-y-6">
      <div class="flex flex-col md:flex-row justify-between items-center gap-4">
          <h1 class="text-2xl font-bold">Parent Dashboard</h1>
          <StudentSelect />
      </div>

      <div v-if="parentStore.currentChildOverview" class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- KPI Cards -->
          <div class="stats shadow bg-base-100">
              <div class="stat">
                  <div class="stat-title">Average (GPA)</div>
                  <div class="stat-value text-primary">{{ parentStore.currentChildOverview.gpa }}</div>
                  <div class="stat-desc">Latest term</div>
              </div>
          </div>
          <div class="stats shadow bg-base-100">
              <div class="stat">
                  <div class="stat-title">Attendance</div>
                  <div class="stat-value" :class="parentStore.currentChildOverview.attendance < 90 ? 'text-warning' : 'text-success'">
                      {{ parentStore.currentChildOverview.attendance }}%
                  </div>
              </div>
          </div>
          <div class="stats shadow bg-base-100 cursor-pointer hover:bg-base-200">
              <div class="stat">
                  <div class="stat-title">Next Colloquium</div>
                  <div class="stat-value text-lg">{{ parentStore.currentChildOverview.nextColloquium || 'None' }}</div>
                  <div class="stat-desc text-info">Book Now -></div>
              </div>
          </div>
          
          <!-- Recent Marks Preview -->
          <div class="card bg-base-100 shadow col-span-1 md:col-span-2">
              <div class="card-body">
                  <div class="flex justify-between">
                     <h2 class="card-title">Latest Marks</h2>
                     <button class="btn btn-link btn-xs" @click="$router.push('/parent/marks')">View All</button>
                  </div>
                  <MarksView :marks="parentStore.currentChildOverview.recentMarks" />
              </div>
          </div>

          <!-- Shortcuts -->
          <ul class="menu bg-base-100 w-full rounded-box shadow">
              <li class="menu-title">Actions</li>
              <li><router-link to="/parent/absences">Justify Absences</router-link></li>
              <li><router-link to="/parent/messages">Message Teachers</router-link></li>
              <li><router-link to="/parent/documents">Report Cards</router-link></li>
          </ul>
      </div>
  </div>
</template>
