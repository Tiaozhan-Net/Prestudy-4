<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { registerUser } from '../composables/auth'

const router = useRouter()
const name = ref('')
const password = ref('')

const handleRegister = async () => {
    const trimmedName = name.value.trim()
    if (!trimmedName || !password.value) {
        ElMessage.warning('请输入用户名和密码')
        return
    }

    try {
        await registerUser(trimmedName, password.value)
        ElMessage.success('注册成功')
        router.push('/login')
    } catch (error) {
        ElMessage.error(error instanceof Error ? error.message : '注册失败')
    }
}
</script>

<template>
    <main class="auth-wrap">
        <div class="auth-card">
            <div class="brand">
                <span class="brand-inner">
                    <img src="/book.svg" />
                </span>
                Tenzors
            </div>
            <p class="auth-sub">注册后查看所有笔记</p>
            <form class="form-grid" @submit.prevent="handleRegister">
                <label class="form-field">
                    <span>用户名</span>
                    <input v-model="name" class="form-input" type="text" placeholder="请输入用户名" />
                </label>
                <label class="form-field">
                    <span>密码</span>
                    <input v-model="password" class="form-input" type="password" placeholder="请输入密码" />
                </label>
                <button class="form-submit" type="submit">注册</button>
            </form>
            <p class="link">已有账号？<RouterLink to="/login">去登录</RouterLink></p>
        </div>
    </main>
</template>

<style scoped>
.auth-wrap {
    min-height: 100vh;
    display: grid;
    place-items: center;
    padding: 24px;
}

.auth-card {
    width: min(100%, 380px);
    padding: 30px 28px;
    border: 1px solid #f0eeee;
    border-radius: 20px;
    background: var(--bg-surface);
    box-shadow: 0 10px 28px #422d300c;
}

.brand {
    display: flex;
    align-items: center;
    justify-content: flex-start;
    color: var(--text-main);
    font-size: 23px;
    font-weight: 700;
    text-decoration: none;
    margin-bottom: 6px;
}

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

.auth-sub {
    margin: 0 0 20px;
    color: var(--text-mute);
    font-size: 12px;
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

.form-input {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #f0e7e7;
    border-radius: 10px;
    color: var(--text-main);
    font: inherit;
    font-size: 13px;
    background: var(--bg-soft);
    transition: border-color 0.2s;
}

.form-input:focus {
    border-color: var(--accent-border-2);
    outline: 1px solid #fff1f3;
}

.form-input.is-invalid {
    border-color: #e0313f;
    background: var(--danger-soft);
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

.form-error {
    margin: -4px 0 0;
    color: var(--danger);
    font-size: 12px;
}

.form-submit:hover {
    filter: brightness(1.07);
}

.link {
    margin-top: 12px;
    color: var(--text-mute);
    font-size: 12px;
}

.link a {
    color: var(--accent);
    text-decoration: none;
}
</style>