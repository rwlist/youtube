import { createApp } from 'vue'
import App from './App.vue'
import { useToast } from 'vue-toast-notification'
import 'vue-toast-notification/dist/theme-sugar.css'
import { updateUser } from './stores/user'

export const app = createApp(App)

export const toast = useToast({})
app.provide('$toast', toast)

export const root = app.mount('#app')

updateUser()
