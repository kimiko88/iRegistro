import { setActivePinia, createPinia } from 'pinia';
import { useStudentStore } from '../student';
import { describe, it, expect, beforeEach, vi } from 'vitest';

// Mock UI Store
vi.mock('../ui', () => ({
    useUIStore: () => ({
        setLoading: vi.fn(),
    })
}));

// Mock Student API
vi.mock('@/services/student', () => ({
    default: {
        getOverview: vi.fn(),
        getMarks: vi.fn(),
    }
}));

import studentApi from '@/services/student';

describe('Student Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
        vi.clearAllMocks();
    });

    it('fetches overview', async () => {
        const store = useStudentStore();

        (studentApi.getOverview as any).mockResolvedValue({
            data: {
                gpa: 8.1,
                recentMarks: [1, 2]
            }
        });

        await store.fetchOverview();

        expect(store.overview).not.toBeNull();
        expect(store.overview.gpa).toBe(8.1);
        expect(store.overview.recentMarks).toHaveLength(2);
    });
});
