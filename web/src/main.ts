import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth.ts'
import { setUnauthorizedHandler } from './api/client.ts'

const app = createApp(App)

app.use(createPinia())

const auth = useAuthStore()
auth.restore()

setUnauthorizedHandler(() => {
  auth.logout()
  router.push({ name: 'login' })
})

app.use(router)

app.mount('#app')
