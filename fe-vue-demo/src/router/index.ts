import { createRouter, createWebHistory } from 'vue-router'
import { ElMessage } from 'element-plus'

import { currentUser } from '../composables/auth'
import { articles } from '../composables/data/articles'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Home',
      meta: { requiresAuth: true },
      component: () => import('../views/HomePage.vue'),
    },
    {
      path: '/message',
      name: 'Message',
      meta: { requiresAuth: true },
      component: () => import('../views/MessagePage.vue'),
    },
    {
      path: '/profile',
      name: 'Profile',
      meta: { requiresAuth: true },
      component: () => import('../views/ProfilePage.vue'),
    },
    {
      path: '/editing',
      name: 'Editing',
      meta: { requiresAuth: true },
      component: () => import('../views/EditingPage.vue'),
    },
    {
      path: '/article/:id',
      name: 'Article',
      meta: { requiresAuth: true },
      component: () => import('../views/DetailsPage.vue'),
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/LoginPage.vue'),
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('../views/RegisterPage.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('../views/404Page.vue'),
    }
  ],
})

router.beforeEach((to, from) => {
  const isAuthPage = to.path === '/login' || to.path === '/register'
  const cameFromAuthPage = from.path === '/login' || from.path === '/register'

  if (isAuthPage && currentUser.value) {
    return '/'
  }

  if (to.meta.requiresAuth && !currentUser.value) {
    if (cameFromAuthPage) {
      ElMessage.warning('请先登录')
    }
    return {
      path: '/login',
      query: to.fullPath === '/' ? undefined : { redirect: to.fullPath },
    }
  }
})

router.afterEach((to) => {
  const article = to.name === 'Article'
    ? articles.value.find(item => item.id === Number(to.params.id))
    : undefined
  const pageTitles: Record<string, string> = {
    Home: '首页',
    Message: '消息',
    Profile: currentUser.value?.name || '个人中心',
    Editing: '发布笔记',
    Article: article?.title || '文章',
    Login: '登录',
    Register: '注册'
  }

  document.title = `${pageTitles[String(to.name)] || ''} | Tenzors`
})

export default router
