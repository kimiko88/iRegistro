import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import AbsencesView from '../absences/AbsencesView.vue';
import { useAuthStore } from '@/stores/auth';

// Mock child components
vi.mock('@/components/AbsenceCalendar.vue', () => ({
    default: { template: '<div>Calendar</div>' }
}));

describe('AbsencesView.vue', () => {
    it('renders absences list', () => {
        const wrapper = mount(AbsencesView, {
            global: {
                plugins: [createTestingPinia({
                    initialState: {
                        auth: { user: { role: 'parent' } },
                        parent: {
                            absences: [
                                { id: 1, date: '2025-01-10', type: 'Illness', justified: false }
                            ]
                        }
                    }
                })],
            }
        });

        expect(wrapper.text()).toContain('Absences & Attendance');
        expect(wrapper.text()).toContain('Illness');
        expect(wrapper.text()).toContain('Unjustified');
        // Parent should see Justify button
        expect(wrapper.text()).toContain('Justify');
    });
});
