import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './assets/main.css'

const app = createApp(App)

// Initialize Pinia
const pinia = createPinia()

app.use(pinia)

app.use(router)

app.mount('#app')
