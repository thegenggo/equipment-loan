import type { LoginResponse, User } from "@/types";
import { request } from "./client";

export function login(email: string, password: string): Promise<LoginResponse> {
  return request<LoginResponse>('/auth/login', {
    method: 'POST',
    body: { email, password },
    isPublic: true
  })
}

export function fetchMe(): Promise<Pick<User, 'id' | 'role'>> {
  return request<Pick<User, 'id' | 'role'>>('/me')
}