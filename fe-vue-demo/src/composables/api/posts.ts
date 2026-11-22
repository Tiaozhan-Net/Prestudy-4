import { users } from '../data/users'

import { articles } from '../data/articles'

export const Posts = () => {
    const userMap = new Map(users.value.map(u => [u.id, u]))
    return articles.value.map(article => ({
        ...article,
        author: userMap.get(article.authorId)?.name || '未知作者',
        avatar: userMap.get(article.authorId)?.avatar || ''
    }))
}