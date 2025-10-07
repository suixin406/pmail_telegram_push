import { get, post, ApiPrefix } from '@/api/api'
import type { ApiResponse } from '@/utils/requests'

export interface Setting {
  chat_id: string
  show_content: boolean
  spoiler_content: boolean
  send_attachments: boolean
  disable_link_preview: boolean
}

export function getSettingInfo(): Promise<ApiResponse<Setting>> {
  return get(ApiPrefix + 'settings')
}

export function saveSettingInfo(setting: Setting): Promise<ApiResponse<void>> {
  return post(ApiPrefix + 'submit', setting)
}
