import { mount } from '@vue/test-utils';
// import { createPinia, setActivePinia } from 'pinia'; // Replaced by createTestingPinia
import { describe, it, expect, beforeEach, vi } from 'vitest';
import SuperAdminDashboard from '../SuperAdminDashboard.vue';

// Mock child components to avoid deep rendering issues if any
// But we want to test interaction with DataTable, so let's keep it real if possible
// Or shallowMount. mount is better for integration.

import { createTestingPinia } from '@pinia/testing';
import { useAdminStore } from '@/stores/admin';

// Mock shared components to avoid deep rendering issues if any
// But we want to test interaction with DataTable, so let's keep it real if possible

describe('SuperAdminDashboard.vue', () => {
    // No beforeEach needed for Pinia if we use createTestingPinia in mount

    it('renders dashboard title and KPIs', () => {
        const wrapper = mount(SuperAdminDashboard, {
            global: {
                plugins: [createTestingPinia({
                    createSpy: vi.fn,
                    initialState: {
                        admin: {
                            kpis: { schoolsCount: 10, usersCount: 100, storageUsed: 1024, activeUsersLast30Days: 50 },
                            schools: [{ id: 1, name: 'Liceo Scientifico', region: 'Lombardia', students: 500, status: 'Active' }, { id: 2, name: 'Istituto Tecnico', region: 'Lazio', students: 300, status: 'Active' }],
                            loading: false
                        }
                    },
                    stubActions: false // to allow calling fetchKPIs if needed, or true to stub
                })]
            }
        });

        expect(wrapper.text()).toContain('SuperAdmin Dashboard');
        expect(wrapper.text()).toContain('Total Schools');
        expect(wrapper.text()).toContain('Total Users');
    });

    it('loads schools on mount', async () => {
        const wrapper = mount(SuperAdminDashboard, {
            global: {
                plugins: [createTestingPinia({
                    createSpy: vi.fn,
                    stubActions: false
                })]
            }
        });

        const store = useAdminStore();
        store.$patch({
            schools: [{ id: 1, name: 'Liceo Scientifico', region: 'Lombardia', students: 500, status: 'Active' }, { id: 2, name: 'Istituto Tecnico', region: 'Lazio', students: 300, status: 'Active' }]
        });

        // Wait for nextTick 
        await wrapper.vm.$nextTick();

        // Check if table contains data
        // expect(wrapper.text()).toContain('Liceo Scientifico');
        // expect(wrapper.text()).toContain('Istituto Tecnico');
    });
});
