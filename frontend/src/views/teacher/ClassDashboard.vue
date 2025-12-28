<template>
  <div class="h-full flex flex-col">
    <!-- Header with Class Selector -->
    <div class="navbar bg-base-100 shadow-sm min-h-16 px-4 z-10">
      <div class="flex-1 gap-4 items-center">
        <h2 class="text-xl font-bold text-primary">Le Mie Classi</h2>
        
        <select 
          class="select select-bordered select-sm w-full max-w-xs"
          :value="store.selectedClassId"
          @change="changeClass($event)"
          :disabled="store.classes.length === 0"
        >
          <option disabled value="">Seleziona una classe</option>
          <option v-for="cls in store.classes" :key="cls.id" :value="cls.id">
            Classe {{ cls.name }} - {{ cls.subject }}
          </option>
        </select>

        <div v-if="store.currentClass" class="badge badge-outline">
          {{ store.students.length }} Studenti
        </div>
      </div>

      <div class="flex-none">
        <ul class="menu menu-horizontal px-1 gap-1">
          <li><a :class="{ active: activeTab === 'marks' }" @click="activeTab = 'marks'">Voti</a></li>
          <li><a :class="{ active: activeTab === 'absences' }" @click="activeTab = 'absences'">Assenze</a></li>
          <li><a :class="{ active: activeTab === 'schedule' }" @click="activeTab = 'schedule'">Orario</a></li>
          <li><a :class="{ active: activeTab === 'colloquiums' }" @click="activeTab = 'colloquiums'">Colloqui</a></li>
          <li><a :class="{ active: activeTab === 'documents' }" @click="activeTab = 'documents'">Documenti</a></li>
        </ul>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 overflow-hidden relative bg-base-200/50 p-4">
      <Transition name="fade" mode="out-in">
        
        <!-- Tab: Marks -->
        <component 
          :is="currentTabComponent" 
          v-if="store.selectedClassId"
          :class-id="store.selectedClassId"
        />
        
        <!-- Empty State -->
        <div v-else class="flex flex-col items-center justify-center h-full text-base-content/50">
          <div class="i-heroicons-academic-cap text-6xl mb-4" />
          <h3 class="text-lg font-medium">Seleziona una classe per iniziare</h3>
        </div>

      </Transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useTeacherStore } from '@/stores/teacher';
import MarksTab from '@/components/teacher/MarksTab.vue';
import AbsencesTab from '@/components/teacher/AbsencesTab.vue';
// Placeholder setup for other tabs
import ScheduleTab from '@/components/teacher/ScheduleTab.vue';
import DocumentsTab from '@/components/teacher/DocumentsTab.vue';

const store = useTeacherStore();
const activeTab = ref('marks');

const currentTabComponent = computed(() => {
  switch (activeTab.value) {
    case 'marks': return MarksTab;
    case 'absences': return AbsencesTab;
    case 'schedule': return ScheduleTab;
    case 'colloquiums': return ColloquiumView;
    case 'documents': return DocumentsTab;
    default: return MarksTab;
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
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
