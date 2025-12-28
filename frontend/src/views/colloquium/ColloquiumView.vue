<template>
  <div class="p-4 space-y-6">
    <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100">Colloquiums</h1>
    
    <!-- Active Colloquiums -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4">
      <h3 class="text-xl font-semibold mb-4 text-gray-800 dark:text-gray-200">Upcoming Appointments</h3>
      <div v-if="appointments.length === 0" class="text-gray-500 italic">No upcoming colloquiums.</div>
      <div v-else class="space-y-4">
        <div v-for="c in appointments" :key="c.id" class="border p-4 rounded-lg flex justify-between items-center dark:border-gray-700">
           <div>
             <h4 class="font-bold">{{ c.teacherName }} - {{ c.subject }}</h4>
             <p class="text-sm text-gray-600 dark:text-gray-400">{{ new Date(c.date).toLocaleString() }}</p>
             <p class="text-sm text-gray-600 dark:text-gray-400">Type: {{ c.type }}</p>
           </div>
           <span class="px-3 py-1 rounded bg-green-100 text-green-800 text-sm">Confirmed</span>
        </div>
      </div>
    </div>

    <!-- Booking Section (Parent Only) -->
    <div v-if="isParent" class="bg-white dark:bg-gray-800 rounded-lg shadow p-4 mt-6">
      <h3 class="text-xl font-semibold mb-4 text-gray-800 dark:text-gray-200">Book New Colloquium</h3>
      
      <!-- List of Teachers/Subjects -->
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead class="bg-gray-50 dark:bg-gray-700">
             <tr>
               <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Teacher</th>
               <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Subject</th>
               <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Action</th>
             </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
            <tr v-for="teacher in teachers" :key="teacher.id">
              <td class="px-6 py-4">{{ teacher.name }}</td>
              <td class="px-6 py-4">{{ teacher.subject }}</td>
              <td class="px-6 py-4">
                <button 
                  @click="openBooking(teacher.id)"
                  class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 text-sm"
                >
                  Book
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <BookingModal 
      :isOpen="isBookingModalOpen" 
      :teacherId="selectedTeacherId"
      :fetchSlots="fetchSlots"
      @close="isBookingModalOpen = false"
      @book="handleWebBooking"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useParentStore } from '@/stores/parent';
import { useStudentStore } from '@/stores/student';
import { useAuthStore } from '@/stores/auth';
import BookingModal from '@/components/ColloquiumBookingModal.vue';
import parentApi from '@/services/parent';
import studentApi from '@/services/student';

const auth = useAuthStore();
const parentStore = useParentStore();
const studentStore = useStudentStore();

const isBookingModalOpen = ref(false);
const selectedTeacherId = ref<number>(0);

const isParent = computed(() => auth.user?.role === 'parent');

const appointments = computed(() => {
   return isParent.value ? parentStore.colloquiums : studentStore.colloquiums;
});

// Mock teachers list - ideally fetched from store based on child's class
const teachers = ref([
 { id: 1, name: "Mario Rossi", subject: "Math" },
 { id: 2, name: "Luigi Verdi", subject: "Science" },
]);

const openBooking = (teacherId: number) => {
  selectedTeacherId.value = teacherId;
  isBookingModalOpen.value = true;
};

const fetchSlots = async (teacherId: number) => {
   if (isParent.value) return (await parentApi.getBookableSlots(teacherId)).data;
   return (await studentApi.getBookableSlots(teacherId)).data;
};

const handleWebBooking = async (slotId: number) => {
  try {
    if (isParent.value && parentStore.selectedChildId) {
      await parentApi.bookColloquium(parentStore.selectedChildId, slotId);
      await parentStore.fetchColloquiums();
    } else {
       await studentApi.bookColloquium(slotId);
       await studentStore.fetchColloquiums();
    }
    isBookingModalOpen.value = false;
  } catch (e) {
    alert("Booking failed");
  }
};

onMounted(() => {
   if (isParent.value) parentStore.fetchColloquiums();
   else studentStore.fetchColloquiums();
});
</script>
