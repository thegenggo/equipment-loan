<script setup lang="ts">
import { RouterLink, RouterView, useRouter } from 'vue-router';

import { useAuthStore } from './stores/auth';

const auth = useAuthStore()
const router = useRouter()

function onLogout(): void {
  auth.logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <header v-if="auth.isAuthenticated">
    <nav>
      <RouterLink :to="{ name: 'equipments' }">อุปกรณ์</RouterLink>
      <RouterLink :to="{ name: 'my-loans' }">คำขอของฉัน</RouterLink>
      <template v-if="auth.isAdmin">
        <RouterLink :to="{ name: 'admin-loans' }">อนุมัติคำขอ</RouterLink>
        <RouterLink :to="{ name: 'admin-equipments' }">จัดการอุปกรณ์</RouterLink>
      </template>
    </nav>
    <button type="button" @click="onLogout">ออกจากระบบ</button>
  </header>

  <RouterView />
</template>
