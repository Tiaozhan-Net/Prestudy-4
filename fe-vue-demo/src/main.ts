import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'
import router from './router'
import { initializeTheme } from './composables/theme'
import { loadCurrentUser } from './composables/auth'
import './assets/main.css'
import './assets/dark.css'

const bootstrap = async () => {
  initializeTheme()

  const app = createApp(App)
  const pinia = createPinia()
  app.use(pinia)

  // 登录态从后端 Cookie + /auth/profile 恢复，不再依赖本地假用户
  await loadCurrentUser()

  for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
  }

  app.use(ElementPlus)
  app.use(router)
  app.mount('#app')
}

bootstrap()
