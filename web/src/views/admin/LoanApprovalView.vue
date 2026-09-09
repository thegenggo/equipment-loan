<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { toMessage } from '@/api/client'
import { approveLoan, listLoans, rejectLoan } from '@/api/loan'
import StateWrapper from '@/components/StateWrapper.vue'
import { useAsyncData } from '@/composables/useAsyncData'
import type { LoanStatus } from '@/types'
import { formatDateTime, loanStatusLabel } from '@/utils/format'

const statusFilter = ref<LoanStatus | ''>('pending')

const {
  data: loans,
  isLoading,
  errorMessage,
  load
} = useAsyncData(() => listLoans(statusFilter.value))

const decidingId = ref<number | null>(null)
const actionError = ref('')

async function decide(id: number, action: (id: number) => Promise<unknown>): Promise<void> {
  decidingId.value = id
  actionError.value = ''

  try {
    await action(id)
  } catch (error) {
    actionError.value = toMessage(error)
  } finally {
    decidingId.value = null
    await load()
  }
}

onMounted(load)
</script>

<template>
  <h1>อนุมัติคำขอ</h1>

  <label for="status">สถานะ</label>
  <select id="status" v-model="statusFilter" @change="load">
    <option value="">ทั้งหมด</option>
    <option value="pending">รออนุมัติ</option>
    <option value="approved">กำลังยืม</option>
    <option value="rejected">ไม่อนุมัติ</option>
    <option value="returned">คืนแล้ว</option>
  </select>

  <p v-if="actionError" class="error" role="alert"{{ actionError }}></p>

  <StateWrapper
    :is-loading="isLoading"
    :error-message="errorMessage"
    :is-empty="loans?.length === 0"
    empty-text="ไม่มีคำขอในสถานะนี้"
    @retry="load"
  >
    <table>
      <thead>
        <tr>
          <th>ผู้ขอ</th>
          <th>อุปกรณ์</th>
          <th>วัตถุประสงค์</th>
          <th>สถานะ</th>
          <th>วันที่ขอ</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="loan in loans" :key="loan.id">
          <td>{{ loan.user_name }}</td>
          <td>{{ loan.equipment_code }} - {{ loan.equipment_name }}</td>
          <td>{{ loan.purpose }}</td>
          <td>{{ loanStatusLabel[loan.status] }}</td>
          <td>{{ formatDateTime(loan.requested_at) }}</td>
          <td>
            <template v-if="loan.status === 'pending'">
              <button
                type="button"
                :disabled="decidingId === loan.id"
                @click="decide(loan.id, approveLoan)"
              >
                อนุมัติ
              </button>
              <button
                type="button"
                :disabled="decidingId === loan.id"
                @click="decide(loan.id, rejectLoan)"
              >
                ปฏิเสธ
              </button>
            </template>
          </td>
        </tr>
      </tbody>
    </table>
  </StateWrapper>
</template>