<template>
  <dialog class="modal" :class="{ 'modal-open': isOpen }">
    <div class="modal-box">
      <h3 class="font-bold text-lg" :class="{'text-error': type === 'danger'}">{{ title }}</h3>
      <p class="py-4">{{ message }}</p>
      <div class="modal-action">
        <button class="btn" @click="$emit('cancel')">Cancel</button>
        <button 
          class="btn" 
          :class="confirmBtnClass"
          @click="$emit('confirm')"
        >
          {{ confirmText }}
        </button>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop" v-if="isOpen">
      <button @click="$emit('cancel')">close</button>
    </form>
  </dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  isOpen: boolean;
  title?: string;
  message?: string;
  confirmText?: string;
  type?: 'info' | 'warning' | 'danger';
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Confirm Action',
  message: 'Are you sure?',
  confirmText: 'Confirm',
  type: 'info'
});

defineEmits(['confirm', 'cancel']);

const confirmBtnClass = computed(() => {
  switch (props.type) {
    case 'danger': return 'btn-error';
    case 'warning': return 'btn-warning';
    default: return 'btn-primary';
  }
});
</script>
