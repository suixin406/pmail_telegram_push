import { get, post } from '@/api'
import type { Setting, ApiResponse } from '@/types'

export function getSettingInfo(): Promise<ApiResponse<Setting>> {
  return get('settings')
}

export function saveSettingInfo(setting: Setting): Promise<ApiResponse<void>> {
  return post('submit', setting)
}
