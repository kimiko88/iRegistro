<template>
  <div class="card bg-base-100 shadow-xl h-full flex flex-col">
    <div class="card-body p-4 gap-4 h-full">
      <div class="flex justify-between items-center">
        <h3 class="card-title text-sm opacity-70">Registro Assenze</h3>
        <div class="join">
          <button class="btn btn-sm join-item" @click="shiftDate(-1)">❮</button>
          <input type="date" class="input input-sm join-item input-bordered" v-model="currentDate" @change="fetchAbsences" />
          <button class="btn btn-sm join-item" @click="shiftDate(1)">❯</button>
        </div>
      </div>

      <div class="overflow-auto border rounded-box flex-1">
        <table class="table w-full">
          <thead>
            <tr>
              <th>Studente</th>
              <th>Stato Presenza</th>
              <th>Note</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="student in store.students" :key="student.id">
              <td class="font-bold">{{ student.last_name }} {{ student.first_name }}</td>
              <td>
                <div class="join">
                  <input 
                    class="join-item btn btn-sm " 
                    :class="getBtnClass(student.id, 'PRESENT')"
                    type="radio" 
                    :name="`status-${student.id}`" 
                    aria-label="P"
                    :checked="getAbsenceType(student.id) === 'PRESENT'"
                    @click="setAbsence(student.id, 'PRESENT')"
                  />
                  <input 
                    class="join-item btn btn-sm" 
                    :class="getBtnClass(student.id, 'ABSENT')"
                    type="radio" 
                    :name="`status-${student.id}`" 
                    aria-label="A"
                    :checked="getAbsenceType(student.id) === 'ABSENT'"
                    @click="setAbsence(student.id, 'ABSENT')"
                  />
                  <input 
                    class="join-item btn btn-sm" 
                    :class="getBtnClass(student.id, 'LATE')"
                    type="radio" 
                    :name="`status-${student.id}`" 
                    aria-label="R"
                    :checked="getAbsenceType(student.id) === 'LATE'"
                    @click="setAbsence(student.id, 'LATE')"
                  />
                </div>
              </td>
              <td>
                <input 
                  type="text" 
                  class="input input-xs input-bordered w-full max-w-xs" 
                  placeholder="Note..." 
                  :value="getNote(student.id)"
                  @change="updateNote(student.id, ($event.target as HTMLInputElement).value)"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="card-actions justify-end mt-2">
        <button class="btn btn-primary" @click="saveDailyAbsences">Salva Registro</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useTeacherStore } from '@/stores/teacher';
import teacherApi from '@/services/teacher'; // Direct API for specific calls or extend store

const store = useTeacherStore();
const currentDate = ref(new Date().toISOString().split('T')[0]);

// Local state for the day's absences before save
// Map<studentId, { type: string, note: string }>
const dailyAbsences = ref<Map<number, any>>(new Map());

const fetchAbsences = async () => {
  if (!store.selectedClassId) return;
  // This would typically fetch from API. 
  // For now, simulate or use store if extended.
  try {
    const res = await teacherApi.getAbsences(store.selectedClassId, currentDate.value);
    // Populate map
    dailyAbsences.value.clear();
    res.data.forEach((a: any) => {
      dailyAbsences.value.set(a.student_id, { type: a.type, note: a.note, id: a.id });
    });
  } catch (e) {
    dailyAbsences.value.clear();
  }
};

watch(() => store.selectedClassId, fetchAbsences);
onMounted(fetchAbsences);

const shiftDate = (days: number) => {
  const date = new Date(currentDate.value);
  date.setDate(date.getDate() + days);
  currentDate.value = date.toISOString().split('T')[0];
  fetchAbsences();
};

const getAbsenceType = (sid: number) => {
  return dailyAbsences.value.get(sid)?.type || 'PRESENT';
};

const getNote = (sid: number) => {
  return dailyAbsences.value.get(sid)?.note || '';
};

const setAbsence = (sid: number, type: string) => {
  const current = dailyAbsences.value.get(sid) || {};
  dailyAbsences.value.set(sid, { ...current, type });
};

const updateNote = (sid: number, note: string) => {
  const current = dailyAbsences.value.get(sid) || { type: 'PRESENT' };
  dailyAbsences.value.set(sid, { ...current, note });
};

const getBtnClass = (sid: number, type: string) => {
  const currentType = getAbsenceType(sid);
  if (currentType === type) {
    if (type === 'PRESENT') return 'btn-success text-white';
    if (type === 'ABSENT') return 'btn-error text-white';
    if (type === 'LATE') return 'btn-warning text-white';
  }
  return 'btn-ghost';
};

const saveDailyAbsences = async () => {
  if (!store.selectedClassId) return;
  
  const payload = Array.from(dailyAbsences.value.entries()).map(([sid, data]) => ({
    student_id: sid,
    date: currentDate.value,
    type: data.type === 'PRESENT' ? null : data.type, // Assuming null/delete for present, or specific API logic
    note: data.note
  })).filter(a => a.type !== null || a.note); // Only send actual events

  await teacherApi.saveAbsences(store.selectedClassId, payload);
  // Toast success
};
</script>
