<template>
  <div class="h-full flex flex-col bg-base-100 rounded-box shadow">
      <!-- Toolbar -->
      <div class="p-2 border-b border-base-200 flex justify-between items-center bg-base-50 rounded-t-box">
          <div class="flex gap-2 items-center">
              <button class="btn btn-sm btn-ghost" @click="prevWeek">‹</button>
              <span class="font-bold text-sm min-w-[150px] text-center">
                  {{ weekStartStr }} - {{ weekEndStr }}
              </span>
              <button class="btn btn-sm btn-ghost" @click="nextWeek">›</button>
              <button class="btn btn-sm btn-ghost" @click="today">Today</button>
          </div>
          <div class="flex gap-2">
             <button class="btn btn-sm btn-outline btn-success" @click="markAllPresent">
                  Mark All Present
             </button>
          </div>
      </div>

      <!-- Grid -->
      <div class="flex-1 overflow-auto">
          <table class="table table-pin-rows table-pin-cols w-full">
              <thead>
                  <tr>
                      <th class="bg-base-200 z-20 w-12 text-center">#</th>
                      <th class="bg-base-200 z-20 w-48">Student</th>
                      <th v-for="day in weekDays" :key="day.dateStr" class="bg-base-200 min-w-[6rem] text-center">
                          <div class="flex flex-col items-center">
                             <span class="text-xs font-bold">{{ day.name }}</span>
                             <span class="text-[10px] opacity-60">{{ day.shortDate }}</span>
                          </div>
                      </th>
                      <th class="bg-base-200 text-center min-w-[4rem]">Total</th>
                  </tr>
              </thead>
              <tbody>
                  <tr v-for="(student, idx) in students" :key="student.id" class="hover">
                      <th class="bg-base-100 text-center text-xs opacity-50">{{ idx + 1 }}</th>
                      <td class="bg-base-100 font-medium whitespace-nowrap">
                          {{ student.last_name }} {{ student.first_name }}
                      </td>
                      
                      <td v-for="day in weekDays" :key="day.dateStr" class="p-1 border border-base-200 text-center">
                           <!-- Simple Select for Absence Type -->
                           <!-- In real app, maybe a popover or detailed modal -->
                           <select 
                             class="select select-xs select-ghost w-full max-w-[80px] h-8 text-center" 
                             :class="getAbsenceClass(student.id, day.dateStr)"
                             :value="getAbsenceType(student.id, day.dateStr)"
                             @change="setAbsence(student.id, day.dateStr, ($event.target as HTMLSelectElement).value)"
                           >
                              <option value="">●</option>
                              <option value="ABSENT">A</option>
                              <option value="LATE">L</option>
                              <option value="EXCUSED">E</option>
                           </select>
                      </td>
                      
                      <td class="text-center font-bold text-base-content/70">
                          {{ countAbsences(student.id) }}
                      </td>
                  </tr>
              </tbody>
          </table>
      </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

const props = defineProps<{
    students: any[],
    absences: any[]
}>();

const emit = defineEmits(['save']);

const currentDate = ref(new Date());

const getStartOfWeek = (date: Date) => {
    const d = new Date(date);
    const day = d.getDay();
    const diff = d.getDate() - day + (day === 0 ? -6 : 1); // Adjust when day is sunday
    return new Date(d.setDate(diff));
};

const weekStart = computed(() => getStartOfWeek(currentDate.value));

const weekDays = computed(() => {
    const start = weekStart.value;
    const days = [];
    for (let i = 0; i < 6; i++) { // Mon-Sat
        const d = new Date(start);
        d.setDate(start.getDate() + i);
        days.push({
            dateStr: d.toISOString().split('T')[0],
            name: d.toLocaleDateString('en-US', { weekday: 'short' }),
            shortDate: d.toLocaleDateString(undefined, { day: 'numeric', month: 'numeric' })
        });
    }
    return days;
});

const weekStartStr = computed(() => weekStart.value.toLocaleDateString());
const weekEndStr = computed(() => {
    const d = new Date(weekStart.value);
    d.setDate(d.getDate() + 5);
    return d.toLocaleDateString();
});

const prevWeek = () => {
    const d = new Date(currentDate.value);
    d.setDate(d.getDate() - 7);
    currentDate.value = d;
};
const nextWeek = () => {
    const d = new Date(currentDate.value);
    d.setDate(d.getDate() + 7);
    currentDate.value = d;
};
const today = () => { currentDate.value = new Date(); };

const getAbsence = (sid: number, dateStr: string) => {
    // Check props.absences
    // Absence date might be ISO string with time
    return props.absences.find(a => a.student_id === sid && a.date.startsWith(dateStr));
};

const getAbsenceType = (sid: number, dateStr: string) => {
    const abs = getAbsence(sid, dateStr);
    return abs ? abs.type : '';
};

const getAbsenceClass = (sid: number, dateStr: string) => {
    const type = getAbsenceType(sid, dateStr);
    switch(type) {
        case 'ABSENT': return 'text-error font-bold';
        case 'LATE': return 'text-warning font-bold';
        case 'EXCUSED': return 'text-success font-bold';
        default: return 'text-base-content/20'; // Present (Dot)
    }
};

const setAbsence = (sid: number, dateStr: string, type: string) => {
    emit('save', {
        student_id: sid,
        date: dateStr, // In real app, maybe specific hour? For now Full Day
        type: type || null // null implies delete or present
    });
};

const countAbsences = (sid: number) => {
    return props.absences.filter(a => a.student_id === sid && a.type === 'ABSENT').length;
};

const markAllPresent = () => {
    // For current week/day? Usually confirms "Mark all students present for Today"
    // Implementation left as exercise or confirmation modal
};
</script>
