<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const visible = computed(() => route.path !== '/')

const goHome = () => {
    router.push('/')
}
</script>

<template>
    <button v-if="visible" class="back-home" type="button" aria-label="返回首页" title="返回首页" @click="goHome">
        <el-icon>
            <HomeFilled />
        </el-icon>
    </button>
</template>

<style scoped>
.back-home {
    position: fixed;
    right: 24px;
    bottom: 24px;
    z-index: 20;
    display: inline-flex;
    width: 44px;
    height: 44px;
    align-items: center;
    justify-content: center;
    border: 0;
    border-radius: 50%;
    color: #fff;
    background: var(--accent);
    box-shadow: 0 8px 18px #412b2b26;
    cursor: pointer;
    transition: transform 0.2s, background 0.2s;
}

.back-home:hover {
    background: var(--accent-strong);
    transform: translateY(-3px);
}

.back-home:active {
    transform: translateY(-1px);
}
</style>
