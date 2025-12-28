import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import DirectorDashboard from '../DirectorDashboard.vue';

// Mock the store
const fetchDashboardDataSpy = vi.fn();
vi.mock('@/stores/director', () => ({
    useDirectorStore: () => ({
        kpis: {
            totalStudents: 999,
            totalTeachers: 50,
            totalClasses: 20,
            averageGrade: 8,
            attendanceRate: 98
        },
        documentsToSign: [],
        pendingRequests: [],
        loading: false,
        fetchDashboardData: fetchDashboardDataSpy
    })
}));

describe('DirectorDashboard.vue', () => {
    it('renders KPIs correctly', async () => {
        const wrapper = mount(DirectorDashboard);

        expect(fetchDashboardDataSpy).toHaveBeenCalled();
        expect(wrapper.text()).toContain('999');
        expect(wrapper.text()).toContain('Total Students');
    });
});
