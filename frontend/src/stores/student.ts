import { defineStore } from 'pinia';
import studentApi from '@/services/student';
import { useUIStore } from './ui';

export const useStudentStore = defineStore('student', {
    state: () => ({
        overview: null as any | null,
        marks: [] as any[],
        absences: [] as any[],
        colloquiums: [] as any[],
        documents: [] as any[],
        messages: [] as any[],
    }),

    actions: {
        async fetchOverview() {
            const ui = useUIStore();
            ui.setLoading(true);
            try {
                const res = await studentApi.getOverview();
                this.overview = res.data;
            } catch (e) {
                console.error("Failed to fetch overview", e);
            } finally {
                ui.setLoading(false);
            }
        },
        async fetchMarks(params?: any) {
            try {
                const res = await studentApi.getMarks(params);
                this.marks = res.data;
            } catch (e) {
                console.error("Failed to fetch marks");
            }
        },
        async fetchAbsences(params?: any) {
            try {
                const res = await studentApi.getAbsences(params);
                this.absences = res.data;
            } catch (e) { console.error(e); }
        },
        async fetchColloquiums() {
            try {
                const res = await studentApi.getColloquiums();
                this.colloquiums = res.data;
            } catch (e) { console.error(e); }
        },
        async fetchDocuments() {
            try {
                const res = await studentApi.getDocuments();
                this.documents = res.data;
            } catch (e) { console.error(e); }
        },
        async fetchMessages(params?: any) {
            try {
                const res = await studentApi.getMessages(params);
                this.messages = res.data;
            } catch (e) { console.error(e); }
        }
    },
});
