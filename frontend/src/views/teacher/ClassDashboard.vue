<template>
  <div class="h-full flex flex-col bg-gray-50 dark:bg-gray-900">
    <!-- Header with Class Selector -->
    <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4 shadow-sm">
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div class="flex items-center gap-4">
          <div class="flex flex-col">
              <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-100 flex items-center gap-2">
                  <GraduationCap class="w-7 h-7 text-indigo-600 dark:text-indigo-400"/>
                  Class Dashboard
              </h1>
              <span class="text-sm text-gray-500 dark:text-gray-400 mt-1" v-if="store.currentClass">
                  {{ store.currentClass.name }} - {{ store.currentClass.subjectName }}
              </span>
          </div>
        </div>
        
        <!-- Class Selector & Student Count -->
        <div class="flex items-center gap-3">
          <select 
            class="select select-bordered select-sm w-full md:w-64 bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600"
            :value="store.selectedClassId"
            @change="changeClass($event)"
            :disabled="store.classes.length === 0"
          >
            <option disabled value="">Select Class</option>
            <option v-for="cls in store.classes" :key="cls.id" :value="cls.id">
              {{ cls.name }} ({{ cls.subjectName }})
            </option>
          </select>
          
          <div v-if="store.currentClass" class="badge badge-lg bg-indigo-50 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300 border-indigo-200 dark:border-indigo-700 gap-2 px-4">
              <Users class="w-4 h-4"/>
              {{ store.students.length }} Students
          </div>
        </div>
      </div>
    </div>
    
    <!-- Tab Navigation -->
    <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 overflow-x-auto">
      <div role="tablist" class="tabs tabs-bordered py-2">
        <a 
          v-for="tab in tabs" 
          :key="tab.id"
          role="tab" 
          class="tab gap-2 transition-all duration-200 text-gray-600 dark:text-gray-400 hover:text-indigo-600 dark:hover:text-indigo-400" 
          :class="{ 'tab-active !text-indigo-600 dark:!text-indigo-400 !border-indigo-600 dark:!border-indigo-400 font-semibold': activeTab === tab.id }"
          @click="activeTab = tab.id"
        >
            <component :is="tab.icon" class="w-4 h-4" />
            {{ tab.label }}
        </a>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 overflow-hidden relative">
      <div v-if="!store.selectedClassId" class="flex flex-col items-center justify-center h-full">
          <div class="text-center p-8 bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-200 dark:border-gray-700 max-w-md">
              <School class="w-20 h-20 mx-auto mb-4 text-gray-300 dark:text-gray-600"/>
              <p class="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-2">No Class Selected</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">Please select a class from the dropdown above to begin</p>
          </div>
      </div>

      <Transition name="fade" mode="out-in" v-else>
        <component 
          :is="currentTabComponent" 
          class="h-full w-full p-6 overflow-y-auto bg-gray-50 dark:bg-gray-900"
          :class-id="store.selectedClassId"
        />
      </Transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useTeacherStore } from '@/stores/teacher';
import MarksTab from './MarksTab.vue';
import AbsencesTab from './AbsencesTab.vue';
import ScheduleTab from './ScheduleTab.vue';
import MessagesTab from './MessagesTab.vue';
import ColloquiumTab from './ColloquiumTab.vue';
import DocumentsTab from './DocumentsTab.vue';
import CoordinatorTab from './CoordinatorTab.vue';
import { 
  GraduationCap, 
  Users, 
  School,
  FileBarChart,
  CalendarX,
  Calendar,
  MessageSquare,
  FileText,
  UserCog
} from 'lucide-vue-next';

const store = useTeacherStore();
const activeTab = ref('marks');

const tabs = computed(() => {
    const list = [
        { id: 'marks', label: 'Marks', icon: FileBarChart },
        { id: 'absences', label: 'Absences', icon: CalendarX },
        { id: 'schedule', label: 'Schedule', icon: Calendar },
        { id: 'messages', label: 'Messages', icon: MessageSquare },
        { id: 'colloquiums', label: 'Colloquiums', icon: Users },
        { id: 'documents', label: 'Documents', icon: FileText },
    ];
    if (store.isCoordinator) {
        list.push({ id: 'coordinator', label: 'Coordinator', icon: UserCog });
    }
    return list;
});

const currentTabComponent = computed(() => {
  switch (activeTab.value) {
    case 'marks': return MarksTab;
    case 'absences': return AbsencesTab;
    case 'schedule': return ScheduleTab;
    case 'messages': return MessagesTab;
    case 'colloquiums': return ColloquiumTab;
    case 'documents': return DocumentsTab;
    case 'coordinator': return CoordinatorTab;
    default: return { template: `<div class="p-10 text-center capitalize">${activeTab.value} Tab Coming Soon</div>` };
  }
});

onMounted(() => {
  store.fetchClasses();
});

const changeClass = (event: Event) => {
  const target = event.target as HTMLSelectElement;
  if (target.value) {
    store.selectClass(Number(target.value));
  }
};
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
