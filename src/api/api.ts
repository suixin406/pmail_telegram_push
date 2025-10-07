import service from '@/utils/requests'
import type { ApiResponse } from '@/utils/requests'

export const ApiPrefix = '/api/plugin/settings/pmail_telegram_push/'

export function get<T>(url: string, params?: object): Promise<ApiResponse<T>> {
  return service.get(url, { params })
}

export function post<T>(url: string, data?: object): Promise<ApiResponse<T>> {
  return service.post(url, data)
}
