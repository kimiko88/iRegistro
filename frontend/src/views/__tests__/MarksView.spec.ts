import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import MarksView from '../marks/MarksView.vue';
import { useStudentStore } from '@/stores/student';
import { useAuthStore } from '@/stores/auth';

vi.mock('chart.js/auto', () => ({
    default: class {
        destroy() { }
    }
}));

describe('MarksView.vue', () => {
    it('renders marks list', () => {
        const wrapper = mount(MarksView, {
            global: {
                plugins: [createTestingPinia({
                    initialState: {
                        auth: { user: { role: 'student' } },
                        student: {
                            marks: [
                                { id: 1, date: '2025-01-01', subject: 'Math', grade: 9, teacher: 'Rossi', notes: 'Good' }
                            ]
                        }
                    }
                })],
                stubs: {
                    MarksChart: true
                }
            }
        });

        expect(wrapper.text()).toContain('Marks & Grades');
        expect(wrapper.text()).toContain('Math');
        expect(wrapper.text()).toContain('9');
    });
});
