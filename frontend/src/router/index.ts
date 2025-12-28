import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import Layout from '@/components/common/Layout.vue';
import Login from '@/views/Login.vue';
import Dashboard from '@/views/Dashboard.vue';
import SuperAdminDashboard from '@/views/admin/SuperAdminDashboard.vue';
import UserManagement from '@/views/admin/UserManagement.vue';
import ParentDashboard from '@/views/parent/ParentDashboard.vue';
import StudentDashboard from '@/views/student/StudentDashboard.vue';
import DocumentManagement from '@/views/secretary/DocumentManagement.vue';
import Archive from '@/views/secretary/Archive.vue';

import MarksView from '@/views/marks/MarksView.vue';
import AbsencesView from '@/views/absences/AbsencesView.vue';
import MessagesView from '@/views/messages/MessagesView.vue';
import ColloquiumView from '@/views/colloquiums/ColloquiumView.vue';

import DirectorDashboard from '@/views/director/DirectorDashboard.vue';
import DocumentSigning from '@/views/director/DocumentSigning.vue';
import RequestApprovals from '@/views/director/RequestApprovals.vue';
import SchoolReports from '@/views/director/SchoolReports.vue';

import TeacherDashboard from '@/views/teacher/ClassDashboard.vue';

const SchoolSettings = { template: '<div>Settings</div>' }; // Stub if missing
// ClassDashboard stub removed

// For simplicity, defining explicit routes for parent sub-views inside parent children

const routes: Array<RouteRecordRaw> = [
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: { guest: true }
    },
    {
        path: '/',
        component: Layout,
        meta: { requiresAuth: true },
        children: [
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: Dashboard,
                alias: ''
            },
            {
                path: 'superadmin',
                name: 'SuperAdminDashboard',
                component: SuperAdminDashboard,
                meta: { requiresRole: 'SuperAdmin' }
            },
            {
                path: 'admin/users',
                name: 'UserManagement',
                component: UserManagement,
                meta: { requiresRole: 'Admin' }
            },
            {
                path: 'admin/settings',
                name: 'SchoolSettings',
                component: SchoolSettings,
                meta: { requiresRole: 'Admin' }
            },
            {
                path: 'teacher',
                name: 'TeacherDashboard',
                component: TeacherDashboard,
                meta: { requiresRole: 'Insegnante' } // Updated role to match backend seed (Insegnante)
            },
            {
                path: 'parent',
                component: ParentDashboard,
                meta: { requiresRole: 'Parent' },
                children: [
                    { path: '', name: 'ParentDashboard', component: { template: '' } }, // Default view logic handled in dashboard
                    { path: 'marks', component: MarksView },
                    { path: 'absences', component: AbsencesView },
                    { path: 'messages', component: MessagesView },
                    { path: 'colloquiums', component: ColloquiumView },
                    { path: 'documents', component: Archive }
                ]
            },
            {
                path: 'student',
                component: StudentDashboard,
                meta: { requiresRole: 'Student' },
                children: [
                    { path: '', name: 'StudentDashboard', component: { template: '' } },
                    { path: 'marks', component: MarksView },
                    { path: 'absences', component: AbsencesView },
                    { path: 'documents', component: Archive }
                ]
            },
            {
                path: 'secretary',
                name: 'SecretaryInbox',
                component: DocumentManagement,
                meta: { requiresRole: 'Secretary' }
            },
            {
                path: 'secretary/archive',
                name: 'SecretaryArchive',
                component: Archive,
                meta: { requiresRole: 'Secretary' }
            },
            {
                path: 'director',
                name: 'DirectorDashboard',
                component: DirectorDashboard,
                meta: { requiresRole: 'Director' }
            },
            {
                path: 'director/signing',
                name: 'DocumentSigning',
                component: DocumentSigning,
                meta: { requiresRole: 'Director' }
            },
            {
                path: 'director/approvals',
                name: 'RequestApprovals',
                component: RequestApprovals,
                meta: { requiresRole: 'Director' }
            },
            {
                path: 'director/reports',
                name: 'SchoolReports',
                component: SchoolReports,
                meta: { requiresRole: 'Director' }
            },
            // Add other routes here
        ]
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((to, _from, next) => {
    const auth = useAuthStore();

    if (to.meta.requiresAuth && !auth.isAuthenticated) {
        next('/login');
    } else if (to.meta.guest && auth.isAuthenticated) {
        next('/dashboard');
    } else {
        next();
    }
});

export default router;
