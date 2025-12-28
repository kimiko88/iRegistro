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
            // ui.setLoading(true);
            try {
                // const res = await studentApi.getOverview();
                // this.overview = res.data;

                // Mock
                this.overview = {
                    gpa: 8.1,
                    attendance: 98,
                    nextColloquium: null,
                    recentMarks: [
                        { subject: 'Italian', value: 9, date: '2023-12-21' },
                        { subject: 'English', value: 8, date: '2023-12-19' }
                    ]
                };
            } finally {
                // ui.setLoading(false);
            }
        },
        async fetchMarks() {
            // Mock
            this.marks = []; // Populate with more data if needed
        }
    },
});
