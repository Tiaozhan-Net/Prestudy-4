<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { CURRENT_USER_STORAGE_KEY, currentUser, logoutUser } from '@/composables/auth'
import type { User } from '@/composables/type'
import PostCard from '@/components/PostCard.vue'
import { Posts } from '@/composables/api/posts'
import { saveUsers, users } from '@/composables/data/users'
import { setStoredValue } from '@/composables/data/storage'
import { compressImage, formatFileSize, IMAGE_FILE_MAX_SIZE, isImageFile } from '@/composables/image'

const router = useRouter()
const isEditing = ref(false)
const avatarInput = ref<HTMLInputElement | null>(null)
const profileName = ref('')
const profileDescription = ref('')
const profileAbout = ref('')
const profileLocation = ref('')
const profileInterests = ref('')
const profilePassword = ref('')

const user = computed<User>(() => {
    if (!currentUser.value) {
        return {
            id: 0,
            name: '未登录',
            avatar: '/avatar.png',
            createTime: '',
            description: '请登录后查看信息',
            interests: []
        }
    }

    return currentUser.value
})

const profileAvatar = ref('')
const savedAvatar = ref('')

const startEditing = () => {
    profileName.value = user.value.name
    profileDescription.value = user.value.description || ''
    profileAbout.value = user.value.about || ''
    profileLocation.value = user.value.location || ''
    profileInterests.value = (user.value.interests || []).join(' ')
    profilePassword.value = ''
    isEditing.value = true
}

const openAvatarPicker = () => {
    if (isEditing.value) avatarInput.value?.click()
}

const updateAvatar = async (event: Event) => {
    const input = event.target as HTMLInputElement
    const file = input.files?.[0]
    // 先取出文件再清空，保证校验失败后再次选择同一张图片仍能触发 change
    input.value = ''

    if (!file) return

    if (!isImageFile(file)) {
        ElMessage.warning('请选择图片文件')
        return
    }

    if (file.size > IMAGE_FILE_MAX_SIZE) {
        ElMessage.warning(`头像图片不能超过 ${formatFileSize(IMAGE_FILE_MAX_SIZE)}，当前 ${formatFileSize(file.size)}`)
        return
    }

    try {
        profileAvatar.value = await compressImage(file, 256)
    } catch {
        ElMessage.error('图片读取失败，请重试')
    }
}

const saveProfile = async () => {
    try {
        if (currentUser.value) {
            const interests = [...new Set(profileInterests.value
                .trim()
                .split(/\s+/)
                .map(interest => interest.trim())
                .filter(Boolean))]
            const updatedUser = {
                ...currentUser.value,
                name: profileName.value.trim() || currentUser.value.name,
                description: profileDescription.value.trim(),
                about: profileAbout.value.trim(),
                location: profileLocation.value.trim(),
                interests,
                ...(profilePassword.value.trim() ? { password: profilePassword.value.trim() } : {}),
                ...(profileAvatar.value ? { avatar: profileAvatar.value } : {})
            }

            if (updatedUser.name !== currentUser.value.name
                || updatedUser.description !== (currentUser.value.description || '')
                || updatedUser.about !== (currentUser.value.about || '')
                || updatedUser.location !== (currentUser.value.location || '')
                || JSON.stringify(updatedUser.interests) !== JSON.stringify(currentUser.value.interests || [])
                || Boolean(profilePassword.value.trim())
                || profileAvatar.value) {
                await setStoredValue(CURRENT_USER_STORAGE_KEY, updatedUser)
                currentUser.value = updatedUser
                if (profileAvatar.value) savedAvatar.value = profileAvatar.value

                const userIndex = users.value.findIndex(item => item.id === updatedUser.id)
                if (userIndex !== -1) {
                    users.value[userIndex] = updatedUser
                    await saveUsers()
                }
            }
        }
    } catch {
        ElMessage.error('修改失败')
        return
    }

    profileAvatar.value = ''
    profilePassword.value = ''
    isEditing.value = false
    ElMessage.success('个人信息已保存')
}

