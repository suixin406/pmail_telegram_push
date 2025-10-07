import { get, ApiPrefix } from '@/api/api'

export interface BotInfo {
  username: string
  first_name: string
  bot_link: string
}

export const getBotInfo = () => {
  return get<BotInfo>(ApiPrefix + 'bot')
}
