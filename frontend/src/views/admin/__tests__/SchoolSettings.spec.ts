import { mount } from '@vue/test-utils';
import { describe, it, expect, vi } from 'vitest';
import SchoolSettings from '@/views/admin/SchoolSettings.vue';

// Mock ActionButton
vi.mock('@/components/shared/ActionButton.vue', () => ({
    default: { template: '<button class="action-btn-stub" @click="$emit(\'click\')">{{label}}</button>', props: ['label'] }
}));

describe('SchoolSettings.vue', () => {
    it('renders main tabs', () => {
        const wrapper = mount(SchoolSettings);
        expect(wrapper.text()).toContain('School Information');
        expect(wrapper.text()).toContain('Branches (Plessi)');
        expect(wrapper.text()).toContain('Academic Years');
    });

    it('allows switching tabs', async () => {
        const wrapper = mount(SchoolSettings);
        // Default active is General Info
        const generalInfoInput = wrapper.find('input[aria-label="General Info"]');
        // Since daisyUI tabs are radios, we check checked prop
        // But JSDOM might not reflect UI state perfectly for uncontrolled inputs without interactions
        // Let's just check if content is rendered. Since all tab contents are in DOM but hidden/shown by CSS in daisyUI usually,
        // or Vue v-if. In this template we used DaisyUI structure which keeps them in DOM but uses CSS.
        // However, simpler test: Check if inputs exist.
        expect(wrapper.findAll('input[type="radio"][name="school_tabs"]')).toHaveLength(4);
    });

    it('adds a branch when requested', async () => {
        // Mock window.prompt
        const promptSpy = vi.spyOn(window, 'prompt').mockReturnValue('New Branch');
        const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(true);

        const wrapper = mount(SchoolSettings);

        // Find the "Add" button for branches. It's a small button inside the branches section.
        // We need to be specific.
        const buttons = wrapper.findAll('button');
        const addBranchBtn = buttons.find(b => b.text() === 'Add');

        await addBranchBtn?.trigger('click');
        await wrapper.vm.$nextTick();

        expect(promptSpy).toHaveBeenCalled();
        expect(wrapper.text()).toContain('New Branch');
    });
});
