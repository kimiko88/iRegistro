import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import UserEditModal from '../UserEditModal.vue';
import { useAuthStore } from '@/stores/auth';

describe('UserEditModal', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    describe('allowedRoles computed property', () => {
        it('should show only Teacher, Student, Parent for Secretary', async () => {
            const authStore = useAuthStore();
            authStore.user = {
                id: 1,
                role: 'Secretary',
                email: 'secretary@test.com',
                schoolId: 1
            };

            const wrapper = mount(UserEditModal, {
                props: {
                    isOpen: true,
                    showSchoolSelect: false
                },
                global: {
                    stubs: {
                        FormModal: {
                            template: '<div><slot /></div>'
                        }
                    }
                }
            });

            // Use vm to access component instance
            const allowedRoles = (wrapper.vm as any).allowedRoles;
            expect(allowedRoles).toEqual(['Teacher', 'Student', 'Parent']);
            expect(allowedRoles).not.toContain('Admin');
            expect(allowedRoles).not.toContain('SuperAdmin');
        });

        it('should show all roles except SuperAdmin for Admin', () => {
            const authStore = useAuthStore();
            authStore.user = {
                id: 1,
                role: 'Admin',
                email: 'admin@test.com',
                schoolId: 1
            };

            const wrapper = mount(UserEditModal, {
                props: {
                    isOpen: true,
                    showSchoolSelect: false
                },
                global: {
                    stubs: {
                        FormModal: {
                            template: '<div><slot /></div>'
                        }
                    }
                }
            });

            const allowedRoles = (wrapper.vm as any).allowedRoles;
            expect(allowedRoles).toContain('Teacher');
            expect(allowedRoles).toContain('Secretary');
            expect(allowedRoles).not.toContain('SuperAdmin');
        });

        it('should show all roles for SuperAdmin', () => {
            const authStore = useAuthStore();
            authStore.user = {
                id: 1,
                role: 'SuperAdmin',
                email: 'superadmin@test.com'
            };

            const wrapper = mount(UserEditModal, {
                props: {
                    isOpen: true,
                    showSchoolSelect: true
                },
                global: {
                    stubs: {
                        FormModal: {
                            template: '<div><slot /></div>'
                        }
                    }
                }
            });

            const allowedRoles = (wrapper.vm as any).allowedRoles;
            expect(allowedRoles).toContain('SuperAdmin');
            expect(allowedRoles).toContain('Admin');
            expect(allowedRoles).toContain('Secretary');
        });
    });

    describe('Subject selection for Teachers', () => {
        it('should show subject selection when role is Teacher and school is selected', async () => {
            const authStore = useAuthStore();
            authStore.user = {
                id: 1,
                role: 'Admin',
                email: 'admin@test.com',
                schoolId: 1
            };

            const wrapper = mount(UserEditModal, {
                props: {
                    isOpen: true,
                    showSchoolSelect: false,
                    preselectedSchoolId: 1
                },
                global: {
                    stubs: {
                        FormModal: {
                            template: '<div><slot /></div>'
                        }
                    }
                }
            });

            // Set role to Teacher
            await wrapper.setData({
                form: {
                    ...wrapper.vm.form,
                    role: 'Teacher',
                    schoolId: 1
                }
            });

            await wrapper.vm.$nextTick();

            // Subject selection should be visible
            expect(wrapper.html()).toContain('Assigned Subjects');
        });
    });
});