const cancelEditing = () => {
    profileName.value = user.value.name
    profileDescription.value = user.value.description || ''
    profileAbout.value = user.value.about || ''
    profileLocation.value = user.value.location || ''
    profileInterests.value = (user.value.interests || []).join(' ')
    profilePassword.value = ''
    profileAvatar.value = ''
    isEditing.value = false
}

const joinDate = computed(() => {
    if (!user.value.createTime) return '暂无'

    const date = new Date(user.value.createTime)
    if (Number.isNaN(date.getTime())) return '暂无'

    return `${date.getUTCFullYear()} 年 ${date.getUTCMonth() + 1} 月`
})

const userPosts = computed(() => {
    return Posts().filter(post => post.authorId === user.value.id)
})

const sortOrder = ref<'date-desc' | 'date-asc' | 'likes-desc'>('date-desc')

const sortedUserPosts = computed(() => {
    return [...userPosts.value].sort((firstPost, secondPost) => {
        if (sortOrder.value === 'likes-desc') {
            return secondPost.likes - firstPost.likes
        }

        const firstTime = new Date(firstPost.createTime).getTime()
        const secondTime = new Date(secondPost.createTime).getTime()

        return sortOrder.value === 'date-desc'
            ? secondTime - firstTime
            : firstTime - secondTime
    })
})

const userLikes = computed(() => {
    return userPosts.value.reduce((total, post) => total + post.likes, 0)
})

const userCategories = computed(() => {
    return new Set(userPosts.value.map(post => post.category)).size
})

const latestPostDate = computed(() => {
    const latestPost = userPosts.value.reduce<typeof userPosts.value[number] | null>((latest, post) => {
        if (!latest || new Date(post.createTime).getTime() > new Date(latest.createTime).getTime()) {
            return post
        }

        return latest
    }, null)

    if (!latestPost) return '暂无'

    const date = new Date(latestPost.createTime)
    if (Number.isNaN(date.getTime())) return '暂无'

    return `${date.getUTCFullYear()} 年 ${date.getUTCMonth() + 1} 月 ${date.getUTCDate()} 日`
})

const logout = async () => {
    await logoutUser()
    ElMessage.success('已退出登录')
    await router.replace('/login')
}
</script>

