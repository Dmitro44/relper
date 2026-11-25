import { createApp } from 'vue'
import 'maplibre-gl/dist/maplibre-gl.css'
import './style.css'
import App from './App.vue'
import { initTheme } from './theme'

initTheme()
createApp(App).mount('#app')
