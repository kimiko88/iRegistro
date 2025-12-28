import { describe, it, expect, vi } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import DocumentManagement from '../DocumentManagement.vue';
import { useSecretaryStore } from '@/stores/secretary';

// Mock components
vi.mock('@/components/secretary/DocumentReceiver.vue', () => ({
    default: { template: '<div>Receiver</div>' }
}));

describe('DocumentManagement.vue', () => {
    it('fetches stats on mount and renders them', async () => {
        const wrapper = mount(DocumentManagement, {
            global: {
                plugins: [createTestingPinia({
                    createSpy: vi.fn,
                    initialState: {
                        secretary: {
                            stats: {
                                new_documents: 42,
                                processed_today: 10,
                                delivery_issues: 5
                            },
                            inbox: []
                        }
                    }
                })],
                stubs: {
                    RouterLink: true
                }
            }
        });

        const store = useSecretaryStore();
        await flushPromises();
        expect(store.fetchStats).toHaveBeenCalled();

        // Check if stats are rendered
        expect(wrapper.text()).toContain('42'); // New
        expect(wrapper.text()).toContain('10'); // Processed
        expect(wrapper.text()).toContain('5');  // Issues
    });
});
