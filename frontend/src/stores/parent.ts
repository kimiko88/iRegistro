import { defineStore } from 'pinia';
import parentApi from '@/services/parent';
import { useUIStore } from './ui';

export const useParentStore = defineStore('parent', {
    state: () => ({
        children: [] as any[],
        selectedChildId: null as number | null,
        currentChildOverview: null as any | null,
    }),

    getters: {
        selectedChild: (state) => state.children.find(c => c.id === state.selectedChildId),
    },

    actions: {
        async fetchChildren() {
            const ui = useUIStore();
            // ui.setLoading(true);
            try {
                // const res = await parentApi.getChildren();
                // this.children = res.data;

                // Mock
                this.children = [
                    { id: 101, name: 'Alice Rossi', class: '1A', school: 'Liceo Scientifico' },
                    { id: 205, name: 'Marco Rossi', class: '3C', school: 'Media' }
                ];

                if (!this.selectedChildId && this.children.length > 0) {
                    this.selectChild(this.children[0].id);
                }
            } finally {
                // ui.setLoading(false);
            }
        },
        selectChild(id: number) {
            this.selectedChildId = id;
            this.fetchChildOverview(id);
        },
        async fetchChildOverview(childId: number) {
            // Mock
            this.currentChildOverview = {
                gpa: 7.8,
                attendance: 92,
                nextColloquium: '2024-01-15 10:00',
                recentMarks: [
                    { subject: 'Math', value: 8, date: '2023-12-20' },
                    { subject: 'History', value: 7.5, date: '2023-12-18' }
                ]
            };
        }
    },
});
