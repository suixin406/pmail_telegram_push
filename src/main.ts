import { createApp } from 'vue'
import App from './App.vue'
import 'element-plus/dist/index.css'
import '@/assets/fonts/iconfont.css'

if (import.meta.env.DEV) {
  await import('@/mock')
}

createApp(App).mount('#app')
