<template>
  <div class="p-6 space-y-8 bg-gray-50 min-h-screen dark:bg-gray-900">
    <div class="flex flex-col md:flex-row justify-between items-center gap-4">
      <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100 flex items-center gap-3">
        <LayoutDashboard class="w-8 h-8 text-indigo-600 dark:text-indigo-400" />
        SuperAdmin Dashboard
      </h1>
      <ActionButton 
        label="New School" 
        variant="primary" 
        @click="showCreateSchoolModal = true"
        :icon="Plus"
        class="bg-indigo-600 hover:bg-indigo-700 text-white border-none shadow-sm shadow-indigo-200 dark:shadow-none"
      />
    </div>

    <!-- KPI Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="stat bg-white dark:bg-gray-800 shadow-md rounded-2xl p-6 border border-gray-100 dark:border-gray-700">
        <div class="stat-figure text-indigo-600 dark:text-indigo-400">
          <div class="w-12 h-12 rounded-full bg-indigo-50 dark:bg-indigo-900/30 flex items-center justify-center">
             <School class="w-6 h-6" />
          </div>
        </div>
        <div class="stat-title font-medium text-gray-500">Total Schools</div>
        <div class="stat-value text-indigo-600 dark:text-indigo-400 text-4xl mt-1">{{ kpis?.schoolsCount || 0 }}</div>
        <div class="stat-desc text-gray-400 mt-1">Active across system</div>
      </div>
      
      <div class="stat bg-white dark:bg-gray-800 shadow-md rounded-2xl p-6 border border-gray-100 dark:border-gray-700">
        <div class="stat-figure text-purple-600 dark:text-purple-400">
          <div class="w-12 h-12 rounded-full bg-purple-50 dark:bg-purple-900/30 flex items-center justify-center">
             <Users class="w-6 h-6" />
          </div>
        </div>
        <div class="stat-title font-medium text-gray-500">Total Users</div>
        <div class="stat-value text-purple-600 dark:text-purple-400 text-4xl mt-1">{{ kpis?.usersCount || 0 }}</div>
        <div class="stat-desc text-green-500 mt-1 flex items-center gap-1">
          <TrendingUp class="w-3 h-3"/> {{ kpis?.activeUsersLast30Days || 0 }} active (30d)
        </div>
      </div>
      
      <div class="stat bg-white dark:bg-gray-800 shadow-md rounded-2xl p-6 border border-gray-100 dark:border-gray-700">
        <div class="stat-figure text-blue-600 dark:text-blue-400">
           <div class="w-12 h-12 rounded-full bg-blue-50 dark:bg-blue-900/30 flex items-center justify-center">
             <Database class="w-6 h-6" />
          </div>
        </div>
        <div class="stat-title font-medium text-gray-500">Storage Used</div>
        <div class="stat-value text-blue-600 dark:text-blue-400 text-4xl mt-1">{{ formatBytes(kpis?.storageUsed || 0) }}</div>
        <div class="stat-desc text-gray-400 mt-1">of 1TB Quota</div>
      </div>

       <div class="stat bg-white dark:bg-gray-800 shadow-md rounded-2xl p-6 border border-gray-100 dark:border-gray-700">
        <div class="stat-figure text-cyan-500">
           <div class="w-12 h-12 rounded-full bg-cyan-50 dark:bg-cyan-900/30 flex items-center justify-center">
             <Activity class="w-6 h-6" />
          </div>
        </div>
        <div class="stat-title font-medium text-gray-500">System Health</div>
        <div class="stat-value text-cyan-500 text-4xl mt-1">98%</div>
        <div class="stat-desc text-gray-400 mt-1">Operational</div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- School Map / Distribution -->
      <div class="lg:col-span-2 bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 h-[450px] flex flex-col">
        <h3 class="font-bold text-lg mb-4 text-gray-800 dark:text-gray-100 flex items-center gap-2">
            <Map class="w-5 h-5 text-gray-400"/>
            School Distribution
        </h3>
        <div class="flex-grow flex items-center justify-center bg-gray-50 dark:bg-gray-700/30 rounded-xl border border-dashed border-gray-200 dark:border-gray-600">
           <!-- Placeholder for Map -->
           <div class="text-center text-gray-500">
             <Globe class="w-16 h-16 mx-auto mb-3 opacity-30"/>
             <p class="font-medium">Interactive Map Component</p>
             <p class="text-sm mt-1">North: 45% • Center: 30% • South: 25%</p>
           </div>
        </div>
      </div>

      <!-- Quick Actions / Status -->
      <div class="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 flex flex-col">
         <h3 class="font-bold text-lg mb-6 text-gray-800 dark:text-gray-100 flex items-center gap-2">
            <Server class="w-5 h-5 text-gray-400"/>
            System Status
         </h3>
         <ul class="space-y-4 mb-8">
           <li class="flex justify-between items-center p-3 bg-gray-50 dark:bg-gray-700/30 rounded-lg">
             <span class="text-gray-600 dark:text-gray-300 font-medium">Database</span>
             <span class="badge badge-success gap-1 pl-1 pr-3 py-3 text-white"><CheckCircle class="w-3 h-3"/> Healthy</span>
           </li>
           <li class="flex justify-between items-center p-3 bg-gray-50 dark:bg-gray-700/30 rounded-lg">
             <span class="text-gray-600 dark:text-gray-300 font-medium">API Gateway</span>
             <span class="badge badge-success gap-1 pl-1 pr-3 py-3 text-white"><CheckCircle class="w-3 h-3"/> Healthy</span>
           </li>
           <li class="flex justify-between items-center p-3 bg-gray-50 dark:bg-gray-700/30 rounded-lg">
             <span class="text-gray-600 dark:text-gray-300 font-medium">Storage Service</span>
             <span class="badge badge-success gap-1 pl-1 pr-3 py-3 text-white"><CheckCircle class="w-3 h-3"/> Healthy</span>
           </li>
           <li class="flex justify-between items-center p-3 bg-gray-50 dark:bg-gray-700/30 rounded-lg">
             <span class="text-gray-600 dark:text-gray-300 font-medium">Backups</span>
             <span class="badge badge-warning gap-1 pl-1 pr-3 py-3 text-yellow-900"><Loader2 class="w-3 h-3 animate-spin"/> Checking...</span>
           </li>
         </ul>
         
         <div class="divider my-0"></div>
         
         <div class="mt-6">
             <h3 class="font-bold text-lg mb-4 text-gray-800 dark:text-gray-100 flex items-center gap-2">
                <FileText class="w-5 h-5 text-gray-400"/>
                Recent Audits
             </h3>
              <div class="text-sm space-y-3">
                 <div v-for="i in 3" :key="i" class="flex justify-between items-center text-gray-600 dark:text-gray-400 p-2 hover:bg-gray-50 dark:hover:bg-gray-700/30 rounded transition-colors">
                   <div class="flex items-center gap-2">
                       <User class="w-4 h-4 text-gray-400"/>
                       <span>User login (admin)</span>
                   </div>
                   <span class="text-xs font-mono bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">2m ago</span>
                 </div>
              </div>
              <div class="mt-6 text-center">
                <button class="btn btn-sm btn-ghost w-full">View All Logs</button>
              </div>
         </div>
      </div>
    </div>

    <!-- School List -->
    <div class="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700">
      <h3 class="font-bold text-lg mb-6 text-gray-800 dark:text-gray-100 flex items-center gap-2">
          <Building2 class="w-5 h-5 text-gray-400"/>
          Managed Schools
      </h3>
      <div class="overflow-hidden rounded-xl border border-gray-200 dark:border-gray-700">
          <DataTable
            :columns="schoolColumns"
            :data="schools"
            :loading="loading"
            :totalItems="100" 
            @search="fetchSchools"
          >
            <template #item-actions="{ item }">
              <div class="flex gap-2">
                 <button class="btn btn-xs btn-ghost text-indigo-600 hover:bg-indigo-50"><Eye class="w-3 h-3 mr-1"/> View</button>
                 <button class="btn btn-xs btn-ghost text-amber-600 hover:bg-amber-50"><Edit2 class="w-3 h-3 mr-1"/> Edit</button>
              </div>
            </template>
          </DataTable>
      </div>
    </div>

    <!-- Create School Modal -->
    <FormModal
      :isOpen="showCreateSchoolModal"
      title="Create New School"
      @close="showCreateSchoolModal = false"
      @submit="handleCreateSchool"
    >
      <div class="form-control w-full space-y-4">
        <div>
            <label class="label"><span class="label-text font-medium">School Name</span></label>
            <input type="text" class="input input-bordered w-full" v-model="newSchool.name" required placeholder="e.g. Liceo Scientifico Galilei" />
        </div>
        <div>
            <label class="label"><span class="label-text font-medium">Region</span></label>
            <select class="select select-bordered w-full" v-model="newSchool.region">
              <option value="" disabled selected>Select Region</option>
              <option value="Lombardia">Lombardia</option>
              <option value="Lazio">Lazio</option>
              <option value="Campania">Campania</option>
            </select>
        </div>
      </div>
    </FormModal>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAdminStore } from '@/stores/admin';
import { storeToRefs } from 'pinia';
import ActionButton from '@/components/shared/ActionButton.vue';
import DataTable from '@/components/shared/DataTable.vue';
import FormModal from '@/components/shared/FormModal.vue';
import { 
  School, 
  Users, 
  Database, 
  Activity, 
  Map, 
  Plus, 
  LayoutDashboard,
  TrendingUp,
  Globe,
  Server,
  CheckCircle,
  Loader2,
  FileText,
  User,
  Building2,
  Eye,
  Edit2
} from 'lucide-vue-next';

const adminStore = useAdminStore();
const { kpis, schools, loading } = storeToRefs(adminStore);
const { fetchKPIs, fetchSchools } = adminStore;

const showCreateSchoolModal = ref(false);
const newSchool = ref({ name: '', region: '' });

const schoolColumns = [
  { key: 'name', label: 'Name', sortable: true },
  { key: 'region', label: 'Region', sortable: true },
  { key: 'students', label: 'Students', sortable: true },
  { key: 'status', label: 'Status' },
];

const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const handleCreateSchool = async () => {
  // adminStore.createSchool(newSchool.value);
  showCreateSchoolModal.value = false;
};

onMounted(() => {
  fetchKPIs();
  fetchSchools();
});
</script>
