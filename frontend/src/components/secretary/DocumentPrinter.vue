<script setup lang="ts">
// Handles batch printing selection
const props = defineProps<{
    selectedIds: number[];
}>();

const emit = defineEmits(['print', 'cancel']);

function print() {
    emit('print');
}
</script>

<template>
  <div class="modal modal-open" v-if="selectedIds.length > 0">
      <div class="modal-box">
          <h3 class="font-bold text-lg">Print Documents</h3>
          <p class="py-4">You are about to generate a print job for {{ selectedIds.length }} documents.</p>
          <div class="modal-action">
              <button class="btn btn-primary" @click="print">Confirm Print</button>
              <button class="btn" @click="$emit('cancel')">Cancel</button>
          </div>
      </div>
  </div>
</template>
