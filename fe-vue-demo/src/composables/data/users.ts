import { ref } from 'vue'

import type { User } from '../type'
import { getStoredValue, setStoredValue } from './storage'

export const USERS_STORAGE_KEY = 'users'

export const users = ref<User[]>([])

export const loadUsers = async () => {
    users.value = (await getStoredValue<User[]>(USERS_STORAGE_KEY)) || []
}

export const saveUsers = async () => {
    await setStoredValue(USERS_STORAGE_KEY, users.value)
}

export const toggleFollowing = async (userId: number, authorId: number) => {
    const userIndex = users.value.findIndex(user => user.id === userId)
    if (userIndex === -1) return false

    const user = users.value[userIndex]
    if (!user) return false
    const following = new Set(user.following || [])
    const isFollowing = following.has(authorId)

    if (isFollowing) {
        following.delete(authorId)
    } else {
        following.add(authorId)
    }

    users.value[userIndex] = { ...user, following: [...following] }
    await saveUsers()
    return !isFollowing
}