<script setup lang="ts">
const value = defineModel<string>({ default: 'all' })
const options = [
    { value: 'all', label: '全部' },
    { value: '音乐', label: '音乐' },
    { value: '美食', label: '美食' },
    { value: '萌宠', label: '萌宠' },
    { value: '旅行', label: '旅行' },
    { value: '游戏', label: '游戏' },
    { value: '影视', label: '影视' },
    { value: '读书', label: '读书' },
    { value: '摄影', label: '摄影' },
    { value: '编程', label: '编程' },
    { value: '其他', label: '其他' }
]
</script>

<template>
    <el-select class="custom-select" v-model="value" style="width: 110px" :teleported="false">
        <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
    </el-select>
</template>

<style scoped>
.custom-select :deep(.el-select__wrapper) {
    min-width: 110px;
    height: 42px;
    padding: 0 17px;
    border: 1px solid #eee;
    border-radius: 21px;
    color: var(--text-soft);
    font: inherit;
    background: var(--bg-surface);
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
</style>