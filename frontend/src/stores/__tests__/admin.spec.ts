import { setActivePinia, createPinia } from 'pinia';
import { useAdminStore } from '../admin';
import { describe, it, expect, beforeEach, vi, beforeAll } from 'vitest';

// Mock UI Store
vi.mock('../ui', () => ({
    useUIStore: () => ({
        setLoading: vi.fn()
    })
}));

// Mock Admin API
vi.mock('@/services/admin', () => ({
    default: {
        getSchools: vi.fn(),
        getUsers: vi.fn(),
        createSchool: vi.fn()
    }
}));

describe('Admin Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    beforeAll(() => {
        // Mock localStorage/Document if likely accessed by dependencies
        global.document.documentElement.setAttribute = () => { };
    });

    it('fetches schools and populates state', async () => {
        const store = useAdminStore();
        expect(store.schools).toHaveLength(0);

        await store.fetchSchools();

        // Based on the mock implementation in store (which ignores API response for now)
        expect(store.schools).toHaveLength(2);
        expect(store.schools[0].name).toBe('Liceo Scientifico');
    });

    it('fetches users and populates state', async () => {
        const store = useAdminStore();
        expect(store.users).toHaveLength(0);

        await store.fetchUsers();

        expect(store.users).toHaveLength(2);
        expect(store.users[0].name).toBe('Mario Rossi');
    });
});
