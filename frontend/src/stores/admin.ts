import { defineStore } from 'pinia';
import { ref } from 'vue';
import adminService, { KPIStats } from '../services/admin';
import { useNotificationStore } from './notification';

export const useAdminStore = defineStore('admin', () => {
    const kpis = ref<KPIStats | null>(null);
    const users = ref<any[]>([]);
    const schools = ref<any[]>([]);
    const auditLogs = ref<any[]>([]);
    const loading = ref(false);
    const error = ref<string | null>(null);
    const totalUsers = ref(0);
    const totalSchools = ref(0);

    const notificationStore = useNotificationStore();

    const fetchKPIs = async () => {
        try {
            loading.value = true;
            const response = await adminService.getKPIs();
            kpis.value = response.data;
        } catch (err) {
            error.value = 'Failed to fetch KPIs';
            notificationStore.error('Failed to fetch KPIs');
        } finally {
            loading.value = false;
        }
    };

    const fetchUsers = async (params: any = {}) => {
        try {
            loading.value = true;
            const response = await adminService.getUsers(params);
            users.value = response.data.data || response.data; // Handle pagination or list
            totalUsers.value = response.data.total || users.value.length;
        } catch (err) {
            notificationStore.error('Failed to fetch users');
        } finally {
            loading.value = false;
        }
    };

    const fetchSchools = async (params: any = {}) => {
        try {
            loading.value = true;
            const response = await adminService.getSchools(params);
            schools.value = response.data.data || response.data;
            totalSchools.value = response.data.total || schools.value.length;
        } catch (err) {
            // notificationStore.error('Failed to fetch schools');
            // Silencing for now as backend might not have this endpoint fully ready or mock
        } finally {
            loading.value = false;
        }
    };

    const fetchAuditLogs = async (params: any = {}) => {
        try {
            loading.value = true;
            const response = await adminService.getAuditLogs(params);
            auditLogs.value = response.data.data || response.data;
        } catch (err) {
            notificationStore.error('Failed to load audit logs');
        } finally {
            loading.value = false;
        }
    };

    const createSchool = async (data: any) => {
        try {
            loading.value = true;
            await adminService.createSchool(data);
            notificationStore.success('School created successfully');
            await fetchSchools();
        } catch (err) {
            notificationStore.error('Failed to create school');
            throw err;
        } finally {
            loading.value = false;
        }
    };

    const updateSchool = async (id: string, data: any) => {
        try {
            loading.value = true;
            await adminService.updateSchool(id, data);
            notificationStore.success('School updated successfully');
            await fetchSchools();
        } catch (err) {
            notificationStore.error('Failed to update school');
            throw err;
        } finally {
            loading.value = false;
        }
    };

    const createUser = async (data: any) => {
        try {
            loading.value = true;
            await adminService.createUser(data); // This maps to backend "CreateUser"
            notificationStore.success('User created successfully');
            await fetchUsers();
        } catch (err) {
            notificationStore.error('Failed to create user');
            throw err;
        } finally {
            loading.value = false;
        }
    };

    const updateUser = async (id: string, data: any) => {
        try {
            loading.value = true;
            await adminService.updateUser(id, data);
            // notificationStore.success('User updated'); // View handles success msg or we do? View did it.
        } catch (err) {
            throw err;
        } finally {
            loading.value = false;
        }
    };

    return {
        kpis,
        users,
        schools,
        auditLogs,
        loading,
        error,
        totalUsers,
        totalSchools,
        fetchKPIs,
        fetchUsers,
        fetchSchools,
        fetchAuditLogs,
        createSchool,
        updateSchool,
        createUser,
        updateUser
    };
});
