import { defineStore } from 'pinia';
import adminApi from '@/services/admin';
import { useUIStore } from './ui';

export const useAdminStore = defineStore('admin', {
    state: () => ({
        schools: [] as any[],
        users: [] as any[],
        auditLogs: [] as any[],
        settings: {} as any,
    }),

    actions: {
        async fetchSchools() {
            const ui = useUIStore();
            ui.setLoading(true);
            try {
                // const res = await adminApi.getSchools();
                // this.schools = res.data;

                // Mock
                this.schools = [
                    { id: 1, name: 'Liceo Scientifico', type: 'High School', users: 120, storage: '45GB' },
                    { id: 2, name: 'Istituto Tecnico', type: 'Technical', users: 85, storage: '12GB' }
                ];
            } finally {
                ui.setLoading(false);
            }
        },
        async fetchUsers(params = {}) {
            const ui = useUIStore();
            ui.setLoading(true);
            try {
                const res = await adminApi.getUsers(params);
                this.users = res.data.map((u: any) => ({
                    id: u.id,
                    name: `${u.first_name} ${u.last_name}`,
                    email: u.email,
                    role: u.role,
                    status: u.locked_until ? 'Locked' : 'Active'
                }));
            } catch (error) {
                console.error('Failed to fetch users:', error);
                ui.addNotification({ type: 'error', message: 'Failed to load users' });
            } finally {
                ui.setLoading(false);
            }
        },
        async createSchool(data: any) {
            await adminApi.createSchool(data);
            await this.fetchSchools();
        },
        async fetchAuditLogs() {
            const ui = useUIStore();
            ui.setLoading(true);
            try {
                const res = await adminApi.getAuditLogs({});
                this.auditLogs = res.data.logs || [];
            } catch (e) {
                // Mock
                this.auditLogs = [
                    { id: 1, action: 'LOGIN', user_id: 1, timestamp: new Date().toISOString() }
                ];
            } finally {
                ui.setLoading(false);
            }
        }
        // More actions...
    },
});
