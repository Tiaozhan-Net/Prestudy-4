// 用户表
export interface User {
    id: number
    name: string
    avatar: string
    createTime: string
    description?: string
    about?: string
    location?: string
    interests?: string[]
    following?: number[]
    likeArticles?: number[]
    likeComments?: number[]
    password?: string
}

// 文章表
export interface Article {
    id: number
    authorId: number
    cover: string
    title: string
    images: string[]
    content: string
    tags: string[]
    category: '音乐' | '美食' | '萌宠' | '旅行' | '游戏' | '影视' | '读书' | '摄影' | '编程' | '其他'
    likes: number
    comments?: Comment[]
    createTime: string
}

// 单条评论
export interface Comment {
    id: number
    articleId: number
    userId: number
    content: string
    likes: number
    createTime: string
}