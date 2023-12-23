import './assets/main.css'
import 'vant/lib/index.css';

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import {ActionSheet, Toast, FloatingPanel, Dialog } from 'vant';
const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ActionSheet);
app.use(Toast);
app.use(FloatingPanel);
app.use(Dialog);
app.mount('#app')
