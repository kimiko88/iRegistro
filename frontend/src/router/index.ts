import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import Layout from '@/components/common/Layout.vue';
import Login from '@/views/Login.vue';
import Dashboard from '@/views/Dashboard.vue';
import SuperAdminDashboard from '@/views/admin/SuperAdminDashboard.vue';
import UserManagement from '@/views/admin/UserManagement.vue';
import ParentDashboard from '@/views/parent/ParentDashboard.vue';
import StudentDashboard from '@/views/student/StudentDashboard.vue';
import AbsencesView from '@/views/absences/AbsencesView.vue'; // Full view wrapper needed? Or link directly 

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
                component: ClassDashboard,
                meta: { requiresRole: 'Teacher' } // In real map this to 'Insegnante'
            },
            {
                path: 'parent',
                name: 'ParentDashboard',
                component: ParentDashboard,
                meta: { requiresRole: 'Parent' }
            },
            {
                path: 'student',
                name: 'StudentDashboard',
                component: StudentDashboard,
                meta: { requiresRole: 'Student' }
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
