import { defineStore } from 'pinia';
import secretaryApi from '@/services/secretary';
import { useUIStore } from './ui';

export const useSecretaryStore = defineStore('secretary', {
    state: () => ({
        inbox: [] as any[],
        archive: [] as any[],
        deliveryReports: [] as any[],
        stats: { new_documents: 0, processed_today: 0, delivery_issues: 0 } as any,
    }),

    actions: {
        async fetchStats() {
            try {
                const res = await secretaryApi.getDashboardStats();
                this.stats = res.data;
            } catch (e) { console.error(e); }
        },
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
            if (docIndex !== -1) {
                // Move to archive strictly for UI feel if needed, or just remove from inbox
                this.inbox.splice(docIndex, 1);
            }
            await secretaryApi.approveDocument(id);
        },
        async rejectDocument(id: number, reason: string = '') {
            this.inbox = this.inbox.filter(d => d.id !== id);
            await secretaryApi.rejectDocument(id, reason);
        },
        async fetchDeliveryReports(params: any = {}) {
            const ui = useUIStore();
            // ui.setLoading(true); // Optional: don't block full UI
            try {
                const res = await secretaryApi.getDeliveryReports(params);
                this.deliveryReports = res.data;
            } catch (e) { console.error(e); }
            // finally { ui.setLoading(false); }
        },
        async markAsDelivered(docId: number, recipientId: number) {
            await secretaryApi.markDelivered(docId, recipientId);
            // Optimistic update
            const report = this.deliveryReports.find(r => r.docId === docId && r.recipientId === recipientId);
            if (report) report.delivered = true;
        },
        async bulkExport(ids: number[]) {
            const res = await secretaryApi.exportArchive(ids);
            // Convert blob to download
            const url = window.URL.createObjectURL(new Blob([res.data]));
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', `archive_export_${Date.now()}.zip`);
            document.body.appendChild(link);
            link.click();
            link.remove();
        },
        async printDocuments(ids: number[]) {
            const res = await secretaryApi.printDocuments(ids);
            const url = window.URL.createObjectURL(new Blob([res.data], { type: 'application/pdf' }));
            // print via iframe or new window
            const printWindow = window.open(url);
            if (printWindow) {
                printWindow.onload = () => { printWindow.print(); };
            }
        }
    },
});

