import type { Equipment, EquipmentInput, EquipmentStatus, EquipmentUpdateInput } from "@/types"
import { request } from './client'

export function listEquipments(status: EquipmentStatus | ''): Promise<Equipment[]> {
  const params = new URLSearchParams()
  if (status !== '') {
    params.set('status', status)
  }

  const query = params.toString()
  return request<Equipment[]>(query === '' ? '/equipments' : `/equipments?${query}`)
}

export function createEquipment(input: EquipmentInput): Promise<Equipment> {
  return request<Equipment>('/equipments', { method: 'POST', body: input })
}

export function updateEquipment(id: number, input: EquipmentUpdateInput): Promise<Equipment> {
  return request<Equipment>(`/equipments/${id}`, { method: 'PUT', body: input })
}

export function deleteEquipment(id: number): Promise<void> {
  return request<void>(`/equipments/${id}`, { method: 'DELETE' })
}