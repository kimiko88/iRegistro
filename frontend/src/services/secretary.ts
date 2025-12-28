import api from './api';

export default {
    getInbox() {
        return api.get('/secretary/documents/inbox');
    },
    getArchive(params: any) {
        return api.get('/secretary/documents/archive', { params });
    },
    approveDocument(docId: number) {
        return api.post(`/secretary/documents/${docId}/approve`);
    },
    rejectDocument(docId: number, reason: string) {
        return api.post(`/secretary/documents/${docId}/reject`, { reason });
    },
    printDocuments(docIds: number[]) {
        // Theoretically triggers backend PDF generation for download or sends to printer
        // For frontend, likely returns a PDF blob
        return api.post('/secretary/documents/print-batch', { ids: docIds }, { responseType: 'blob' });
    }
};
