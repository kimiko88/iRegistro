import { mount } from '@vue/test-utils';
import { describe, it, expect } from 'vitest';
import FormModal from '@/components/shared/FormModal.vue';

describe('FormModal.vue', () => {
    it('renders correctly when open', () => {
        const wrapper = mount(FormModal, {
            props: { isOpen: true, title: 'Test Modal' }
        });
        expect(wrapper.find('dialog').classes()).toContain('modal-open');
        expect(wrapper.text()).toContain('Test Modal');
    });

    it('does not render content visible when closed', () => {
        const wrapper = mount(FormModal, {
            props: { isOpen: false, title: 'Test Modal' }
        });
        expect(wrapper.find('dialog').classes()).not.toContain('modal-open');
    });

    it('emits close event when cancel button is clicked', async () => {
        const wrapper = mount(FormModal, {
            props: { isOpen: true, title: 'Test Modal' }
        });
        const buttons = wrapper.findAll('button');
        // We have: close X (top right), Cancel (bottom), Submit (bottom)
        // Cancel is usually the one with text "Cancel"
        const cancelBtn = buttons.find(b => b.text() === 'Cancel');
        await cancelBtn?.trigger('click');
        expect(wrapper.emitted()).toHaveProperty('close');
    });

    it('emits submit event when form is submitted', async () => {
        const wrapper = mount(FormModal, {
            props: { isOpen: true, title: 'Test Modal' }
        });
        await wrapper.find('form').trigger('submit.prevent');
        expect(wrapper.emitted()).toHaveProperty('submit');
    });

    it('shows loading state on submit button', () => {
        const wrapper = mount(FormModal, {
            props: { isOpen: true, title: 'Test Modal', loading: true }
        });
        const submitBtn = wrapper.find('button[type="submit"]');
        expect(submitBtn.classes()).toContain('loading');
        expect(submitBtn.attributes()).toHaveProperty('disabled');
    });
});
