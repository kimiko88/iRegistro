<template>
  <div class="card bg-base-100 shadow-xl h-full flex flex-col">
    <div class="card-body p-4 gap-4 h-full">
      <!-- Toolbar -->
      <div class="flex justify-between items-center">
        <h3 class="card-title text-sm opacity-70">Registro Voti</h3>
        <button class="btn btn-primary btn-sm gap-2" @click="openAddMarkModal">
          <span class="i-heroicons-plus-circle w-4 h-4" />
          Nuovo Voto
        </button>
      </div>

      <!-- Scrollable Grid Container -->
      <div class="flex-1 overflow-auto border rounded-box relative">
        <table class="table table-pin-rows table-pin-cols w-full">
          <thead>
            <tr>
              <th class="bg-base-200 z-20 w-48">Studente</th>
              <!-- Render existing marks columns or dynamic dates? 
                   Usually teachers see a grid of dates. For simplicity, we list marks per student row -->
              <th class="bg-base-200">Voti (Media)</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="student in store.students" :key="student.id" class="hover">
              <td class="font-medium bg-base-100 z-10">{{ student.last_name }} {{ student.first_name }}</td>
              <td class="flex flex-wrap gap-2 items-center">
                <!-- Marks Loop -->
                <div 
                  v-for="mark in getStudentMarks(student.id)" 
                  :key="mark.id" 
                  class="badge badge-lg tooltip cursor-pointer hover:scale-110 transition-transform"
                  :class="getMarkColor(mark.value)"
                  :data-tip="`${mark.date} - ${mark.type}`"
                >
                  {{ mark.value.toFixed(1) }}
                </div>
                
                <!-- Average Badge -->
                <div class="divider divider-horizontal mx-1"></div>
                <div class="font-bold text-base-content/70">
                  Avg: {{ calculateAverage(student.id).toFixed(1) }}
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add Mark Modal -->
    <dialog id="add_mark_modal" class="modal" ref="markModal">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">Aggiungi Voto</h3>
        
        <div class="form-control w-full mb-2">
          <label class="label"><span class="label-text">Studente</span></label>
          <select class="select select-bordered" v-model="newMark.student_id">
            <option v-for="s in store.students" :key="s.id" :value="s.id">
              {{ s.last_name }} {{ s.first_name }}
            </option>
          </select>
        </div>

        <div class="grid grid-cols-2 gap-4 mb-2">
          <div class="form-control">
            <label class="label"><span class="label-text">Voto</span></label>
            <input type="number" step="0.25" min="1" max="10" class="input input-bordered" v-model.number="newMark.value" />
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text">Data</span></label>
            <input type="date" class="input input-bordered" v-model="newMark.date" />
          </div>
        </div>

        <div class="form-control mb-4">
          <label class="label"><span class="label-text">Tipo</span></label>
          <select class="select select-bordered" v-model="newMark.type">
            <option value="Written">Scritto</option>
            <option value="Oral">Orale</option>
            <option value="Practical">Pratico</option>
          </select>
        </div>

        <div class="modal-action">
          <form method="dialog">
            <button class="btn btn-ghost">Annulla</button>
            <button class="btn btn-primary ml-2" @click.prevent="saveMark">Salva</button>
          </form>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useTeacherStore } from '@/stores/teacher';

const store = useTeacherStore();
const markModal = ref<HTMLDialogElement | null>(null);

const newMark = reactive({
  student_id: null as number | null,
  value: 6.0,
  date: new Date().toISOString().split('T')[0],
  type: 'Oral',
  subject_id: null as number | null // Need to inject or use store getter
});

const getStudentMarks = (studentId: number) => {
  return store.marks.filter((m: any) => m.student_id === studentId);
};

const calculateAverage = (studentId: number) => {
  const marks = getStudentMarks(studentId);
  if (marks.length === 0) return 0;
  const sum = marks.reduce((acc: number, m: any) => acc + m.value, 0);
  return sum / marks.length;
};

const getMarkColor = (value: number) => {
  if (value >= 6) return 'badge-success text-white';
  if (value >= 5) return 'badge-warning text-white';
  return 'badge-error text-white';
};

const openAddMarkModal = () => {
  if (markModal.value) {
    markModal.value.showModal();
    // Default to first student if none selected
    if (!newMark.student_id && store.students.length > 0) {
      newMark.student_id = store.students[0].id;
    }
  }
};

const saveMark = async () => {
  const selectedClass = store.classes.find(c => c.id === store.selectedClassId);
  if (!selectedClass) return;

  await store.saveMark({
    ...newMark,
    class_id: store.selectedClassId,
    subject_id: selectedClass.subjectId
  });
  
  if (markModal.value) markModal.value.close();
};
</script>
