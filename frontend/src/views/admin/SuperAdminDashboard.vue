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
        @click="openCreateModal"
        :icon="Plus"
        class="bg-indigo-600 hover:bg-indigo-700 text-white border-none shadow-sm shadow-indigo-200 dark:shadow-none"
      />
    </div>

    <!-- KPI Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="stat bg-white dark:bg-gray-800 shadow-md rounded-2xl p-6 border border-gray-100 dark:border-gray-700 hover:shadow-lg transition-shadow">
        <div class="stat-figure text-indigo-600 dark:text-indigo-400">
          <div class="w-12 h-12 rounded-full bg-indigo-50 dark:bg-indigo-900/30 flex items-center justify-center">
             <School class="w-6 h-6" />
          </div>
        </div>
        <div class="stat-title font-medium text-gray-500">Total Schools</div>
        <div class="stat-value text-indigo-600 dark:text-indigo-400 text-4xl mt-1">{{ kpis?.schoolsCount || 0 }}</div>
        <div class="stat-desc text-gray-400 mt-1">Active across system</div>
      </div>
      
      <div class="stat bg-white dark:bg-gray-800 shadow-md rounded-2xl p-6 border border-gray-100 dark:border-gray-700 hover:shadow-lg transition-shadow">
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
      
      <div class="stat bg-white dark:bg-gray-800 shadow-md rounded-2xl p-6 border border-gray-100 dark:border-gray-700 hover:shadow-lg transition-shadow">
        <div class="stat-figure text-blue-600 dark:text-blue-400">
           <div class="w-12 h-12 rounded-full bg-blue-50 dark:bg-blue-900/30 flex items-center justify-center">
             <Database class="w-6 h-6" />
          </div>
        </div>
        <div class="stat-title font-medium text-gray-500">Storage Used</div>
        <div class="stat-value text-blue-600 dark:text-blue-400 text-4xl mt-1">{{ formatBytes(kpis?.storageUsed || 0) }}</div>
        <div class="stat-desc text-gray-400 mt-1">of 1TB Quota</div>
      </div>

       <div class="stat bg-white dark:bg-gray-800 shadow-md rounded-2xl p-6 border border-gray-100 dark:border-gray-700 hover:shadow-lg transition-shadow">
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

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- School Map / Distribution -->
      <div class="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 h-[400px] flex flex-col">
        <h3 class="font-bold text-lg mb-4 text-gray-800 dark:text-gray-100 flex items-center gap-2">
            <Map class="w-5 h-5 text-gray-400"/>
            School Distribution by Region
        </h3>
        <div class="flex-grow">
           <BarChart :data="distributionData" />
        </div>
      </div>

       <!-- User Growth -->
      <div class="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 h-[400px] flex flex-col">
        <h3 class="font-bold text-lg mb-4 text-gray-800 dark:text-gray-100 flex items-center gap-2">
            <TrendingUp class="w-5 h-5 text-gray-400"/>
            User Growth (Last 6 Months)
        </h3>
        <div class="flex-grow">
           <LineChart :data="userGrowthData" />
        </div>
      </div>
    </div>

     <!-- Quick Actions / Status -->
      <div class="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 flex flex-col">
         <h3 class="font-bold text-lg mb-6 text-gray-800 dark:text-gray-100 flex items-center gap-2">
            <Server class="w-5 h-5 text-gray-400"/>
            System Status Details
         </h3>
         
         <div class="mt-6">
             <h3 class="font-bold text-lg mb-4 text-gray-800 dark:text-gray-100 flex items-center gap-2">
                <FileText class="w-5 h-5 text-gray-400"/>
                Recent Audits
             </h3>
              <div class="text-sm space-y-3">
                 <div v-for="log in auditLogs || []" :key="log.id" class="flex justify-between items-center text-gray-600 dark:text-gray-400 p-2 hover:bg-gray-50 dark:hover:bg-gray-700/30 rounded transition-colors">
                   <div class="flex items-center gap-2">
                       <User class="w-4 h-4 text-gray-400"/>
                       <span>{{ log.action }} ({{ log.entity }})</span>
                   </div>
                   <span class="text-xs font-mono bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">{{ new Date(log.timestamp).toLocaleTimeString() }}</span>
                 </div>
                 <div v-if="!auditLogs || auditLogs.length === 0" class="text-gray-400 italic text-center text-xs">No recent activity</div>
              </div>
              <div class="mt-6 text-center">
                <button class="btn btn-sm btn-ghost w-full">View All Logs</button>
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
            @search="handleSearch"
          >
            <template #item-actions="{ item }">
              <div class="flex gap-2">
                 <button @click="handleViewSchool(item)" class="btn btn-xs btn-ghost text-indigo-600 hover:bg-indigo-50"><Eye class="w-3 h-3 mr-1"/> View</button>
                 <button @click="openEditModal(item)" class="btn btn-xs btn-ghost text-amber-600 hover:bg-amber-50"><Edit2 class="w-3 h-3 mr-1"/> Edit</button>
                 <button @click="manageUsers(item)" class="btn btn-xs btn-ghost text-green-600 hover:bg-green-50"><Users class="w-3 h-3 mr-1"/> Users</button>
              </div>
            </template>
          </DataTable>
      </div>
    </div>

    <!-- Create/Edit School Modal -->
    <FormModal
      :isOpen="showCreateSchoolModal"
      :title="isEditing ? 'Edit School' : 'Create New School'"
      @close="showCreateSchoolModal = false"
      @submit="handleSaveSchool"
    >
      <div class="form-control w-full space-y-4">
        
        <!-- Search Input with Dropdown (Only for Create) -->
         <div v-if="!isEditing" class="relative">
            <label class="label"><span class="label-text font-medium">Search Italian School</span></label>
            <input 
              type="text" 
              class="input input-bordered w-full" 
              v-model="searchQuery" 
              placeholder="Start typing... (e.g. Galilei)" 
              @input="showSuggestions = true"
            />
            
             <ul v-if="showSuggestions && filteredSchools.length > 0" class="absolute z-10 menu p-2 shadow bg-base-100 rounded-box w-full mt-1 max-h-48 overflow-y-auto">
                <li v-for="school in filteredSchools" :key="school.code">
                    <a @click="selectSchool(school)">
                         <span class="font-bold">{{ school.name }}</span>
                         <span class="text-xs text-gray-500">{{ school.city }} ({{ school.code }})</span>
                    </a>
                </li>
             </ul>
        </div>

        <div class="divider text-xs text-gray-400" v-if="!isEditing">OR ENTER MANUALLY</div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
             <div>
                <label class="label"><span class="label-text font-medium">School Name</span></label>
                <input type="text" class="input input-bordered w-full" v-model="formSchool.name" required />
            </div>
             <div>
                <label class="label"><span class="label-text font-medium">Mechanicographic Code</span></label>
                <input type="text" class="input input-bordered w-full" v-model="formSchool.code" required placeholder="e.g. MIPS01000N" />
            </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="label"><span class="label-text font-medium">City</span></label>
                <input type="text" class="input input-bordered w-full" v-model="formSchool.city" />
            </div>
            <div>
                <label class="label"><span class="label-text font-medium">Region</span></label>
                <div class="relative">
                    <input 
                        type="text" 
                        class="input input-bordered w-full pr-10" 
                        v-model="formSchool.region"
                        placeholder="Select or Type Region"
                        @focus="showRegionSuggestions = true"
                        @input="showRegionSuggestions = true"
                    />
                    <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none text-gray-400">
                        <ChevronDown class="w-4 h-4"/>
                    </div>
                    <ul v-if="showRegionSuggestions" class="absolute z-20 menu p-2 shadow bg-base-100 rounded-box w-full mt-1 max-h-48 overflow-y-auto">
                        <li v-for="region in filteredRegions" :key="region">
                            <a @click="selectRegion(region)" class="font-medium">{{ region }}</a>
                        </li>
                        <li v-if="filteredRegions.length === 0" class="text-gray-400 p-2 text-sm text-center">No matches found</li>
                    </ul>
                </div>
            </div>
        </div>

        <div>
            <label class="label"><span class="label-text font-medium">Address</span></label>
            <input type="text" class="input input-bordered w-full" v-model="formSchool.address" />
        </div>
         <div>
            <label class="label"><span class="label-text font-medium">Principal Email</span></label>
            <input type="email" class="input input-bordered w-full" v-model="formSchool.email" />
        </div>
      </div>
    </FormModal>

    <!-- View School Modal -->
    <dialog id="view_school_modal" class="modal" :class="{ 'modal-open': showViewModal }">
      <div class="modal-box relative bg-white dark:bg-gray-800 max-w-2xl">
        <button @click="showViewModal = false" class="btn btn-sm btn-circle absolute right-2 top-2">✕</button>
        
        <div class="flex items-center gap-4 mb-6">
            <div class="w-16 h-16 rounded-full bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center text-indigo-600 dark:text-indigo-400">
                <School class="w-8 h-8" />
            </div>
            <div>
              <h3 class="text-xl font-bold">{{ schoolToView?.name }}</h3>
              <p class="text-gray-500">{{ schoolToView?.code }}</p>
            </div>
        </div>

        <div v-if="schoolToView" class="grid grid-cols-1 md:grid-cols-2 gap-4">
             <div>
                 <span class="block text-sm font-medium text-gray-500">Region</span>
                 <span>{{ schoolToView.region }}</span>
             </div>
             <div>
                 <span class="block text-sm font-medium text-gray-500">City</span>
                 <span>{{ schoolToView.city }}</span>
             </div>
             <div class="col-span-1 md:col-span-2">
                 <span class="block text-sm font-medium text-gray-500">Address</span>
                 <span>{{ schoolToView.address }}</span>
             </div>
              <div class="col-span-1 md:col-span-2">
                 <span class="block text-sm font-medium text-gray-500">Principal Email</span>
                 <span>{{ schoolToView.email }}</span>
             </div>
        </div>
      </div>
    </dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useAdminStore } from '@/stores/admin';
import { useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import ActionButton from '@/components/shared/ActionButton.vue';
import DataTable from '@/components/shared/DataTable.vue';
import FormModal from '@/components/shared/FormModal.vue';
import BarChart from '@/components/charts/BarChart.vue';
import LineChart from '@/components/charts/LineChart.vue';
import { 
  School, Users, Database, Activity, Map, Plus, LayoutDashboard, TrendingUp,
  Globe, Server, CheckCircle, Loader2, FileText, User, Building2, Eye, Edit2,
  ChevronDown
} from 'lucide-vue-next';
import italianSchools from '@/assets/schools_it.json';
import italianRegions from '@/assets/regions_it.json';

const adminStore = useAdminStore();
const router = useRouter(); 
const { kpis, schools, loading, auditLogs } = storeToRefs(adminStore);
const { fetchKPIs, fetchSchools, fetchAuditLogs } = adminStore;

const showCreateSchoolModal = ref(false);
const showViewModal = ref(false);
const isEditing = ref(false);
const editingId = ref<string | null>(null);
const schoolToView = ref<any>(null);

// Mock Data for Charts
const distributionData = ref({
  labels: ['North', 'Center', 'South', 'Islands'],
  datasets: [{
    label: 'Schools by Region',
    data: [12, 19, 8, 5],
    backgroundColor: ['rgba(99, 102, 241, 0.5)', 'rgba(168, 85, 247, 0.5)', 'rgba(59, 130, 246, 0.5)', 'rgba(236, 72, 153, 0.5)'],
    borderColor: ['rgb(99, 102, 241)', 'rgb(168, 85, 247)', 'rgb(59, 130, 246)', 'rgb(236, 72, 153)'],
    borderWidth: 1
  }]
});

const userGrowthData = ref({
  labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
  datasets: [{
    label: 'New Users',
    data: [65, 59, 80, 81, 56, 120],
    fill: true,
    borderColor: 'rgb(75, 192, 192)',
    backgroundColor: 'rgba(75, 192, 192, 0.2)',
    tension: 0.4
  }]
});

// Form Data with new fields
const formSchool = ref({
  name: '',
  region: '',
  city: '',
  address: '',
  code: '',
  email: ''
});

// Search State (Schools)
const searchQuery = ref('');
const showSuggestions = ref(false);

const filteredSchools = computed(() => {
  if (!searchQuery.value || searchQuery.value.length < 3) return [];
  const q = searchQuery.value.toLowerCase();
  return italianSchools.filter(s => 
    s.name.toLowerCase().includes(q) || 
    s.city.toLowerCase().includes(q) ||
    s.code.toLowerCase().includes(q)
  );
});

// Search State (Regions)
const showRegionSuggestions = ref(false);
const filteredRegions = computed(() => {
  if (!formSchool.value.region) return italianRegions;
  const q = formSchool.value.region.toLowerCase();
  return italianRegions.filter(r => r.toLowerCase().includes(q));
});

const selectSchool = (school: any) => {
  formSchool.value = {
    name: school.name,
    region: school.region,
    city: school.city,
    address: school.address,
    code: school.code,
    email: school.email
  };
  searchQuery.value = school.name;
  showSuggestions.value = false;
};

const selectRegion = (region: string) => {
    formSchool.value.region = region;
    showRegionSuggestions.value = false;
}

// Reset Form
const resetForm = () => {
    formSchool.value = { name: '', region: '', city: '', address: '', code: '', email: '' };
    searchQuery.value = '';
    isEditing.value = false;
    editingId.value = null;
    showRegionSuggestions.value = false;
}

const openCreateModal = () => {
    resetForm();
    showCreateSchoolModal.value = true;
}

const openEditModal = (item: any) => {
    formSchool.value = { ...item }; // Map fields
    isEditing.value = true;
    editingId.value = item.id;
    showCreateSchoolModal.value = true;
};

const manageUsers = (item: any) => {
    router.push({ 
        name: 'SchoolUserManagement', 
        params: { schoolId: item.id },
        query: { schoolName: item.name }
    });
};

const handleSaveSchool = async () => {
  try {
    if (isEditing.value && editingId.value) {
        await adminStore.updateSchool(editingId.value, formSchool.value);
    } else {
        await adminStore.createSchool(formSchool.value);
    }
    showCreateSchoolModal.value = false;
    resetForm();
    await fetchSchools(); 
  } catch(e) {
    console.error(e);
  }
};

const handleViewSchool = (item: any) => {
    schoolToView.value = item;
    showViewModal.value = true;
};

const handleSearch = (q: string) => {
    fetchSchools({ q });
};

const schoolColumns = [
  { key: 'name', label: 'Name', sortable: true },
  { key: 'city', label: 'City', sortable: true },
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

// Macro-regions mapping
const getMacroRegion = (region: string) => {
    const north = ['Lombardia', 'Piemonte', 'Veneto', 'Trentino-Alto Adige', 'Friuli Venezia Giulia', 'Liguria', 'Emilia-Romagna', "Valle d'Aosta"];
    const center = ['Toscana', 'Lazio', 'Marche', 'Umbria'];
    // All others South/Islands
    if (north.includes(region)) return 'North';
    if (center.includes(region)) return 'Center';
    return 'South';
};

const distributionStats = computed(() => {
    if (!schools.value || schools.value.length === 0) return { north: 0, center: 0, south: 0 };
    
    const initial = { north: 0, center: 0, south: 0 };
    const counts = schools.value.reduce((acc: any, school: any) => {
        const macro = getMacroRegion(school.region || '');
        if (macro === 'North') acc.north++;
        else if (macro === 'Center') acc.center++;
        else acc.south++;
        return acc;
    }, initial);

    const total = schools.value.length;
    return {
        north: Math.round((counts.north / total) * 100),
        center: Math.round((counts.center / total) * 100),
        south: Math.round((counts.south / total) * 100)
    };
});

onMounted(() => {
  fetchKPIs();
  fetchSchools();
  fetchAuditLogs(); 
});
</script>
