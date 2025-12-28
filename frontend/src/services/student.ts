import api from './api';

export default {
    getOverview() {
        // Returns the overview for the currently logged-in student
        return api.get('/student/overview');
    },
    // Add other methods as needed: getMarks, getAbsences specific to 'me'
    getMarks(subjectId?: number) {
        return api.get('/student/marks', { params: { subjectId } });
    },
    getAbsences() {
        return api.get('/student/absences');
    }
};
