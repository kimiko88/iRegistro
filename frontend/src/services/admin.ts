import api from './api';

// const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

/* const api = axios.create({
    baseURL: API_URL,
    withCredentials: true,
}); */

// Types can be moved to a shared types file later
export interface School {
    id: string;
    name: string;
    region: string;
    address: string;
    city: string;
    code: string;
}

export interface KPIStats {
    schoolsCount: number;
    usersCount: number;
    storageUsed: number;
    activeUsersLast30Days: number;
}

export default {
    // KPIs & Dashboard
    getKPIs() {
        return api.get<KPIStats>('/admin/kpis');
    },
    getSchoolDistribution() {
        return api.get('/admin/schools/distribution');
    },

    // User Management
    getUsers(params: any) {
        return api.get('/admin/users', { params });
    },
    createUser(userData: any) {
        return api.post('/admin/users', userData);
    },
    updateUser(id: string, userData: any) {
        return api.put(`/admin/users/${id}`, userData);
    },
    deleteUser(id: string) {
        return api.delete(`/admin/users/${id}`);
    },
    importUsers(formData: FormData) {
        return api.post('/admin/users/import', formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
    },

    // School Management
    getSchools(params: any) {
        return api.get('/admin/schools', { params });
    },
    createSchool(schoolData: any) {
        return api.post('/superadmin/schools', schoolData);
    },
    updateSchool(id: string, schoolData: any) {
        return api.put(`/superadmin/schools/${id}`, schoolData);
    },

    // Audits & Backups
    getAuditLogs(params: any) {
        return api.get('/admin/audit-logs', { params });
    },
    getBackups() {
        return api.get('/admin/backups');
    },
    createBackup() {
        return api.post('/admin/backups');
    },

    // Export
    exportData(request: any) {
        return api.post('/admin/export', request, { responseType: 'blob' });
    }
};
