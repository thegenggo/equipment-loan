<script setup lang="ts">
  import { RouterLink, useRouter } from 'vue-router'

  import { useAuthStore } from '@/stores/auth';

  const auth = useAuthStore()
  const router = useRouter()

  function onLogout(): void {
    auth.logout()
    router.push({ name: 'login' })
  }
</script>

<template>
  <header v-if="auth.isAuthenticated" class="layout__header">
    <nav>
      <RouterLink :to="{ name: 'equipments' }">อุปกรณ์</RouterLink>
      <RouterLink :to="{ name: 'my-loans' }">{{ auth.isAdmin ? 'คำขอทั้งหมด' : 'คำขอของฉัน' }}</RouterLink>
      <template v-if="auth.isAdmin">
        <RouterLink :to="{ name: 'admin-loans' }">อนุมัติคำขอ</RouterLink>
        <RouterLink :to="{ name: 'admin-equipments' }">จัดการอุปกรณ์</RouterLink>
      </template>
    </nav>

    <div>
      <span>{{ auth.user?.name }}</span>
      <button type="button" @click="onLogout">ออกจากระบบ</button>
    </div>
  </header>

  <main class="layout__main">
    <slot />
  </main>
</template>