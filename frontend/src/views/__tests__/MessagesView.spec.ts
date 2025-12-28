import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import MessagesView from '../messages/MessagesView.vue';

describe('MessagesView.vue', () => {
    it('renders messages list', () => {
        const wrapper = mount(MessagesView, {
            global: {
                plugins: [createTestingPinia({
                    initialState: {
                        auth: { user: { role: 'student' } },
                        student: {
                            messages: [
                                { id: 1, date: '2025-01-01', sender: 'Teacher', subject: 'Homework', body: 'Do it', read: false }
                            ]
                        }
                    }
                })],
            }
        });

        expect(wrapper.text()).toContain('Messages');
        expect(wrapper.text()).toContain('Homework');
        expect(wrapper.text()).toContain('Teacher');
    });
});
