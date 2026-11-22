<script setup lang="ts">
import { defineProps } from 'vue'

const props = defineProps<{
    item: {
        id: number
        cover: string
        title: string
        category: string
    }
}>()
</script>

<template>
    <article class="carousel-slide">
        <RouterLink :to="`/article/${item.id}`" class="note-open" :aria-label="`查看文章：${item.title}`" />
        <img :src="item.cover" />
        <div class="slide-overlay">
            <span class="slide-tag">{{ item.category }}</span>
            <h2>{{ item.title }}</h2>
        </div>
    </article>
</template>

<style scoped>
.carousel-slide {
    position: relative;
    width: 100%;
    height: 100%;
    border-radius: 16px;
    overflow: hidden;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    cursor: pointer;
    background-color: #f5f5f5;
}

.note-open {
    position: absolute;
    inset: 0;
    z-index: 10;
}

/* 3. 图片样式 */
.carousel-slide img {
    width: 100%;
    height: 100%;
    display: block;
    object-fit: cover;
    transition: transform 0.3s ease;
}

.carousel-slide:hover img {
    transform: scale(1.03);
}

.slide-overlay {
    position: absolute;
    inset: auto 0 0;
    padding: 46px 32px 28px;
    color: #fff;
    background: linear-gradient(transparent, #000000ad);
}

.slide-tag {
    display: inline-block;
    margin-bottom: 8px;
    padding: 3px 10px;
    border-radius: 20px;
    font-size: 12px;
    background: #ffffff40;
}

.slide-overlay h2 {
    margin: 0;
    font-size: 23px;
}
</style>