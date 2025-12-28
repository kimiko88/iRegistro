<template>
  <div class="h-full flex flex-col gap-4">
      <AbsenceCalendar
        :students="store.students"
        :absences="store.absences"
        @save="handleSave"
        class="flex-1"
      />
  </div>
</template>

<script setup lang="ts">
import { onMounted, watch } from 'vue';
import { useTeacherStore } from '@/stores/teacher';
import AbsenceCalendar from '@/components/teacher/AbsenceCalendar.vue';

const props = defineProps<{
    classId: number
}>();

const store = useTeacherStore();

const loadData = () => {
    if (props.classId) {
        store.fetchAbsences(props.classId);
    }
};

onMounted(() => {
    loadData();
});

watch(() => props.classId, () => {
    loadData();
});

const handleSave = async (absenceData: any) => {
    // absenceData: { student_id, date, type }
    // If type is empty/null, it might mean delete. Backend needs to handle this or front sends "delete" action.
    // For MVP, assuming saveAbsence handles Upsert, and if type is null it deletes? 
    // Let's assume sending type="" means remove absence (make present).
    
    // We need ClassID
    await store.saveAbsence({
        ...absenceData,
        class_id: props.classId,
        hour: 0 // Full day default
    });
    // saveAbsence pushes to store, but for delete we might need to reload or filter locally
    // Ideally reload or smarter store update
    if (!absenceData.type) {
        // Reload list to clear removed one
         store.fetchAbsences(props.classId);
    }
};
</script>
