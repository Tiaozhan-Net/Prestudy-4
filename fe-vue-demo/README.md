# Tenzors 课件样例

保留最初版本 `fe-vue-demo` 的完整站点外观：

- 首页：搜索、分类、推荐笔记
- 消息页
- 个人中心
- 发布笔记
- 文章详情
- 侧边栏、主题切换

只把**文章数据**清空了，所以登录后像一个“还没有人上传过内容”的 Tenzors。

## 已接通

- `POST /api/auth/register`：注册
- `POST /api/auth/login`：登录，后端下发 `tz-sessions` Cookie
- `GET /api/auth/profile`：恢复登录用户
- `DELETE /api/auth/logout`：退出

## 启动

```bash
pnpm install
pnpm dev
```

前端：`http://127.0.0.1:5173`  
后端：`http://127.0.0.1:8088`

先启动后端，再注册/登录。

## 当前数据状态

- `src/composables/data/articles.ts` 里的 `articles` 初始为空
- `main.ts` 不再从 IndexedDB 读取旧文章
- 首页轮播图自动隐藏，显示空状态
- 个人中心、消息页、发布页保持原版界面
- 文章详情页没有文章时显示“文章不存在”
- 文章列表、详情、发布 API 留给课上接入

## 和原版的差异

| 位置 | 改动 |
| --- | --- |
| `src/composables/auth.ts` | 改为调用后端登录/注册/资料接口，再同步到 `currentUser` |
| `src/stores/auth.ts` | 保存真实登录态 |
| `src/main.ts` | 安装 Pinia，启动时恢复登录态；不再加载旧文章数据 |
| `src/views/HomePage.vue` | 空文章时显示“还没有人发布笔记” |
| `src/views/ProfilePage.vue` | 使用后端登录用户资料 |
| `src/views/LoginPage.vue` / `RegisterPage.vue` | 调用真实接口 |

原版的 `Sidebar`、`PostCard`、`CarouselCard`、`SearchBox`、`CategorySelect`、消息页、发布页、详情页都保留。
