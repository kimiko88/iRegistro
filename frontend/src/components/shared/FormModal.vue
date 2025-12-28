<template>
  <dialog class="modal" :class="{ 'modal-open': isOpen }">
    <div class="modal-box w-11/12 max-w-5xl">
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" @click="$emit('close')">✕</button>
      <h3 class="font-bold text-lg mb-4">{{ title }}</h3>
      
      <form @submit.prevent="onSubmit">
        <slot></slot>
        
        <div class="modal-action">
          <button type="button" class="btn" @click="$emit('close')">Cancel</button>
          <button 
            type="submit" 
            class="btn btn-primary"
            :class="{ 'loading': loading }"
            :disabled="loading || disabled"
          >
            {{ submitLabel }}
          </button>
        </div>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop" v-if="isOpen">
      <button @click="$emit('close')">close</button>
    </form>
  </dialog>
</template>

<script setup lang="ts">
interface Props {
  isOpen: boolean;
  title: string;
  loading?: boolean;
  disabled?: boolean;
  submitLabel?: string;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  disabled: false,
  submitLabel: 'Save',
});

const emit = defineEmits(['close', 'submit']);

const onSubmit = () => {
  emit('submit');
};
</script>
