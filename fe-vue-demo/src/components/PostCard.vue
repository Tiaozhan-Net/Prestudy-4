<script setup lang="ts">
import { RouterLink } from 'vue-router'

const props = defineProps<{
  clickable?: boolean
  item: {
    id: number
    author: string
    avatar: string
    cover: string
    title: string
    images: string[]
    tags: string[]
    category: string
    likes: number
  }
}>()

const clickable = props.clickable ?? true
</script>

<template>
  <RouterLink :to="`/article/${item.id}`" custom v-slot="{ navigate }">
    <article class="feed-card card-tall" :class="{ 'is-clickable': clickable }"
      @click="clickable ? navigate() : undefined">
      <div class="feed-card-cover">
        <img v-if="item.cover" :src="item.cover" />
        <span v-else class="cover-placeholder">无封面</span>
      </div>
      <div class="feed-card-body">
        <h3>{{ item.title }}</h3>
        <div class="tag-group">
          <span v-for="tag in item.tags" :key="tag" class="note-tag">
            # {{ tag }}
          </span>
        </div>
        <div class="feed-card-meta">
          <span class="note-author"><img :src="item.avatar" :alt="item.author" />{{ item.author }}</span>
          <button class="note-like" type="button" data-action="like">
            <span class="like-heart">
              <el-icon>
                <StarFilled />
              </el-icon>
            </span>
            <span data-likes>{{ item.likes }}</span>
          </button>
        </div>
      </div>
    </article>
  </RouterLink>
</template>

<style scoped>
.feed-card {
  position: relative;
  overflow: hidden;
  border: 1px solid #f0eeee;
  border-radius: 16px;
  background: var(--bg-surface);
  box-shadow: 0 7px 20px #412b2b0c;
  transition: 0.2s;
}

.feed-card:hover {
  transform: translateY(-3px);
}

.feed-card.is-clickable {
  cursor: pointer;
}

.feed-card img {
  width: 100%;
  display: block;
  object-fit: cover;
}

.feed-card-cover {
  display: flex;
  align-items: center;
  justify-content: center;
  aspect-ratio: 1 / 1;
  overflow: hidden;
  background: var(--bg-soft);
}

.cover-placeholder {
  color: var(--text-mute);
  font-size: 13px;
}

.feed-card-cover img {
  width: 100%;
  height: 100%;
}

.feed-card-body {
  padding: 12px 13px 13px;
}

.feed-card h3 {
  display: -webkit-box;
  min-height: 40px;
  margin: 0;
  overflow: hidden;
  color: var(--text-main);
  font-size: 14px;
  line-height: 1.45;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.feed-card-meta {
  display: flex;
  align-items: center;
  /* 第 1 课在 qzone-carousel.css 里补的 */
  justify-content: space-between;
  margin-top: 9px;
  color: var(--text-mute);
  font-size: 11px;
}

.feed-card-meta span:last-child {
  color: #f36a7d;
}

.tag-group {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 6px;
}

.note-tag {
  margin-top: 5px;
  overflow: hidden;
  color: var(--highlight);
  font-size: 11px;
  line-height: 1.25;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.note-tag {
  display: inline-block;
  width: fit-content;
  padding: 2px 8px;
  margin-top: 0;
  border-radius: 8px;
  background: var(--highlight-soft);
}

.note-author {
  display: inline-flex;
  min-width: 0;
  align-items: center;
  gap: 6px;
  overflow: hidden;
  color: var(--text-mute);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.note-author img {
  width: 18px;
  height: 18px;
  flex: 0 0 18px;
  border: 1px solid #f1eaea;
  border-radius: 50%;
  object-fit: cover;
}

.note-like {
  flex: 0 0 auto;
  color: #f36a7d;
  display: inline-flex;
  align-items: center;
  gap: 3px;
  padding: 2px 6px;
  border: 0;
  border-radius: 8px;
  color: #f36a7d;
  font: inherit;
  font-size: 11px;
  background: transparent;
  cursor: pointer;
  transition:
    background 0.2s,
    transform 0.15s;
}

.note-like:hover {
  background: var(--accent-soft);
}

.note-like:active {
  transform: scale(0.94);
}

.like-heart {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
</style>