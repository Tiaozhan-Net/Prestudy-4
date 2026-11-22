import { ref } from 'vue'

import type { Comment } from '../type'
import { getStoredValue, setStoredValue } from './storage'

export const COMMENTS_STORAGE_KEY = 'comments'

export const comments = ref<Comment[]>([])

export const loadComments = async () => {
    comments.value = (await getStoredValue<Comment[]>(COMMENTS_STORAGE_KEY)) || []
}

export const saveComments = async () => {
    await setStoredValue(COMMENTS_STORAGE_KEY, comments.value)
}

export const addComment = async (comment: Omit<Comment, 'id' | 'createTime'>) => {
    const nextId = comments.value.reduce((maxId, item) => Math.max(maxId, item.id), 0) + 1
    const newComment: Comment = {
        ...comment,
        id: nextId,
        createTime: new Date().toISOString()
    }

    comments.value.push(newComment)
    await saveComments()
    return newComment
}

export const deleteComment = async (commentId: number) => {
    comments.value = comments.value.filter(comment => comment.id !== commentId)
    await saveComments()
}