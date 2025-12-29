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

      <FormModal
        :isOpen="showAddModal"
        title="Add New Mark"
        :loading="saving"
        @close="showAddModal = false"
        @submit="handleSaveMark"
      >
        <div class="grid grid-cols-1 gap-4">
            <div class="form-control">
                <label class="label"><span class="label-text">Student</span></label>
                <select class="select select-bordered" v-model="markForm.student_id" required>
                    <option value="" disabled>Select Student</option>
                    <option v-for="student in store.students" :key="student.id" :value="student.id">
                        {{ student.lastName }} {{ student.firstName }}
                    </option>
                </select>
            </div>
            
            <div class="grid grid-cols-2 gap-4">
                <div class="form-control">
                    <label class="label"><span class="label-text">Mark (1-10)</span></label>
                    <input type="number" step="0.25" min="1" max="10" class="input input-bordered" v-model="markForm.value" required />
                </div>
                <div class="form-control">
                    <label class="label"><span class="label-text">Type</span></label>
                    <select class="select select-bordered" v-model="markForm.type">
                        <option value="ORAL">Oral</option>
                        <option value="WRITTEN">Written</option>
                        <option value="PRACTICAL">Practical</option>
                    </select>
                </div>
            </div>

            <div class="form-control">
                <label class="label"><span class="label-text">Date</span></label>
                <input type="date" class="input input-bordered" v-model="markForm.date" required />
            </div>

             <div class="form-control">
                <label class="label"><span class="label-text">Notes</span></label>
                <textarea class="textarea textarea-bordered" v-model="markForm.description"></textarea>
            </div>
        </div>
      </FormModal>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, reactive } from 'vue';
import { useTeacherStore } from '@/stores/teacher';
import MarksGrid from '@/components/teacher/MarksGrid.vue';
import FormModal from '@/components/shared/FormModal.vue';
import { useNotificationStore } from '@/stores/notification';

const store = useTeacherStore();
const notificationStore = useNotificationStore();

const classAverage = computed(() => {
    if (!store.marks.length) return '-';
    const sum = store.marks.reduce((acc, m) => acc + m.value, 0);
    return (sum / store.marks.length).toFixed(2);
});

const showAddModal = ref(false);
const saving = ref(false);
const markForm = reactive({
    student_id: '',
    value: null as number | null,
    type: 'ORAL',
    date: new Date().toISOString().split('T')[0],
    description: ''
});

const openAddMarkModal = () => {
    markForm.student_id = '';
    markForm.value = null;
    markForm.type = 'ORAL';
    markForm.date = new Date().toISOString().split('T')[0];
    markForm.description = '';
    showAddModal.value = true;
};

const handleSaveMark = async () => {
    if (!markForm.student_id || !markForm.value) {
        notificationStore.error('Please select a student and enter a mark');
        return;
    }
    
    saving.value = true;
    try {
        await store.saveMark({
             student_id: Number(markForm.student_id),
             subject_id: store.currentSubjectId || 1, // Fallback needs to be addressed if 0
             teacher_id: 3, // TODO: Get from auth store
             class_id: store.selectedClassId,
             value: Number(markForm.value),
             type: markForm.type,
             date: new Date(markForm.date).toISOString()
        });
        notificationStore.success('Mark added successfully');
        showAddModal.value = false;
    } catch (err) {
        notificationStore.error('Failed to save mark');
    } finally {
        saving.value = false;
    }
};

const handleBulkSave = async (changes: any[]) => {
    for (const ch of changes) {
        const [date, type] = ch.key.split('-'); 
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
    notificationStore.success('Marks saved successfully');
};
</script>
