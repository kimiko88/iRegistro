import { defineStore } from 'pinia';
import directorService, { DirectorKPIs, DocumentToSign, PendingRequest } from '../services/director';

export const useDirectorStore = defineStore('director', {
    state: () => ({
        kpis: null as DirectorKPIs | null,
        documentsToSign: [] as DocumentToSign[],
        pendingRequests: [] as PendingRequest[],
        loading: false,
        error: null as string | null,
    }),

    actions: {
        async fetchDashboardData() {
            this.loading = true;
            this.error = null;
            try {
                this.kpis = await directorService.getKPIs();
                this.documentsToSign = await directorService.getDocumentsToSign(); // Fetching basic count usually, but list fits mock
                this.pendingRequests = await directorService.getPendingRequests();
            } catch (err: any) {
                this.error = 'Failed to load dashboard data';
                console.error(err);
            } finally {
                this.loading = false;
            }
        },

        async fetchDocumentsToSign() {
            this.loading = true;
            try {
                this.documentsToSign = await directorService.getDocumentsToSign();
            } catch (err) {
                this.error = 'Failed to fetch documents';
            } finally {
                this.loading = false;
            }
        },

        async signDocument(id: number, pin: string) {
            this.loading = true;
            try {
                await directorService.signDocument(id, pin);
                // Remove from local list or refresh
                this.documentsToSign = this.documentsToSign.filter(d => d.id !== id);
            } catch (err) {
                this.error = 'Failed to sign document. Invalid PIN?';
                throw err;
            } finally {
                this.loading = false;
            }
        },

        async fetchRequests() {
            this.loading = true;
            try {
                this.pendingRequests = await directorService.getPendingRequests();
            } catch (err) {
                this.error = 'Failed to fetch requests';
            } finally {
                this.loading = false;
            }
        },

        async approveRequest(id: number) {
            try {
                await directorService.approveRequest(id);
                this.pendingRequests = this.pendingRequests.filter(r => r.id !== id);
            } catch (err) {
                this.error = 'Failed to approve request';
            }
        },

        async rejectRequest(id: number) {
            try {
                await directorService.rejectRequest(id);
                this.pendingRequests = this.pendingRequests.filter(r => r.id !== id);
            } catch (err) {
                this.error = 'Failed to reject request';
            }
        }
    }
});
