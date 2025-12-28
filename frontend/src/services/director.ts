import api from './api';

export interface DirectorKPIs {
    totalStudents: number;
    totalTeachers: number;
    totalClasses: number;
    averageGrade: number;
    attendanceRate: number;
}

export interface DocumentToSign {
    id: number;
    title: string;
    type: string;
    studentName: string;
    date: string;
    status: 'PENDING' | 'SIGNED';
}

export interface PendingRequest {
    id: number;
    type: 'GRADE_APPEAL' | 'TRANSFER' | 'LEAVE';
    requester: string;
    details: string;
    date: string;
    status: 'PENDING' | 'APPROVED' | 'REJECTED';
}

export default {
    async getKPIs(): Promise<DirectorKPIs> {
        // return api.get('/director/kpi');
        // Mock data
        return new Promise(resolve => setTimeout(() => resolve({
            totalStudents: 1250,
            totalTeachers: 85,
            totalClasses: 42,
            averageGrade: 7.8,
            attendanceRate: 94.5
        }), 500));
    },

    async getDocumentsToSign(): Promise<DocumentToSign[]> {
        // return api.get('/director/documents/sign');
        return new Promise(resolve => setTimeout(() => resolve([
            { id: 1, title: 'Pagella 1 Quadrimestre', type: 'ReportCard', studentName: 'Mario Rossi', date: '2025-01-20', status: 'PENDING' },
            { id: 2, title: 'Piano Didattico Personalizzato', type: 'PDP', studentName: 'Luigi Verdi', date: '2024-11-15', status: 'PENDING' },
            { id: 3, title: 'Documento 15 Maggio', type: 'Official', studentName: 'Classe 5A', date: '2025-05-15', status: 'PENDING' }
        ]), 600));
    },

    async signDocument(id: number, pin: string): Promise<void> {
        // return api.post(`/director/documents/${id}/sign`, { pin });
        console.log(`Signing doc ${id} with pin ${pin}`);
        return new Promise(resolve => setTimeout(resolve, 800));
    },

    async getPendingRequests(): Promise<PendingRequest[]> {
        // return api.get('/director/requests');
        return new Promise(resolve => setTimeout(() => resolve([
            { id: 101, type: 'GRADE_APPEAL', requester: 'Genitore: Giulia Bianchi', details: 'Contestazione voto matematica', date: '2025-02-10', status: 'PENDING' },
            { id: 102, type: 'TRANSFER', requester: 'Studente: Marco Neri', details: 'Nulla osta trasferimento', date: '2025-03-01', status: 'PENDING' }
        ]), 500));
    },

    async approveRequest(id: number): Promise<void> {
        // return api.post(`/director/requests/${id}/approve`);
        return new Promise(resolve => setTimeout(resolve, 500));
    },

    async rejectRequest(id: number): Promise<void> {
        // return api.post(`/director/requests/${id}/reject`);
        return new Promise(resolve => setTimeout(resolve, 500));
    }
};
