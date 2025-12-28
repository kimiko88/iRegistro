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
        },
        async fetchChildOverview(childId: number) {
            const ui = useUIStore();
            // Optional: loading for sub-component?
            try {
                const res = await parentApi.getChildOverview(childId);
                this.currentChildOverview = res.data;
            } catch (e) {
                console.error("Failed to fetch child overview", e);
            }
        }
    },
});
