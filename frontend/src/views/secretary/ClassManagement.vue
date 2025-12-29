<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
// import { useAuthStore } from '@/stores/auth';
import api from '@/services/api'; 
import Modal from '@/components/shared/Modal.vue';
import { useUIStore } from '@/stores/ui';

const ui = useUIStore();
// const auth = useAuthStore();
const schoolId = 1; // Simplification: should get from user context

const classes = ref<any[]>([]);
const teachers = ref<any[]>([]);
const subjects = ref<any[]>([]);

const filteredTeachers = computed(() => {
    if (!newAssignment.value.subjectId) return [];
    return teachers.value.filter(t => {
        // Check if teacher has this subject assigned
        // Assuming t.subjects is an array of objects {id, ...}
        if (!t.subjects || !Array.isArray(t.subjects)) return false; // Or true if we want to show all if no subjects assigned? Request implies strict filter.
        return t.subjects.some((s: any) => s.id === newAssignment.value.subjectId);
    });
});


const isCreateClassModalOpen = ref(false);
const isAssignModalOpen = ref(false);

const activeTab = ref('classes');

const newClass = ref({
    grade: 1,
    section: '',
    year: '2024-25',
    curriculumId: 1 // hardcoded for demo
});

const newAssignment = ref({
    classId: null as number | null,
    subjectId: null as number | null,
    teacherId: null as number | null
});

onMounted(async () => {
    await fetchClasses();
});

const fetchClasses = async () => {
    ui.setLoading(true);
    try {
        const res = await api.get(`/schools/${schoolId}/classes`);
        classes.value = res.data;
    } catch (e) {
        console.error(e);
        // Fallback for demo if API fails
        if (classes.value.length === 0) {
             classes.value = [
                 { id: 101, grade: 1, section: 'A', year: '2024-25', assignments: ['Matematica', 'Italiano'] },
                 { id: 102, grade: 2, section: 'B', year: '2024-25', assignments: ['Storia', 'Inglese'] },
             ];
        }
    } finally {
        ui.setLoading(false);
    }
};

// Fetch teachers/subjects only when needed (e.g. opening modal)
const fetchMetadata = async () => {
     try {
         // Mocking or needing endpoints for all teachers/subjects
         // Assuming we have basic list endpoints or using admin ones if permitted
         // For Secretary, we might need specific endpoints.
         // Let's assume generic access for now or mocked.
         const tRes = await api.get('/admin/users?role=Insegnante'); 
         teachers.value = tRes.data || [];
         
         // Subjects endpoint might be needed.
         // const sRes = await api.get(`/schools/${schoolId}/subjects`);
         // subjects.value = sRes.data;
         subjects.value = [
             {id: 1, name: 'Matematica', code: 'MAT01'},
             {id: 2, name: 'Italiano', code: 'ITA01'},
             {id: 3, name: 'Storia', code: 'STO01'},
             {id: 4, name: 'Inglese', code: 'ING01'},
             {id: 5, name: 'Scienze', code: 'SCI01'}
         ];
     } catch(e) { console.error(e); }
};

const openCreateClassModal = () => { isCreateClassModalOpen.value = true; };
const closeCreateClassModal = () => { isCreateClassModalOpen.value = false; };

const createClass = async () => {
    try {
        await api.post(`/schools/${schoolId}/classes`, {
            curriculum_id: newClass.value.curriculumId,
            grade: Number(newClass.value.grade),
            section: newClass.value.section.toUpperCase(),
            year: newClass.value.year
        });
        ui.addNotification({ type: 'success', message: 'Class created' });
        closeCreateClassModal();
        fetchClasses();
    } catch (e) {
        ui.addNotification({ type: 'error', message: 'Failed to create class' });
    }
};

const openAssignModal = async (cls: any) => {
    newAssignment.value.classId = cls.id;
    newAssignment.value.subjectId = null;
    newAssignment.value.teacherId = null;
    await fetchMetadata();
    isAssignModalOpen.value = true;
};
const closeAssignModal = () => { isAssignModalOpen.value = false; };

const assignSubject = async () => {
     try {
        await api.post(`/schools/${schoolId}/assignments`, {
            class_id: newAssignment.value.classId,
            subject_id: newAssignment.value.subjectId,
            teacher_id: newAssignment.value.teacherId
        });
        ui.addNotification({ type: 'success', message: 'Assignment created successfully' });
        closeAssignModal();
    } catch (e) {
        ui.addNotification({ type: 'error', message: 'Failed to assign' });
    }
};
</script>

