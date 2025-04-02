import { createApp } from 'vue'
import App from './App.vue'
import ModalPlugin from './plugins/modal'
import router from './router/router'
import './style.css'

const app = createApp(App)
app.use(router)
app.use(ModalPlugin)
app.mount('#app')
