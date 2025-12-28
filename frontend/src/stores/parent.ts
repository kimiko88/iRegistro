import { defineStore } from 'pinia';
import parentApi from '@/services/parent';
import { useUIStore } from './ui';

export const useParentStore = defineStore('parent', {
    state: () => ({
        children: [] as any[],
        selectedChildId: null as number | null,
        currentChildOverview: null as any | null,
        marks: [] as any[],
        absences: [] as any[],
        colloquiums: [] as any[],
        documents: [] as any[],
        messages: [] as any[],
    }),

    getters: {
        selectedChild: (state) => state.children.find(c => c.id === state.selectedChildId),
    },

    actions: {
        async fetchChildren() {
            const ui = useUIStore();
            ui.setLoading(true);
            try {
                const res = await parentApi.getChildren();
                this.children = res.data;

                if (!this.selectedChildId && this.children.length > 0) {
                    this.selectChild(this.children[0].id);
                }
            } catch (e) {
                console.error("Failed to fetch children", e);
            } finally {
                ui.setLoading(false);
            }
        },
        selectChild(id: number) {
            this.selectedChildId = id;
            this.fetchChildOverview(id);
            // Clear or refetch other data
            this.marks = [];
            this.absences = [];
            this.colloquiums = [];
        },
        async fetchChildOverview(childId: number) {
            const ui = useUIStore();
            try {
                const res = await parentApi.getChildOverview(childId);
                this.currentChildOverview = res.data;
            } catch (e) {
                console.error("Failed to fetch child overview", e);
            }
        },
        async fetchMarks(params?: any) {
            if (!this.selectedChildId) return;
            try {
                const res = await parentApi.getMarks(this.selectedChildId, params);
                this.marks = res.data;
            } catch (e) { console.error(e); }
        },
        async fetchAbsences(params?: any) {
            if (!this.selectedChildId) return;
            try {
                const res = await parentApi.getAbsences(this.selectedChildId, params);
                this.absences = res.data;
            } catch (e) { console.error(e); }
        },
        async fetchColloquiums() {
            if (!this.selectedChildId) return;
            try {
                const res = await parentApi.getColloquiums(this.selectedChildId);
                this.colloquiums = res.data;
            } catch (e) { console.error(e); }
        },
        async fetchDocuments() {
            if (!this.selectedChildId) return;
            try {
                const res = await parentApi.getDocuments(this.selectedChildId);
                this.documents = res.data;
            } catch (e) { console.error(e); }
        },
        async justifyAbsence(absenceId: number, reason: string) {
            if (!this.selectedChildId) return;
            try {
                await parentApi.justifyAbsence(this.selectedChildId, absenceId, reason);
                await this.fetchAbsences(); // refresh
            } catch (e) { console.error(e); throw e; }
        }
    },
});
