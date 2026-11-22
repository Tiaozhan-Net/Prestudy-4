import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import { authApi, type LoginPayload } from '@/api/auth'
import type { UserProfile } from '@/types/auth'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<UserProfile | null>(null)
  const loading = ref(false)
  const error = ref('')
  const initialized = ref(false)

  const isLoggedIn = computed(() => user.value !== null)
  const displayName = computed(() => user.value?.name ?? '游客')

  /** 登录：先拿 Session，再用 Session 换用户资料 */
  async function login(payload: LoginPayload) {
    loading.value = true
    error.value = ''

    try {
      await authApi.login(payload)
      user.value = await authApi.profile()
    } catch (err) {
      error.value = err instanceof Error ? err.message : '登录失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  function clearError() {
    error.value = ''
  }

  /** 注册：先创建账号，再自动登录 */
  async function register(payload: LoginPayload) {
    loading.value = true
    error.value = ''

    try {
      await authApi.register(payload)
      await login(payload)
    } catch (err) {
      error.value = err instanceof Error ? err.message : '注册失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  /** 刷新页面后：带着浏览器里的 Cookie 恢复登录态 */
  async function restore() {
    if (initialized.value) return

    try {
      user.value = await authApi.profile()
    } catch {
      user.value = null
    } finally {
      initialized.value = true
    }
  }

  /** 退出：先通知后端清 Session，再清本地状态 */
  async function logout() {
    try {
      await authApi.logout()
    } catch {
      // 后端暂时连不上时，也允许本地退出
    } finally {
      user.value = null
      error.value = ''
    }
  }

  return {
    user,
    loading,
    error,
    initialized,
    isLoggedIn,
    displayName,
    clearError,
    login,
    register,
    restore,
    logout,
  }
})
