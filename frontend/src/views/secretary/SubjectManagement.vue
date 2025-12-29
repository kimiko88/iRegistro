<template>
  <div class="p-6 space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-3xl font-bold">Gestione Materie</h1>
      <button @click="showModal = true" class="btn btn-primary gap-2">
        <Plus class="w-5 h-5" />
        Nuova Materia
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
            <p class="text-sm text-gray-600" v-if="subject.competition_class">Classe di concorso: {{ subject.competition_class }}</p>
            <div class="card-actions justify-end mt-4">
              <button @click="editSubject(subject)" class="btn btn-sm btn-ghost">
                <Edit class="w-4 h-4" /> Modifica
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Subject Modal -->
    <FormModal
      :isOpen="showModal"
      :title="editingSubject ? 'Modifica Materia' : 'Crea Nuova Materia'"
      :submitLabel="editingSubject ? 'Salva' : 'Crea'"
      :loading="saving"
      @close="closeModal"
      @submit="handleSubmit"
    >
      <div class="space-y-4">
        <div class="form-control">
          <label class="label"><span class="label-text">Nome Materia</span></label>
          <input type="text" v-model="formData.name" class="input input-bordered" placeholder="es. Matematica" required />
        </div>
        
        <div class="form-control">
          <label class="label"><span class="label-text">Codice</span></label>
          <input type="text" v-model="formData.code" class="input input-bordered" placeholder="es. MAT01" required />
        </div>
        
        <div class="form-control">
          <label class="label"><span class="label-text">Classe di Concorso</span></label>
          <input type="text" v-model="formData.competition_class" class="input input-bordered" placeholder="es. A026, A027" required />
          <label class="label"><span class="label-text-alt">Codice della classe di concorso ministeriale</span></label>
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
import { Plus, Edit } from 'lucide-vue-next';

const authStore = useAuthStore();
const notificationStore = useNotificationStore();

const subjects = ref<any[]>([]);
const loading = ref(false);
const saving = ref(false);
const showModal = ref(false);
const editingSubject = ref<any>(null);

const formData = ref({
  name: '',
  code: '',
  competition_class: ''
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

const editSubject = (subject: any) => {
  editingSubject.value = subject;
  formData.value = { ...subject };
  showModal.value = true;
};

const closeModal = () => {
  showModal.value = false;
  editingSubject.value = null;
  formData.value = { name: '', code: '', competition_class: '' };
};

const handleSubmit = async () => {
  saving.value = true;
  try {
    const schoolId = authStore.user?.schoolId;
    if (!schoolId) {
      notificationStore.error('School ID non trovato');
      return;
    }
    
    if (editingSubject.value) {
      // Update existing subject
      await api.put(`/schools/${schoolId}/subjects/${editingSubject.value.id}`, formData.value);
      notificationStore.success('Materia aggiornata con successo');
    } else {
      // Create new subject
      await api.post(`/schools/${schoolId}/subjects`, formData.value);
      notificationStore.success('Materia creata con successo');
    }
    closeModal();
    fetchSubjects();
  } catch (err) {
    notificationStore.error(editingSubject.value ? 'Errore nell\'aggiornamento' : 'Errore nella creazione');
  } finally {
    saving.value = false;
  }
};

onMounted(() => {
  fetchSubjects();
});
</script>
