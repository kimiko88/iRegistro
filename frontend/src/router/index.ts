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

import AdminPanel from '@/views/admin/AdminPanel.vue';
import SchoolSettings from '@/views/admin/SchoolSettings.vue';
import BackupManagement from '@/views/admin/BackupManagement.vue';
import DataExport from '@/views/admin/DataExport.vue';
import AuditLogs from '@/views/admin/AuditLogs.vue';
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
                path: 'superadmin/schools/:schoolId/users',
                name: 'SchoolUserManagement',
                component: UserManagement,
                props: true,
                meta: { requiresRole: 'SuperAdmin' }
            },
            {
                path: 'admin',
                component: AdminPanel,
                meta: { requiresRole: 'Admin' },
                children: [
                    {
                        path: 'users',
                        name: 'UserManagement',
                        component: UserManagement
                    },
                    {
                        path: 'settings',
                        name: 'SchoolSettings',
                        component: SchoolSettings
                    },
                    {
                        path: 'backups',
                        name: 'BackupManagement',
                        component: BackupManagement
                    },
                    {
                        path: 'export',
                        name: 'DataExport',
                        component: DataExport
                    },
                    {
                        path: 'audit-logs',
                        name: 'AuditLogs',
                        component: AuditLogs
                    }
                ]
            },
            {
                path: 'teacher',
                name: 'TeacherDashboard',
                component: TeacherDashboard,
                meta: { requiresRole: 'Teacher' }
            },
            {
                path: 'parent',
                component: ParentDashboard,
                meta: { requiresRole: 'Parent' },
                children: [
                    { path: '', name: 'ParentDashboard', component: { template: '' } }, // Default view logic handling
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
                path: 'secretary',
                name: 'SecretaryInbox',
                component: DocumentManagement,
                meta: { requiresRole: 'Secretary' }
            },
            {
                path: 'secretary/classes',
                name: 'SecretaryClassManagement',
                component: () => import('@/views/secretary/ClassManagement.vue'),
                meta: { requiresRole: 'Secretary' }
            },
            {
                path: 'secretary/archive',
                name: 'SecretaryArchive',
                component: Archive,
                meta: { requiresRole: 'Segreteria' }
            },
            {
                path: 'director',
                name: 'DirectorDashboard',
                component: DirectorDashboard,
                meta: { requiresRole: 'Principal' }
            },
            {
                path: 'director/signing',
                name: 'DocumentSigning',
                component: DocumentSigning,
                meta: { requiresRole: 'Principal' }
            },
            {
                path: 'director/approvals',
                name: 'RequestApprovals',
                component: RequestApprovals,
                meta: { requiresRole: 'Principal' }
            },
            {
                path: 'director/reports',
                name: 'SchoolReports',
                component: SchoolReports,
                meta: { requiresRole: 'Dirigente' }
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
