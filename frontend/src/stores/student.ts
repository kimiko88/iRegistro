import { defineStore } from 'pinia';
import studentApi from '@/services/student';
import { useUIStore } from './ui';

export const useStudentStore = defineStore('student', {
    state: () => ({
        overview: null as any | null,
        marks: [] as any[],
        absences: [] as any[],
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
        async fetchMarks() {
            const ui = useUIStore();
            try {
                const res = await studentApi.getMarks();
                this.marks = res.data;
            } catch (e) {
                console.error("Failed to fetch marks");
            }
        }
    },
});
