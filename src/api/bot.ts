import { get } from '@/api'
import type { BotInfo } from '@/types'

export const getBotInfo = () => {
  return get<BotInfo>('bot')
}
