import axios, { type AxiosRequestConfig } from 'axios'
import type { ApiResult } from '@/types/auth'

/**
 * 带业务码的请求错误：
 * - status: HTTP 状态码，例如 401
 * - code:   后端业务码，例如 2001 / 2003
 */
export class ApiError extends Error {
  status?: number
  code?: number

  constructor(message: string, options?: { status?: number; code?: number }) {
    super(message)
    this.name = 'ApiError'
    this.status = options?.status
    this.code = options?.code
  }
}

const baseURL = import.meta.env.VITE_API_BASE_URL

if (!baseURL) {
  throw new Error('缺少 VITE_API_BASE_URL，请检查 .env.development')
}

export const http = axios.create({
  baseURL,
  timeout: 8000,
  withCredentials: true, // 跨域请求也要带上 tz-sessions Cookie
})

http.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.code === 'ECONNABORTED') {
      return Promise.reject(new ApiError('请求超时，请检查后端是否启动'))
    }

    if (!error.response) {
      return Promise.reject(new ApiError('连不上后端，请检查 8088 端口和 CORS 配置'))
    }

    const data = error.response.data as Partial<ApiResult<unknown>> | undefined
    return Promise.reject(
      new ApiError(data?.message || `请求失败（${error.response.status}）`, {
        status: error.response.status,
        code: data?.code,
      }),
    )
  },
)

/** 统一判断 code，并把 data 剥出来；失败时抛 ApiError */
export async function request<T>(config: AxiosRequestConfig): Promise<T> {
  const response = await http.request<ApiResult<T>>(config)
  const result = response.data

  if (result.code !== 0) {
    throw new ApiError(result.message || '请求失败', {
      status: response.status,
      code: result.code,
    })
  }

  return result.data
}
