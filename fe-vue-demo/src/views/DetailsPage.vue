<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'

import type { User } from '@/composables/type'
import { CURRENT_USER_STORAGE_KEY, currentUser } from '@/composables/auth'
import { deleteArticle, articles, saveArticles } from '@/composables/data/articles'
import { addComment, comments } from '@/composables/data/comments'
import { toggleFollowing, users, saveUsers } from '@/composables/data/users'
import { setStoredValue } from '@/composables/data/storage'

const route = useRoute()
const router = useRouter()
const commentInput = ref('')
const commentSort = ref('likes')

const article = computed(() => {
    const articleId = Number(route.params.id)
    return articles.value.find(item => item.id === articleId)
})

const author = computed(() => {
    return users.value.find(item => item.id === article.value?.authorId) || {
        id: 0,
        name: '未知作者',
        avatar: '/avatar.png',
        createTime: 'NULL'
    } satisfies User
})

const articleComments = computed(() => {
    if (!article.value) return []

    return comments.value
        .filter(comment => comment.articleId === article.value?.id)
        .sort((first, second) => commentSort.value === 'likes'
            ? second.likes - first.likes
            : new Date(second.createTime).getTime() - new Date(first.createTime).getTime())
})

const articleImages = computed(() => {
    if (!article.value) return []

    return article.value.images.length > 0 ? article.value.images : [article.value.cover]
})

const publishedDate = computed(() => {
    if (!article.value) return ''

    const date = new Date(article.value.createTime)

    return `${date.getUTCMonth() + 1} 月 ${date.getUTCDate()} 日`
})

const isAuthor = computed(() => currentUser.value?.id === article.value?.authorId)

const isFollowing = computed(() => {
    const authorId = article.value?.authorId
    return authorId !== undefined && (currentUser.value?.following || []).includes(authorId)
})

const isLiked = computed(() => {
    const articleId = article.value?.id
    return articleId !== undefined && (currentUser.value?.likeArticles || []).includes(articleId)
})

const handleDelete = async () => {
    if (!article.value) return

    try {
        await ElMessageBox.confirm('确定删除这篇文章吗？', '删除文章', {
            type: 'warning',
            confirmButtonText: '删除',
            cancelButtonText: '取消'
        })
        await deleteArticle(article.value.id)
        ElMessage.success('文章已删除')
        await router.push('/')
    } catch {

    }
}

const handleFollow = async () => {
    if (!currentUser.value) {
        ElMessage.warning('请先登录')
        return
    }

    if (!article.value || isAuthor.value) return

    const following = await toggleFollowing(currentUser.value.id, article.value.authorId)
    const updatedUser = users.value.find(user => user.id === currentUser.value?.id)
    if (updatedUser) {
        currentUser.value = updatedUser
        await setStoredValue(CURRENT_USER_STORAGE_KEY, updatedUser)
    }
    ElMessage.success(following ? '关注成功' : '已取消关注')
}

const handleLike = async () => {
    if (!currentUser.value || !article.value) {
        ElMessage.warning('请先登录')
        return
    }

    const userIndex = users.value.findIndex(user => user.id === currentUser.value?.id)
    const articleIndex = articles.value.findIndex(item => item.id === article.value?.id)
    const user = userIndex >= 0 ? users.value[userIndex] : undefined
    const targetArticle = articleIndex >= 0 ? articles.value[articleIndex] : undefined

    if (!user || !targetArticle) return

    const likedArticles = new Set(user.likeArticles || [])
    if (likedArticles.has(targetArticle.id)) {
        likedArticles.delete(targetArticle.id)
        targetArticle.likes = Math.max(0, targetArticle.likes - 1)
    } else {
        likedArticles.add(targetArticle.id)
        targetArticle.likes += 1
    }

    const updatedUser = { ...user, likeArticles: [...likedArticles] }
    users.value[userIndex] = updatedUser
    currentUser.value = updatedUser
    await Promise.all([
        saveUsers(),
        saveArticles(),
        setStoredValue(CURRENT_USER_STORAGE_KEY, updatedUser)
    ])
}

const handleCommentSubmit = async () => {
    if (!currentUser.value || !article.value) {
        ElMessage.warning('请先登录')
        return
    }

    const content = commentInput.value.trim()
    if (!content) {
        ElMessage.warning('请输入评论内容')
        return
    }

    await addComment({
        articleId: article.value.id,
        userId: currentUser.value.id,
        content,
        likes: 0
    })
    commentInput.value = ''
    ElMessage.success('评论发布成功')
}
</script>

