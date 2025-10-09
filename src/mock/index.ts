import Mock from 'mockjs'
import type { ApiResponse, BotInfo, Setting } from '@/types'
import service from '@/utils/requests'

const baseurl = service.defaults.baseURL
const Random = Mock.Random

Mock.setup({
  timeout: '500-1000',
})

export const botInfoMock = Mock.mock(baseurl + 'bot', 'get', () => {
  const botInfo: ApiResponse<BotInfo> = {
    code: 0,
    message: 'success',
    data: {
      username: Random.string('lower', 10),
      first_name: Random.string('lower', 10),
      bot_link: Random.string('lower', 10),
    },
  }
  console.log('botInfo mock response', botInfo)
  return botInfo
})

export const settingMock = Mock.mock(baseurl + 'settings', 'get', () => {
  const setting: ApiResponse<Setting> = {
    code: 0,
    message: 'success',
    data: {
      chat_id: Random.string('number', 9),
      show_content: Random.boolean(),
      spoiler_content: Random.boolean(),
      send_attachments: Random.boolean(),
      disable_link_preview: Random.boolean(),
    },
  }
  console.log('setting mock response', setting)
  return setting
})

export const submitMock = Mock.mock(baseurl + 'submit', 'post', (options) => {
  console.log('submit options', options)
  const code = Random.integer(0, -1)
  const message = code === 0 ? 'success' : 'error'
  const setting: ApiResponse<Setting> = {
    code,
    message,
  }
  console.log('submit mock response', setting)
  return setting
})

export default {
  botInfoMock,
  settingMock,
  submitMock,
}
