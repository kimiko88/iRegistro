import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router'

import { useAuthStore } from './stores/auth'

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)

const authStore = useAuthStore()
if (authStore.token) {
    authStore.fetchUser().finally(() => {
        app.mount('#app')
    })
} else {
    app.mount('#app')
}
