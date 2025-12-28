<template>
  <div class="toast toast-top toast-end z-50">
    <transition-group name="fade">
      <div 
        v-for="notification in notifications" 
        :key="notification.id"
        class="alert shadow-lg"
        :class="getAlertClass(notification.type)"
      >
        <component :is="getIcon(notification.type)" class="w-6 h-6" />
        <div>
          <h3 v-if="notification.title" class="font-bold">{{ notification.title }}</h3>
          <div class="text-xs">{{ notification.message }}</div>
        </div>
        <button v-if="notification.dismissible" class="btn btn-sm btn-ghost" @click="remove(notification.id)">
          ✕
        </button>
      </div>
    </transition-group>
  </div>
</template>

<script setup lang="ts">
import { useNotificationStore } from '@/stores/notification';
import { storeToRefs } from 'pinia';
import { AlertCircle, CheckCircle, Info, XCircle } from 'lucide-vue-next';

const store = useNotificationStore();
const { notifications } = storeToRefs(store);
const { remove } = store;

const getAlertClass = (type: string) => {
  switch (type) {
    case 'success': return 'alert-success';
    case 'error': return 'alert-error';
    case 'warning': return 'alert-warning';
    case 'info': return 'alert-info';
    default: return '';
  }
};

const getIcon = (type: string) => {
  switch (type) {
    case 'success': return CheckCircle;
    case 'error': return XCircle;
    case 'warning': return AlertCircle;
    default: return Info;
  }
};
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
