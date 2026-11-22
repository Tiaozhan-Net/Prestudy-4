<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import { currentUser } from '../composables/auth'
import { isDark, toggleTheme } from '../composables/theme'

const route = useRoute()

const defaultUser = {
    name: '未登录',
    avatar: '/avatar.png',
    description: '请登录后查看信息'
}

const user = computed(() => currentUser.value
    ? {
        name: currentUser.value.name,
        avatar: currentUser.value.avatar,
        description: currentUser.value.description || '这个人很懒，什么都没留下'
    }
    : defaultUser)
</script>

<template>
    <aside class="sidebar">
        <div class="brand">
            <span class="brand-inner">
                <img src="/book.svg" />
            </span>
            Tenzors
        </div>
        <RouterLink class="sidebar-user" :to="currentUser ? '/profile' : '/login'">
            <span class="user-revamp">
                <img class="user-avatar" data-me="avatar" alt="avatar" :src="user.avatar" />
                <span class="user-info">
                    <strong data-me="nickname">{{ user.name }}</strong>
                    <small>{{ user.description }}</small>
                </span>
                <i class="user-arrow">
                    <el-icon>
                        <ArrowRightBold />
                    </el-icon>
                </i>
            </span>
        </RouterLink>
        <nav class="sidebar-nav">
            <RouterLink to="/">
                <a class="nav-item" :class="{ active: route.path === '/' }">
                    <span class="nav-icon">
                        <el-icon>
                            <HomeFilled />
                        </el-icon>
                    </span>
                    首页
                </a>
            </RouterLink>
            <RouterLink to="/message">
                <a class="nav-item" :class="{ active: route.path === '/message' }">
                    <span class="nav-icon">
                        <el-icon>
                            <BellFilled />
                        </el-icon>
                    </span>
                    消息
                </a>
            </RouterLink>
            <RouterLink :to="currentUser ? '/profile' : '/login'">
                <a class="nav-item" :class="{ active: route.path === '/profile' }">
                    <span class="nav-icon">
                        <el-icon>
                            <UserFilled />
                        </el-icon>
                    </span>
                    个人中心
                </a>
            </RouterLink>
        </nav>
        <div class="sidebar-bottom">
            <button class="theme-toggle" type="button" @click="toggleTheme">
                <span class="theme-icon">
                    <template v-if="isDark">
                        <el-icon>
                            <Sunny />
                        </el-icon>
                    </template>
                    <template v-else>
                        <el-icon>
                            <Moon />
                        </el-icon>
                    </template>
                </span>
                <span class="theme-text">{{ isDark ? '白天模式' : '黑夜模式' }}</span>
            </button>
            <RouterLink :to="currentUser ? '/editing' : '/login'">
                <a class="publish">
                    <span class="publish-icon">
                        <el-icon>
                            <Plus />
                        </el-icon>
                    </span>
                    <span class="publish-text">发布笔记</span>
                </a>
            </RouterLink>
        </div>
    </aside>
</template>

<style scoped>
/* 总样式 */
.sidebar {
    position: fixed;
    inset: 0 auto 0 0;
    width: 220px;
    padding: 36px 22px 28px;
    display: flex;
    flex-direction: column;
    background: var(--bg-surface);
    border-right: 1px solid #eee;
}


/* 四段样式 */
.brand {
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-main);
    font-size: 23px;
    font-weight: 700;
    text-decoration: none;
}

.sidebar-user {
    display: block;
    margin-top: 32px;
    padding: 10px;
    border: 1px solid #f1eeee;
    border-radius: 16px;
    color: var(--text-main);
    text-decoration: none;
    background: linear-gradient(135deg, var(--bg-surface), #fff8f8);
    box-shadow: 0 5px 16px #412b2b08;
    transition: 0.2s;
}

.sidebar-nav {
    display: grid;
    gap: 12px;
    margin-top: 24px;
}

.sidebar-bottom {
    margin-top: auto;
}


/* Brand */
.brand-inner {
    width: 29px;
    height: 29px;
    margin-right: 8px;
}

.brand-inner img {
    width: 100%;
    height: 100%;
    display: block;
}


/* SidebarUser */
.sidebar-user:hover {
    border-color: var(--border-strong);
    background: var(--bg-hover);
    transform: translateY(-1px);
}

.user-revamp {
    display: flex;
    align-items: center;
    gap: 10px;
}

.user-avatar {
    width: 44px;
    height: 44px;
    flex: 0 0 auto;
    border: 1px solid #eee;
    border-radius: 50%;
    object-fit: cover;
}

.user-info {
    flex: 1;
    min-width: 0;
}

.sidebar-user strong,
.sidebar-user small {
    display: block;
}

.sidebar-user strong {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.sidebar-user small {
    margin-top: 4px;
    color: var(--text-mute);
    font-size: 10px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.user-arrow {
    display: inline-flex;
    flex: 0 0 auto;
    align-items: center;
}


/* SidebarNav */
.sidebar-nav .nav-item:hover {
    color: #fff;
    background: var(--accent);
    transform: translateY(-3px);
}

.sidebar-nav .nav-item.active {
    color: var(--accent-text);
    font-weight: 700;
    background: var(--accent-soft);
}

.sidebar-nav .nav-item.active:hover {
    color: #fff;
    background: var(--accent);
    transform: translateY(-3px);
}

.nav-item {
    display: flex;
    align-items: center;
    padding: 14px 15px;
    border-radius: 12px;
    color: var(--text-soft);
    font-size: 16px;
    text-decoration: none;
    transition: 0.2s;
}

.nav-icon {
    display: inline-flex;
    align-items: center;
    width: 27px;
}


/* SidebarBottom */
.theme-toggle {
    width: 100%;
    margin-bottom: 8px;
    padding: 11px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    border: 1px solid var(--accent-border);
    border-radius: 10px;
    color: var(--accent-text);
    font: inherit;
    font-size: 14px;
    background: var(--bg-surface);
    cursor: pointer;
    transition:
        border-color 0.2s,
        color 0.2s,
        background 0.2s,
        transform 0.2s;
}

.theme-icon {
    font-size: 15px;
    line-height: 1;
}

.theme-toggle:hover {
    border-color: var(--accent);
    background: var(--bg-hover);
    transform: translateY(-2px);
}

.theme-toggle:active {
    transform: translateY(0);
}

.publish {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    line-height: 1.2;
    white-space: nowrap;
    margin-bottom: 8px;
    padding: 13px;
    border: 0;
    border-radius: 10px;
    color: #fff;
    font-size: 15px;
    font-weight: 700;
    text-align: center;
    text-decoration: none;
    background: #ff3858;
    transition: 0.2s;
}

.publish-icon {
    display: inline-flex;
    align-items: center;
}

a.publish:hover {
    transform: translateY(-3px);
}

a.publish:active {
    transform: translateY(-1px);
}
</style>