<template>
  <div class="p-6 space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-3xl font-bold">Gestione Materie</h1>
      <button @click="showCreateModal = true" class="btn btn-primary gap-2">
        <Plus class="w-5 h-5" />
        Nuova Mat eria
      </button>
    </div>

    <div class="bg-base-100 p-4 rounded-lg shadow">
      <div v-if="loading" class="text-center py-8">Caricamento...</div>
      <div v-else-if="subjects.length === 0" class="text-center py-8 text-gray-500">
        Nessuna materia creata. Clicca "Nuova Materia" per iniziare.
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="subject in subjects" :key="subject.id" 
             class="card bg-base-200 shadow-sm hover:shadow-md transition-shadow">
          <div class="card-body">
            <h3 class="card-title text-lg">{{ subject.name }}</h3>
            <p class="text-sm text-gray-600">Codice: {{ subject.code }}</p>
            <p class="text-sm text-gray-600">Ore settimanali: {{ subject.hours_per_week }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Subject Modal -->
    <FormModal
      :isOpen="showCreateModal"
      title="Crea Nuova Materia"
      submitLabel="Crea"
      :loading="creating"
      @close="showCreateModal = false"
      @submit="handleCreate"
    >
      <div class="space-y-4">
        <div class="form-control">
          <label class="label"><span class="label-text">Nome Materia</span></label>
          <input type="text" v-model="newSubject.name" class="input input-bordered" placeholder="es. Matematica" required />
        </div>
        
        <div class="form-control">
          <label class="label"><span class="label-text">Codice</span></label>
          <input type="text" v-model="newSubject.code" class="input input-bordered" placeholder="es. MAT01" required />
        </div>
        
        <div class="form-control">
          <label class="label"><span class="label-text">Ore Settimanali</span></label>
          <input type="number" v-model.number="newSubject.hours_per_week" class="input input-bordered" min="1" max="10" required />
        </div>
      </div>
    </FormModal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useNotificationStore } from '@/stores/notification';
import api from '@/services/api';
import FormModal from '@/components/shared/FormModal.vue';
import { Plus } from 'lucide-vue-next';

const authStore = useAuthStore();
const notificationStore = useNotificationStore();

const subjects = ref<any[]>([]);
const loading = ref(false);
const creating = ref(false);
const showCreateModal = ref(false);

const newSubject = ref({
  name: '',
  code: '',
  hours_per_week: 3
});

const fetchSubjects = async () => {
  loading.value = true;
  try {
    const schoolId = authStore.user?.schoolId;
    if (!schoolId) {
      notificationStore.error('School ID non trovato');
      return;
    }
    const res = await api.get(`/schools/${schoolId}/subjects`);
    subjects.value = res.data || [];
  } catch (err) {
    notificationStore.error('Errore nel caricamento delle materie');
  } finally {
    loading.value = false;
  }
};

const handleCreate = async () => {
  creating.value = true;
  try {
    const schoolId = authStore.user?.schoolId;
    if (!schoolId) {
      notificationStore.error('School ID non trovato');
      return;
    }
    
    await api.post(`/schools/${schoolId}/subjects`, newSubject.value);
    notificationStore.success('Materia creata con successo');
    showCreateModal.value = false;
    newSubject.value = { name: '', code: '', hours_per_week: 3 };
    fetchSubjects();
  } catch (err) {
    notificationStore.error('Errore nella creazione della materia');
  } finally {
    creating.value = false;
  }
};

onMounted(() => {
  fetchSubjects();
});
</script>
