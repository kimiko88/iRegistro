<script lang="ts">
export default {
  inheritAttrs: false
}
</script>

<script setup lang="ts">
import { ref, computed, useAttrs } from 'vue';
import ConfirmDialog from './ConfirmDialog.vue';

// Props
interface Props {
  label?: string;
  type?: 'button' | 'submit' | 'reset';
  variant?: 'primary' | 'secondary' | 'accent' | 'ghost' | 'link' | 'error' | 'warning' | 'success' | 'info';
  size?: 'lg' | 'md' | 'sm' | 'xs';
  disabled?: boolean;
  loading?: boolean;
  icon?: any;
  requiresConfirmation?: boolean;
  confirmTitle?: string;
  confirmMessage?: string;
}

const props = withDefaults(defineProps<Props>(), {
  type: 'button',
  variant: 'primary',
  size: 'md',
  disabled: false,
  loading: false,
  requiresConfirmation: false,
  confirmTitle: 'Confirm Action',
  confirmMessage: 'Are you sure you want to proceed?',
});

const emit = defineEmits(['click', 'confirmed']);
const attrs = useAttrs(); // Use attrs if needed, but v-bind="$attrs" on button is enough

const showConfirm = ref(false);

const variantClass = computed(() => {
  switch (props.variant) {
    case 'primary': return 'btn-primary';
    case 'secondary': return 'btn-secondary';
    case 'accent': return 'btn-accent';
    case 'ghost': return 'btn-ghost';
    case 'link': return 'btn-link';
    case 'error': return 'btn-error';
    case 'warning': return 'btn-warning';
    case 'success': return 'btn-success';
    case 'info': return 'btn-info';
    default: return '';
  }
});

const sizeClass = computed(() => {
  switch (props.size) {
    case 'lg': return 'btn-lg'; // DaisyUI classes
    case 'sm': return 'btn-sm';
    case 'xs': return 'btn-xs';
    default: return '';
  }
});

const handleClick = (e: Event) => {
  if (props.requiresConfirmation) {
    e.preventDefault();
    showConfirm.value = true;
  } else {
    emit('click', e);
  }
};

const onConfirm = () => {
  showConfirm.value = false;
  emit('click');
  emit('confirmed');
};
</script>

<template>
  <button
    :class="[
      'btn',
      variantClass,
      sizeClass,
      { 'loading': loading, 'btn-disabled': disabled || loading }
    ]"
    v-bind="$attrs"
    @click="handleClick"
    :disabled="disabled || loading"
    :type="type"
  >
    <component
      v-if="icon && !loading"
      :is="icon"
      class="w-4 h-4 mr-2"
    />
    <span v-if="loading" class="loading loading-spinner"></span>
    <slot>{{ label }}</slot>
  </button>

  <ConfirmDialog
    v-if="requiresConfirmation"
    :is-open="showConfirm"
    :title="confirmTitle"
    :message="confirmMessage"
    @confirm="onConfirm"
    @cancel="showConfirm = false"
  />
</template>
