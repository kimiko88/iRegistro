import { defineStore } from 'pinia';
import teacherApi from '@/services/teacher';
import { useUIStore } from './ui';

export const useTeacherStore = defineStore('teacher', {
    state: () => ({
        classes: [] as any[],
        selectedClassId: null as number | null,
        currentClass: null as any | null,
        students: [] as any[],
        marks: [] as any[],
        schedule: [] as any[],
    }),

    getters: {
        selectedClass: (state) => state.classes.find(c => c.id === state.selectedClassId),
    },

    actions: {
        async fetchClasses() {
            const ui = useUIStore();
            ui.setLoading(true);
            try {
                const res = await teacherApi.getClasses();

                // Backend returns Assignment = { class: {...}, subject: {...} }
                // We Map it so Select works
                this.classes = res.data.map((a: any) => ({
                    id: a.class.id,
                    name: `${a.class.grade}${a.class.section}`,
                    subject: a.subject.name,
                    subjectId: a.subject.id
                }));

                if (!this.selectedClassId && this.classes.length > 0) {
                    this.selectedClassId = this.classes[0].id;
                    this.fetchStudents(this.classes[0].id);
                }
            } finally {
                ui.setLoading(false);
            }
        },
        selectClass(id: number) {
            this.selectedClassId = id;
            this.fetchStudents(id);
            const cls = this.classes.find(c => c.id === id);
            if (cls && cls.subjectId) {
                this.fetchMarks(cls.subjectId);
            }
        },
        async fetchStudents(classId: number) {
            const ui = useUIStore();
            ui.setLoading(true);
            try {
                const res = await teacherApi.getStudents(classId);
                this.students = res.data;
            } finally {
                ui.setLoading(false);
            }
        },
        async fetchMarks(subjectId: number) {
            if (!this.selectedClassId) return;
            try {
                const res = await teacherApi.getMarks(this.selectedClassId, subjectId);
                this.marks = res.data;
            } catch (e) {
                console.error("Failed to fetch marks");
            }
        },
        async saveMark(mark: any) {
            const ui = useUIStore();
            try {
                const res = await teacherApi.saveMark(mark);
                this.marks.push(res.data);
            } catch (e) {
                ui.addNotification({ type: 'error', message: 'Failed to save mark' });
            }
        }
    },
});
