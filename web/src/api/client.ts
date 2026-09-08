import type { ApiError } from "@/types";

const BASE_URL = '/api/v1'

let authToken: string | null = null
let onUnauthorized: () => void = () => { }

export function setAuthToken(token: string | null): void {
  authToken = token
}

export function setUnauthorizedHandler(handler: () => void): void {
  onUnauthorized = handler
}

export class ApiRequestError extends Error {
  constructor(
    readonly status: number,
    readonly code: string,
    message: string,
  ) {
    super(message)
    this.name = 'ApiRequestError'
  }
}

interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  body?: unknown
  isPublic?: boolean
}

export async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = 'GET', body, isPublic = false } = options

  const headers: Record<string, string> = {}
  if (body !== undefined) {
    headers['Content-Type'] = 'application/json'
  }
  if (!isPublic && authToken !== null) {
    headers.Authorization = `Bearer ${authToken}`
  }

  const response = await fetch(BASE_URL + path, {
    method,
    headers,
    body: body === undefined ? undefined : JSON.stringify(body),
  })

  if (response.status === 401 && !isPublic) {
    onUnauthorized()
  }

  if (!response.ok) {
    const payload = (await response.json().catch(() => null)) as ApiError | null
    throw new ApiRequestError(
      response.status,
      payload?.error ?? 'unknown_error',
      payload?.message ?? 'ระบบขัดข้อง ลองใหม่อีกครั้ง',
    )
  }

  if (response.status === 204) {
    return undefined as T
  }

  return (await response.json()) as T
}

export function toMessage(error: unknown): string {
  return error instanceof ApiRequestError
    ? error.message
    : 'เชื่อมต่อเซิร์ฟเวอร์ไม่ได้ ลองใหม่อีกครั้ง'
}