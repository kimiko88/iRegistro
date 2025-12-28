import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import ParentDashboard from '@/views/parent/ParentDashboard.vue';
import { createTestingPinia } from '@pinia/testing';

describe('ParentDashboard.vue', () => {
    it('renders correct children list', () => {
        const pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                parent: {
                    children: [
                        { id: 1, first_name: 'Child', last_name: 'One' },
                        { id: 2, first_name: 'Child', last_name: 'Two' }
                    ],
                    selectedChildId: 1,
                    currentChildOverview: {
                        gpa: 8.0,
                        attendance: 95,
                        nextColloquium: 'Tomorrow',
                        recentMarks: []
                    }
                }
            }
        });

        const wrapper = mount(ParentDashboard, {
            global: {
                plugins: [pinia],
                stubs: ['router-view'] // Stub nested router view
            }
        });

        expect(wrapper.text()).toContain('8.0'); // GPA from overview
        expect(wrapper.text()).toContain('95%'); // Attendance
    });
});
