import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import ColloquiumView from '../ColloquiumView.vue';

describe('ColloquiumView', () => {
    it('renders teacher list', () => {
        const wrapper = mount(ColloquiumView);
        expect(wrapper.text()).toContain('Prof. Rossi');
        expect(wrapper.text()).toContain('Prof. Verdi');
    });

    it('opens modal on View Slots click', async () => {
        const wrapper = mount(ColloquiumView);

        // Find button for teacher with slots (Prof Rossi)
        const buttons = wrapper.findAll('button');
        const viewSlotsBtn = buttons[0]; // First button

        await viewSlotsBtn.trigger('click');

        const modal = wrapper.find('dialog');
        expect(modal.classes()).toContain('modal-open');
        expect(wrapper.text()).toContain('Available Slots for Prof. Rossi');
    });

    it('handles booking confirmation', async () => {
        const wrapper = mount(ColloquiumView);
        window.confirm = vi.fn(() => true);
        window.alert = vi.fn();

        // Open modal
        await wrapper.findAll('button')[0].trigger('click');

        // Click Book
        const bookBtn = wrapper.find('.modal-box button.btn-accent');
        if (bookBtn.exists()) {
            await bookBtn.trigger('click');
            expect(window.confirm).toHaveBeenCalled();
            expect(window.alert).toHaveBeenCalledWith('Booking confirmed!');
        }
    });
});
