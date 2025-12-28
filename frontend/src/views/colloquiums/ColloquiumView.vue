<script setup lang="ts">
import { ref } from 'vue';

import communicationApi from '@/services/communication';
import { useUIStore } from '@/stores/ui';
import { onMounted } from 'vue';

const teachers = ref<any[]>([]); // This view logic might need adjustment if API returns flat slots
// Adapting: Fetch available slots and group by teacher? 
// Or fetching teachers first? 
// For simplicity, let's assume we fetch a list of teachers with slots, 
// OR we fetch slots and group them locally.
// Let's implement fetchSlots and grouping for now.

const ui = useUIStore();

onMounted(async () => {
    ui.setLoading(true);
    try {
        // Mocking structure match since API likely returns flat list of slots
        // We simulate fetching "Teachers with Slots" or we fetch slots and transform.
        // For this existing UI, let's try to group functionality.
        const res = await communicationApi.getAvailableSlots(); 
        // Assuming res.data is [ {id, teacher_id, teacher_name, time...} ]
        
        // Grouping logic (simplified for integration):
        const slots = res.data || [];
        const groups: Record<number, any> = {};
        
        slots.forEach((s: any) => {
            if (!groups[s.teacher_id]) {
                groups[s.teacher_id] = { id: s.teacher_id, name: s.teacher_name || 'Teacher ' + s.teacher_id, slots: [] };
            }
            groups[s.teacher_id].slots.push(s);
        });
        teachers.value = Object.values(groups);

    } catch (e) {
        console.error("Failed to fetch slots", e);
    } finally {
        ui.setLoading(false);
    }
});

const selectedTeacher = ref<any>(null);

async function bookSlot(slot: any) {
    if (confirm(`Book appointment with ${selectedTeacher.value.name} at ${slot.time}?`)) {
        try {
            await communicationApi.bookSlot(slot.id);
            alert('Booking confirmed!');
            // remove slot locally
            selectedTeacher.value.slots = selectedTeacher.value.slots.filter((s: any) => s.id !== slot.id);
        } catch (e) {
            alert('Booking failed');
        }
    }
}
</script>

<template>
  <div class="p-4 space-y-6">
      <h1 class="text-2xl font-bold">Colloquiums</h1>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div v-for="teacher in teachers" :key="teacher.id" class="card bg-base-100 shadow-xl">
              <div class="card-body">
                  <h2 class="card-title">{{ teacher.name }}</h2>
                  <p v-if="teacher.slots.length > 0">{{ teacher.slots.length }} slots available</p>
                  <p v-else class="text-error">No slots available</p>
                  
                  <div class="card-actions justify-end mt-4">
                      <button class="btn btn-primary" @click="selectedTeacher = teacher" :disabled="teacher.slots.length === 0">
                          View Slots
                      </button>
                  </div>
              </div>
          </div>
      </div>

      <!-- Slots Modal -->
      <dialog id="slots_modal" class="modal" :class="{ 'modal-open': selectedTeacher }">
          <div class="modal-box" v-if="selectedTeacher">
              <h3 class="font-bold text-lg">Available Slots for {{ selectedTeacher.name }}</h3>
              <ul class="py-4 space-y-2">
                  <li v-for="slot in selectedTeacher.slots" :key="slot.id" class="flex justify-between items-center bg-base-200 p-2 rounded">
                      <span>{{ slot.time }}</span>
                      <button class="btn btn-sm btn-accent" @click="bookSlot(slot)">Book</button>
                  </li>
              </ul>
              <div class="modal-action">
                  <button class="btn" @click="selectedTeacher = null">Close</button>
              </div>
          </div>
      </dialog>
  </div>
</template>
