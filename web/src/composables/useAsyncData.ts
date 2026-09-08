import { ref } from 'vue'

import { toMessage } from '@/api/client'

export function useAsyncData<T>(loader: () => Promise<T>) {
  const data = ref<T | null>(null)
  const isLoading = ref(false)
  const errorMessage = ref('')

  async function load(): Promise<void> {
    isLoading.value = true
    errorMessage.value = ''

    try {
      data.value = await loader()
    } catch (error) {
      errorMessage.value = toMessage(error)
    } finally {
      isLoading.value = false
    }
  }

  return { data, isLoading, errorMessage, load }
}