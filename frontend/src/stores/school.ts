import { defineStore } from 'pinia';
// import api from '@/services/api';

export const useSchoolStore = defineStore('school', {
    state: () => ({
        currentSchool: null as any | null,
        classes: [] as any[],
        teachers: [] as any[],
        students: [] as any[],
    }),

    actions: {
        async fetchSchoolDetails(schoolId: number) {
            // Implementation pending backend integration
            console.log("Fetching details for school", schoolId);
        },
        // Add other actions for managing school data
    },
});