<template>
  <div class="p-4 space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold">Class Management</h1>
      <button class="btn btn-primary" @click="openCreateClassModal">+ Create New Class</button>
    </div>

    <!-- Stats or Filters -->
    <div class="stats shadow w-full">
        <div class="stat">
            <div class="stat-title">Total Classes</div>
            <div class="stat-value">{{ classes.length }}</div>
        </div>
        <div class="stat">
            <div class="stat-title">Academic Year</div>
            <div class="stat-value text-secondary">2024-25</div>
        </div>
    </div>

    <div class="overflow-x-auto bg-base-100 rounded-box shadow">
      <table class="table w-full">
        <thead class="bg-base-200">
          <tr>
            <th>Class</th>
            <th>Year</th>
            <th>Assignments</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="cls in classes" :key="cls.id" class="hover">
            <td>
                <div class="font-bold text-lg">{{ cls.grade }}{{ cls.section }}</div>
            </td>
            <td>{{ cls.year }}</td>
            <td>
                <div class="text-sm opacity-70">
                    <span v-if="!cls.assignments || cls.assignments.length === 0">No subjects assigned</span>
                    <div v-else class="flex flex-wrap gap-1">
                        <span v-for="a in cls.assignments" :key="a" class="badge badge-ghost badge-xs">{{ a }}</span>
                    </div>
                </div>
            </td>
            <td>
              <button class="btn btn-sm btn-outline btn-info gap-2" @click="openAssignModal(cls)">
                   <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>
                   Assign Teachers
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Class Modal -->
    <Modal title="Create Class" :isOpen="isCreateClassModalOpen" @close="closeCreateClassModal">
         <div class="space-y-4">
             <div class="alert alert-info shadow-sm text-xs">
                 <span>Creating a class will automatically set up the default curriculum structure.</span>
             </div>
             <div class="flex gap-4">
                 <div class="form-control w-1/2">
                    <label class="label">Grade (1-5)</label>
                    <input type="number" v-model="newClass.grade" class="input input-bordered" min="1" max="5" />
                 </div>
                 <div class="form-control w-1/2">
                    <label class="label">Section</label>
                    <input type="text" v-model="newClass.section" class="input input-bordered" placeholder="A" maxlength="2" />
                 </div>
             </div>
             <div class="form-control">
                <label class="label">Academic Year</label>
                <input type="text" v-model="newClass.year" class="input input-bordered" />
             </div>
         </div>
         <template #actions>
             <button class="btn" @click="closeCreateClassModal">Cancel</button>
             <button class="btn btn-primary" @click="createClass">Create</button>
         </template>
    </Modal>
    
    <!-- Assign Teacher Modal -->
    <Modal title="Assign Subject & Teacher" :isOpen="isAssignModalOpen" @close="closeAssignModal">
         <div class="space-y-4">
             <div class="form-control">
                <label class="label">Subject</label>
                <select class="select select-bordered" v-model="newAssignment.subjectId">
                    <option disabled selected :value="null">Select Subject</option>
                    <option v-for="s in subjects" :key="s.id" :value="s.id">{{ s.name }} ({{ s.code }})</option>
                </select>
             </div>
             <div class="form-control">
                <label class="label">Teacher</label>
                <select class="select select-bordered" v-model="newAssignment.teacherId" :disabled="!newAssignment.subjectId">
                     <option disabled selected :value="null">Select Teacher</option>
                    <option v-for="t in filteredTeachers" :key="t.id" :value="t.id">{{ t.first_name }} {{ t.last_name }}</option>
                </select>
                <label class="label" v-if="newAssignment.subjectId && filteredTeachers.length === 0">
                    <span class="label-text-alt text-warning">No teachers found for this subject.</span>
                </label>
             </div>
         </div>
         <template #actions>
             <button class="btn" @click="closeAssignModal">Cancel</button>
             <button class="btn btn-primary" @click="assignSubject" :disabled="!newAssignment.teacherId || !newAssignment.subjectId">
                 Confirm Assignment
             </button>
         </template>
    </Modal>

  </div>
</template>
