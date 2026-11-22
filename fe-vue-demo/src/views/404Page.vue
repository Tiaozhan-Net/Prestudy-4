<script setup lang="ts">

</script>

<template>
    <section class="messages-page">
        <header class="messages-heading">
            <div>
                <p>ERROR</p>
                <h1>错误</h1>
            </div>
        </header>
        <div class="messages-empty">
            <span class="messages-empty-icon"><el-icon>
                    <Warning />
                </el-icon></span>
            <h2>404 Not Found</h2>
            <p>未找到页面</p>
        </div>
    </section>
</template>

<style scoped>
.messages-page {
    width: 100%;
}

.messages-heading {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
}

.messages-heading p {
    margin: 3px 0 6px;
    color: #c2b8b8;
    font-size: 10px;
    font-weight: 700;
    letter-spacing: 0.12em;
}

.messages-heading h1 {
    margin: 0;
    color: var(--text-main);
    font-size: 24px;
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