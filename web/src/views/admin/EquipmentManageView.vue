<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { toMessage } from '@/api/client'
import {
  createEquipment,
  deleteEquipment,
  listEquipments,
  updateEquipment,
} from '@/api/equipment'
import StateWrapper from '@/components/StateWrapper.vue'
import { useAsyncData } from '@/composables/useAsyncData'
import type { EditableEquipmentStatus, Equipment } from '@/types'
import { equipmentStatusLabel } from '@/utils/format'

const { data: equipments, isLoading, errorMessage, load } = useAsyncData(() => listEquipments(''))

const editing = ref<Equipment | null>(null)
const code = ref('')
const name = ref('')
const category = ref('')
const status = ref<EditableEquipmentStatus>('available')

const isSubmitting = ref(false)
const formError = ref('')
const actionError = ref('')
const deletingId = ref<number | null>(null)

const isEditMode = computed(() => editing.value !== null)
const canSubmit = computed(
  () =>
    code.value.trim() !== '' &&
    name.value.trim() !== '' &&
    category.value.trim() !== '' &&
    !isSubmitting.value,
)

function startCreate(): void {
  editing.value = null
  code.value = ''
  name.value = ''
  category.value = ''
  status.value = 'available'
  formError.value = ''
}

function startEdit(equipment: Equipment): void {
  editing.value = equipment
  code.value = equipment.code
  name.value = equipment.name
  category.value = equipment.category
  status.value = equipment.status === 'repair' ? 'repair' : 'available'
  formError.value = ''
}

async function onSubmit(): Promise<void> {
  isSubmitting.value = true
  formError.value = ''

  const input = {
    code: code.value.trim(),
    name: name.value.trim(),
    category: category.value.trim(),
  }

  try {
    if (editing.value === null) {
      await createEquipment(input)
    } else {
      await updateEquipment(editing.value.id, { ...input, status: status.value })
    }
    startCreate()
    await load()
  } catch (error) {
    formError.value = toMessage(error)
  } finally {
    isSubmitting.value = false
  }
}

async function onDelete(equipment: Equipment): Promise<void> {
  if (!window.confirm(`ลบ ${equipment.code} - ${equipment.name}?`)) {
    return
  }

  deletingId.value = equipment.id
  actionError.value = ''

  try {
    await deleteEquipment(equipment.id)
    if (editing.value?.id === equipment.id) {
      startCreate()
    }
  } catch (error) {
    actionError.value = toMessage(error)
  } finally {
    deletingId.value = null
    await load()
  }
}

onMounted(load)
</script>

<template>
  <h1>จัดการอุปกรณ์</h1>

  <form @submit.prevent="onSubmit">
    <h2>{{ isEditMode ? `แก้ไข ${editing?.code}` : 'เพิ่มอุปกรณ์' }}</h2>

    <label for="code">รหัส</label>
    <input id="code" v-model="code" required maxlength="50" />

    <label for="name">ชื่อ</label>
    <input id="name" v-model="name" required maxlength="150" />

    <label for="category">หมวด</label>
    <input id="category" v-model="category" required maxlength="50" />

    <template v-if="isEditMode">
      <label for="equipment-status">สถานะ</label>
      <select id="equipment-status" v-model="status">
        <option value="available">ว่าง</option>
        <option value="repair">ซ่อมบำรุง</option>
      </select>
    </template>

    <p v-if="formError" class="error" role="alert">{{ formError }}</p>

    <button type="submit" :disabled="!canSubmit">
      {{ isSubmitting ? 'กำลังบันทึก...' : isEditMode ? 'บันทึก' : 'เพิ่ม' }}
    </button>
    <button v-if="isEditMode" type="button" @click="startCreate">ยกเลิก</button>
  </form>

  <p v-if="actionError" class="error" role="alert">{{ actionError }}</p>

  <StateWrapper
    :is-loading="isLoading"
    :error-message="errorMessage"
    :is-empty="equipments?.length === 0"
    empty-text="ยังไม่มีอุปกรณ์ในระบบ"
    @retry="load"
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
          <td>{{ equipmentStatusLabel[equipment.status]}}</td>
          <td>
            <template v-if="equipment.status === 'borrowed'">
              <span>ถูกยืมอยู่ แก้ไขไม่ได้</span>
            </template>
            <template v-else>
              <button type="button" @click="startEdit(equipment)">แก้ไข</button>
              <button
                type="button"
                :disabled="deletingId === equipment.id"
                @click="onDelete(equipment)"
              >
                ลบ
              </button>
            </template>
          </td>
        </tr>
      </tbody>
    </table>
  </StateWrapper>
</template>