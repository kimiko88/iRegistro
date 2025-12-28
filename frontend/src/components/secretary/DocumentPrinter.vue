<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{
    selectedIds: number[];
}>();

const emit = defineEmits(['print', 'cancel']);

const selectedTemplate = ref('school_official');

function print() {
    emit('print', selectedTemplate.value);
}
</script>

<template>
  <div class="modal modal-open">
      <div class="modal-box">
          <h3 class="font-bold text-lg">Print Documents</h3>
          <p class="py-4">You are about to generate a print job for {{ selectedIds.length }} documents.</p>
          
          <div class="form-control w-full max-w-xs mb-4">
            <label class="label">
                <span class="label-text">Select Template</span>
            </label>
            <select class="select select-bordered" v-model="selectedTemplate">
                <option value="school_official">Official School Letterhead</option>
                <option value="simple_report">Simple Report (Internal)</option>
                <option value="certificate">Certificate Format</option>
            </select>
          </div>

          <div class="modal-action">
              <button class="btn" @click="$emit('cancel')">Cancel</button>
              <button class="btn btn-primary" @click="print">
                  <span class="loading loading-spinner loading-xs hidden"></span>
                  Confirm Print
              </button>
          </div>
      </div>
  </div>
</template>
