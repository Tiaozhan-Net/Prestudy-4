<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { UploadInstance, UploadUserFile } from 'element-plus'

import { currentUser } from '@/composables/auth'
import type { Article } from '@/composables/type'
import { addArticle } from '@/composables/data/articles'
import { compressImage, formatFileSize, IMAGE_FILE_MAX_SIZE, isImageFile, revokeBlobUrl } from '@/composables/image'
import PostCard from '@/components/PostCard.vue'

const router = useRouter()

const categories: Article['category'][] = ['音乐', '美食', '萌宠', '旅行', '游戏', '影视', '读书', '摄影', '编程', '其他']

const articleDraft = reactive<Omit<Pick<Article, 'title' | 'category' | 'images' | 'cover' | 'tags' | 'content'>, 'category'> & {
    category: Article['category'] | ''
}>({
    title: '',
    category: '',
    images: [],
    cover: '',
    tags: [],
    content: ''
})

const uploadFiles = ref<UploadUserFile[]>([])
const uploadRef = ref<UploadInstance>()
const tagInput = ref('')

/** 上传区是唯一的图片来源，预览卡片从它派生，避免出现"上传区有图、卡片没封面" */
const syncDraftImages = () => {
    articleDraft.images = uploadFiles.value
        .filter(file => file.status === 'success' && Boolean(file.url))
        .map(file => file.url as string)
    articleDraft.cover = articleDraft.images[0] || ''
}

const removeUploadFile = (target: UploadUserFile) => {
    const uid = target.uid

    uploadFiles.value = uploadFiles.value.filter(file => file !== target && (uid === undefined || file.uid !== uid))
}

/** 校验不通过的文件必须立刻从列表移除，否则它仍留在上传区，看起来像"大小限制没生效" */
const rejectUploadFile = (uploadFile: UploadUserFile, message: string) => {
    ElMessage.warning(message)
    revokeBlobUrl(uploadFile.url)
    removeUploadFile(uploadFile)
    syncDraftImages()
}

const updateImages = async (uploadFile: UploadUserFile) => {
    const file = uploadFile.raw

    if (!file) {
        rejectUploadFile(uploadFile, '图片读取失败，请重试')
        return
    }

    if (!isImageFile(file)) {
        rejectUploadFile(uploadFile, '只能上传图片文件')
        return
    }

    if (file.size > IMAGE_FILE_MAX_SIZE) {
        rejectUploadFile(uploadFile, `图片不能超过 ${formatFileSize(IMAGE_FILE_MAX_SIZE)}，当前 ${formatFileSize(file.size)}`)
        return
    }

    try {
        const image = await compressImage(file)
        // 压缩是异步的，期间用户可能已经删除或清空了图片
        const uploaded = uploadFiles.value.find(item => item.uid === uploadFile.uid)
        if (!uploaded) return

        revokeBlobUrl(uploaded.url)
        uploaded.url = image
        uploaded.status = 'success'
        syncDraftImages()
    } catch {
        ElMessage.error('图片读取失败，请重试')
        revokeBlobUrl(uploadFile.url)
        removeUploadFile(uploadFile)
        syncDraftImages()
    }
}

const removeImage = (uploadFile: UploadUserFile) => {
    revokeBlobUrl(uploadFile.url)
    removeUploadFile(uploadFile)
    syncDraftImages()
}

const updateTags = (value: string) => {
    tagInput.value = value
    articleDraft.tags = value.trim() ? value.trim().split(/\s+/) : []
}

const submitArticle = async () => {
    if (!currentUser.value) {
        ElMessage.warning('请先登录')
        return
    }

    if (!articleDraft.title.trim() || !articleDraft.content.trim() || !articleDraft.category) {
        ElMessage.warning('请填写标题、分类和正文')
        return
    }

    try {
        await addArticle({
            authorId: currentUser.value.id,
            cover: articleDraft.images[0] || '',
            title: articleDraft.title.trim(),
            images: [...articleDraft.images],
            content: articleDraft.content.trim(),
            tags: [...articleDraft.tags],
            category: articleDraft.category,
            likes: 0,
            comments: []
        })
        ElMessage.success('发布成功')
        resetDraft()
        await router.replace('/')
    } catch {
        ElMessage.error('发布失败，请重试')
    }
}

const resetDraft = () => {
    articleDraft.title = ''
    articleDraft.category = ''
    articleDraft.images = []
    articleDraft.cover = ''
    articleDraft.tags = []
    articleDraft.content = ''
    tagInput.value = ''
    uploadFiles.value.forEach(file => revokeBlobUrl(file.url))
    uploadFiles.value = []
}
</script>

