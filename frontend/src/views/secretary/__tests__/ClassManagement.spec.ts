import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import ClassManagement from '../ClassManagement.vue';

describe('ClassManagement', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    describe('filteredTeachers computed property', () => {
        it('should return empty array when no subject is selected', () => {
            const wrapper = mount(ClassManagement, {
                global: {
                    stubs: {
                        Modal: true,
                        DataTable: true
                    }
                }
            });

            // Access computed property
            const filteredTeachers = (wrapper.vm as any).filteredTeachers;
            expect(filteredTeachers).toEqual([]);
        });

        it('should filter teachers by selected subject', async () => {
            const wrapper = mount(ClassManagement, {
                global: {
                    stubs: {
                        Modal: true,
                        DataTable: true
                    }
                }
            });

            // Setup test data
            const teachers = [
                {
                    id: 1,
                    first_name: 'John',
                    last_name: 'Doe',
                    subjects: [{ id: 1, name: 'Math' }, { id: 2, name: 'Science' }]
                },
                {
                    id: 2,
                    first_name: 'Jane',
                    last_name: 'Smith',
                    subjects: [{ id: 3, name: 'History' }]
                },
                {
                    id: 3,
                    first_name: 'Bob',
                    last_name: 'Johnson',
                    subjects: [{ id: 1, name: 'Math' }]
                }
            ];

            // Set component data
            (wrapper.vm as any).teachers = teachers;
            (wrapper.vm as any).newAssignment.subjectId = 1; // Select Math

            await wrapper.vm.$nextTick();

            // Only teachers with Math should be filtered
            const filteredTeachers = (wrapper.vm as any).filteredTeachers;
            expect(filteredTeachers).toHaveLength(2);
            expect(filteredTeachers[0].id).toBe(1);
            expect(filteredTeachers[1].id).toBe(3);
        });

        it('should return empty array when teacher has no subjects', async () => {
            const wrapper = mount(ClassManagement, {
                global: {
                    stubs: {
                        Modal: true,
                        DataTable: true
                    }
                }
            });

            const teachers = [
                {
                    id: 1,
                    first_name: 'John',
                    last_name: 'Doe',
                    subjects: null // No subjects assigned
                }
            ];

            (wrapper.vm as any).teachers = teachers;
            (wrapper.vm as any).newAssignment.subjectId = 1;

            await wrapper.vm.$nextTick();

            const filteredTeachers = (wrapper.vm as any).filteredTeachers;
            expect(filteredTeachers).toEqual([]);
        });
    });
});
