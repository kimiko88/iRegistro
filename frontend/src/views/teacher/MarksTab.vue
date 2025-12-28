<template>
  <div class="h-full flex flex-col gap-4">
      <!-- Toolbar -->
      <div class="flex justify-between items-center bg-base-100 p-4 rounded-box shadow-sm">
          <div class="flex gap-4 items-center">
             <div class="stats stats-horizontal shadow compact">
                <div class="stat px-4 py-2">
                    <div class="stat-title text-xs">Class Average</div>
                    <div class="stat-value text-lg text-primary">{{ classAverage }}</div>
                </div>
                 <div class="stat px-4 py-2">
                    <div class="stat-title text-xs">Total Marks</div>
                    <div class="stat-value text-lg">{{ store.marks.length }}</div>
                </div>
             </div>
          </div>

          <button class="btn btn-primary gap-2" @click="openAddMarkModal">
              <span>+</span> Add Mark
          </button>
      </div>

      <!-- Grid -->
      <MarksGrid 
        :students="store.students" 
        :marks="store.marks"
        @save="handleBulkSave"
        class="flex-1"
      />

      <!-- Add Mark Modal stub -->
      <!-- Ideally we reuse Modal.vue -->
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTeacherStore } from '@/stores/teacher';
import MarksGrid from '@/components/teacher/MarksGrid.vue';
import { useUIStore } from '@/stores/ui';

const store = useTeacherStore();
const ui = useUIStore();

const classAverage = computed(() => {
    if (!store.marks.length) return '-';
    const sum = store.marks.reduce((acc, m) => acc + m.value, 0);
    return (sum / store.marks.length).toFixed(2);
});

const openAddMarkModal = () => {
    // TODO: proper modal
    // Creating a dummy mark for now to test reactivity
    if(store.students.length) {
        store.saveMark({
             student_id: store.students[0].id,
             subject_id: store.currentSubjectId || 1,
             teacher_id: 3, // Current user
             class_id: store.selectedClassId,
             value: 8.5,
             type: 'NUMERIC',
             date: new Date().toISOString()
        });
    }
};

const handleBulkSave = async (changes: any[]) => {
    // Process bulk changes
    // In real app, verify backend has bulk endpoint or loop
    // Demo: Loop
    for (const ch of changes) {
        // ch.key is "date-type" string from Grid
        const [date, type] = ch.key.split('-'); // Simple parsing
        await store.saveMark({
            student_id: ch.student_id,
            subject_id: store.currentSubjectId || 1,
            teacher_id: 3,
            class_id: store.selectedClassId,
            value: ch.value,
            type: type || 'NUMERIC',
            date: date || new Date().toISOString()
        });
    }
};
</script>
