import type { Equipment, EquipmentStatus } from "@/types"
import { request } from './client'

export function listEquipments(status: EquipmentStatus | ''): Promise<Equipment[]> {
  const params = new URLSearchParams()
  if (status !== '') {
    params.set('status', status)
  }

  const query = params.toString()
  return request<Equipment[]>(query === '' ? '/equipments' : `/equipments?${query}`)
}

