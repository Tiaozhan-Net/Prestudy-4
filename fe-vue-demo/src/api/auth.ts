import { request } from './http'
import type { UserProfile } from '@/types/auth'

export interface LoginPayload {
  name: string
  password: string
}

export type RegisterPayload = LoginPayload

/**
 * 登录相关的接口都放在这里。
 * 页面和 store 不直接写 Axios，只调用 authApi。
 */
export const authApi = {
  /** POST /api/auth/register：注册新账号 */
  register(data: RegisterPayload) {
    return request<null>({
      url: '/auth/register',
      method: 'POST',
      data,
    })
  },

  /** POST /api/auth/login：验证账号密码，后端下发 tz-sessions Cookie */
  login(data: LoginPayload) {
    return request<null>({
      url: '/auth/login',
      method: 'POST',
      data,
    })
  },

  /** GET /api/auth/profile：带着 Cookie 问“我是谁” */
  profile() {
    return request<UserProfile>({
      url: '/auth/profile',
      method: 'GET',
      timeout: 4000,
    })
  },

  /** DELETE /api/auth/logout：清除后端 Session */
  logout() {
    return request<null>({
      url: '/auth/logout',
      method: 'DELETE',
    })
  },
}
