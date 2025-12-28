<template>
  <div class="p-6 space-y-8 bg-gray-50 min-h-screen dark:bg-gray-900">
    <div class="flex flex-col md:flex-row justify-between items-center gap-4">
      <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100 flex items-center gap-3">
        <LayoutDashboard class="w-8 h-8 text-indigo-600 dark:text-indigo-400" />
        Parent Dashboard
      </h1>
      <!-- Child Selector as a more prominent visual element if needed, currently component handles it -->
    </div>

    <!-- Child Selector -->
    <div v-if="store.children.length > 1" class="bg-white dark:bg-gray-800 p-4 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700">
        <StudentSelect 
          :children="store.children" 
          :selectedId="store.selectedChildId"
          @select="store.selectChild"
        />
    </div>

    <div v-if="store.currentChildOverview" class="space-y-8">
      <!-- Main Student Card -->
      <StudentCard :student="store.currentChildOverview" class="shadow-sm" />

      <!-- KPI Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div class="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 relative overflow-hidden group hover:shadow-md transition-all">
          <div class="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
            <GraduationCap class="w-16 h-16 text-indigo-600" />
          </div>
          <h3 class="text-gray-500 text-sm font-medium uppercase tracking-wider mb-2">GPA</h3>
          <p class="text-3xl font-bold text-indigo-600 dark:text-indigo-400">{{ store.currentChildOverview.averageGrade?.toFixed(1) || '-' }}</p>
          <p class="text-xs text-gray-400 mt-2">Overall Average</p>
        </div>

        <div class="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 relative overflow-hidden group hover:shadow-md transition-all">
          <div class="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
             <UserX class="w-16 h-16 text-red-500" />
          </div>
          <h3 class="text-gray-500 text-sm font-medium uppercase tracking-wider mb-2">Absences (Month)</h3>
          <p class="text-3xl font-bold text-red-500">{{ store.currentChildOverview.absencesMonth || 0 }}</p>
          <p class="text-xs text-gray-400 mt-2">Days missed this month</p>
        </div>

        <div class="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 relative overflow-hidden group hover:shadow-md transition-all">
          <div class="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
            <Percent class="w-16 h-16 text-green-500" />
          </div>
          <h3 class="text-gray-500 text-sm font-medium uppercase tracking-wider mb-2">Attendance %</h3>
          <p class="text-3xl font-bold text-green-500">{{ store.currentChildOverview.attendancePercentage || 0 }}%</p>
          <p class="text-xs text-gray-400 mt-2">Year to date</p>
        </div>

        <div class="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 relative overflow-hidden group hover:shadow-md transition-all">
          <div class="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
            <MessageSquare class="w-16 h-16 text-yellow-500" />
          </div>
          <h3 class="text-gray-500 text-sm font-medium uppercase tracking-wider mb-2">Next Colloquium</h3>
          <p class="text-lg font-bold text-yellow-600 dark:text-yellow-400 truncate">{{ store.currentChildOverview.nextColloquium || 'None Scheduled' }}</p>
          <p class="text-xs text-gray-400 mt-2">Upcoming meeting</p>
        </div>
      </div>

      <!-- Quick Links -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <router-link :to="{ name: 'MarksView' }" 
          class="flex flex-col items-center justify-center p-6 bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 hover:border-indigo-200 hover:shadow-md hover:-translate-y-1 transition-all group">
          <div class="w-12 h-12 rounded-full bg-indigo-50 dark:bg-indigo-900/30 flex items-center justify-center mb-3 group-hover:bg-indigo-100 dark:group-hover:bg-indigo-900/50 transition-colors">
            <FileBarChart class="w-6 h-6 text-indigo-600 dark:text-indigo-400" />
          </div>
          <span class="font-semibold text-gray-700 dark:text-gray-200">View Grades</span>
        </router-link>

        <router-link :to="{ name: 'AbsencesView' }" 
          class="flex flex-col items-center justify-center p-6 bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 hover:border-red-200 hover:shadow-md hover:-translate-y-1 transition-all group">
          <div class="w-12 h-12 rounded-full bg-red-50 dark:bg-red-900/30 flex items-center justify-center mb-3 group-hover:bg-red-100 dark:group-hover:bg-red-900/50 transition-colors">
            <CalendarX class="w-6 h-6 text-red-500" />
          </div>
          <span class="font-semibold text-gray-700 dark:text-gray-200">View Absences</span>
        </router-link>

        <router-link :to="{ name: 'MessagesView' }" 
          class="flex flex-col items-center justify-center p-6 bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 hover:border-green-200 hover:shadow-md hover:-translate-y-1 transition-all group">
          <div class="w-12 h-12 rounded-full bg-green-50 dark:bg-green-900/30 flex items-center justify-center mb-3 group-hover:bg-green-100 dark:group-hover:bg-green-900/50 transition-colors">
            <Mail class="w-6 h-6 text-green-500" />
          </div>
          <span class="font-semibold text-gray-700 dark:text-gray-200">Messages</span>
        </router-link>

         <router-link :to="{ name: 'ColloquiumView' }" 
          class="flex flex-col items-center justify-center p-6 bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 hover:border-yellow-200 hover:shadow-md hover:-translate-y-1 transition-all group">
          <div class="w-12 h-12 rounded-full bg-yellow-50 dark:bg-yellow-900/30 flex items-center justify-center mb-3 group-hover:bg-yellow-100 dark:group-hover:bg-yellow-900/50 transition-colors">
            <Users class="w-6 h-6 text-yellow-500" />
          </div>
          <span class="font-semibold text-gray-700 dark:text-gray-200">Colloquiums</span>
        </router-link>
      </div>

      <!-- Recent Marks Preview -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6">
        <div class="flex items-center gap-2 mb-6">
           <TrendingUp class="w-6 h-6 text-gray-400" />
           <h3 class="text-xl font-bold text-gray-800 dark:text-gray-100">Performance Trend</h3>
        </div>
        <MarksChart :marks="store.currentChildOverview.latestMarks || []" class="max-h-64" />
      </div>

    </div>
    <div v-else class="flex flex-col items-center justify-center py-20 text-gray-500">
      <div class="loading loading-spinner loading-lg text-primary mb-4"></div>
      <p>Loading child data...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useParentStore } from '@/stores/parent';
import StudentSelect from './StudentSelect.vue';
import StudentCard from '@/components/StudentCard.vue';
import MarksChart from '@/components/MarksChart.vue';
import { 
  LayoutDashboard, 
  GraduationCap, 
  UserX, 
  Percent, 
  MessageSquare, 
  FileBarChart, 
  CalendarX, 
  Mail, 
  Users,
  TrendingUp 
} from 'lucide-vue-next';

const store = useParentStore();

onMounted(() => {
  store.fetchChildren();
});
</script>
