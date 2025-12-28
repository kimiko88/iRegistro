import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import DocumentsView from '../documents/DocumentsView.vue';

describe('DocumentsView.vue', () => {
    it('renders documents list', () => {
        const wrapper = mount(DocumentsView, {
            global: {
                plugins: [createTestingPinia({
                    initialState: {
                        auth: { user: { role: 'parent' } },
                        parent: {
                            documents: [
                                { id: 1, date: '2025-01-01', name: 'Report Card', type: 'PDF' }
                            ]
                        }
                    }
                })],
            }
        });

        expect(wrapper.text()).toContain('Documents');
        expect(wrapper.text()).toContain('Report Card');
    });
});
