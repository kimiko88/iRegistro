import { mount } from '@vue/test-utils';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ParentDashboard from '../ParentDashboard.vue';
import { createPinia, setActivePinia } from 'pinia';
import { createRouter, createWebHistory } from 'vue-router';

describe('ParentDashboard.vue', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    // Mock Router
    const router = createRouter({
        history: createWebHistory(),
        routes: [{ path: '/', component: { template: '' } }]
    });

    it('renders dashboard with children selector', async () => {
        const wrapper = mount(ParentDashboard, {
            global: {
                plugins: [router],
                stubs: {
                    MarksView: true, // Stub child components
                }
            }
        });

        // Trigger onMounted fetch
        await wrapper.vm.$nextTick();

        // Check static text
        expect(wrapper.text()).toContain('Parent Dashboard');

        // Check if stats are rendered (based on default Mock data in store)
        // Wait for store to update
        await new Promise(resolve => setTimeout(resolve, 10));
        await wrapper.vm.$nextTick();

        const stats = wrapper.findAll('.stat-value');
        expect(stats.length).toBeGreaterThan(0);
        // GPA 7.8 from mock
        expect(wrapper.text()).toContain('7.8');
    });
});
