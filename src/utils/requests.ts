import axios from 'axios'
import type { AxiosResponse, AxiosInstance } from 'axios'

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

const service: AxiosInstance = axios.create({
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json;charset=utf-8',
  },
})

service.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data
    if (res.code === 0) {
      return res
    }
    return Promise.reject(res.message || '请求失败')
  },
  (error) => {
    return Promise.reject(error)
  },
)

export default service
