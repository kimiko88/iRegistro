import api from './api';


export default {
    getOverview() {
        return api.get('/student/overview');
    },
    getMarks(params?: any) {
        return api.get('/student/marks', { params });
    },
    getAbsences(params?: any) {
        return api.get('/student/absences', { params });
    },
    getColloquiums() {
        return api.get('/student/colloquiums');
    },
    getBookableSlots(teacherId: number) {
        return api.get('/student/colloquiums/slots', { params: { teacherId } });
    },
    bookColloquium(slotId: number) {
        return api.post('/student/colloquiums', { slotId });
    },
    getDocuments() {
        return api.get('/student/documents');
    },
    getMessages(params?: any) {
        return api.get('/student/messages', { params });
    },
    markMessageRead(messageId: number) {
        return api.post(`/student/messages/${messageId}/read`);
    }
};
