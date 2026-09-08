<script setup lang="ts">
import { toMessage } from '@/api/client';
import { listLoans, returnLoan } from '@/api/loan';
import StateWrapper from '@/components/StateWrapper.vue';
import { useAsyncData } from '@/composables/useAsyncData';
import { formatDateTime, loanStatusLabel } from '@/utils/format';
import { onMounted, ref } from 'vue';

const { data: loans, isLoading, errorMessage, load } = useAsyncData(listLoans)

const returningId = ref<number | null>(null)
const actionError = ref('')

async function onReturn(id: number): Promise<void> {
  returningId.value = id
  actionError.value = ''

  try {
    await returnLoan(id)
  } catch (error) {
    actionError.value = toMessage(error)
  } finally {
    returningId.value = null
    await load()
  }
}

onMounted(load)
</script>

<template>
  <h1>คำขอของฉัน</h1>

  <p v-if="actionError" class="error" role="alert">{{ actionError }}</p>

  <StateWrapper
    :is-loading="isLoading"
    :error-message="errorMessage"
    :is-empty="loans?.length === 0"
    empty-text="ยังไม่เคยขอยืมอุปกรณ์"
    @retry="load"
  >
    <table>
      <thead>
        <tr>
          <th>อุปกรณ์</th>
          <th>วัตถุประสงค์</th>
          <th>สถานะ</th>
          <th>วันที่ขอ</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="loan in loans" :key="loan.id">
          <td>{{ loan.equipment_code }} - {{ loan.equipment_name }}</td>
          <td>{{ loan.purpose }}</td>
          <td>{{ loanStatusLabel[loan.status] }}</td>
          <td>{{ formatDateTime(loan.requested_at) }}</td>
          <td>
            <button
              v-if="loan.status === 'approved'"
              type="button"
              :disabled="returningId === loan.id"
              @click="onReturn(loan.id)"
            >
              {{ returningId === loan.id ? 'กำลังคืน...' : 'คืน' }}
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </StateWrapper>
</template>