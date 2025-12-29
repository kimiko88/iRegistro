<script setup lang="ts">
import { ref, computed } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';
import { 
  LayoutDashboard, School, Users, Settings, FileText, Archive, 
  GraduationCap, BookOpen, LogOut, Menu, ChevronLeft, ChevronRight,
  Shield, UserCog, Home, Wallet
} from 'lucide-vue-next';

const auth = useAuthStore();
const router = useRouter();
const isSidebarOpen = ref(true);

const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value;
};

const logout = () => {
  auth.logout();
  router.push('/login');
};

const menuItems = computed(() => {
  const role = auth.user?.role;
  const items = [
    { label: 'Dashboard', to: '/dashboard', icon: LayoutDashboard, roles: ['*'] },
    
    // SuperAdmin
    { label: 'Manage Schools', to: '/superadmin', icon: School, roles: ['SuperAdmin'] },
    
    // Admin
    { label: 'Users', to: '/admin/users', icon: Users, roles: ['Admin'] },
    { label: 'Settings', to: '/admin/settings', icon: Settings, roles: ['Admin'] },
    { label: 'Audit Logs', to: '/admin/audit-logs', icon: Shield, roles: ['Admin'] },

    // Teacher
    { label: 'Registro Docente', to: '/teacher', icon: BookOpen, roles: ['Teacher'] },
    
    // Parent
    { label: 'My Children', to: '/parent', icon: Users, roles: ['Parent'] },
    
    // Student
    { label: 'My Dashboard', to: '/student', icon: GraduationCap, roles: ['Student'] },
    
    // Secretary
    { label: 'Inbox', to: '/secretary', icon: FileText, roles: ['Secretary'] },
    { label: 'Archive', to: '/secretary/archive', icon: Archive, roles: ['Secretary'] },
    { label: 'Classes', to: '/secretary/classes', icon: BookOpen, roles: ['Secretary'] },
    { label: 'Users', to: '/secretary/users', icon: Users, roles: ['Secretary'] },

    // Director
    { label: 'School Overview', to: '/director', icon: Home, roles: ['Principal'] },
  ];

  return items.filter(item => item.roles.includes('*') || (role && item.roles.includes(role)));
});
</script>

<template>
  <div class="drawer lg:drawer-open">
    <input id="my-drawer-2" type="checkbox" class="drawer-toggle" />
    
    <div class="drawer-content flex flex-col bg-gray-50 dark:bg-gray-900 min-h-screen transition-all duration-300" 
         :class="{ 'lg:ml-20': !isSidebarOpen, 'lg:ml-0': isSidebarOpen }">
      
      <!-- Navbar (Always Visible) -->
      <div class="navbar bg-white dark:bg-gray-800 shadow-sm sticky top-0 z-30 px-4 h-16 border-b border-gray-200 dark:border-gray-700">
        <div class="flex-none">
          <!-- Mobile Toggle -->
          <label for="my-drawer-2" class="btn btn-square btn-ghost lg:hidden text-gray-600 dark:text-gray-300">
            <Menu class="w-6 h-6" />
          </label>
          
          <!-- Desktop Toggle -->
          <button @click="toggleSidebar" class="btn btn-square btn-ghost hidden lg:inline-flex text-gray-600 dark:text-gray-300">
            <Menu class="w-6 h-6" />
          </button>
        </div>
        
        <div class="flex-1 px-4 items-center gap-4">
           <h1 class="text-xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent hidden md:block">
             iRegistro
           </h1>
           <div v-if="auth.user?.school" class="hidden md:flex items-center gap-2 text-sm text-gray-500 bg-gray-100 dark:bg-gray-700 px-3 py-1 rounded-full">
                <School class="w-4 h-4" />
                <span class="font-medium">{{ auth.user.school.name }}</span>
           </div>
        </div>
        
        <div class="flex-none gap-4">
           <div class="text-right hidden md:block mr-2">
              <p class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ auth.user?.name || 'User' }}</p>
              <p class="text-xs text-gray-500 capitalize">{{ auth.user?.role }}</p>
           </div>
           
           <div class="dropdown dropdown-end">
              <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar placeholder">
                <div class="bg-indigo-100 dark:bg-indigo-900 text-indigo-700 dark:text-indigo-300 rounded-full w-10">
                  <span class="text-lg font-bold">{{ auth.user?.name?.charAt(0) || 'U' }}</span>
                </div>
              </div>
              <ul tabindex="0" class="mt-3 z-[1] p-2 shadow menu menu-sm dropdown-content bg-base-100 rounded-box w-52">
                <li><a>Profile</a></li>
                <li><a>Settings</a></li>
                <li><a @click="logout" class="text-red-600">Logout</a></li>
              </ul>
           </div>
        </div>
      </div>
      
      <!-- Main Content -->
      <main class="p-6">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div> 
    
    <!-- Sidebar -->
    <div class="drawer-side z-40">
      <label for="my-drawer-2" aria-label="close sidebar" class="drawer-overlay"></label> 
      <aside 
        class="bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 min-h-screen transition-all duration-300 flex flex-col"
        :class="[isSidebarOpen ? 'w-64' : 'w-20 hidden lg:flex']"
      >
        <!-- Sidebar Header -->
        <div class="h-16 flex items-center justify-center border-b border-gray-200 dark:border-gray-700 px-4">
           <div v-if="isSidebarOpen" class="flex items-center gap-2 font-bold text-2xl text-indigo-600">
             <BookOpen class="w-8 h-8" />
             <span>iRegistro</span>
           </div>
           <BookOpen v-else class="w-8 h-8 text-indigo-600" />
        </div>

        <!-- Menu Items -->
        <ul class="menu p-2 flex-grow gap-1 overflow-y-auto w-full">
           <li v-for="item in menuItems" :key="item.to">
             <router-link 
               :to="item.to" 
               active-class="active-link"
               class="flex items-center gap-3 px-4 py-3 rounded-lg text-gray-600 dark:text-gray-400 hover:bg-indigo-50 dark:hover:bg-indigo-900/20 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
               :class="{ 'justify-center': !isSidebarOpen }"
             >
                <component :is="item.icon" class="w-5 h-5 flex-shrink-0" />
                <span v-if="isSidebarOpen" class="font-medium whitespace-nowrap">{{ item.label }}</span>
             </router-link>
           </li>
        </ul>

        <!-- Footer / Logout -->
        <div class="p-2 border-t border-gray-200 dark:border-gray-700">
           <button 
             @click="logout" 
             class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
             :class="{ 'justify-center': !isSidebarOpen }"
           >
              <LogOut class="w-5 h-5 flex-shrink-0" />
              <span v-if="isSidebarOpen" class="font-medium whitespace-nowrap">Logout</span>
           </button>
           
           <div v-if="isSidebarOpen" class="text-xs text-center text-gray-400 mt-4 pb-2">
             v1.0.0
           </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.active-link {
  @apply bg-indigo-600 text-white hover:bg-indigo-700 hover:text-white dark:hover:text-white !important;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
