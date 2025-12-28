import { defineStore } from 'pinia';
import teacherApi from '@/services/teacher';
import { useUIStore } from './ui';

export const useTeacherStore = defineStore('teacher', {
    state: () => ({
        classes: [] as any[],
        selectedClassId: null as number | null,
        students: [] as any[],
        marks: [] as any[],
        absences: [] as any[],
        // Mapping of classId -> subjectId
        classSubjectMap: {} as Record<number, number>,
        // Mock Coordinator Role
        isCoordinator: true,
    }),

    getters: {
        currentClass: (state) => state.classes.find(c => c.id === state.selectedClassId),
        currentSubjectId: (state) => state.selectedClassId ? state.classSubjectMap[state.selectedClassId] : null
    },

    actions: {
        async fetchClasses() {
            const ui = useUIStore();
            ui.setLoading(true);
            try {
                const res = await teacherApi.getClasses();
                // Backend returns list of { class: {...}, subject: {...} } assignments
                // We map this to a clean class list and store the subject mapping
                this.classes = res.data.map((item: any) => {
                    this.classSubjectMap[item.class.id] = item.subject.id;
                    return {
                        id: item.class.id,
                        name: `${item.class.grade}${item.class.section}`,
                        year: item.class.year,
                        subjectName: item.subject.name
                    };
                });

                // Auto-select first class if none selected
                if (!this.selectedClassId && this.classes.length > 0) {
                    this.selectClass(this.classes[0].id);
                }
            } catch (err) {
                ui.addNotification({ type: 'error', message: 'Failed to load classes' });
            } finally {
                ui.setLoading(false);
            }
        },

        selectClass(classId: number) {
            this.selectedClassId = classId;
            this.fetchStudents(classId);
            const subjectId = this.classSubjectMap[classId];
            if (subjectId) {
                this.fetchMarks(classId, subjectId);
            }
        },

        async fetchStudents(classId: number) {
            const ui = useUIStore();
            // Don't block UI fully, maybe just grid loader
            try {
                const res = await teacherApi.getStudents(classId);
                this.students = res.data;
            } catch (err) {
                ui.addNotification({ type: 'error', message: 'Failed to load students' });
            }
        },

        async fetchMarks(classId: number, subjectId: number) {
            try {
                const res = await teacherApi.getMarks(classId, subjectId);
                this.marks = res.data;
            } catch (err) {
                console.error("Error fetching marks", err);
            }
        },

        async saveMark(markData: any) {
            const ui = useUIStore();
            try {
                const res = await teacherApi.saveMark(markData);
                this.marks.push(res.data);
                ui.addNotification({ type: 'success', message: 'Mark saved' });
                return true;
            } catch (err) {
                ui.addNotification({ type: 'error', message: 'Failed to save mark' });
                return false;
            }
        },

        async fetchAbsences(classId: number) {
            const ui = useUIStore();
            // ui.setLoading(true); // maybe subtle loading
            try {
                const res = await teacherApi.getAbsences(classId);
                this.absences = res.data;
            } catch (e) {
                console.error(e);
                ui.addNotification({ type: 'error', message: 'Failed to fetch absences' });
            } finally {
                // ui.setLoading(false);
            }
        },

        async saveAbsence(absenceData: any) {
            const ui = useUIStore();
            try {
                const res = await teacherApi.saveAbsence(absenceData);
                this.absences.push(res.data);
                // ui.addNotification({ type: 'success', message: 'Absence saved' }); 
                // Silent success for quick marking? Or toast
                return true;
            } catch (e) {
                ui.addNotification({ type: 'error', message: 'Failed to save absence' });
                return false;
            }
        }
    }
});
