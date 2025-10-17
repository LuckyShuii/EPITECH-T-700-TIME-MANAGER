import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import '@/assets/css/tailwind.css';
import { setupCalendar, Calendar } from 'v-calendar'
import 'v-calendar/style.css'
import { createPinia } from 'pinia';
import { useAuthStore } from '@/store/AuthStore';  // ✅ AJOUT

const pinia = createPinia();
const app = createApp(App);

app.use(pinia);

//Initialiser l'auth AVANT le router
const authStore = useAuthStore();

authStore.initAuth().finally(() => {
    app.use(router);
    app.use(setupCalendar, {})
    app.component('VCalendar', Calendar)
    app.mount('#app');
});