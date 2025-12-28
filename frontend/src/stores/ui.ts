import { defineStore } from 'pinia';

interface Notification {
    id: string;
    type: 'info' | 'success' | 'warning' | 'error';
    message: string;
    timeout?: number;
}

export const useUIStore = defineStore('ui', {
    state: () => ({
        theme: localStorage.getItem('theme') || 'corporate',
        isLoading: false,
        notifications: [] as Notification[],
    }),

    actions: {
        setTheme(theme: string) {
            this.theme = theme;
            localStorage.setItem('theme', theme);
            document.documentElement.setAttribute('data-theme', theme);
        },
        setLoading(loading: boolean) {
            this.isLoading = loading;
        },
        addNotification(notification: Omit<Notification, 'id'>) {
            const id = Date.now().toString();
            this.notifications.push({ ...notification, id });

            if (notification.timeout !== 0) {
                setTimeout(() => {
                    this.removeNotification(id);
                }, notification.timeout || 3000);
            }
        },
        removeNotification(id: string) {
            this.notifications = this.notifications.filter((n) => n.id !== id);
        },
    },
});