<template>
    <section class="profile-page">
        <div class="profile-notes-heading">
            <div>
                <p>NEW POST</p>
                <h1>发布笔记</h1>
                <p id="author-line" style="margin: 6px 0 0; color: var(--text-mute); font-size: 12px">
                    作者：{{ currentUser?.name || '未登录' }}
                </p>
            </div>
        </div>
        <div class="edit-columns">
            <div class="profile-card">
                <div class="form-grid">
                    <label class="form-field">
                        <span>标题 <b class="required">*</b></span>
                        <el-input v-model="articleDraft.title" class="custom-input" maxlength="40" placeholder="请输入标题"
                            show-word-limit type="text" word-limit-position="outside" />
                    </label>
                    <label class="form-field">
                        <span>分类 <b class="required">*</b></span>
                        <el-select v-model="articleDraft.category" class="custom-select" style="width: 100%"
                            :teleported="false" placeholder="请选择分类">
                            <el-option v-for="category in categories" :key="category" :label="category"
                                :value="category" />
                        </el-select>
                    </label>
                    <div class="form-field is-media">
                        <span class="field-title">
                            图片
                            <small class="form-count">
                                <span class="cover-path-hint">第一张当作封面</span>
                                <b class="cover-count"></b>
                            </small>
                        </span>
                        <el-upload ref="uploadRef" v-model:file-list="uploadFiles" accept="image/*"
                            class="custom-upload" action="#" list-type="picture-card" :auto-upload="false"
                            :on-change="updateImages" :on-remove="removeImage">
                            <el-icon>
                                <Plus />
                            </el-icon> <template #file="{ file }">
                                <div>
                                    <img class="el-upload-list__item-thumbnail" :src="file.url" />
                                    <span class="el-upload-list__item-actions">
                                        <span class="el-upload-list__item-delete"
                                            @click.stop="uploadRef?.handleRemove(file)">
                                            <el-icon>
                                                <Delete />
                                            </el-icon>
                                        </span>
                                    </span>
                                </div>
                            </template>
                        </el-upload>
                    </div>
                    <label class="form-field">
                        <span>标签</span>
                        <el-input class="custom-input" style="width: 100%" placeholder="请输入标签" maxlength="40"
                            :model-value="tagInput" @update:model-value="updateTags" />
                        <small class="form-count">多个标签用空格分开</small>
                    </label>
                    <label class="form-field">
                        <span>正文 <b class="required">*</b></span>
                        <el-input v-model="articleDraft.content" class="custom-textarea" maxlength="500"
                            placeholder="请输入正文" show-word-limit type="textarea" word-limit-position="outside" />
                    </label>
                </div>
                <div class="form-actions">
                    <button class="btn-ghost" type="button" @click="resetDraft">
                        清空
                    </button>
                    <button class="form-submit form-submit-inline" type="button" @click="submitArticle">
                        发布笔记
                    </button>
                </div>
            </div>
            <aside>
                <p class="preview-label">实时预览</p>
                <PostCard :clickable="false" :item="{
                    id: 0,
                    author: currentUser?.name || '未登录',
                    avatar: currentUser?.avatar || '/avatar.png',
                    cover: articleDraft.cover,
                    title: articleDraft.title || '未命名笔记',
                    images: articleDraft.images,
                    tags: articleDraft.tags,
                    category: articleDraft.category || '其他',
                    likes: 0
                }" />
                <p class="preview-body" :class="{ 'is-placeholder': !articleDraft.content }">
                    {{ articleDraft.content || '正文预览会显示在这里' }}
                </p>
            </aside>
        </div>
    </section>
</template>

<style scoped>
.profile-page {
    width: 100%;
}

.profile-notes-heading {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin: 4px 0 17px;
}

.profile-notes-heading p {
    margin: 0 0 5px;
    color: #c2b8b8;
    font-size: 10px;
    font-weight: 700;
    letter-spacing: 0.12em;
}

.profile-notes-heading h1 {
    margin: 0;
    color: var(--text-main);
    font-size: 24px;
}

.profile-notes-heading h2 {
    margin: 0;
    color: var(--text-main);
    font-size: 21px;
}

.profile-notes-heading h2 span {
    margin-left: 4px;
    color: var(--text-mute);
    font-size: 13px;
    font-weight: 400;
}

.edit-columns {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 300px;
    gap: 20px;
    margin-top: 20px;
    align-items: start;
}

