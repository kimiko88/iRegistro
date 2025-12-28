<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';

const auth = useAuthStore();
const router = useRouter();

const logout = () => {
  auth.logout();
  router.push('/login');
};
</script>

<template>
  <div class="drawer lg:drawer-open">
    <input id="my-drawer-2" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content flex flex-col bg-base-200">
      <!-- Navbar -->
      <div class="w-full navbar bg-base-100 lg:hidden">
        <div class="flex-none lg:hidden">
          <label for="my-drawer-2" aria-label="open sidebar" class="btn btn-square btn-ghost">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-6 h-6 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
          </label>
        </div>
        <div class="flex-1 px-2 mx-2">iRegistro</div>
      </div>
      
      <!-- Content -->
      <main class="p-6">
        <router-view />
      </main>
    </div> 
    
    <!-- Sidebar -->
    <div class="drawer-side z-20">
      <label for="my-drawer-2" aria-label="close sidebar" class="drawer-overlay"></label> 
      <ul class="menu p-4 w-80 min-h-full bg-base-100 text-base-content flex flex-col justify-between">
        <div>
          <li class="mb-4 text-2xl font-bold px-4">iRegistro</li>
          <li><router-link to="/dashboard">Dashboard</router-link></li>
          
          <!-- Admin Links (Should check role) -->
          <li v-if="auth.user?.role === 'SuperAdmin'"><router-link to="/superadmin">Manage Schools</router-link></li>
          <li v-if="auth.user?.role === 'Admin'">
             <details open>
                <summary>Administration</summary>
                <ul>
                    <li><router-link to="/admin/users">Users</router-link></li>
                    <li><router-link to="/admin/settings">Settings</router-link></li>
                </ul>
             </details>
          </li>
          
          <li v-if="auth.user?.role === 'Teacher'"><router-link to="/teacher">Registro Docente</router-link></li>
          <li v-if="auth.user?.role === 'Parent'"><router-link to="/parent">My Children</router-link></li>
          <li v-if="auth.user?.role === 'Student'"><router-link to="/student">My Dashboard</router-link></li>
          <li v-if="auth.user?.role === 'Secretary'">
             <details open>
                <summary>Secretary</summary>
                <ul>
                    <li><router-link to="/secretary">Inbox</router-link></li>
                    <li><router-link to="/secretary/archive">Archive</router-link></li>
                    <li><router-link to="/secretary/classes">Classes</router-link></li>
                </ul>
             </details>
          </li>
          <li v-if="auth.user?.role === 'Principal'"><router-link to="/director">School Overview</router-link></li>

          <li v-if="auth.isAuthenticated"><a @click="logout">Logout</a></li>
        </div>
        <div class="text-xs text-center opacity-50">v0.1.0</div>
      </ul>
    </div>
  </div>
</template>
