import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import MessagesView from '../MessagesView.vue';

describe('MessagesView', () => {
    it('renders message list', () => {
        const wrapper = mount(MessagesView);
        expect(wrapper.findAll('li').length).toBeGreaterThan(0);
        expect(wrapper.text()).toContain('Prof. Bianchi');
    });

    it('displays message detail on selection', async () => {
        const wrapper = mount(MessagesView);

        const firstMsg = wrapper.find('li');
        await firstMsg.trigger('click');

        expect(wrapper.find('h2').exists()).toBe(true);
        expect(wrapper.text()).toContain('Math Homework'); // Subject of first mock msg
    });

    it('shows placeholder when no message selected', () => {
        const wrapper = mount(MessagesView);
        expect(wrapper.text()).toContain('Select a message to read');
    });
});
