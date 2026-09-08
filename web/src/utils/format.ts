import type { EquipmentStatus, LoanStatus } from "@/types"

export function formatDateTime(value: string): string {
  return new Date(value).toLocaleString('th-TH', {
    dateStyle: 'medium',
    timeStyle: 'short',
  })
}

export const equipmentStatusLabel: Record<EquipmentStatus, string> = {
  available: 'ว่าง',
  borrowed: 'ถูกยืม',
  repair: 'ซ่อมบำรุง',
}

export const loanStatusLabel: Record<LoanStatus, string> = {
  pending: 'รออนุมัติ',
  approved: 'กำลังยืม',
  rejected: 'ไม่อนุมัติ',
  returned: 'คืนแล้ว',
}