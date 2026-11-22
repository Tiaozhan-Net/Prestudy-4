<script setup lang="ts">
import { ref, computed } from 'vue'

import SearchBox from '../components/SearchBox.vue'
import CategorySelect from '../components/CategorySelect.vue'
import CarouselCard from '../components/CarouselCard.vue'
import PostCard from '../components/PostCard.vue'

import { Posts } from '../composables/api/posts'

const PostsList = ref(Posts())
const selectedCategory = ref('all')
const searchQuery = ref('')

const carouselPosts = computed(() => {
    const sortedList = [...PostsList.value].sort((a, b) => b.likes - a.likes)
    return sortedList.slice(0, 6)
})

const displayedPosts = computed(() => {
    const sortedList = [...PostsList.value].sort((a, b) => b.likes - a.likes)
    const categoryPosts = selectedCategory.value === 'all'
        ? sortedList
        : sortedList.filter(item => item.category === selectedCategory.value)
    const query = searchQuery.value.trim().toLowerCase()
    const filteredList = query
        ? categoryPosts.filter(item => [
            item.title,
            item.content,
            item.category,
            item.author,
            ...item.tags
        ].some(value => value.toLowerCase().includes(query)))
        : categoryPosts

    return selectedCategory.value === 'all' && !query ? filteredList.slice(0, 4) : filteredList
})

const status = ref(0)
const updateStatus = () => {
    status.value = selectedCategory.value === 'all' && !searchQuery.value.trim() ? 0 : 1
}

const selectCategory = (category: string) => {
    selectedCategory.value = category
    updateStatus()
}

const searchPosts = (query: string) => {
    searchQuery.value = query
    updateStatus()
}

const displayMore = () => {
    status.value = 1 - status.value;
}

const emptyText = computed(() => {
    if (searchQuery.value.trim()) return '没有找到符合条件的笔记，换个关键词试试'
    if (selectedCategory.value !== 'all') return '这个分类下还没有笔记'
    return '还没有人发布笔记，来写下第一篇吧'
})
</script>

<template>
    <header class="topbar">
        <SearchBox v-model="searchQuery" @update:model-value="searchPosts" />
        <CategorySelect v-model="selectedCategory" @update:model-value="selectCategory" />
    </header>
    <div class="carousel-wrapper" v-if="!status && carouselPosts.length > 0">
        <el-carousel :interval="3000" type="card" arrow="never" indicator-position="none">
            <el-carousel-item v-for="item in carouselPosts" :key="item.id">
                <CarouselCard :item="item" />
            </el-carousel-item>
        </el-carousel>
    </div>
    <section class="feed">
        <div class="feed-heading">
            <div>
                <p class="feed-eyebrow">EXPLORE MORE</p>
                <h2>推荐笔记</h2>
            </div>
            <button
                v-if="displayedPosts.length > 0 && !status"
                class="feed-more"
                type="button"
                @click="displayMore"
            >
                发现更多<el-icon>
                    <Right />
                </el-icon>
            </button>
            <button
                v-else-if="displayedPosts.length > 0"
                class="feed-more"
                type="button"
                @click="displayMore"
            >
                <el-icon>
                    <Back />
                </el-icon>返回
            </button>
        </div>
        <div class="note-grid">
            <PostCard v-for="item in displayedPosts" :key="item.id" :item="item" clickable />
            <p v-if="displayedPosts.length === 0" class="empty-state">
                {{ emptyText }}
            </p>
        </div>
    </section>
</template>

<style scoped>
/* 三段 */
.topbar {
    height: 51px;
    display: flex;
    align-items: center;
    gap: 13px;
}

.carousel-wrapper {
    width: 100%;
    height: clamp(270px, 36vh, 350px);
    margin: 30px auto 0 auto;
    box-sizing: border-box;
}

.feed {
    margin: clamp(22px, 3vh, 32px) 0 0 0;
    display: flex;
    flex-direction: column;
}


/* 走马灯 */
.carousel-wrapper :deep(.el-carousel__mask) {
    display: none !important;
}

.carousel-wrapper :deep(.el-carousel__item) {
    border-radius: 16px;
    overflow: hidden;
    background: transparent !important;
}

.carousel-wrapper :deep(.el-carousel__item .carousel-slide) {
    position: relative;
    transition: transform 0.4s ease;
}

.carousel-wrapper :deep(.el-carousel__item .carousel-slide::after) {
    content: '';
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.45);
    opacity: 0;
    transition: opacity 0.4s ease;
    pointer-events: none;
    z-index: 10;
}

.carousel-wrapper :deep(.el-carousel__item:not(.is-active) .carousel-slide::after) {
    opacity: 1;
}

.carousel-wrapper :deep(.el-carousel__item:not(.is-active) .carousel-slide) {
    transform: scale(0.92);
}

.el-carousel__item h3 {
    color: #475669;
    opacity: 0.75;
    line-height: 200px;
    margin: 0;
    text-align: center;
}

.el-carousel__item:nth-child(2n) {
    background-color: #99a9bf;
}

.el-carousel__item:nth-child(2n + 1) {
    background-color: #d3dce6;
}


/* FEED */
.feed-heading {
    display: flex;
    align-items: end;
    justify-content: space-between;
    margin-bottom: 18px;
}

.feed-eyebrow {
    margin: 0 0 5px;
    color: #c2b8b8;
    font-size: 10px;
    font-weight: 700;
    letter-spacing: 0.12em;
}

.feed-heading h2 {
    margin: 0;
    font-size: 22px;
}

.feed-more {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    margin: 0;
    padding: 0;
    border: 0;
    color: var(--text-mute);
    font: inherit;
    font-size: 13px;
    cursor: pointer;
    background: none;
    transition: color 0.2s;
}

.feed-more:hover {
    color: var(--accent);
}

.note-grid {
    display: grid;
    grid-template-columns: repeat(4, minmax(190px, 1fr));
    gap: 18px;
}

.empty-state {
    grid-column: 1 / -1;
    margin: 0;
    padding: 40px 0;
    color: var(--text-mute);
    font-size: 14px;
    text-align: center;
}
</style>