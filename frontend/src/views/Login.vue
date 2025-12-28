<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useUIStore } from '@/stores/ui';
import { useRouter } from 'vue-router';
import { useForm } from 'vee-validate';
import * as yup from 'yup';
import { LogIn, Mail, Lock, School } from 'lucide-vue-next';

const authStore = useAuthStore();
const uiStore = useUIStore();
const router = useRouter();

const schema = yup.object({
  email: yup.string().required('Email richiesta').email('Email non valida'),
  password: yup.string().required('Password richiesta').min(6, 'Password deve essere almeno 6 caratteri'),
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
    uiStore.addNotification({ type: 'success', message: 'Accesso effettuato con successo' });
    router.push('/dashboard');
  } else {
    uiStore.addNotification({ type: 'error', message: 'Credenziali non valide' });
  }
});
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-purple-50 dark:from-gray-900 dark:via-gray-800 dark:to-indigo-950 flex items-center justify-center p-4">
    <div class="w-full max-w-6xl grid lg:grid-cols-2 gap-8 items-center">
      <!-- Left Side - Welcome Message -->
      <div class="hidden lg:flex flex-col justify-center p-12 space-y-6">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-16 h-16 bg-indigo-600 rounded-2xl flex items-center justify-center shadow-lg shadow-indigo-200 dark:shadow-none">
            <School class="w-10 h-10 text-white" />
          </div>
          <div>
            <h2 class="text-3xl font-bold text-gray-800 dark:text-gray-100">iRegistro</h2>
            <p class="text-sm text-gray-600 dark:text-gray-400">Registro Elettronico</p>
          </div>
        </div>
        
        <h1 class="text-5xl font-bold text-gray-900 dark:text-white leading-tight">
          Benvenuto nel tuo
          <span class="bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">Registro Digitale</span>
        </h1>
        
        <p class="text-lg text-gray-600 dark:text-gray-300">
          Accedi alla piattaforma per gestire la tua scuola in modo semplice ed efficace.
        </p>

        <div class="grid grid-cols-2 gap-4 pt-6">
          <div class="flex items-start gap-3 p-4 bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700">
            <div class="w-10 h-10 rounded-full bg-green-100 dark:bg-green-900/30 flex items-center justify-center shrink-0">
              <svg class="w-5 h-5 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
              </svg>
            </div>
            <div>
              <h3 class="font-semibold text-gray-900 dark:text-white text-sm">Sicuro</h3>
              <p class="text-xs text-gray-500 dark:text-gray-400">Dati protetti</p>
            </div>
          </div>
          
          <div class="flex items-start gap-3 p-4 bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700">
            <div class="w-10 h-10 rounded-full bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center shrink-0">
              <svg class="w-5 h-5 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
              </svg>
            </div>
            <div>
              <h3 class="font-semibold text-gray-900 dark:text-white text-sm">Veloce</h3>
              <p class="text-xs text-gray-500 dark:text-gray-400">Accesso rapido</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Side - Login Form -->
      <div class="w-full">
        <div class="bg-white dark:bg-gray-800 rounded-3xl shadow-2xl shadow-gray-200/50 dark:shadow-none border border-gray-100 dark:border-gray-700 p-8 md:p-12">
          <!-- Mobile Logo -->
          <div class="lg:hidden flex items-center justify-center gap-3 mb-8">
            <div class="w-12 h-12 bg-indigo-600 rounded-xl flex items-center justify-center">
              <School class="w-7 h-7 text-white" />
            </div>
            <div>
              <h2 class="text-2xl font-bold text-gray-800 dark:text-gray-100">iRegistro</h2>
            </div>
          </div>

          <div class="mb-8">
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Accedi al tuo account</h2>
            <p class="text-gray-600 dark:text-gray-400">Inserisci le tue credenziali per continuare</p>
          </div>

          <form @submit.prevent="onSubmit" class="space-y-6">
            <!-- Email Field -->
            <div class="form-control">
              <label class="label">
                <span class="label-text font-medium text-gray-700 dark:text-gray-300">Email</span>
              </label>
              <div class="relative">
                <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                  <Mail class="w-5 h-5 text-gray-400" />
                </div>
                <input
                  v-model="email"
                  v-bind="emailAttrs"
                  type="email"
                  placeholder="mario.rossi@scuola.it"
                  class="input input-bordered w-full pl-12 bg-gray-50 dark:bg-gray-700 border-gray-200 dark:border-gray-600 focus:border-indigo-500 dark:focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 dark:focus:ring-indigo-900/30"
                  :class="{ 'input-error border-red-500': errors.email }"
                />
              </div>
              <label v-if="errors.email" class="label">
                <span class="label-text-alt text-error">{{ errors.email }}</span>
              </label>
            </div>

            <!-- Password Field -->
            <div class="form-control">
              <label class="label">
                <span class="label-text font-medium text-gray-700 dark:text-gray-300">Password</span>
              </label>
              <div class="relative">
                <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                  <Lock class="w-5 h-5 text-gray-400" />
                </div>
                <input
                  v-model="password"
                  v-bind="passwordAttrs"
                  type="password"
                  placeholder="••••••••"
                  class="input input-bordered w-full pl-12 bg-gray-50 dark:bg-gray-700 border-gray-200 dark:border-gray-600 focus:border-indigo-500 dark:focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 dark:focus:ring-indigo-900/30"
                  :class="{ 'input-error border-red-500': errors.password }"
                />
              </div>
              <label v-if="errors.password" class="label">
                <span class="label-text-alt text-error">{{ errors.password }}</span>
              </label>
              <label class="label">
                <a href="#" class="label-text-alt link link-hover text-indigo-600 dark:text-indigo-400 hover:text-indigo-700 dark:hover:text-indigo-300">Password dimenticata?</a>
              </label>
            </div>

            <!-- Submit Button -->
            <div class="form-control mt-8">
              <button 
                type="submit"
                class="btn btn-primary w-full bg-indigo-600 hover:bg-indigo-700 border-none text-white shadow-lg shadow-indigo-200 dark:shadow-none h-12 text-base font-semibold" 
                :disabled="uiStore.isLoading"
              >
                <span v-if="uiStore.isLoading" class="loading loading-spinner"></span>
                <LogIn v-else class="w-5 h-5 mr-2" />
                {{ uiStore.isLoading ? 'Accesso in corso...' : 'Accedi' }}
              </button>
            </div>
          </form>

          <!-- Divider -->
          <div class="divider text-gray-400 dark:text-gray-500 text-sm my-8">oppure</div>

          <!-- SPID/CIE Buttons -->
          <div class="grid grid-cols-2 gap-4">
            <button class="btn btn-outline border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 gap-2">
              <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2L2 7v10c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V7l-10-5z"/>
              </svg>
              <span class="text-sm">SPID</span>
            </button>
            <button class="btn btn-outline border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 gap-2">
              <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
                <path d="M21 16v-2l-8-5V3.5c0-.83-.67-1.5-1.5-1.5S10 2.67 10 3.5V9l-8 5v2l8-2.5V19l-2 1.5V22l3.5-1 3.5 1v-1.5L13 19v-5.5l8 2.5z"/>
              </svg>
              <span class="text-sm">CIE</span>
            </button>
          </div>

          <!-- Footer Note -->
          <p class="text-center text-xs text-gray-500 dark:text-gray-400 mt-8">
            Accedendo, accetti i nostri 
            <a href="#" class="link link-hover text-indigo-600 dark:text-indigo-400">Termini di Servizio</a> e 
            <a href="#" class="link link-hover text-indigo-600 dark:text-indigo-400">Privacy Policy</a>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
