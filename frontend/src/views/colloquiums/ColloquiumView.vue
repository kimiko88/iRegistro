<script setup lang="ts">
import { ref } from 'vue';

const teachers = ref([
    { id: 1, name: 'Prof. Rossi (Math)', slots: [{ id: 101, time: '2023-11-01 10:00' }, { id: 102, time: '2023-11-01 10:30' }] },
    { id: 2, name: 'Prof. Verdi (History)', slots: [] },
]);

const selectedTeacher = ref<any>(null);

function bookSlot(slot: any) {
    if (confirm(`Book appointment with ${selectedTeacher.value.name} at ${slot.time}?`)) {
        alert('Booking confirmed!');
        // remove slot locally
        selectedTeacher.value.slots = selectedTeacher.value.slots.filter((s: any) => s.id !== slot.id);
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
