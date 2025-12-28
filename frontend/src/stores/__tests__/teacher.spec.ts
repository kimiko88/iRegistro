import { setActivePinia, createPinia } from 'pinia';
import { useTeacherStore } from '../teacher';
import { describe, it, expect, beforeEach, vi, beforeAll } from 'vitest';

// Mock UI Store
vi.mock('../ui', () => ({
    useUIStore: () => ({
        setLoading: vi.fn(),
        addNotification: vi.fn()
    })
}));

// Mock Teacher API
vi.mock('@/services/teacher', () => ({
    default: {
        getClasses: vi.fn(),
        getStudents: vi.fn(),
        getMarks: vi.fn(),
        saveMark: vi.fn()
    }
}));

import teacherApi from '@/services/teacher';

describe('Teacher Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
        vi.clearAllMocks();
    });

    it('fetches classes and maps them correctly', async () => {
        const store = useTeacherStore();

        // Mock backend response structure
        const mockResponse = {
            data: [
                {
                    class: { id: 1, grade: 1, section: 'A' },
                    subject: { id: 10, name: 'Math' }
                },
                {
                    class: { id: 2, grade: 2, section: 'B' },
                    subject: { id: 11, name: 'Physics' }
                }
            ]
        };
        (teacherApi.getClasses as any).mockResolvedValue(mockResponse);
        // Mock getStudents because fetchClasses calls it internally for the first class
        (teacherApi.getStudents as any).mockResolvedValue({ data: [] });

        await store.fetchClasses();

        expect(store.classes).toHaveLength(2);
        expect(store.classes[0].name).toBe('1A');
        expect(store.classes[0].subject).toBe('Math');
        expect(store.classes[0].subjectId).toBe(10);

        // Should auto-select first class
        expect(store.selectedClassId).toBe(1);
    });

    it('fetches students for selected class', async () => {
        const store = useTeacherStore();

        (teacherApi.getStudents as any).mockResolvedValue({
            data: [{ id: 100, name: 'Student 1' }]
        });

        await store.fetchStudents(1);
        expect(store.students).toHaveLength(1);
        expect(store.students[0].name).toBe('Student 1');
    });

    it('optimistically saves mark', async () => {
        const store = useTeacherStore();
        const newMark = { studentId: 100, value: 9 };

        (teacherApi.saveMark as any).mockResolvedValue({ data: newMark });

        await store.saveMark(newMark);

        expect(store.marks).toHaveLength(1);
        expect(store.marks[0].value).toBe(9);
    });
});
