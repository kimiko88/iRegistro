import { setActivePinia, createPinia } from 'pinia';
import { useStudentStore } from '../student';
import { describe, it, expect, beforeEach, vi } from 'vitest';

vi.mock('../ui', () => ({
    useUIStore: () => ({
        setLoading: vi.fn(),
    })
}));

describe('Student Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('fetches overview', async () => {
        const store = useStudentStore();
        await store.fetchOverview();

        expect(store.overview).not.toBeNull();
        expect(store.overview.gpa).toBe(8.1);
        expect(store.overview.recentMarks).toHaveLength(2);
    });
});
