<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-md p-6">
      <h3 class="text-lg font-bold mb-4 text-gray-900 dark:text-white">Book Colloquium</h3>
      
      <div v-if="loading" class="text-center py-4">Loading slots...</div>
      
      <div v-else class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Select Slot</label>
          <select v-model="selectedSlotId" class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 dark:bg-gray-700 dark:border-gray-600">
            <option v-for="slot in slots" :key="slot.id" :value="slot.id">
              {{ new Date(slot.startTime).toLocaleString() }} - {{ new Date(slot.endTime).toLocaleTimeString() }}
            </option>
          </select>
        </div>
      </div>

      <div class="mt-6 flex justify-end gap-3">
        <button @click="$emit('close')" class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200">Cancel</button>
        <button 
          @click="confirmBooking" 
          :disabled="!selectedSlotId || loading"
          class="px-4 py-2 text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 disabled:opacity-50"
        >
          Confirm
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{
  isOpen: boolean;
  teacherId: number;
  fetchSlots: (teacherId: number) => Promise<any[]>;
}>();

const emit = defineEmits(['close', 'book']);

const slots = ref<any[]>([]);
const loading = ref(false);
const selectedSlotId = ref<number | null>(null);

const loadSlots = async () => {
  if (!props.teacherId) return;
  loading.value = true;
  try {
    slots.value = await props.fetchSlots(props.teacherId);
  } finally {
    loading.value = false;
  }
};

watch(() => props.isOpen, (newVal) => {
  if (newVal) loadSlots();
});

const confirmBooking = () => {
  if (selectedSlotId.value) {
    emit('book', selectedSlotId.value);
  }
};
</script>
