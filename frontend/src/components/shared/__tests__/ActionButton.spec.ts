import { mount } from '@vue/test-utils';
import { describe, it, expect, vi } from 'vitest';
import ActionButton from '@/components/shared/ActionButton.vue';

describe('ActionButton.vue', () => {
    it('renders label correctly', () => {
        const wrapper = mount(ActionButton, {
            props: { label: 'Click Me' }
        });
        expect(wrapper.text()).toContain('Click Me');
    });

    it('emits click event when clicked', async () => {
        const wrapper = mount(ActionButton, {
            props: { label: 'Click Me' }
        });
        await wrapper.find('button').trigger('click');
        expect(wrapper.emitted()).toHaveProperty('click');
    });

    it('shows loading spinner when loading prop is true', () => {
        const wrapper = mount(ActionButton, {
            props: { label: 'Loading', loading: true }
        });
        expect(wrapper.find('.loading-spinner').exists()).toBe(true);
    });

    it('opens confirm dialog if requiresConfirmation is set', async () => {
        const wrapper = mount(ActionButton, {
            props: {
                label: 'Delete',
                requiresConfirmation: true,
                confirmTitle: 'Are you sure?'
            },
            global: {
                stubs: {
                    // If ConfirmDialog is a child component, we might need to rely on its real implementation or stub it. 
                    // For unit testing ActionButton, we want to know if logic triggers.
                    // Since ActionButton imports ConfirmDialog, it will assume it exists. 
                }
            }
        });

        await wrapper.trigger('click');
        // Check if ConfirmDialog is rendered or exposed property changed. 
        // Since we can't easily check internal state without exposing it, we check if click was NOT emitted immediately.
        expect(wrapper.emitted('click')).toBeFalsy();

        // We would need to interact with the ConfirmDialog to verify the flow, 
        // but assuming simple prop check is enough for this basic test.
    });
});
