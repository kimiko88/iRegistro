import { mount } from '@vue/test-utils';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { createTestingPinia } from '@pinia/testing';
import UserManagement from '@/views/admin/UserManagement.vue';
import { useAdminStore } from '@/stores/admin';

// Mock shared components to simplify view testing
vi.mock('@/components/shared/DataTable.vue', () => ({
    default: { template: '<div class="data-table-stub"><slot name="item-actions" :item="{id: 1}"></slot></div>' }
}));
vi.mock('@/components/admin/UserEditModal.vue', () => ({
    default: { template: '<div class="user-edit-modal-stub" v-if="isOpen"></div>', props: ['isOpen'] }
}));
vi.mock('@/components/admin/UserImportModal.vue', () => ({
    default: { template: '<div class="user-import-modal-stub" v-if="isOpen"></div>', props: ['isOpen'] }
}));

describe('UserManagement.vue', () => {
    let wrapper: any;

    beforeEach(() => {
        wrapper = mount(UserManagement, {
            global: {
                plugins: [createTestingPinia({
                    createSpy: vi.fn,
                    initialState: {
                        admin: {
                            users: [
                                { id: 1, firstName: 'John', lastName: 'Doe', role: 'Admin', email: 'john@example.com' }
                            ],
                            totalUsers: 1,
                            loading: false
                        }
                    }
                })]
            }
        });
    });

    it('renders correctly', () => {
        expect(wrapper.find('h1').text()).toBe('User Management');
    });

    it('calls fetchUsers on mount', () => {
        const store = useAdminStore();
        expect(store.fetchUsers).toHaveBeenCalled();
    });

    it('opens create modal when "Create User" is clicked', async () => {
        const createBtn = wrapper.findAllComponents({ name: 'ActionButton' })
            .find((c: any) => c.props('label') === 'Create User');

        await createBtn.find('button').trigger('click');
        expect(wrapper.find('.user-edit-modal-stub').exists()).toBe(true);
    });

    it('opens import modal when "Import CSV" is clicked', async () => {
        const importBtn = wrapper.findAllComponents({ name: 'ActionButton' })
            .find((c: any) => c.props('label') === 'Import CSV');

        await importBtn.find('button').trigger('click');
        const modal = wrapper.findComponent({ name: 'UserImportModal' });
        // Since we mocked it with an object that might not have name, let's rely on import or order. 
        // Actually, let's just check the stub class again but ensure we wait.
        // Or check the prop on the matching component.
        // Finding by the stub class should work if v-if is true.
        // Let's rely on finding all properties.
        const modalComponent = wrapper.findAllComponents({ name: 'UserImportModal' })[0]
            || wrapper.findComponent({ name: 'UserImportModal' }); // name might be missing in mock object

        // Better: Find by importing the mock? 
        // Let's try checking the prop on the component found by ref or type.
        // Since we can't easily import the mock reference here perfectly as it's hoisted,
        // Let's use the class check but verify why it failed. 
        // Maybe ActionButton wasn't found.

        // Let's try nextTick
        await wrapper.vm.$nextTick();
        expect(wrapper.find('.user-import-modal-stub').exists()).toBe(true);
    });
});
