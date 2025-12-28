import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import ParentDashboard from '../parent/ParentDashboard.vue';
import { useParentStore } from '@/stores/parent';

// Mock Router
vi.mock('vue-router', () => ({
    useRouter: () => ({
        push: vi.fn(),
    }),
    RouterLink: {
        template: '<a><slot /></a>',
    }
}));

// Mock Chart.js to avoid canvas errors
vi.mock('chart.js/auto', () => ({
    default: class {
        destroy() { }
    }
}));

describe('ParentDashboard.vue', () => {
    it('renders correctly', () => {
        const wrapper = mount(ParentDashboard, {
            global: {
                plugins: [createTestingPinia({
                    createSpy: vi.fn,
                    initialState: {
                        parent: {
                            children: [],
                            selectedChildId: null,
                            currentChildOverview: null
                        }
                    }
                })],
                stubs: {
                    StudentSelect: true,
                    StudentCard: true,
                    MarksChart: true
                }
            }
        });

        expect(wrapper.text()).toContain('Parent Dashboard');
    });

    it('fetches children on mount', () => {
        const wrapper = mount(ParentDashboard, {
            global: {
                plugins: [createTestingPinia({
                    createSpy: vi.fn,
                    stubActions: false
                })],
                stubs: {
                    StudentSelect: true,
                    StudentCard: true,
                    MarksChart: true
                }
            }
        });

        const store = useParentStore();
        expect(store.fetchChildren).toHaveBeenCalled();
    });
});
