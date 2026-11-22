import { ref } from 'vue'

import type { Article } from '../type'
import { getStoredValue, removeStoredValue, setStoredValue } from './storage'

export const ARTICLES_STORAGE_KEY = 'articles'
const USER_ARTICLES_STORAGE_KEY = 'user-articles'

export const articles = ref<Article[]>([])

export const loadArticles = async () => {
    const storedArticles = await getStoredValue<Article[]>(USER_ARTICLES_STORAGE_KEY)
    articles.value = storedArticles || []

    await removeStoredValue(ARTICLES_STORAGE_KEY)
}

export const saveArticles = async () => {
    await setStoredValue(USER_ARTICLES_STORAGE_KEY, articles.value)
}

export const addArticle = async (article: Omit<Article, 'id' | 'createTime'>) => {
    const nextId = articles.value.reduce((maxId, item) => Math.max(maxId, item.id), 0) + 1
    const newArticle: Article = {
        ...article,
        id: nextId,
        createTime: new Date().toISOString()
    }

    articles.value.unshift(newArticle)
    await saveArticles()
    return newArticle
}

export const deleteArticle = async (articleId: number) => {
    articles.value = articles.value.filter(article => article.id !== articleId)
    await saveArticles()
}