import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import MarksTab from '@/components/teacher/MarksTab.vue';
import { createTestingPinia } from '@pinia/testing';
import { useTeacherStore } from '@/stores/teacher';

describe('MarksTab.vue', () => {
    it('renders students and marks correctly', async () => {
        const pinia = createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                teacher: {
                    students: [
                        { id: 1, first_name: 'Mario', last_name: 'Rossi' },
                        { id: 2, first_name: 'Luigi', last_name: 'Verdi' }
                    ],
                    marks: [
                        { id: 101, student_id: 1, value: 8.5, type: 'Oral', date: '2024-10-10' }
                    ],
                    classes: [{ id: 1, subjectId: 5 }]
                }
            }
        });

        const store = useTeacherStore();
        store.selectedClassId = 1; // Manually select class since component assumes it

        const wrapper = mount(MarksTab, {
            global: {
                plugins: [pinia]
            }
        });

        // Check student names
        expect(wrapper.text()).toContain('Rossi Mario');
        expect(wrapper.text()).toContain('Verdi Luigi');

        // Check mark rendered
        const marks = wrapper.findAll('.badge');
        expect(marks.length).toBeGreaterThan(0);
        expect(marks[0].text()).toContain('8.5');
    });

    it('opens modal on click', async () => {
        const pinia = createTestingPinia({ createSpy: vi.fn });
        const wrapper = mount(MarksTab, {
            global: { plugins: [pinia] }
        });

        // Mock modal showModal which is not implemented in JSDOM usually
        const modal = wrapper.find('dialog').element as HTMLDialogElement;
        modal.showModal = vi.fn();

        await wrapper.find('button.btn-primary').trigger('click');
        expect(modal.showModal).toHaveBeenCalled();
    });
});
