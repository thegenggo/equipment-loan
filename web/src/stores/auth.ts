import { computed, ref } from "vue"
import { defineStore } from "pinia"

import * as authApi from '@/api/auth'
import { setAuthToken } from "@/api/client"
import type { Role, User } from "@/types"

const TOKEN_KEY = 'equipment-loan.token'
const USER_KEY = 'equipment-loan.user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => token.value !== null)
  const role = computed<Role | null>(() => user.value?.role ?? null)
  const isAdmin = computed(() => role.value === 'admin')

  function restore(): void {
    const storedToken = localStorage.getItem(TOKEN_KEY)
    const storedUser = localStorage.getItem(USER_KEY)
    if (storedToken === null || storedUser === null) {
      return
    }

    try {
      user.value = JSON.parse(storedUser) as User
    } catch {
      clear()
      return
    }

    token.value = storedToken
    setAuthToken(storedToken)
  }

  async function login(email: string, password: string): Promise<void> {
    const result = await authApi.login(email, password)

    token.value = result.token
    user.value = result.user
    setAuthToken(result.token)

    localStorage.setItem(TOKEN_KEY, result.token)
    localStorage.setItem(USER_KEY, JSON.stringify(result.user))
  }

  function clear(): void {
    token.value = null
    user.value = null
    setAuthToken(null)
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  }

  return { token, user, isAuthenticated, role, isAdmin, restore, login, logout: clear }
})