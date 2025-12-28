import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import SuperAdminDashboard from '../SuperAdminDashboard.vue';

// Mock child components to avoid deep rendering issues if any
// But we want to test interaction with DataTable, so let's keep it real if possible
// Or shallowMount. mount is better for integration.

describe('SuperAdminDashboard.vue', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('renders dashboard title and KPIs', () => {
        const wrapper = mount(SuperAdminDashboard, {
            global: {
                stubs: {
                    // Stub complex sub-components if needed
                    // 'DataTable': true 
                }
            }
        });

        expect(wrapper.text()).toContain('SuperAdmin Dashboard');
        expect(wrapper.text()).toContain('Total Schools');
        expect(wrapper.text()).toContain('Total Users');
    });

    it('loads schools on mount', async () => {
        // Since store is mocked in store logic (hardcoded in action), 
        // mounting triggers fetchSchools which populates state.
        const wrapper = mount(SuperAdminDashboard);

        // Wait for lifecycle hooks
        await wrapper.vm.$nextTick();

        // Check if table contains data from mock
        expect(wrapper.text()).toContain('Liceo Scientifico');
        expect(wrapper.text()).toContain('Istituto Tecnico');
    });
});
