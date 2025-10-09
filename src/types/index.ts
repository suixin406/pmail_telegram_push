export interface ApiResponse<T> {
  code: number
  message: string
  data?: T
}

export interface BotInfo {
  username: string
  first_name: string
  bot_link: string
}

export interface Setting {
  chat_id: string
  show_content: boolean
  spoiler_content: boolean
  send_attachments: boolean
  disable_link_preview: boolean
}
