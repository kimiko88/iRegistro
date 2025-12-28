import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import StudentDashboard from '@/views/student/StudentDashboard.vue';
import { createTestingPinia } from '@pinia/testing';

describe('StudentDashboard.vue', () => {
    it('renders student stats overview', () => {
        const pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                student: {
                    overview: {
                        average_grade: 8.5,
                        total_absences: 3,
                        next_assignment: 'Math Test'
                    },
                    loading: false
                }
            }
        });

        const wrapper = mount(StudentDashboard, {
            global: {
                plugins: [pinia],
                stubs: ['router-link', 'router-view']
            }
        });

        expect(wrapper.text()).toContain('8.5');
        expect(wrapper.text()).toContain('3');
        expect(wrapper.text()).toContain('Math Test');
    });
});
