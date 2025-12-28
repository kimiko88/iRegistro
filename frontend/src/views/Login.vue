<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useUIStore } from '@/stores/ui';
import { useRouter } from 'vue-router';
import { useForm } from 'vee-validate';
import * as yup from 'yup';

const authStore = useAuthStore();
const uiStore = useUIStore();
const router = useRouter();

const schema = yup.object({
  email: yup.string().required().email(),
  password: yup.string().required().min(6),
});

const { handleSubmit, errors, defineField } = useForm({
  validationSchema: schema,
});

const [email, emailAttrs] = defineField('email');
const [password, passwordAttrs] = defineField('password');

const onSubmit = handleSubmit(async (values) => {
  uiStore.setLoading(true);
  const success = await authStore.login(values);
  uiStore.setLoading(false);

  if (success) {
    uiStore.addNotification({ type: 'success', message: 'Login successful' });
    router.push('/dashboard');
  } else {
    uiStore.addNotification({ type: 'error', message: 'Invalid credentials' });
  }
});
</script>

<template>
  <div class="hero min-h-screen bg-base-200">
    <div class="hero-content flex-col lg:flex-row-reverse">
      <div class="text-center lg:text-left">
        <h1 class="text-5xl font-bold">Login now!</h1>
        <p class="py-6">Access your school management dashboard securely.</p>
      </div>
      <div class="card shrink-0 w-full max-w-sm shadow-2xl bg-base-100">
        <form class="card-body" @submit.prevent="onSubmit">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Email</span>
            </label>
            <input
              v-model="email"
              v-bind="emailAttrs"
              type="email"
              placeholder="email"
              class="input input-bordered"
              :class="{ 'input-error': errors.email }"
            />
            <label v-if="errors.email" class="label">
              <span class="label-text-alt text-error">{{ errors.email }}</span>
            </label>
          </div>
          <div class="form-control">
            <label class="label">
              <span class="label-text">Password</span>
            </label>
            <input
              v-model="password"
              v-bind="passwordAttrs"
              type="password"
              placeholder="password"
              class="input input-bordered"
              :class="{ 'input-error': errors.password }"
            />
            <label v-if="errors.password" class="label">
              <span class="label-text-alt text-error">{{ errors.password }}</span>
            </label>
            <label class="label">
              <a href="#" class="label-text-alt link link-hover">Forgot password?</a>
            </label>
          </div>
          <div class="form-control mt-6">
            <button class="btn btn-primary" :disabled="uiStore.isLoading">
              <span v-if="uiStore.isLoading" class="loading loading-spinner"></span>
              Login
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
