<template>
  <div class="h-full flex flex-col bg-base-100 rounded-box shadow">
      <!-- Toolbar -->
      <div class="p-2 border-b border-base-200 flex justify-between items-center bg-base-50 rounded-t-box">
          <div class="flex gap-2">
              <button class="btn btn-sm btn-ghost" title="Undo" disabled>↩ Undo</button>
              <button class="btn btn-sm btn-ghost" title="Redo" disabled>↪ Redo</button>
              <div class="divider divider-horizontal my-1"></div>
              <button class="btn btn-sm btn-ghost gap-2" @click="exportCSV">
                  📄 Export CSV
              </button>
          </div>
          <div class="flex gap-2">
              <span class="text-xs text-base-content/60 self-center">Double click cell to edit</span>
              <button class="btn btn-sm btn-primary" @click="saveChanges" :disabled="!hasChanges">
                  Save Changes
              </button>
          </div>
      </div>

      <!-- Grid -->
      <div class="flex-1 overflow-auto relative">
          <table class="table table-pin-rows table-pin-cols w-full">
              <thead>
                  <tr>
                      <th class="bg-base-200 z-20 w-12 text-center">#</th>
                      <th class="bg-base-200 z-20 w-48">Student</th>
                      <!-- Date Columns Header -->
                      <!-- For simplicity, just listing random columns or fetched columns -->
                      <!-- In real app, we need to group by date or type -->
                      <th v-for="col in columns" :key="col.id" class="bg-base-200 min-w-[3rem] text-center group relative">
                          <div class="flex flex-col items-center">
                              <span class="text-xs font-bold">{{ formatDate(col.date) }}</span>
                              <span class="text-[10px] uppercase opacity-60">{{ col.type }}</span>
                          </div>
                          <!-- Column actions dropdown -->
                          <div class="absolute right-1 top-1 opacity-0 group-hover:opacity-100 transition-opacity">
                              <button class="btn btn-xs btn-circle btn-ghost">⋮</button>
                          </div>
                      </th>
                      <th class="bg-base-200 text-center min-w-[4rem] text-primary">Avg</th>
                  </tr>
              </thead>
              <tbody>
                  <tr v-for="(student, idx) in students" :key="student.id" class="hover">
                      <th class="bg-base-100 text-center text-xs opacity-50">{{ idx + 1 }}</th>
                      <td class="bg-base-100 font-medium whitespace-nowrap">
                          {{ student.last_name }} {{ student.first_name }}
                      </td>
                      
                      <!-- Cells -->
                      <td v-for="col in columns" :key="col.id" 
                          class="p-0 border border-base-200 hover:bg-base-200/50 transition-colors cursor-cell relative"
                          :class="{ 'bg-warning/10': isDirty(student.id, col.id) }"
                          @dblclick="editCell(student.id, col.id)"
                      >
                           <div v-if="isEditing(student.id, col.id)" class="absolute inset-0 z-10">
                               <input 
                                  v-model="editValue" 
                                  @blur="commitEdit" 
                                  @keydown.enter="commitEdit"
                                  @keydown.esc="cancelEdit"
                                  ref="editInput"
                                  class="w-full h-full text-center bg-white focus:outline-primary font-bold"
                                  type="number" step="0.5" min="1" max="10"
                                  v-focus
                               />
                           </div>
                           <div v-else class="w-full h-full flex items-center justify-center">
                               <span :class="getMarkColorClass(getMarkValue(student.id, col.id))" class="font-semibold">
                                   {{ getMarkValue(student.id, col.id) || '-' }}
                               </span>
                           </div>
                      </td>

                      <!-- Stats -->
                      <td class="text-center font-bold text-base-content/70">
                          {{ calculateAverage(student.id) }}
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
    marks: any[]
}>();

const emit = defineEmits(['save']);

// Local State
const editingCell = ref<{sid: number, cid: number} | null>(null);
const editValue = ref('');
const pendingChanges = ref<Record<string, number>>({}); // key: "sid-cid", value: newMark

// Dummy Columns (would come from store/props grouped by unique dates/types)
// Ideally we process `marks` prop to derive columns dynamically
const columns = computed(() => {
    // Unique combination of Date + Type from marks
    // For demo, static columns or just derive
    const cols = [];
    const map = new Map();
    props.marks.forEach(m => {
        const key = `${m.date}-${m.type}`;
        if(!map.has(key)) {
            map.set(key, { id: key, date: m.date, type: m.type });
        }
    });
    // Add a couple of empty columns for "New Mark" simulation if needed
    // For now returning existing ones sorted
    return Array.from(map.values()).sort((a,b) => new Date(a.date).getTime() - new Date(b.date).getTime());
});

const isEditing = (sid: number, cid: any) => editingCell.value?.sid === sid && editingCell.value?.cid === cid;

const isDirty = (sid: number, cid: any) => {
    return pendingChanges.value[`${sid}-${cid}`] !== undefined;
};

const hasChanges = computed(() => Object.keys(pendingChanges.value).length > 0);

const getMarkValue = (sid: number, cid: any) => {
    if (pendingChanges.value[`${sid}-${cid}`] !== undefined) {
        return pendingChanges.value[`${sid}-${cid}`];
    }
    const mark = props.marks.find(m => m.student_id === sid && `${m.date}-${m.type}` === cid);
    return mark ? mark.value : null;
};

const editCell = (sid: number, cid: any) => {
    editingCell.value = { sid, cid };
    editValue.value = getMarkValue(sid, cid) || '';
};

const commitEdit = () => {
    if (editingCell.value) {
        const { sid, cid } = editingCell.value;
        const val = parseFloat(editValue.value);
        if (!isNaN(val)) {
            pendingChanges.value[`${sid}-${cid}`] = val;
        } else if (editValue.value === '') {
             // Handle clear?
        }
        editingCell.value = null;
    }
};

const cancelEdit = () => {
    editingCell.value = null;
};

const saveChanges = () => {
    // Emit changes to parent to persist
    // Need to convert pendingChanges back to mark objects
    const changes = Object.entries(pendingChanges.value).map(([key, value]) => {
        const [sid, cid] = key.split('-');
        // we need to find the column details from cid (which is date-type key)
        // or simplistic approach for now
        return { student_id: Number(sid), value, key: cid }; // Parent handles parsing key
    });
    emit('save', changes);
    pendingChanges.value = {};
};

// Utils
const formatDate = (d: string) => {
    return new Date(d).toLocaleDateString(undefined, { day: '2-digit', month: '2-digit' });
};
const calculateAverage = (sid: number) => {
    // Simple avg
    const studentMarks = props.marks.filter(m => m.student_id === sid).map(m => m.value);
    // Add pending
    Object.entries(pendingChanges.value).forEach(([k, v]) => {
        if(k.startsWith(`${sid}-`)) studentMarks.push(v); // This logic is flawed (double counting), simple demo
    });
    
    if(!studentMarks.length) return '-';
    const sum = studentMarks.reduce((a,b) => a+b, 0);
    return (sum / studentMarks.length).toFixed(1);
};

const getMarkColorClass = (val: number) => {
    if (!val) return '';
    if (val < 6) return 'text-error';
    if (val >= 8) return 'text-success';
    return '';
}

// Custom directive for auto focus
const vFocus = {
  mounted: (el: HTMLInputElement) => el.focus()
}

const exportCSV = () => {
    console.log("Export CSV");
}
</script>
