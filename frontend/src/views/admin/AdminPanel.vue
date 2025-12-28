<template>
  <div class="flex h-full">
    <!-- Admin Sidebar (if not using main layout sidebar, or as a secondary sidebar) -->
    <!-- Assuming this is used inside the main Layout, maybe just a wrapper for router-view with some top tabs or breadcrumbs? -->
    <!-- Or if this IS the main text for /admin routes, it needs a sidebar. -->
    <!-- Let's assume it renders the sub-navigation for Admin actions -->
    
    <div class="w-64 bg-base-200 h-full p-4 hidden lg:block">
      <ul class="menu bg-base-100 w-56 rounded-box">
        <li class="menu-title">Admin Console</li>
        <li><router-link to="/superadmin">Dashboard</router-link></li>
        <li><router-link to="/admin/users">User Management</router-link></li>
        <li><router-link to="/admin/settings">School Settings</router-link></li>
        
        <li class="menu-title mt-4">Advanced</li>
        <li><router-link to="/admin/backups">Backups</router-link></li>
        <li><router-link to="/admin/export">Data Export</router-link></li>
        <li><router-link to="/admin/audit-logs">Audit Logs</router-link></li>
      </ul>
    </div>

    <div class="flex-1 overflow-y-auto">
      <router-view></router-view>
    </div>
  </div>
</template>

<script setup lang="ts">
// Logic for active menu highlighting is handled by router-link-active class usually
</script>