<template>
    <section class="profile-page">
        <header class="profile-hero">
            <div class="profile-avatar-wrap" :class="{ 'is-editing': isEditing }" @click="openAvatarPicker">
                <img class="profile-avatar" :src="profileAvatar || savedAvatar || user.avatar" />
                <input ref="avatarInput" type="file" accept="image/*" hidden @change="updateAvatar" />
            </div>
            <div class="profile-identity">
                <div class="profile-fields">
                    <div class="profile-name-row">
                        <span v-if="isEditing" class="edit-field-label">昵称</span>
                        <input v-if="isEditing" v-model="profileName" class="profile-name-input" maxlength="20" />
                        <div v-else class="nickname-measure">{{ user.name }}</div>
                        <span v-if="!isEditing" class="title-label">个人创作者</span>
                    </div>
                    <div v-if="isEditing" class="profile-description-field">
                        <span class="edit-field-label">简介</span>
                        <input v-model="profileDescription" class="profile-description-input"
                            placeholder="这个人很懒，什么都没留下" maxlength="32" />
                    </div>
                    <p v-else class="editable-bio">{{ user.description || '这个人很懒，什么都没留下' }}</p>
                </div>
                <div class="profile-stats">
                    <span><b id="stat-articles">{{ userPosts.length }}</b> 笔记</span>
                    <span><b id="stat-likes">{{ userLikes }}</b> 获赞</span>
                    <span><b id="stat-categories">{{ userCategories }}</b> 分类</span>
                </div>
            </div>
            <button v-if="!isEditing" class="profile-edit" type="button" @click="startEditing">
                编辑资料
            </button>
            <button v-if="isEditing" class="profile-edit is-primary" type="button" @click="saveProfile">
                保存
            </button>
            <button v-if="isEditing" class="profile-edit" type="button" @click="cancelEditing">
                取消
            </button>

            <button v-if="!isEditing" class="profile-edit" type="button" @click="logout">
                退出登录
            </button>
        </header>
        <div class="profile-content">
            <aside class="profile-sidebar-content">
                <section class="profile-card about-card">
                    <p class="profile-card-eyebrow">ABOUT ME</p>
                    <h2>关于我</h2>
                    <textarea v-if="isEditing" v-model="profileAbout" class="profile-about-input"
                        placeholder="介绍一下自己吧" maxlength="120" rows="1"></textarea>
                    <p v-else class="editable-about">{{ user.about || '还没有填写个人介绍' }}</p>
                    <dl class="profile-details">
                        <div>
                            <dt>所在地</dt>
                            <dd>
                                <input v-if="isEditing" v-model="profileLocation" class="profile-location-input"
                                    maxlength="40" />
                                <span v-else id="location-text">{{ user.location }}</span>
                            </dd>
                        </div>
                        <div>
                            <dt>加入时间</dt>
                            <dd>{{ joinDate }}</dd>
                        </div>
                        <div>
                            <dt>最近发布</dt>
                            <dd>{{ latestPostDate }}</dd>
                        </div>

                    </dl>
                </section>
                <section class="profile-card interest-card">
                    <p class="profile-card-eyebrow">MY INTERESTS</p>
                    <h2>我的兴趣</h2>
                    <textarea v-if="isEditing" v-model="profileInterests" class="profile-interests-input"
                        placeholder="例如：摄影 阅读旅行" maxlength="200"></textarea>
                    <div v-else class="interest-tags">
                        <span class="interest-tag" v-for="interest in user.interests" :key="interest">{{ interest
                            }}</span>
                        <span v-if="!user.interests?.length" class="interest-placeholder">暂无兴趣</span>
                    </div>
                </section>
                <section v-if="isEditing" class="profile-card">
                    <p class="profile-card-eyebrow">PASSWORD</p>
                    <h2>修改密码</h2>
                    <input v-model="profilePassword" type="password" autocomplete="new-password"
                        placeholder="为空则不修改密码"
                        class="profile-description-input" style="margin-top: 13px;" maxlength="20" />
                </section>
            </aside>
            <section class="profile-notes">
                <div class="profile-notes-heading">
                    <div>
                        <p>MY POSTS</p>
                        <h2>我的笔记 <span id="note-total">{{ userPosts.length }}</span></h2>
                    </div>
                    <el-select v-model="sortOrder" class="sort-select" :teleported="false">
                        <el-option label="最新发布" value="date-desc" />
                        <el-option label="最早发布" value="date-asc" />
                        <el-option label="点赞最多" value="likes-desc" />
                    </el-select>
                </div>
                <div class="profile-note-columns">
                    <PostCard v-for="post in sortedUserPosts" :key="post.id" :item="post" clickable />
                    <p v-if="sortedUserPosts.length === 0" class="profile-empty">
                        还没有发布笔记，去写下第一篇吧
                    </p>
                </div>
            </section>
        </div>
    </section>
</template>

<style scoped>
.profile-page {
    width: 100%;
}

