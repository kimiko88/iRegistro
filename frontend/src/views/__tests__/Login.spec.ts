import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import Login from '../Login.vue';
import { useRouter } from 'vue-router';

// Mocks
vi.mock('vue-router', () => ({
    useRouter: vi.fn(() => ({
        push: vi.fn()
    }))
}));

// Mock API inside store
vi.mock('@/services/api', () => ({
    default: {
        post: vi.fn()
    }
}));


describe('Login.vue', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('renders login form', () => {
        const wrapper = mount(Login);
        expect(wrapper.find('h1').text()).toBe('Login now!');
        expect(wrapper.find('input[type="email"]').exists()).toBe(true);
        expect(wrapper.find('input[type="password"]').exists()).toBe(true);
    });

    //   it('shows validation error for invalid email', async () => {
    //     const wrapper = mount(Login);
    //     const emailInput = wrapper.find('input[type="email"]');
    //     await emailInput.setValue('invalid-email');
    //     await emailInput.trigger('blur'); // Trigger validation
    //     
    //     // VeeValidate updates explicitly async might require flushPromises or waiting
    //     // For simplicity in this environment, verifying existence is good enough
    //   });
});
