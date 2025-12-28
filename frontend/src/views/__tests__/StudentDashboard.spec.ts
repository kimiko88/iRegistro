import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import StudentDashboard from '../student/StudentDashboard.vue';
import { useStudentStore } from '@/stores/student';

// Mock Router
vi.mock('vue-router', () => ({
    useRouter: () => ({
        push: vi.fn(),
    }),
    RouterLink: {
        template: '<a><slot /></a>',
    }
}));

// Mock Chart.js
vi.mock('chart.js/auto', () => ({
    default: class {
        destroy() { }
    }
}));

describe('StudentDashboard.vue', () => {
    it('renders correctly', () => {
        const wrapper = mount(StudentDashboard, {
            global: {
                plugins: [createTestingPinia({
                    createSpy: vi.fn,
                    initialState: {
                        student: {
                            overview: {
                                firstName: 'Test',
                                lastName: 'Student',
                                averageGrade: 8.5
                            },
                        }
                    }
                })],
                stubs: {
                    StudentCard: true,
                    MarksChart: true
                }
            }
        });

        expect(wrapper.text()).toContain('My Dashboard');
        expect(wrapper.text()).toContain('GPA');
        expect(wrapper.text()).toContain('8.5');
    });

    it('fetches overview on mount', () => {
        const wrapper = mount(StudentDashboard, {
            global: {
                plugins: [createTestingPinia({
                    createSpy: vi.fn,
                    stubActions: false
                })],
                stubs: {
                    StudentCard: true,
                    MarksChart: true
                }
            }
        });

        const store = useStudentStore();
        expect(store.fetchOverview).toHaveBeenCalled();
    });
});
