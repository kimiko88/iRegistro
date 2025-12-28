<template>
  <div class="h-full flex flex-col gap-4 p-4 bg-base-100 rounded-box shadow">
      <!-- Header -->
      <div class="flex justify-between items-center border-b border-base-200 pb-4">
          <h2 class="text-xl font-bold">Colloquiums</h2>
          <button class="btn btn-primary" @click="openCreateSlotModal">+ New Slot</button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 h-full overflow-hidden">
          <!-- Upcoming Slots -->
          <div class="flex flex-col h-full">
              <h3 class="font-bold mb-2">My Available Slots</h3>
              <div class="overflow-y-auto bg-base-50 rounded-box border border-base-200 p-2 space-y-2 flex-1">
                  <div v-for="slot in availableSlots" :key="slot.id" class="card bg-white shadow-sm compact bordered hover:shadow-md transition-shadow">
                      <div class="card-body flex-row justify-between items-center p-3">
                          <div>
                              <div class="font-bold">{{ formatDate(slot.start) }}</div>
                              <div class="text-xs opacity-60">{{ formatTime(slot.start) }} - {{ formatTime(slot.end) }}</div>
                          </div>
                          <div class="badge badge-success badge-outline">Available</div>
                      </div>
                  </div>
              </div>
          </div>

          <!-- Bookings -->
          <div class="flex flex-col h-full">
              <h3 class="font-bold mb-2">Booked Colloquiums</h3>
               <div class="overflow-y-auto bg-base-50 rounded-box border border-base-200 p-2 space-y-2 flex-1">
                  <div v-for="booking in bookings" :key="booking.id" class="card bg-white shadow-sm compact bordered border-l-4 border-l-primary">
                      <div class="card-body p-3">
                          <div class="flex justify-between items-start">
                              <div>
                                  <div class="font-bold">{{ booking.parentName }}</div>
                                  <div class="text-xs">Student: {{ booking.studentName }}</div>
                              </div>
                              <div class="text-right">
                                  <div class="font-bold text-sm">{{ formatDate(booking.date) }}</div>
                                  <div class="text-xs opacity-60">{{ formatTime(booking.date) }}</div>
                              </div>
                          </div>
                          <div class="mt-2 text-xs bg-base-200 p-2 rounded">
                              <span class="font-semibold">Note:</span> {{ booking.note }}
                          </div>
                          <div class="card-actions justify-end mt-2">
                              <button class="btn btn-xs btn-outline">Add Note</button>
                          </div>
                      </div>
                  </div>
               </div>
          </div>
      </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

// Mock Data
const availableSlots = ref([
    { id: 1, start: '2024-05-15T15:00:00', end: '2024-05-15T15:15:00' },
    { id: 2, start: '2024-05-15T15:15:00', end: '2024-05-15T15:30:00' },
    { id: 3, start: '2024-05-22T16:00:00', end: '2024-05-22T16:15:00' },
]);

const bookings = ref([
    { id: 1, parentName: 'Mario Rossi', studentName: 'Luigi Rossi', date: '2024-05-10T16:00:00', note: 'Richiesta informazioni andamento matematica.' },
    { id: 2, parentName: 'Anna Verdi', studentName: 'Sofia Verdi', date: '2024-05-12T15:30:00', note: 'Problemi di comportamento.' },
]);

const formatDate = (d: string) => new Date(d).toLocaleDateString();
const formatTime = (d: string) => new Date(d).toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' });

const openCreateSlotModal = () => {
    alert("Create Slot Modal - Coming Soon");
};
</script>
