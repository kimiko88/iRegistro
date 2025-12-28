import api from './api';

export default {
    // Super Admin
    createSchool(data: any) {
        return api.post('/superadmin/schools', data);
    },
    getSchools() {
        return api.get('/superadmin/schools'); // Assuming endpoint exists
    },

    // School Admin
    getUsers(params: any) {
        return api.get('/admin/users', { params });
    },
    createUser(data: any) {
        return api.post('/admin/users', data);
    },
    getAuditLogs(params: any) {
        return api.get('/admin/audit-logs', { params });
    },
    getSettings() {
        return api.get('/admin/settings');
    },
    updateSetting(key: string, value: any) {
        return api.put('/admin/settings', { key, value });
    },
    requestExport(type: string) {
        return api.post('/admin/data-export', { type });
    }
};
