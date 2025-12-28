<script setup lang="ts">
defineProps<{
  title: string;
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();
</script>

<template>
  <dialog class="modal" :class="{ 'modal-open': isOpen }">
    <div class="modal-box">
      <h3 class="font-bold text-lg mb-4">{{ title }}</h3>
      <slot></slot>
      <div class="modal-action">
        <form method="dialog">
          <!-- if there is a button in form, it will close the modal -->
          <slot name="actions">
            <button class="btn" @click="emit('close')">Close</button>
          </slot>
        </form>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button @click="emit('close')">close</button>
    </form>
  </dialog>
</template>
