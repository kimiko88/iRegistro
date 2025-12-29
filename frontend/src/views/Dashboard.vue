<script setup lang="ts">
import { onMounted, computed } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';
import { Shield, User, GraduationCap, School, FileText, Users } from 'lucide-vue-next';

const auth = useAuthStore();
const router = useRouter();

const userRole = computed(() => auth.user?.role || 'Guest');

const availableDashboards = [
  { role: ['SuperAdmin'], route: 'SuperAdminDashboard', label: 'Super Admin Dashboard', icon: Shield },
  { role: ['Admin'], route: 'UserManagement', label: 'Admin Panel', icon: Shield },
  { role: ['Teacher', 'Docente'], route: 'TeacherDashboard', label: 'Teacher Dashboard', icon: User },
  { role: ['Student', 'Studente'], route: 'StudentDashboard', label: 'Student Dashboard', icon: GraduationCap },
  { role: ['Parent', 'Genitore'], route: 'ParentDashboard', label: 'Parent Dashboard', icon: Users },
  { role: ['Secretary', 'Segreteria'], route: 'SecretaryInbox', label: 'Secretary Dashboard', icon: FileText },
  { role: ['Principal', 'Director', 'Dirigente'], route: 'DirectorDashboard', label: 'Director Dashboard', icon: School },
];

onMounted(() => {
  const role = auth.user?.role;
  if (role) {
    const target = availableDashboards.find(d => d.role.includes(role));
    if (target) {
      router.replace({ name: target.route });
    }
  }
});

const navigateTo = (routeName: string) => {
    router.push({ name: routeName });
};
</script>

<template>
  <div class="min-h-[80vh] flex flex-col items-center justify-center p-4">
    <div class="text-center mb-8">
       <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100 mb-2">Welcome, {{ auth.user?.name || 'User' }}</h1>
       <p class="text-gray-500">Select your dashboard to continue</p>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 max-w-4xl w-full">
         <button 
            v-for="dashboard in availableDashboards" 
            :key="dashboard.route"
            @click="navigateTo(dashboard.route)"
            class="flex flex-col items-center justify-center p-8 bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 hover:shadow-md hover:border-indigo-500 dark:hover:border-indigo-400 transition-all gap-4 group"
         >
            <div class="w-16 h-16 rounded-full bg-indigo-50 dark:bg-indigo-900/20 flex items-center justify-center text-indigo-600 dark:text-indigo-400 group-hover:scale-110 transition-transform">
                <component :is="dashboard.icon" class="w-8 h-8" />
            </div>
            <span class="font-semibold text-gray-700 dark:text-gray-200">{{ dashboard.label }}</span>
         </button>
    </div>
  </div>
</template>
