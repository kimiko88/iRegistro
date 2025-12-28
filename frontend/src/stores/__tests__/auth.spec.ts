import { setActivePinia, createPinia } from 'pinia';
import { useAuthStore } from '../auth';
import { describe, it, expect, beforeEach, vi, beforeAll } from 'vitest';
import api from '@/services/api';

// Mock API
vi.mock('@/services/api', () => ({
    default: {
        post: vi.fn()
    }
}));

describe('Auth Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
        vi.clearAllMocks();
    });

    beforeAll(() => {
        global.localStorage = {
            getItem: () => null,
            setItem: () => { },
            removeItem: () => { },
            length: 0,
            key: () => null,
            clear: () => { }
        };
    });

    it('initializes with default state', () => {
        const store = useAuthStore();
        expect(store.isAuthenticated).toBe(false);
        expect(store.user).toBeNull();
    });

    it('login success updates state', async () => {
        const store = useAuthStore();
        (api.post as any).mockResolvedValue({
            data: { token: 'fake-token', role: 'Admin' }
        });

        const success = await store.login({ email: 'test@example.com', password: 'password' });

        expect(success).toBe(true);
        expect(store.token).toBe('fake-token');
        expect(store.isAuthenticated).toBe(true);
        expect(store.user).toEqual({ role: 'Admin' });
    });

    it('login failure returns false', async () => {
        const store = useAuthStore();
        (api.post as any).mockRejectedValue(new Error('Auth failed'));

        const success = await store.login({ email: 'test@example.com', password: 'wrong' });

        expect(success).toBe(false);
        expect(store.isAuthenticated).toBe(false);
    });
});
