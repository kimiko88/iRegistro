import api from './api';

export default {
    getChildren() {
        return api.get('/parent/children');
    },
    getChildOverview(studentId: number) {
        return api.get(`/parent/children/${studentId}/overview`);
    },
    getMarks(studentId: number, params?: any) {
        return api.get(`/parent/children/${studentId}/marks`, { params });
    },
    getAbsences(studentId: number, params?: any) {
        return api.get(`/parent/children/${studentId}/absences`, { params });
    },
    getColloquiums(studentId: number) {
        return api.get(`/parent/children/${studentId}/colloquiums`);
    },
    getBookableSlots(teacherId: number) {
        return api.get(`/parent/colloquiums/slots`, { params: { teacherId } });
    },
    bookColloquium(studentId: number, slotId: number) {
        return api.post(`/parent/children/${studentId}/colloquiums`, { slotId });
    },
    getDocuments(studentId: number) {
        return api.get(`/parent/children/${studentId}/documents`);
    },
    justifyAbsence(studentId: number, absenceId: number, reason: string) {
        return api.post(`/parent/children/${studentId}/absences/${absenceId}/justify`, { reason });
    },
    getMessages(studentId: number, params?: any) {
        return api.get(`/parent/children/${studentId}/messages`, { params });
    }
};
