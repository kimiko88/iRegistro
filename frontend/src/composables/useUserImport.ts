import { ref } from 'vue';
import adminService from '@/services/admin';
import { useNotificationStore } from '@/stores/notification';

export function useUserImport() {
    const file = ref<File | null>(null);
    const parsedData = ref<any[]>([]);
    const headers = ref<string[]>([]);
    const mapping = ref<Record<string, string>>({});
    const loading = ref(false);
    const progress = ref(0);
    const error = ref<string | null>(null);
    const importStats = ref<{ total: number; success: number; failed: number; errors: any[] } | null>(null);

    const notificationStore = useNotificationStore();

    const parseCSV = (fileToParse: File) => {
        loading.value = true;
        error.value = null;
        file.value = fileToParse;

        const reader = new FileReader();
        reader.onload = (e) => {
            try {
                const text = e.target?.result as string;
                const lines = text.split('\n').filter(line => line.trim());
                if (lines.length === 0) throw new Error('Empty CSV file');

                // Simple CSV parser (assumes comma separator, no quoted values with newlines for MVP)
                headers.value = lines[0].split(',').map(h => h.trim());

                parsedData.value = lines.slice(1).map(line => {
                    const values = line.split(',');
                    const row: any = {};
                    headers.value.forEach((header, index) => {
                        row[header] = values[index]?.trim() || '';
                    });
                    return row;
                });

                // Auto-map common fields
                autoMapHeaders();
            } catch (err: any) {
                error.value = 'Failed to parse CSV: ' + err.message;
                notificationStore.error(error.value!);
            } finally {
                loading.value = false;
            }
        };
        reader.onerror = () => {
            error.value = 'Error reading file';
            loading.value = false;
        };
        reader.readAsText(fileToParse);
    };

    const autoMapHeaders = () => {
        const commonMappings: Record<string, string[]> = {
            'email': ['email', 'e-mail', 'mail'],
            'firstName': ['firstname', 'first_name', 'nome', 'name'],
            'lastName': ['lastname', 'last_name', 'cognome', 'surname'],
            'role': ['role', 'ruolo', 'type']
        };

        const newMapping: Record<string, string> = {};

        // Target fields we want
        const targetFields = ['email', 'firstName', 'lastName', 'role'];

        targetFields.forEach(target => {
            const found = headers.value.find(h =>
                commonMappings[target]?.includes(h.toLowerCase())
            );
            if (found) {
                newMapping[target] = found;
            }
        });

        mapping.value = newMapping;
    };

    const importUsers = async () => {
        if (!file.value) return;

        loading.value = true;
        progress.value = 0;

        // Prepare data based on mapping
        const usersToImport = parsedData.value.map(row => {
            const user: any = {};
            Object.entries(mapping.value).forEach(([targetField, sourceHeader]) => {
                user[targetField] = row[sourceHeader];
            });
            // Default password if not mapped? Or handled by backend
            return user;
        });

        // In a real scenario with large datasets, we might batch this or send the file directly with mapping
        // For this task, we'll send the file and the mapping configuration to the backend, 
        // OR send the JSON data. The service `importUsers` takes FormData.
        // Let's send the original file + mapping JSON

        const formData = new FormData();
        formData.append('file', file.value);
        formData.append('mapping', JSON.stringify(mapping.value));

        try {
            // Simulated progress for better UX since axios upload progress is just upload, not processing
            const interval = setInterval(() => {
                if (progress.value < 90) progress.value += 10;
            }, 500);

            const response = await adminService.importUsers(formData);

            clearInterval(interval);
            progress.value = 100;

            importStats.value = response.data; // { total, success, failed, errors }
            notificationStore.success(`Import completed: ${importStats.value?.success} imported, ${importStats.value?.failed} failed`);
        } catch (err: any) {
            error.value = err.response?.data?.message || 'Import failed';
            notificationStore.error(error.value!);
        } finally {
            loading.value = false;
        }
    };

    return {
        file,
        parsedData,
        headers,
        mapping,
        loading,
        progress,
        error,
        importStats,
        parseCSV,
        importUsers
    };
}
