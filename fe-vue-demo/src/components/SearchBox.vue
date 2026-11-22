<script setup lang="ts">
import { Search } from '@element-plus/icons-vue'

const input = defineModel<string>({ default: '' })
const handleClear = () => {
    input.value = ''
}
</script>

<template>
    <el-input v-model="input" class="input-box" placeholder="搜索感兴趣的内容" :prefix-icon="Search" clearable
        @clear="handleClear" />
</template>

<style scoped>
.input-box {
    width: min(540px, 53vw);
}

.input-box :deep(.el-input__wrapper) {
    height: 44px;
    display: flex;
    align-items: center;
    padding: 0 15px;
    border-radius: 22px;
    color: var(--text-soft);
    background: var(--bg-surface);
    --el-input-focus-border-color: #e0313f;
}

.input-box :deep(.el-input__prefix) {
    display: flex !important;
    align-items: center !important;
    justify-content: center !important;
    width: 26px;
    height: 26px;
    padding: 0 !important;
    border-radius: 50%;
    background: transparent;
    transition: background 0.2s;
    cursor: pointer;
    margin-right: 6px;
}

.input-box :deep(.el-input__prefix-inner) {
    display: flex !important;
    align-items: center !important;
    justify-content: center !important;
    width: 100%;
    height: 100%;
    margin: 0 !important;
    padding: 0 !important;
    line-height: 1 !important;
}

.input-box :deep(.el-icon) {
    font-size: 15px !important;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 !important;
    vertical-align: middle;
}

.input-box :deep(.el-input__wrapper.is-focus .el-input__prefix) {
    background: #0000000f;
}
</style>