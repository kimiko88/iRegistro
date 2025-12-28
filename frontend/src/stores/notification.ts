import { defineStore } from 'pinia';
import { ref } from 'vue';

export type NotificationType = 'success' | 'error' | 'warning' | 'info';

export interface Notification {
    id: string;
    type: NotificationType;
    message: string;
    title?: string;
    duration?: number;
    dismissible?: boolean;
}

export const useNotificationStore = defineStore('notification', () => {
    const notifications = ref<Notification[]>([]);

    const add = (notification: Omit<Notification, 'id'>) => {
        const id = Date.now().toString() + Math.random().toString(36).substring(2);
        const newNotification = {
            id,
            duration: 5000,
            dismissible: true,
            ...notification
        };

        notifications.value.push(newNotification);

        if (newNotification.duration > 0) {
            setTimeout(() => {
                remove(id);
            }, newNotification.duration);
        }
    };

    const remove = (id: string) => {
        notifications.value = notifications.value.filter(n => n.id !== id);
    };

    const success = (message: string, title?: string) => add({ type: 'success', message, title });
    const error = (message: string, title?: string) => add({ type: 'error', message, title });
    const warning = (message: string, title?: string) => add({ type: 'warning', message, title });
    const info = (message: string, title?: string) => add({ type: 'info', message, title });

    return {
        notifications,
        add,
        remove,
        success,
        error,
        warning,
        info
    };
});
