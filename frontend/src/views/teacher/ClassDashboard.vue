<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useTeacherStore } from '@/stores/teacher';
import MarksTab from './MarksTab.vue';
import AbsencesTab from './AbsencesTab.vue';
import ScheduleTab from './ScheduleTab.vue';
import ColloquiumTab from './ColloquiumTab.vue';
import MessagesTab from './MessagesTab.vue';
import DocumentsTab from './DocumentsTab.vue';
// import CoordinatorDashboard from './CoordinatorDashboard.vue';

const teacherStore = useTeacherStore();
const activeTab = ref('marks');

onMounted(() => {
    teacherStore.fetchClasses();
    if(teacherStore.selectedClassId) {
        teacherStore.fetchStudents(teacherStore.selectedClassId);
    }
});

const currentTabComponent = computed(() => {
    switch(activeTab.value) {
        case 'marks': return MarksTab;
        case 'absences': return AbsencesTab;
        case 'schedule': return ScheduleTab;
        case 'colloquium': return ColloquiumTab;
        case 'messages': return MessagesTab;
        case 'documents': return DocumentsTab;
        default: return MarksTab;
    }
});
</script>

<template>
  <div class="p-4 space-y-4">
    <!-- Header: Class Selector -->
    <div class="flex flex-col md:flex-row justify-between items-center gap-4 bg-base-100 p-4 rounded-lg shadow">
        <div class="prose">
            <h2 class="m-0">Class Dashboard</h2>
        </div>
        
        <div class="flex items-center gap-2">
            <span class="font-bold">Active Class:</span>
            <select class="select select-bordered select-sm w-full max-w-xs" 
                    v-model="teacherStore.selectedClassId" @change="teacherStore.selectClass($event.target.value)">
                <option v-for="cls in teacherStore.classes" :key="cls.id" :value="cls.id">
                    {{ cls.name }} - {{ cls.subject }}
                </option>
            </select>
        </div>
    </div>

    <!-- Navigation Tabs -->
    <div role="tablist" class="tabs tabs-boxed">
        <a role="tab" class="tab" :class="{ 'tab-active': activeTab === 'marks' }" @click="activeTab = 'marks'">Marks</a>
        <a role="tab" class="tab" :class="{ 'tab-active': activeTab === 'absences' }" @click="activeTab = 'absences'">Absences</a>
        <a role="tab" class="tab" :class="{ 'tab-active': activeTab === 'schedule' }" @click="activeTab = 'schedule'">Schedule</a>
        <a role="tab" class="tab" :class="{ 'tab-active': activeTab === 'colloquium' }" @click="activeTab = 'colloquium'">Colloquiums</a>
        <a role="tab" class="tab" :class="{ 'tab-active': activeTab === 'messages' }" @click="activeTab = 'messages'">Messages</a>
        <a role="tab" class="tab" :class="{ 'tab-active': activeTab === 'documents' }" @click="activeTab = 'documents'">Documents</a>
    </div>

    <!-- Content Area -->
    <div class="bg-base-100 rounded-lg shadow min-h-[500px] p-4">
        <component :is="currentTabComponent" />
    </div>
  </div>
</template>
