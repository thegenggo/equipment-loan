export type Role = 'staff' | 'admin'
export type EquipmentStatus = 'available' | 'borrowed' | 'repair'
export type LoanStatus = 'pending' | 'approved' | 'rejected' | 'returned'

export interface User {
  id: number
  email: string
  name: string
  role: Role
}

export interface Equipment {
  id: number
  code: string
  name: string
  category: string
  status: EquipmentStatus
}

export interface LoanRequest {
  id: number
  user_id: number
  equipment_id: number
  purpose: string
  status: LoanStatus
  requested_at: string
  approved_by: number | null
  approved_at: string | null
  returned_at: string | null
  equipment_code: string
  equipment_name: string
  user_name: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface ApiError {
  error: string
  message: string
}