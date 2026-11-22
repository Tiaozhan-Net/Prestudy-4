/** 后端统一响应结构：{ code, message, data } */
export interface ApiResult<T> {
  code: number
  message: string
  data: T
}

/** GET /api/auth/profile 返回的用户资料 */
export interface UserProfile {
  id: number
  name: string
  avatar: string
  slogan: string
}
