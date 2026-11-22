import { ref } from 'vue'

import type { User } from './type'
import { useAuthStore } from '@/stores/auth'

/**
 * 兼容原项目页面：组件仍然从 currentUser 读取用户。
 * 真正的登录态由 Pinia 的 auth store 维护，再同步到这里。
 */
export const CURRENT_USER_STORAGE_KEY = 'current-user'
export const currentUser = ref<User | null>(null)

function authStore() {
  return useAuthStore()
}

function toUser(profile: {
  id: number
  name: string
  avatar: string
  slogan: string
}): User {
  return {
    id: profile.id,
    name: profile.name,
    avatar: profile.avatar || '/avatar.png',
    createTime: new Date().toISOString(),
    description: profile.slogan || '',
  }
}

function syncFromStore() {
  const profile = authStore().user
  currentUser.value = profile ? toUser(profile) : null
}

export const loadCurrentUser = async () => {
  await authStore().restore()
  syncFromStore()
}

export const registerUser = async (name: string, password: string): Promise<User> => {
  await authStore().register({ name, password })
  syncFromStore()
  if (!currentUser.value) throw new Error('注册失败')
  return currentUser.value
}

export const loginUser = async (name: string, password: string): Promise<User> => {
  await authStore().login({ name, password })
  syncFromStore()
  if (!currentUser.value) throw new Error('登录失败')
  return currentUser.value
}

export const logoutUser = async () => {
  await authStore().logout()
  currentUser.value = null
}
