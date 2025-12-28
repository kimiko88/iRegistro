<template>
  <button
    :class="[
      'btn',
      variantClass,
      sizeClass,
      { 'loading': loading, 'btn-disabled': disabled || loading }
    ]"
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

<script setup lang="ts">
import { ref, computed } from 'vue';
import ConfirmDialog from './ConfirmDialog.vue'; // Cyclic dependency handled by Vue usually fine, or we can use slot
// But better to separate concern or use a simpler approach. 
// For now, let's keep it simple and assume ConfirmDialog is globally registered or imported.
// Actually, to avoid circular deps if ConfirmDialog uses ActionButton, be careful.
// Let's implement ConfirmDialog first in next step or assume it handles its buttons primitively.

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
