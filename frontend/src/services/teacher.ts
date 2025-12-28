import api from './api';

export default {
    getClasses() {
        return api.get('/teacher/classes');
    },
    getStudents(classId: number) {
        return api.get(`/teacher/classes/${classId}/students`);
    },
    getMarks(classId: number, subjectId: number) {
        return api.get(`/teacher/classes/${classId}/subjects/${subjectId}/marks`);
    },
    saveMark(mark: any) {
        return api.post('/teacher/marks', mark);
    },
    getAbsences(classId: number) {
        return api.get(`/teacher/classes/${classId}/absences`);
    },
    saveAbsence(absence: any) {
        return api.post(`/teacher/classes/${absence.class_id}/absences`, absence);
    },
    // Future expansion stubs
    getSchedule(classId: number) {
        // return api.get(`/teacher/classes/${classId}/schedule`);
        return Promise.resolve({ data: [] });
    },
    getSlots() {
        // return api.get('/communication/slots/available');
        return Promise.resolve({ data: [] });
    }
};