.profile-hero {
    display: flex;
    align-items: center;
    gap: 18px;
    padding: 25px 28px;
    border: 1px solid #f0eeee;
    border-radius: 20px;
    background: linear-gradient(112deg, var(--bg-surface) 50%, #fff5f6);
    box-shadow: 0 10px 28px #422d300c;
}

.profile-avatar-wrap {
    position: relative;
    flex: 0 0 auto;
    width: 82px;
    height: 82px;
}

.profile-avatar-wrap.is-editing {
    cursor: pointer;
}

.profile-avatar-wrap.is-editing .profile-avatar {
    border: 1px solid var(--accent-border-2);
}

.profile-avatar {
    width: 82px;
    height: 82px;
    flex: 0 0 auto;
    border: 1px solid var(--bg-surface);
    border-radius: 50%;
    object-fit: cover;
    box-shadow: 0 6px 16px #d7647630;
}

.profile-identity {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
}

.profile-fields {
    display: grid;
    grid-template-columns: 44px minmax(0, 1fr);
    align-items: stretch;
    column-gap: 10px;
    row-gap: 8px;
}

.profile-name-row {
    grid-column: 1 / -1;
    display: flex;
    align-items: center;
    gap: 10px;
}

.nickname-measure {
    width: auto;
    height: 35px;
    box-sizing: border-box;
    padding: 2px 6px;
    border: 1px solid transparent;
    border-radius: 6px;
    white-space: pre;
    font: inherit;
    font-size: 25px;
    font-weight: 700;
}

.profile-name-input {
    width: 100%;
    height: 35px;
    box-sizing: border-box;
    padding: 2px 6px;
    border: 1px solid var(--accent-border-2);
    border-radius: 6px;
    outline: none;
    color: var(--text-main);
    font: inherit;
    font-size: 25px;
    font-weight: 700;
    background: var(--bg-soft);
}

.title-label {
    padding: 3px 8px;
    border-radius: 10px;
    color: #f65369;
    font-size: 11px;
    background: var(--accent-soft);
}

.profile-field-label {
    flex: 0 0 auto;
    color: var(--text-mute);
    font-size: 12px;
}

.edit-field-label {
    flex: 0 0 auto;
    margin-left: 2px;
    color: var(--accent-text);
    font-size: 11px;
    font-weight: 700;
    letter-spacing: 0.04em;
}

.editable-bio {
    grid-column: 1 / -1;
    margin: 0;
    box-sizing: border-box;
    height: 31px;
    padding: 4px 7px;
    border: 1px solid transparent;
    border-radius: 6px;
    color: var(--text-mute);
    font-size: 13px;
    line-height: 1.5;
    resize: none;
    overflow: hidden;
}

.profile-description-input {
    width: 100%;
    height: 31px;
    box-sizing: border-box;
    padding: 5px 7px;
    border: 1px solid var(--accent-border-2);
    border-radius: 6px;
    outline: none;
    resize: vertical;
    color: var(--text-prose);
    font: inherit;
    font-size: 13px;
    line-height: 1.5;
    background: var(--bg-soft);
}

.profile-description-field {
    grid-column: 1 / -1;
    display: flex;
    align-items: center;
    gap: 10px;
}

.profile-description-field .profile-description-input {
    flex: 1;
    min-width: 0;
}

.profile-stats {
    grid-column: 1 / -1;
    display: flex;
    gap: 19px;
    margin-top: 13px;
    color: var(--text-mute);
    font-size: 12px;
}

.profile-stats b {
    margin-right: 3px;
    color: var(--text-main);
    font-size: 14px;
}

.profile-edit {
    margin-left: auto;
    padding: 9px 17px;
    border: 1px solid #ffd6dc;
    border-radius: 18px;
    color: #f65369;
    font: inherit;
    font-size: 13px;
    background: var(--bg-surface);
    cursor: pointer;
}

.profile-content {
    display: grid;
    grid-template-columns: 270px minmax(0, 1fr);
    gap: 26px;
    margin-top: 25px;
    align-items: start;
}

.profile-sidebar-content {
    display: grid;
    gap: 15px;
}

.profile-card {
    min-width: 0;
    padding: 18px;
    border: 1px solid #f0eeee;
    border-radius: 16px;
    background: var(--bg-surface);
    box-shadow: 0 7px 20px #412b2b0a;
}

.profile-notes-heading {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin: 4px 0 17px;
}

.profile-card-eyebrow,
.profile-notes-heading p {
    margin: 0 0 5px;
    color: #c2b8b8;
    font-size: 10px;
    font-weight: 700;
    letter-spacing: 0.12em;
}

.profile-card h2 {
    margin: 0;
    color: var(--text-main);
    font-size: 16px;
}

.profile-notes-heading h2 span {
    margin-left: 4px;
    color: var(--text-mute);
    font-size: 13px;
    font-weight: 400;
}

.about-copy {
    margin: 12px 0 15px;
    color: var(--text-mute);
    font-size: 12px;
    line-height: 1.7;
}

.profile-about-input {
    width: 100%;
    height: 55px;
    box-sizing: border-box;
    margin: 12px 0 15px;
    padding: 7px;
    border: 1px solid var(--accent-border-2);
    border-radius: 6px;
    outline: none;
    resize: vertical;
    color: var(--text-prose);
    font: inherit;
    font-size: 12px;
    line-height: 1.7;
    background: var(--bg-soft);
    resize: none;
}

.editable-about {
    min-width: 0;
    margin: 12px 0 15px;
    box-sizing: border-box;
    height: 55px;
    padding: 7px;
    border: 1px solid transparent;
    border-radius: 6px;
    overflow-wrap: anywhere;
    word-break: break-word;
    color: var(--text-mute);
    font-size: 12px;
    line-height: 1.7;
}

.profile-location-input {
    width: 100%;
    max-width: 120px;
    height: 27px;
    box-sizing: border-box;
    padding: 3px 5px;
    border: 1px solid var(--accent-border-2);
    border-radius: 6px;
    outline: none;
    color: var(--text-prose);
    font: inherit;
    font-size: 12px;
    background: var(--bg-soft);
}

.profile-details {
    display: grid;
    gap: 9px;
    margin: 0;
}

.profile-details div {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    min-width: 0;
    font-size: 12px;
}

.profile-details dt {
    flex: 0 0 auto;
    color: var(--text-mute);
}

.profile-details dd {
    min-width: 0;
    max-width: 130px;
    flex: 1 1 auto;
    height: 27px;
    box-sizing: border-box;
    margin: 0;
    padding: 3px 5px;
    border: 1px solid transparent;
    border-radius: 6px;
    color: var(--text-prose);
    text-align: right;
}

.profile-details dd:has(.profile-location-input) {
    width: 130px;
    max-width: 130px;
}

.interest-tags {
    display: flex;
    flex-wrap: wrap;
    align-items: flex-start;
    gap: 7px;
    min-height: 50px;
    box-sizing: border-box;
    margin-top: 13px;
    padding: 7px;
    border: 1px solid transparent;
    border-radius: 6px;
}

.interest-tag {
    padding: 5px 9px;
    border-radius: 12px;
    color: var(--highlight);
    font-size: 11px;
    background: var(--highlight-soft);
}

.interest-placeholder {
    padding: 5px 9px;
    border: 1px solid transparent;
    border-radius: 12px;
    color: var(--text-mute);
    font-size: 11px;
}

.profile-interests-input {
    width: 100%;
    height: 50px;
    box-sizing: border-box;
    margin-top: 13px;
    padding: 7px;
    border: 1px solid var(--accent-border-2);
    border-radius: 6px;
    outline: none;
    resize: vertical;
    color: var(--text-prose);
    font: inherit;
    font-size: 12px;
    line-height: 1.7;
    background: var(--bg-soft);
    resize: none;
}

.sort-select {
    width: max-content !important;
    min-width: 0 !important;
}

.sort-select :deep(.el-select__wrapper) {
    width: 94px !important;
    min-width: 0 !important;
    padding: 6px 12px;
    border: 1px solid #f0e7e7;
    border-radius: 14px;
    color: var(--text-soft);
    font: inherit;
    font-size: 12px;
    background: var(--bg-surface);
    cursor: pointer;
}

.sort-select :deep(.el-select__wrapper:hover) {
    box-shadow: 0 0 0 1px transparent inset !important;
}

.sort-select :deep(.el-select__wrapper.is-focused) {
    box-shadow: 0 0 0 1px transparent inset !important;
}

.sort-select .el-select-dropdown__item.is-selected {
    color: var(--el-text-color-regular) !important;
}

.profile-empty {
    grid-column: 1 / -1;
    margin: 0;
    padding: 46px 0;
    border: 1px dashed #f0dfe2;
    border-radius: 16px;
    color: var(--text-mute);
    background: var(--bg-surface);
    font-size: 13px;
    text-align: center;
}

.profile-note-columns {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 18px;
}
</style>