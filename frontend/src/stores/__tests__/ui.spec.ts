import { setActivePinia, createPinia } from 'pinia';
import { useUIStore } from '../ui';
import { describe, it, expect, beforeEach, beforeAll } from 'vitest';

describe('UI Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    // Mock localStorage
    beforeAll(() => {
        global.localStorage = {
            getItem: () => null,
            setItem: () => { },
            removeItem: () => { },
            length: 0,
            key: () => null,
            clear: () => { }
        };
        global.document.documentElement.setAttribute = () => { };
    });

    it('sets loading state', () => {
        const store = useUIStore();
        expect(store.isLoading).toBe(false);
        store.setLoading(true);
        expect(store.isLoading).toBe(true);
    });

    it('sets theme', () => {
        const store = useUIStore();
        store.setTheme('dark');
        expect(store.theme).toBe('dark');
    });

    it('adds and removes notifications', async () => {
        const store = useUIStore();
        store.addNotification({ type: 'success', message: 'Test' });
        expect(store.notifications).toHaveLength(1);
        expect(store.notifications[0].message).toBe('Test');

        // Test remove
        const id = store.notifications[0].id;
        store.removeNotification(id);
        expect(store.notifications).toHaveLength(0);
    });
});
