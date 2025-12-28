import api from './api';

export default {
    getClasses() {
        return api.get('/teacher/classes');
    },
    getClassDetails(classId: number) {
        return api.get(`/teacher/classes/${classId}`);
    },
    getStudents(classId: number) {
        return api.get(`/teacher/classes/${classId}/students`);
    },
    getMarks(classId: number, subjectId: number) {
        return api.get(`/teacher/classes/${classId}/subjects/${subjectId}/marks`);
    },
    saveMark(data: any) {
        return api.post('/teacher/marks', data);
    },
    getAbsences(classId: number, date: string) {
        return api.get(`/teacher/classes/${classId}/absences`, { params: { date } });
    },
    saveAbsences(classId: number, absences: any[]) {
        return api.post(`/teacher/classes/${classId}/absences`, { absences });
    },
    getSchedule() {
        return api.get('/teacher/schedule');
    },
    // Colloqui e Messaggi would likely reuse Communication Service endpoints or have wrappers here
};
