<script setup lang="ts">
import { computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

import type { Comment, User } from '../composables/type'
import { CURRENT_USER_STORAGE_KEY, currentUser } from '../composables/auth'
import { comments, deleteComment, saveComments } from '../composables/data/comments'
import { saveUsers, users } from '../composables/data/users'
import { setStoredValue } from '../composables/data/storage'

const props = defineProps<{
    comment: Comment
    user: User
}>()

const user = computed(() => {
    return users.value.find(item => item.id === props.comment.userId) || props.user
})

const publishedDate = computed(() => {
    if (!props.comment) return ''

    const date = new Date(props.comment.createTime)

    return `${date.getUTCMonth() + 1} 月 ${date.getUTCDate()} 日`
})

const isLiked = computed(() => currentUser.value?.likeComments?.includes(props.comment.id) || false)
const isCommentAuthor = computed(() => currentUser.value?.id === props.comment.userId)

const handleLike = async () => {
    if (!currentUser.value) {
        ElMessage.warning('请先登录')
        return
    }

    const userIndex = users.value.findIndex(item => item.id === currentUser.value?.id)
    const user = userIndex >= 0 ? users.value[userIndex] : undefined
    if (!user) return

    const likedComments = new Set(user.likeComments || [])
    if (likedComments.has(props.comment.id)) {
        likedComments.delete(props.comment.id)
        props.comment.likes = Math.max(0, props.comment.likes - 1)
    } else {
        likedComments.add(props.comment.id)
        props.comment.likes += 1
    }

    const updatedUser = { ...user, likeComments: [...likedComments] }
    users.value[userIndex] = updatedUser
    currentUser.value = updatedUser

    await Promise.all([
        saveUsers(),
        saveComments(),
        setStoredValue(CURRENT_USER_STORAGE_KEY, updatedUser)
    ])
}

const handleDelete = async () => {
    if (!isCommentAuthor.value) return

    try {
        await ElMessageBox.confirm('确定删除这条评论吗？', '删除评论', {
            type: 'warning',
            confirmButtonText: '删除',
            cancelButtonText: '取消'
        })
        await deleteComment(props.comment.id)
        ElMessage.success('评论已删除')
    } catch {
        // 取消确认时不提示
    }
}
</script>

<template>
    <article :data-comment-id="props.comment.id">
        <img :src="user.avatar" />
        <div>
            <header>
                <b>{{ user.name }}</b>
                <time>{{ publishedDate }}</time>
            </header>
            <p>{{ props.comment.content }}</p>
            <footer>
                <label class="comment-like" :class="{ 'is-liked': isLiked }" @click="handleLike">
                    <span class="comment-heart">
                        <el-icon>
                            <StarFilled />
                        </el-icon>
                    </span>
                    <span data-likes>{{ props.comment.likes }}</span>
                </label>
                <button v-if="isCommentAuthor" class="comment-delete" type="button" @click="handleDelete">
                    删除
                </button>
            </footer>
        </div>
    </article>
</template>

<style scoped>
article {
    display: flex;
    gap: 10px;
}

img {
    width: 30px;
    height: 30px;
    border-radius: 50%;
    object-fit: cover;
}

div {
    flex: 1;
    min-width: 0;
}

header {
    display: flex;
    align-items: center;
    gap: 8px;
}

b {
    color: var(--text-prose);
    font-size: 12px;
}

time,
footer {
    color: var(--text-mute);
    font-size: 11px;
}

p {
    margin: 7px 0 8px;
    color: var(--text-prose);
    font-size: 13px;
    line-height: 1.55;
}

footer {
    display: flex;
    gap: 13px;
}

.comment-like {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    color: var(--text-mute);
    cursor: pointer;
    user-select: none;
}

.comment-heart,
.comment-heart :deep(.el-icon) {
    display: inline-flex;
    align-items: center;
    justify-content: center;
}

.comment-heart {
    line-height: 0;
}

.comment-heart :deep(.el-icon) {
    position: relative;
}

.comment-delete {
    padding: 0;
    border: 0;
    color: var(--text-mute);
    font: inherit;
    font-size: 11px;
    background: transparent;
    cursor: pointer;
}

.comment-delete:hover {
    color: #f65369;
}
</style>