<template>
    <section class="note-detail-page">
        <div class="note-detail-layout">
            <article v-if="article" class="note-article">
                <div class="note-main">
                    <header class="note-author-header">
                        <img :src="author.avatar" :alt="author.name" />
                        <div>
                            <strong class="detail-author">{{ author.name }}</strong>
                            <span class="detail-date">发布于 {{ publishedDate }}</span>
                        </div>
                        <button v-if="isAuthor" class="follow-btn delete-btn" type="button" @click="handleDelete">
                            删除
                        </button>
                        <button v-else class="follow-btn" type="button" @click="handleFollow">
                            {{ isFollowing ? '已关注' : '关注' }}
                        </button>
                    </header>
                    <h1 class="detail-title">{{ article.title }}</h1>
                    <div class="detail-tags">
                        <span v-for="tag in article.tags" :key="tag">
                            # {{ tag }}
                        </span>
                    </div>
                    <section class="detail-carousel">
                        <el-carousel trigger="click" height="350px">
                            <el-carousel-item v-for="item in articleImages" :key="item">
                                <img :src="item" />
                            </el-carousel-item>
                        </el-carousel>
                    </section>
                    <div class="note-copy">
                        <p>{{ article.content }}</p>
                    </div>
                    <div class="note-actions">
                        <button class="like-action" :class="{ 'is-liked': isLiked }" type="button" @click="handleLike">
                            <span>
                                <el-icon>
                                    <StarFilled />
                                </el-icon>
                            </span>
                            <b>{{ article.likes }}</b> {{ isLiked ? '已点赞' : '点赞' }}
                        </button>
                        <span class="comment-count">
                            <el-icon>
                                <Comment />
                            </el-icon>
                            <b>{{ articleComments.length }}</b> 条评论
                        </span>
                        <button class="share-btn" type="button">分享 <el-icon>
                                <Share />
                            </el-icon></button>
                    </div>
                </div>
                <section class="comment-section">
                    <div class="comment-heading">
                        <h2>评论 <span class="comment-total">{{ articleComments.length }}</span></h2>
                        <el-select v-model="commentSort" class="custom-select" style="width: 90px" :teleported="false">
                            <el-option label="按时间" value="time"></el-option>
                            <el-option label="按热度" value="likes"></el-option>
                        </el-select>
                    </div>
                    <form class="comment-form" @submit.prevent="handleCommentSubmit">
                        <img :src="currentUser?.avatar || '/avatar.png'" :alt="currentUser?.name || 'avatar'" />
                        <input v-model="commentInput" class="comment-input" type="text" maxlength="80"
                            placeholder="说点什么吧…" />
                        <button type="submit">发布</button>
                    </form>
                    <div id="comment-list" class="comment-list">
                        <SingleComment v-for="comment in articleComments" :key="comment.id" :comment="comment"
                            :user="author" />
                    </div>
                </section>
            </article>
            <div v-else class="messages-empty">
                <span class="messages-empty-icon">
                    <el-icon>
                        <Warning />
                    </el-icon>
                </span>
                <h2>没有这篇笔记</h2>
            </div>
        </div>
    </section>
</template>

<style scoped>
.note-detail-page {
    width: 100%;
}

.note-detail-layout {
    display: block;
}

.note-article {
    --detail-top-h: 350px;
    display: flex;
    flex-direction: column;
    padding: 18px;
    border: 1px solid #f0eeee;
    border-radius: 20px;
    background: var(--bg-surface);
    box-shadow: 0 10px 28px #422d300c;
    height: auto;
}

.note-main {
    display: grid;
    grid-template-columns: minmax(380px, 0.92fr) minmax(0, 1fr);
    grid-template-rows: auto auto auto minmax(0, 0.75fr) auto;
    column-gap: 34px;
    row-gap: 0;
    height: var(--detail-top-h);
    overflow: hidden;
}

.note-author-header {
    grid-column: 2;
    grid-row: 1;
    display: flex;
    align-items: center;
    gap: 10px;
}

.note-author-header>img {
    width: 38px;
    height: 38px;
    border: 1px solid #f0e7e7;
    border-radius: 50%;
    object-fit: cover;
}

.note-author-header>div {
    display: grid;
    gap: 3px;
}

.detail-author {
    color: var(--text-main);
    font-size: 13px;
}

.detail-date {
    color: var(--text-mute);
    font-size: 11px;
}

.follow-btn {
    margin-left: auto;
    line-height: 1;
    min-height: 30px;
    white-space: nowrap;
    padding: 6px 12px;
    border: 1px solid #ffd6dc;
    border-radius: 14px;
    color: #f65369;
    font: inherit;
    font-size: 12px;
    background: #fff5f6;
    cursor: pointer;
}

.detail-title {
    grid-column: 2;
    grid-row: 2;
    margin: 23px 0 0;
    color: var(--text-main);
    font-size: 26px;
    line-height: 1.3;
}

.detail-tags {
    grid-column: 2;
    grid-row: 3;
    display: flex;
    flex-wrap: nowrap;
    align-items: flex-start;
    align-self: start;
    gap: 7px;
    margin: 7px 0 12px;
    overflow-x: auto;
    scrollbar-width: none;
}

.detail-tags::-webkit-scrollbar {
    display: none;
}

