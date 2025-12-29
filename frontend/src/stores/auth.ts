import { defineStore } from 'pinia';
import api from '@/services/api';

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: null as any | null,
        token: localStorage.getItem('token') || '',
        isAuthenticated: !!localStorage.getItem('token'),
    }),

    actions: {
        async login(credentials: any) {
            try {
                const response = await api.post('/auth/login', credentials);
                this.token = response.data.access_token;
                this.isAuthenticated = true;
                this.user = { role: response.data.role }; // basic simulation
                localStorage.setItem('token', this.token);
                return true;
            } catch (e) {
                return false;
            }
        },
        async fetchUser() {
            if (!this.token) return;
            try {
                const res = await api.get('/auth/me');
                this.user = res.data;
                this.isAuthenticated = true;
            } catch (e) {
                this.logout();
            }
        },
        logout() {
            this.token = '';
            this.user = null;
            this.isAuthenticated = false;
            localStorage.removeItem('token');
        },
    },
});
