import api from './api';

export default {
    getChildren() {
        // Returns list of students linked to parent
        return api.get('/parent/children');
    },
    getChildOverview(studentId: number) {
        return api.get(`/parent/children/${studentId}/overview`);
    },
    // Reuse specific detail endpoints or wrap them
};