.detail-tags span {
    flex: 0 0 auto;
    padding: 4px 8px;
    border-radius: 12px;
    color: var(--highlight);
    font-size: 11px;
    line-height: 1.2;
    white-space: nowrap;
    text-decoration: none;
    background: var(--highlight-soft);
}

.detail-carousel {
    grid-column: 1;
    grid-row: 1 / span 5;
    width: 100%;
    height: 100%;
    min-height: 0;
    position: relative;
    overflow: hidden;
    border-radius: 14px;
    background: var(--bg-soft);
}

.detail-carousel :deep(.el-carousel__container),
.detail-carousel :deep(.el-carousel__item),
.detail-carousel :deep(.el-carousel__item img) {
    width: 100%;
    height: 100%;
}

.detail-carousel :deep(.el-carousel__item img) {
    display: block;
    min-height: 100%;
    object-fit: cover;
}

.note-copy {
    grid-column: 2;
    grid-row: 4;
    min-height: 0;
    margin-top: 12px;
    overflow-y: auto;
    overscroll-behavior: contain;
    color: var(--text-prose);
    font-size: 13px;
    line-height: 1.8;
}

.note-copy p {
    margin: 0 0 8px;
}

.note-actions {
    grid-column: 2;
    grid-row: 5;
    display: flex;
    align-items: center;
    gap: 10px;
    margin-top: auto;
    padding-top: 15px;
}

.like-action,
.share-btn,
.comment-count {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    line-height: 1;
    min-height: 36px;
    white-space: nowrap;
    padding: 8px 11px;
    border: 1px solid #f0e7e7;
    border-radius: 18px;
    color: var(--text-prose);
    font: inherit;
    font-size: 12px;
    text-decoration: none;
    background: var(--bg-surface);
    cursor: pointer;
}

.like-action>span,
.like-action :deep(.el-icon) {
    display: inline-flex;
    align-items: center;
    justify-content: center;
}

.like-action>span {
    line-height: 0;
}

.like-action :deep(.el-icon) {
    position: relative;
}

.comment-section {
    display: flex;
    flex-direction: column;
    margin-top: 10px;
    padding-top: 10px;
    border-top: 1px solid #f1eeee;
}

.comment-heading {
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.comment-heading h2 {
    margin: 0;
    color: var(--text-main);
    font-size: 18px;
}

.comment-total {
    color: var(--text-mute);
    font-size: 13px;
    font-weight: 400;
}

.custom-select :deep(.el-select__wrapper) {
    padding: 6px 12px;
    border: 1px solid #f0e7e7;
    border-radius: 14px;
    color: var(--text-soft);
    font: inherit;
    font-size: 12px;
    background: var(--bg-surface);
    cursor: pointer;
    --el-border-color: transparent;
    --el-select-input-focus-border-color: black;
}

.custom-select :deep(.el-select__wrapper:hover) {
    box-shadow: 0 0 0 1px transparent inset !important;
}

.custom-select :deep(.el-select__wrapper.is-focused) {
    box-shadow: 0 0 0 1px transparent inset !important;
}

.custom-select .el-select-dropdown__item.is-selected {
    color: var(--el-text-color-regular) !important;
}

.comment-form {
    display: flex;
    align-items: center;
    gap: 10px;
    margin: 11px 0;
    padding: 8px;
    border-radius: 13px;
    background: var(--bg-soft);
}

.comment-form img {
    width: 30px;
    height: 30px;
    border-radius: 50%;
    object-fit: cover;
}

.comment-input {
    flex: 1;
    min-width: 0;
    border: 0;
    outline: 0;
    color: var(--text-prose);
    font: inherit;
    font-size: 13px;
    background: transparent;
}

.comment-form button {
    padding: 7px 13px;
    border: 0;
    border-radius: 14px;
    color: #fff;
    font: inherit;
    font-size: 12px;
    background: var(--accent-border-2);
    cursor: pointer;
}

.comment-list {
    display: grid;
    gap: 14px;
    overflow: visible;
}

.comment-list::-webkit-scrollbar {
    width: 6px;
}

.comment-list::-webkit-scrollbar-track {
    background: transparent;
}

.comment-list::-webkit-scrollbar-thumb {
    border-radius: 999px;
    background: var(--border);
}

.messages-empty {
    min-height: 420px;
    display: grid;
    place-content: center;
    justify-items: center;
    margin-top: 22px;
    border: 1px solid #f0eeee;
    border-radius: 20px;
    color: var(--text-mute);
    background: var(--bg-surface);
    box-shadow: 0 10px 28px #422d300c;
}

.messages-empty-icon {
    width: 52px;
    height: 52px;
    display: grid;
    place-items: center;
    border-radius: 50%;
    color: #f65369;
    font-size: 25px;
    background: var(--accent-soft);
}

.messages-empty h2 {
    margin: 16px 0 5px;
    color: var(--text-main);
    font-size: 17px;
}

.messages-empty p {
    margin: 0;
    font-size: 13px;
}
</style>