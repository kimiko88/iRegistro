import { defineStore } from 'pinia';
import secretaryApi from '@/services/secretary';
import { useUIStore } from './ui';

export const useSecretaryStore = defineStore('secretary', {
    state: () => ({
        inbox: [] as any[],
        archive: [] as any[],
    }),

    actions: {
        async fetchInbox() {
            const ui = useUIStore();
            ui.setLoading(true);
            try {
                const res = await secretaryApi.getInbox();
                this.inbox = res.data;
            } finally {
                ui.setLoading(false);
            }
        },
        async fetchArchive(filters: any = {}) {
            const ui = useUIStore();
            ui.setLoading(true);
            try {
                const res = await secretaryApi.getArchive(filters);
                this.archive = res.data;
            } finally {
                ui.setLoading(false);
            }
        },
        async approveDocument(id: number) {
            // Optimistic update
            const docIndex = this.inbox.findIndex(d => d.id === id);
            if (docIndex !== -1) this.inbox.splice(docIndex, 1);
            await secretaryApi.approveDocument(id);
        },
        async rejectDocument(id: number) {
            this.inbox = this.inbox.filter(d => d.id !== id);
        }
    },
});
