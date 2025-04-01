import './assets/main.css'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import pocketbase from './pocketbase'

const app = createApp(App)
app.provide('$pocketbase', pocketbase)
app.use(router)
app.mount('#app')
