import { ref } from 'vue';
import adminService from '@/services/admin';
import { useNotificationStore } from '@/stores/notification';

export function useDataExport() {
    const loading = ref(false);
    const error = ref<string | null>(null);
    const notificationStore = useNotificationStore();

    const exportData = async (type: string, format: 'csv' | 'json' | 'pdf', filters: any = {}) => {
        loading.value = true;
        error.value = null;

        try {
            const response = await adminService.exportData({
                type,
                format,
                filters
            });

            // Handle blob download
            const url = window.URL.createObjectURL(new Blob([response.data]));
            const link = document.createElement('a');
            link.href = url;
            const contentDisposition = response.headers['content-disposition'];
            let fileName = `${type}_export.${format}`;
            if (contentDisposition) {
                const fileNameMatch = contentDisposition.match(/filename="?(.+)"?/);
                if (fileNameMatch && fileNameMatch.length === 2)
                    fileName = fileNameMatch[1];
            }
            link.setAttribute('download', fileName);
            document.body.appendChild(link);
            link.click();
            link.remove();

            notificationStore.success('Data exported successfully');
        } catch (err: any) {
            error.value = 'Failed to export data';
            notificationStore.error('Failed to export data');
        } finally {
            loading.value = false;
        }
    };

    return {
        loading,
        error,
        exportData
    };
}
