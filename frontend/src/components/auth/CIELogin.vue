<template>
  <div class="cie-login">
    <button @click="loginWithCIE" class="btn btn-cie">
      <img src="@/assets/cie-icon.svg" alt="CIE" class="cie-icon" />
      <span>Entra con CIE</span>
    </button>
    
    <p class="cie-info">
      <small>Autenticati con la tua Carta d'Identità Elettronica</small>
    </p>
  </div>
</template>

<script setup lang="ts">
interface Props {
  schoolId: string
}

const props = defineProps<Props>()

const loginWithCIE = () => {
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const redirectUri = `${window.location.origin}/auth-callback`
  
  const authUrl = `${apiUrl}/auth/cie/login?school_id=${props.schoolId}&redirect_uri=${encodeURIComponent(redirectUri)}`
  
  // Redirect to CIE authentication
  window.location.href = authUrl
}
</script>

<style scoped>
.cie-login {
  margin: 1rem 0;
}

.btn-cie {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1.5rem;
  background: #0073E6;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.3s;
  width: 100%;
  justify-content: center;
}

.btn-cie:hover {
  background: #005BB5;
}

.cie-icon {
  width: 24px;
  height: 24px;
}

.cie-info {
  margin-top: 0.5rem;
  text-align: center;
  color: #666;
}

.cie-info small {
  font-size: 0.875rem;
}
</style>