.profile-card {
    padding: 18px;
    border: 1px solid #f0eeee;
    border-radius: 16px;
    background: var(--bg-surface);
    box-shadow: 0 7px 20px #412b2b0a;
}

.form-grid {
    display: grid;
    gap: 14px;
}

.form-field {
    display: grid;
    gap: 6px;
}

.form-field>span {
    color: var(--text-soft);
    font-size: 12px;
    font-weight: 700;
}

.form-count {
    justify-self: end;
    color: var(--text-mute);
    font-size: 11px;
}

.cover-path-hint {
    color: var(--text-mute);
    font-size: 11px;
}

.cover-count {
    margin-left: 6px;
    color: #f65369;
    font-weight: 700;
}

.form-actions {
    display: flex;
    gap: 10px;
    margin-top: 16px;
}

.btn-ghost {
    padding: 10px 18px;
    border: 1px solid var(--border);
    border-radius: 10px;
    color: var(--text-soft);
    font: inherit;
    font-size: 13px;
    background: var(--bg-surface);
    cursor: pointer;
    transition:
        border-color 0.2s,
        color 0.2s,
        background 0.2s,
        transform 0.2s;
}

.btn-ghost:hover {
    border-color: var(--accent);
    color: #fff;
    background: var(--accent);
    transform: translateY(-3px);
}

.btn-ghost:active {
    border-color: var(--accent-strong);
    background: var(--accent-strong);
    transform: translateY(-1px);
}

.form-submit {
    width: 100%;
    padding: 11px;
    border: 0;
    border-radius: 10px;
    color: #fff;
    font: inherit;
    font-size: 14px;
    font-weight: 700;
    background: #ff3858;
    cursor: pointer;
    transition: filter 0.2s;
}

.form-submit:hover {
    filter: brightness(1.07);
}

.form-submit-inline {
    width: auto;
    padding: 10px 22px;
    transition:
        transform 0.2s ease,
        filter 0.2s ease;
}

.form-submit-inline:hover {
    transform: translateY(-3px);
}

.form-submit-inline:active {
    transform: translateY(-1px);
}

.preview-label {
    margin: 0 0 10px;
    color: #c2b8b8;
    font-size: 10px;
    font-weight: 700;
    letter-spacing: 0.12em;
}

.preview-body {
    margin: 16px 0 0;
    padding: 14px;
    border: 1px solid #f0eeee;
    border-radius: 14px;
    color: var(--text-prose);
    font-size: 13px;
    line-height: 1.8;
    white-space: pre-wrap;
    background: var(--bg-surface);
}

.preview-body.is-placeholder {
    color: var(--text-mute);
    font-style: italic;
}

.required {
    color: #e0313f;
}

.custom-input :deep(.el-input__wrapper) {
    width: 100%;
    padding: 5px 12px;
    border-radius: 10px;
    color: var(--text-main);
    font: inherit;
    font-size: 13px;
    background: var(--bg-soft);
    transition: border-color 0.2s;
    --el-input-focus-border-color: #e0313f;
}

.custom-input :deep(.el-input__count-inner) {
    margin-top: 6px;
    color: var(--text-mute);
    font-size: 11px;
}

.custom-textarea :deep(.el-textarea__inner) {
    width: 100%;
    min-height: 120px !important;
    padding: 10px 12px;
    border-radius: 10px;
    color: var(--text-main);
    font: inherit;
    font-size: 13px;
    background: var(--bg-soft);
    transition: border-color 0.2s;
    resize: none;
    --el-input-focus-border-color: #e0313f;
}

.custom-select :deep(.el-select__wrapper) {
    width: 100%;
    padding: 13px 12px;
    box-shadow: 0 0 0 1px var(--el-border-color) inset;
    border-radius: 10px;
    color: var(--text-main);
    background: var(--bg-soft);
}

.custom-select :deep(.el-select__wrapper.is-focused) {
    box-shadow: 0 0 0 1px #e0313f inset !important;
}

.custom-select .el-select-dropdown__item.is-selected {
    color: var(--el-text-color-regular) !important;
}

.custom-upload :deep(.el-upload--picture-card:hover) {
    border-color: #e0313f;
    color: #e0313f;
}

.custom-upload :deep(.el-upload--picture-card) {
    border-radius: 10px;
}

.custom-upload :deep(.el-upload-list__item),
.custom-upload :deep(.el-upload-list--picture-card) {
    border-radius: 10px;
}
</style>