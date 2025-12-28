<template>
  <div class="spid-login">
    <button @click="loginWithSPID" class="btn btn-spid">
      <img src="@/assets/spid-icon.svg" alt="SPID" class="spid-icon" />
      <span>Accedi con SPID</span>
    </button>
    
    <!-- Optional: Show provider selection -->
    <div v-if="showProviders" class="spid-providers">
      <h3>Scegli il tuo Identity Provider</h3>
      <div class="provider-grid">
        <button
          v-for="provider in spidProviders"
          :key="provider.id"
          @click="loginWithProvider(provider.id)"
          class="provider-btn"
        >
          {{ provider.name }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

interface Props {
  schoolId: string
  showProviderSelection?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showProviderSelection: false
})

const router = useRouter()
const showProviders = ref(props.showProviderSelection)

const spidProviders = [
  { id: 'aruba', name: 'Aruba ID' },
  { id: 'infocert', name: 'Infocert ID' },
  { id: 'poste', name: 'Poste ID' },
  { id: 'tim', name: 'TIM ID' },
  { id: 'lepida', name: 'Lepida ID' },
  { id: 'sielteid', name: 'Sielte ID' },
  { id: 'intesa', name: 'Intesa ID' },
  { id: 'register', name: 'Register SPID' }
]

const loginWithSPID = () => {
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const redirectUri = `${window.location.origin}/auth-callback`
  
  const authUrl = `${apiUrl}/auth/spid/login?school_id=${props.schoolId}&redirect_uri=${encodeURIComponent(redirectUri)}`
  
  // Redirect to SPID authentication
  window.location.href = authUrl
}

const loginWithProvider = (providerId: string) => {
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const redirectUri = `${window.location.origin}/auth-callback`
  
  const authUrl = `${apiUrl}/auth/spid/login?school_id=${props.schoolId}&provider=${providerId}&redirect_uri=${encodeURIComponent(redirectUri)}`
  
  window.location.href = authUrl
}
</script>

<style scoped>
.spid-login {
  margin: 1rem 0;
}

.btn-spid {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1.5rem;
  background: #0066CC;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.3s;
}

.btn-spid:hover {
  background: #0052A3;
}

.spid-icon {
  width: 24px;
  height: 24px;
}

.spid-providers {
  margin-top: 1.5rem;
  padding: 1rem;
  border: 1px solid #E0E0E0;
  border-radius: 8px;
}

.spid-providers h3 {
  margin-bottom: 1rem;
  font-size: 1.125rem;
  color: #333;
}

.provider-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 0.75rem;
}

.provider-btn {
  padding: 0.75rem;
  border: 1px solid #0066CC;
  background: white;
  color: #0066CC;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 0.875rem;
}

.provider-btn:hover {
  background: #0066CC;
  color: white;
}
</style>
