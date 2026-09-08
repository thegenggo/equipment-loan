<script setup lang="ts">
defineProps<{
  isLoading: boolean
  errorMessage: string
  isEmpty: boolean
  emptyText?: string
}>()

defineEmits<{ retry: [] }>()
</script>

<template>
  <p v-if="isLoading" class="state" role="status">กำลังโหลด...</p>

  <div v-else-if="errorMessage" class="state state--error" role="alert">
    <p>{{ errorMessage }}</p>
    <button type="button" @click="$emit('retry')">ลองใหม่</button>
  </div>

  <p v-else-if="isEmpty" class="state">{{ emptyText ?? 'ยังไม่มีข้อมูล' }}</p>

  <slot v-else />
</template>