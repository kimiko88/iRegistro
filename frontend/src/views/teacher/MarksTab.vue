<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useTeacherStore } from '@/stores/teacher';
import MarksGrid from '@/components/teacher/MarksGrid.vue';
import Modal from '@/components/shared/Modal.vue';
import { useUIStore } from '@/stores/ui';

const teacherStore = useTeacherStore();
const uiStore = useUIStore();

const isModalOpen = ref(false);
const editingMark = ref<any>({});

onMounted(() => {
    teacherStore.fetchMarks(1); // Subject ID irrelevant for mock
});

const openAddMarkModal = (studentId: number) => {
    editingMark.value = { studentId, value: null, type: 'Oral', date: new Date().toISOString().split('T')[0] };
    isModalOpen.value = true;
};

const openEditMarkModal = (mark: any) => {
    editingMark.value = { ...mark };
    isModalOpen.value = true;
};

const saveMark = async () => {
    if (!editingMark.value.value) return;
    
    await teacherStore.saveMark(editingMark.value);
    uiStore.addNotification({ type: 'success', message: 'Mark saved successfully' });
    isModalOpen.value = false;
};
</script>

<template>
  <div class="space-y-4">
      <div class="flex justify-between items-center">
          <h3 class="text-xl font-bold">Gradebook - {{ teacherStore.selectedClass?.subject }}</h3>
          <div class="space-x-2">
              <button class="btn btn-sm btn-outline">Analytics</button>
              <button class="btn btn-sm btn-outline">Import CSV</button>
          </div>
      </div>

      <MarksGrid 
          :students="teacherStore.students" 
          :marks="teacherStore.marks"
          @add-mark="openAddMarkModal"
          @edit-mark="openEditMarkModal"
      />

      <Modal :title="editingMark.id ? 'Edit Mark' : 'Add Mark'" :isOpen="isModalOpen" @close="isModalOpen = false">
           <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
               <div class="form-control">
                   <label class="label">Value</label>
                   <input type="number" step="0.5" max="10" min="1" class="input input-bordered" v-model.number="editingMark.value" />
               </div>
               <div class="form-control">
                   <label class="label">Date</label>
                   <input type="date" class="input input-bordered" v-model="editingMark.date" />
               </div>
               <div class="form-control">
                   <label class="label">Type</label>
                   <select class="select select-bordered" v-model="editingMark.type">
                       <option>Oral</option>
                       <option>Written</option>
                       <option>Practical</option>
                   </select>
               </div>
           </div>
           
           <template #actions>
               <button class="btn" @click="isModalOpen = false">Cancel</button>
               <button class="btn btn-primary" @click="saveMark">Save</button>
           </template>
      </Modal>
  </div>
</template>
