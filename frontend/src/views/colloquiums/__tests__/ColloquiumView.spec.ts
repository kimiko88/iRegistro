import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import ColloquiumView from '../ColloquiumView.vue';
import { createPinia, setActivePinia } from 'pinia';

vi.mock('@/stores/ui', () => ({
    useUIStore: () => ({
        setLoading: vi.fn(),
    })
}));

const mockGetAvailableSlots = vi.fn();
const mockBookSlot = vi.fn();

vi.mock('@/services/communication', () => ({
    default: {
        getAvailableSlots: () => mockGetAvailableSlots(),
        bookSlot: (id: any) => mockBookSlot(id),
    }
}));

describe('ColloquiumView', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
        vi.clearAllMocks();
    });

    it('renders teacher list from API slots', async () => {
        // Mock data structure: flat list of slots
        const mockSlots = [
            { id: 101, teacher_id: 1, teacher_name: 'Prof. Rossi', time: '2023-11-01 10:00' },
            { id: 102, teacher_id: 1, teacher_name: 'Prof. Rossi', time: '2023-11-01 10:30' },
            { id: 201, teacher_id: 2, teacher_name: 'Prof. Verdi', time: '2023-11-02 09:00' }
        ];
        mockGetAvailableSlots.mockResolvedValue({ data: mockSlots });

        const wrapper = mount(ColloquiumView);
        await flushPromises();

        expect(mockGetAvailableSlots).toHaveBeenCalled();
        expect(wrapper.text()).toContain('Prof. Rossi');
        expect(wrapper.text()).toContain('Prof. Verdi');
    });

    it('handles booking confirmation via API', async () => {
        const mockSlots = [
            { id: 101, teacher_id: 1, teacher_name: 'Prof. Rossi', time: '2023-11-01 10:00' }
        ];
        mockGetAvailableSlots.mockResolvedValue({ data: mockSlots });
        mockBookSlot.mockResolvedValue({});

        const wrapper = mount(ColloquiumView);
        await flushPromises();

        // Open modal
        await wrapper.findAll('button')[0].trigger('click');

        // Mock window confirm
        window.confirm = vi.fn(() => true);
        window.alert = vi.fn();

        // Click Book
        const bookBtn = wrapper.find('.modal-box button.btn-accent');
        await bookBtn.trigger('click');

        expect(mockBookSlot).toHaveBeenCalledWith(101);
        await flushPromises(); // wait for async bookSlot
        expect(window.alert).toHaveBeenCalledWith('Booking confirmed!');
    });
});
