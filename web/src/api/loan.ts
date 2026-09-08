import type { LoanRequest } from "@/types"
import { request } from "./client"

export function listLoans(): Promise<LoanRequest[]> {
  return request<LoanRequest[]>('/loans')
}

export function createLoan(equipmentId: number, purpose: string): Promise<LoanRequest> {
  return request<LoanRequest>('/loans', {
    method: 'POST',
    body: { equipment_id: equipmentId, purpose },
  })
}

export function returnLoan(id: number): Promise<LoanRequest> {
  return request<LoanRequest>(`/loans/${id}/return`, { method: 'PATCH' })
}