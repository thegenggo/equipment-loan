<script setup lang="ts">
import { toMessage } from '@/api/client';
import { listEquipments } from '@/api/equipment';
import { createLoan, listLoans } from '@/api/loan';
import StateWrapper from '@/components/StateWrapper.vue';
import { useAsyncData } from '@/composables/useAsyncData';
import type { Equipment, EquipmentStatus } from '@/types'
import { equipmentStatusLabel } from '@/utils/format';
import { computed, onMounted, ref } from 'vue'

const statusFilter = ref<EquipmentStatus | ''>('')

const {
  data: equipments,
  isLoading,
  errorMessage,
  load,
} = useAsyncData(() => listEquipments(statusFilter.value))

const requestedIds = ref(new Set<number>())

const selected = ref<Equipment | null>(null)
const purpose = ref('')
const isSubmitting = ref(false)
const formError = ref('')
const successMessage = ref('')

const canSubmit = computed(() => purpose.value.trim().length > 0 && !isSubmitting.value)

async function refresh(): Promise<void> {
  await load()
  
  try {
    const loans = await listLoans()
    requestedIds.value = new Set(
      loans
        .filter((loan) => loan.status === 'pending' || loan.status === 'approved')
        .map((loan) => loan.equipment_id),
    )
  } catch {
    requestedIds.value = new Set()
  }
}

function openForm(equipment: Equipment): void {
  selected.value = equipment
  purpose.value = ''
  formError.value = ''
  successMessage.value = ''
}

function closeForm(): void {
  selected.value = null
  purpose.value = ''
  formError.value = ''
}

async function onSubmit(): Promise<void> {
  if (selected.value === null) {
    return
  }

  isSubmitting.value = true
  formError.value = ''

  try {
    await createLoan(selected.value.id, purpose.value.trim())
    successMessage.value = `ส่งคำขอยืม ${selected.value.code} แล้ว รออนุมัติ`
    closeForm()
    await refresh()
  } catch (error) {
    formError.value = toMessage(error)
    await refresh()
  } finally {
    isSubmitting.value = false
  }
}

onMounted(refresh)
</script>

<template>
  <h1>อุปกรณ์</h1>
  
  <label for="status">สถานะ</label>
  <select id="status" v-model="statusFilter" @change="refresh">
    <option value="">ทั้งหมด</option>
    <option value="available">ว่าง</option>
    <option value="borrowed">ถูกยืม</option>
    <option value="repair">ซ่อมบำรุง</option>
  </select>

  <p v-if="successMessage" class="success" role="status">{{ successMessage }}</p>

  <form v-if="selected" class="loan-form" @submit.prevent="onSubmit">
    <h2>ยืม {{ selected.code }} - {{ selected.name }}</h2>

    <label for="purpose">วัตถุประสงค์</label>
    <textarea 
      id="purpose"
      v-model="purpose"
      required
      maxlength="500"
      rows="3"
      placeholder="ใช้ทำอะไร"
    ></textarea>
    <small>{{ purpose.length }}</small>

    <p v-if="formError" class="error" role="alert">{{ formError }}</p>

    <button type="submit" :disabled="!canSubmit">
      {{ isSubmitting ? 'กำลังส่ง...' : 'ส่งคำขอ' }}
    </button>
    <button type="button" @click="closeForm">ยกเลิก</button>
  </form>

  <StateWrapper
    :is-loading="isLoading"
    :error-message="errorMessage"
    :is-empty="equipments?.length === 0"
    empty-text="ไม่มีอุปกรณ์ตรงกับตัวกรองนี้"
    @retry="refresh"
  >
    <table>
      <thead>
        <tr>
          <th>รหัส</th>
          <th>ชื่อ</th>
          <th>หมวด</th>
          <th>สถานะ</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="equipment in equipments" :key="equipment.id">
          <td>{{ equipment.code }}</td>
          <td>{{ equipment.name }}</td>
          <td>{{ equipment.category }}</td>
          <td>{{ equipmentStatusLabel[equipment.status] }}</td>
          <td>
            <span v-if="requestedIds.has(equipment.id)">คุณขอยืมไว้แล้ว</span>
            <button
              v-else-if="equipment.status === 'available'"
              type="button"
              @click="openForm(equipment)"
            >
              ยืม
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </StateWrapper>
</template>