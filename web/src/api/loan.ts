import type { LoanRequest, LoanStatus } from "@/types"
import { request } from "./client"

export function listLoans(status: LoanStatus | '' = ''): Promise<LoanRequest[]> {
  const params = new URLSearchParams()
  if (status !== '') {
    params.set('status', status)
  }

  const query = params.toString()
  return request<LoanRequest[]>(query === '' ? '/loans' : `/loans?${query}`)
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

export function approveLoan(id: number): Promise<LoanRequest> {
  return request<LoanRequest>(`/loans/${id}/approve`, { method: 'PATCH' })
}

export function rejectLoan(id: number): Promise<LoanRequest> {
  return request<LoanRequest>(`/loans/${id}/reject`, { method: 'PATCH' })
}