import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import MessagesView from '../MessagesView.vue';
import { createPinia, setActivePinia } from 'pinia';

// Mock UI Store
vi.mock('@/stores/ui', () => ({
    useUIStore: () => ({
        setLoading: vi.fn(),
    })
}));

// Mock Communication API
const mockGetConversations = vi.fn();
vi.mock('@/services/communication', () => ({
    default: {
        getConversations: () => mockGetConversations(),
    }
}));

describe('MessagesView', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
        vi.clearAllMocks();
    });

    it('renders message list from API', async () => {
        const mockMessages = [
            { id: 1, sender: 'Prof. Bianchi', subject: 'Math Homework', date: '2023-10-25', preview: 'Please remember...' }
        ];
        mockGetConversations.mockResolvedValue({ data: mockMessages });

        const wrapper = mount(MessagesView);
        await flushPromises();

        expect(mockGetConversations).toHaveBeenCalled();
        expect(wrapper.findAll('li').length).toBe(1);
        expect(wrapper.text()).toContain('Prof. Bianchi');
    });

    it('displays message detail on selection', async () => {
        const mockMessages = [
            { id: 1, sender: 'Prof. Bianchi', subject: 'Math Homework', date: '2023-10-25', preview: 'Please remember...' }
        ];
        mockGetConversations.mockResolvedValue({ data: mockMessages });

        const wrapper = mount(MessagesView);
        await flushPromises();

        const firstMsg = wrapper.find('li');
        await firstMsg.trigger('click');

        expect(wrapper.find('h2').exists()).toBe(true);
        expect(wrapper.text()).toContain('Math Homework');
    });
});
