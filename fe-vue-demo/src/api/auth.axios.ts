import axios from 'axios'
import type { UserProfile } from '@/types/auth'

interface LoginPayload {
  name: string
  password: string
}

interface ApiResult<T> {
  code: number
  message: string
  data: T
}

const axiosConfig = {
  baseURL: '/api',
  timeout: 8000,
  withCredentials: true,
}

function errorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const responseData = error.response?.data as { message?: string } | undefined
    if (responseData?.message) return responseData.message
    if (error.code === 'ECONNABORTED') return '请求超时，请检查后端是否启动'
    if (!error.response) return '连不上后端，请检查 8088 端口和 CORS 配置'
    return `请求失败（${error.response.status}）`
  }

  return error instanceof Error ? error.message : '请求失败'
}

export async function login(name: string, password: string) {
  try {
    const response = await axios.post<ApiResult<null>>(
      '/auth/login',
      { name, password } satisfies LoginPayload,
      axiosConfig,
    )

    if (response.data.code !== 0) {
      console.error(response.data.message)
    } else {
      console.log(response.data)
    }
  } catch (error) {
    console.error(errorMessage(error))
  }
}

export async function register(name: string, password: string) {
  try {
    const response = await axios.post<ApiResult<null>>(
      '/auth/register',
      { name, password } satisfies LoginPayload,
      axiosConfig,
    )

    if (response.data.code !== 0) {
      console.error(response.data.message)
    } else {
      console.log(response.data)
    }
  } catch (error) {
    console.error(errorMessage(error))
  }
}

export async function profile() {
  try {
    const response = await axios.get<ApiResult<UserProfile>>('/auth/profile', {
      ...axiosConfig,
      timeout: 4000,
    })

    if (response.data.code !== 0) {
      console.error(response.data.message)
    } else {
      console.log(response.data)
    }
  } catch (error) {
    console.error(errorMessage(error))
  }
}

export async function logout() {
  try {
    const response = await axios.delete<ApiResult<null>>('/auth/logout', axiosConfig)

    if (response.data.code !== 0) {
      console.error(response.data.message)
    } else {
      console.log(response.data)
    }
  } catch (error) {
    console.error(errorMessage(error))
  }
}

export const authAxios = {
  login,
  register,
  profile,
  logout,
}
