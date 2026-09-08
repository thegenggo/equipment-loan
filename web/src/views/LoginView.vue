<script setup lang="ts">
import { ApiRequestError } from '@/api/client';
import { useAuthStore } from '@/stores/auth';
import { ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const email = ref('')
const password = ref('')
const errorMessage = ref('')
const isSubmitting = ref(false)

function safeRedirect(value: unknown): string {
  return typeof value === 'string' && value.startsWith('/') && !value.startsWith('/') ? value : '/'
}

async function onSubmit(): Promise<void> {
  errorMessage.value = ''
  isSubmitting.value = true

  try {
    await auth.login(email.value, password.value)
    await router.replace(safeRedirect(route.query.redirect))
  } catch (error) {
    errorMessage.value = error instanceof ApiRequestError && error.status === 401 ? 'อีเมลหรือรหัสผ่านไม่ถูกต้อง' : 'เข้าสู่ระบบไม่สำเร็จ ลองใหม่อีกครั้ง'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="login">
    <h1>ระบบยืม-คืนอุปกรณ์</h1>

    <form @submit.prevent="onSubmit">
      <label for="email">อีเมล</label>
      <input id="email" v-model="email" type="email" required autocomplete="email" />

      <label for="password">รหัสผ่าน</label>
      <input 
        id="password"
        v-model="password"
        type="password"
        required
        autocomplete="current-password"
      />

      <p v-if="errorMessage" class="error" role="alert">{{ errorMessage }}</p>

      <button type="submit" :disabled="isSubmitting">
        {{ isSubmitting ? 'กำลังเข้าสู่ระบบ...' : 'เข้าสู่ระบบ' }}
      </button>
    </form>
  </main>
</template>