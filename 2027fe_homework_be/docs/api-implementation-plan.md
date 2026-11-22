# Tenzors Blog API 后端编写方案（controller + service 范围）

> 依据：`Tenzors Blog API.openapi.json`（v2.0.0，26 个操作）
> **改动面约束**：主体只改 `controller/` 与 `service/`；`router/` 允许补 handler 参数与上传图片的静态目录；`model/` 允许**小幅调整**（加索引、改字段长度、加一张表）；`common/form.go` 的 `PagerForm`、`config/logconf.go`（颜色只写控制台）、`middleware/log.go`（二进制 body 不落日志）做过小幅修改。
> 不改：`main.go`、`config/`、`common/` 其它文件、`middleware/`、`logger/`、`pkg/`、`sql/`。

---

## 1. 范围与目标

| 项 | 内容 |
| --- | --- |
| 接口数 | 26（认证 4 / 文章 8 / 评论 5 / 用户 6 / 消息 4 / 上传 1） |
| 认证 | Cookie + Session（gin-contrib cookie store，`tz-sessions`），登出/改密清当前 cookie |
| 响应 | 统一 `{ code, message, data }`，`code=0` 成功 |
| 错误 | 15 个业务码（1001~9001）+ HTTP 语义状态码 |
| 允许改的文件 | `controller/*.go`、`service/*.go`、`router/router.go`（仅 handler 绑定） |
| 允许小改的文件 | `model/*`（索引/长度/新增 `session.go`）、`common/form.go`（仅 `PagerForm`） |
| 禁止改的文件 | 其余全部；尤其 `common/response.go`、`middleware/*`、`config/*` 一行不动 |

**为什么能在不改 middleware 的情况下统一响应**：`middleware/error.go` 只在 `c.Errors` 非空时才接管输出。controller 改为**自己写 JSON、不调用 `c.Error`**，错误响应就完全由 controller 决定，middleware 形同直通。

---

## 2. 现状评估

### 2.1 可直接复用

| 设施 | 位置 | 说明 |
| --- | --- | --- |
| 分层骨架 | `controller/controller.go`、`service/service.go` | 结构体嵌入注册，`ctr.Auth.X` / `srv.AuthSrv.X` 调用 |
| 会话工具 | `controller/session.go` | `SessionGet/Set/Update/Clear`，本方案只换底层实现，调用点不变 |
| 分页结构 | `common/form.go:7` `PagerForm` | 复用；3.5 会给它补 json tag / 放宽 binding / 上限 100 |
| 分页 scope | `model/scopes.go:9` `model.Paginate` | 直接 `Scopes(model.Paginate(p))` |
| 函数式映射 | `pkg/func_prog.go:49` `pkg.Map` | model 切片 → DTO 切片；空输入返回非 nil `[]T` |
| 校验翻译 | `service/validator`（zh） | 已按 `APP_LANGUAGE=zh` 注册，controller 可直接取 `vl.Trans` |
| 模型 | `model/{user,article,comment,like_record,tag,image_url,conversation,message}.go` | 表结构够用；索引/长度小幅调整见 3.9 |
| 密码 | `service/auth.go`（argon2id） | 逻辑保留，错误返回/参数绑定重写 |

### 2.2 本范围内必须解决的 8 个问题

| # | 问题 | 位置 | 处理 |
| --- | --- | --- | --- |
| 1 | **编译失败**：`Service` 嵌入 4 个 `Hello` 方法，`srv.Hello` 歧义 | `service/{comment,conversation,user,image}.go` + 对应 controller 桩 | 删掉 Hello，写真实方法 |
| 2 | 响应体 `{success,...}` 不符 spec | `controller/init.go`（`Response` 仍被 middleware 引用，不能删） | 新增 `ok()/fail()` helper，旧 `Response` 保留但不再使用 |
| 3 | 错误码只有 3~7 分类 | 无 | `service/errs.go` 定义 `AppError` + 15 个业务码 |
| 4 | 登录/注册用 `ShouldBindQuery` | `controller/auth.go:16,50` | 改 `ShouldBindJSON` |
| 5 | 不传 `page` 直接 400（`binding:"min=1"`）、size 上限 20、响应输出 `Page`/`Size` 大写 | `common/form.go:7` | 小幅修改 `PagerForm`（3.5） |
| 6 | 文章列表 SQL 表名/列名错（`users`/`articles`/`likes_count`） | `service/article.go` | 改成 `user`/`article`/`like_count` |
| 7 | 文章 DTO 命名混乱（`Author UserBrief` 与 spec 的字符串 `author` 冲突）、缺 `authorId/avatar/cover`、`likes` 带 omitempty | `service/struct.go` | 按 3.4 的命名规则重写 |
| 8 | 会话是 `tz-sessions` 签名 cookie，无法服务端失效 | `config/default.go`（**不能改**） | 保留 `controller/session.go` 初版，不改 |

### 2.3 本范围解决不了、需要上层配合（详见第 8 章出界项）

| 出界项 | 影响 |
| --- | --- |
| ~~E1 模型缺索引/长度~~ | **已解除**：改为 `model/` 小幅调整（见 3.9），不再需要迁移 hack |
| ~~E2 `PagerForm` 无 json tag~~ | **已解除**：`common/form.go` 允许小幅修改（3.5），直接绑定即可 |
| E3 路由带 `/api` 前缀、`/articles` 实际是 `/api/articles/`（尾斜杠）、`/image` 实际是 `/api/image/` | 需前端 baseURL 配合；本方案按实际路径实现 |
| ~~E4 没有静态目录路由~~ | **已解除**：图片存 `Config.AppUploadDir`（默认 `./update`），`r.Static(service.ImageURLPath, config.Config.AppUploadDir)`，上传返回 `/update/xxx.png` |
| E5 CORS 未开 `AllowCredentials`、`AllowOrigins=*` | 跨域带 cookie 的联调会失败，需改 `config/default.go` |
| E6 `middleware/error.go` 恒 200 + 分类码 | 已被"controller 自己写响应"绕开；panic 兜底仍是旧格式 |

---

## 3. 设计

### 3.1 统一响应与错误码（controller 内闭环）

`service/errs.go` 新增：

```go
type AppError struct {
	Code int    // 业务码，进 body
	HTTP int    // HTTP 状态码
	Msg  string // 中文文案
	Err  error  // 内层原因，只进日志
}

func (e *AppError) Error() string { ... }
func (e *AppError) Unwrap() error { return e.Err }
func (e *AppError) With(err error) *AppError // 返回副本，不改全局变量

var (
	ErrParam             = &AppError{1002, 400, "参数不合法", nil}
	ErrUserNameExists    = &AppError{1001, 409, "用户名已存在", nil}
	ErrCredentials       = &AppError{2001, 401, "用户名或密码错误", nil}
	ErrTooManyAttempts   = &AppError{2002, 429, "尝试次数过多，请稍后再试", nil}
	ErrUnauthenticated   = &AppError{2003, 401, "未登录或登录已过期", nil}
	ErrArticleNotFound   = &AppError{3001, 404, "文章不存在", nil}
	ErrArticleForbidden  = &AppError{3002, 403, "无权限删除该文章", nil}
	ErrCommentForbidden  = &AppError{4001, 403, "无权限删除该评论", nil}
	ErrCommentNotFound   = &AppError{4002, 404, "评论不存在", nil}
	ErrUserNotFound      = &AppError{5001, 404, "用户不存在", nil}
	ErrFollowSelf        = &AppError{6001, 400, "不能关注自己", nil}
	ErrConversationNotFound  = &AppError{7001, 404, "会话不存在", nil}
	ErrConversationForbidden = &AppError{7002, 403, "无权访问该会话", nil}
	ErrChatSelf          = &AppError{7003, 400, "不能和自己聊天", nil}
	ErrInternal          = &AppError{9001, 500, "服务器内部错误", nil}
)
```

`controller/init.go` 新增（不动旧 `Response`，middleware 仍能编译）：

```go
func ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": successMessage, "data": data})
}

func fail(c *gin.Context, err error) {
	var ae *service.AppError
	if errors.As(err, &ae) {
		if ae.Err != nil {
			logger.Errorf("controller: %v", ae.Err)
		}
		c.JSON(ae.HTTP, gin.H{"code": ae.Code, "message": ae.Msg, "data": nil})
		return
	}
	logger.Errorf("controller: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"code": 9001, "message": "服务器内部错误", "data": nil})
}

// 参数校验失败：validator 中文翻译 + 兜底文案
func paramErr(c *gin.Context, err error) {
	if ve, ok := err.(validator.ValidationErrors); ok {
		var msgs []string
		for _, v := range ve.Translate(vl.Trans) {
			msgs = append(msgs, v)
		}
		fail(c, service.ErrParam.WithMsg(strings.Join(msgs, "；")))
		return
	}
	fail(c, service.ErrParam.With(err))
}
```

**铁律**：controller 里**不再出现 `c.Error(...)`**（一旦调用，`middleware/error.go` 就会用旧格式覆盖响应）。

### 3.2 会话：保留 `controller/session.go` 初版（gin-contrib cookie store）

最终决定不做服务端会话存储，`controller/session.go` **保持初版、不再修改**：

- 存储：`config/default.go` 的 `cookie.NewStore(APP_SECRET)` + `sessions.Sessions("tz-sessions", store)`，数据（gob 编码的 `UserSession`）在**签名 cookie** 里；
- cookie：`tz-sessions`，`Max-Age=1800`（30 分钟），`SameSite=Lax`，`HttpOnly` 仅生产模式（均由 `config/default.go` 决定，未改）；
- 读写：`SessionGet/Set/Update/Clear/Delete` 是 `sessions.Default(c)` 的薄封装，调用点（`middleware/role.go`、各 controller）不变；
- 失效：登出、修改密码时 `SessionClear` 清空当前浏览器的 cookie；
- **代价（明确记录）**：spec 里的 `sid` / `HttpOnly` / 服务端吊销**不满足**；旧 cookie 的副本在过期前仍有效（无法踢下线）；因此已删除 `service/session.go`、`model/session.go`；
- 若库里已建过 `session` 表，属于遗留，可手工 `DROP TABLE session;`。

### 3.3 数据访问：表名与列名对齐

model 只做 3.9 的索引/长度调整，表名列名保持现状，所以**新写的查询必须按现有命名**：

- 表名用单数：`user`、`article`、`comment`、`tag`、`image_url`、`like_record`、`conversation`、`message`；`user` 是 MySQL 关键字，**裸 SQL 里必须反引号**。
- 列名：`like_count`（不是 `likes_count`）、`author_id`、`article_id`、`user_low_id`/`user_high_id`。
- 优先用关联 `Preload("Author")` 等，少写 join；关键字搜索需要作者名时：

```go
model.DB.Model(&model.Article{}).
	Joins("JOIN `user` AS u ON u.id = article.author_id AND u.deleted_at IS NULL").
	Where("article.title LIKE ? OR article.content LIKE ? OR article.category LIKE ? OR u.name LIKE ?", like, like, like, like)
```

- 软删除由 `deleted_at` 自动过滤（GORM），删除文章需在事务里手工清理 `comment`/`like_record`/`image_url`/`article_tags`，避免残留。

### 3.4 DTO 命名与时间格式

`service/struct.go` 重写。**不做"万能 DTO"**：每个接口一套 `XxxParams` / `XxxResponse`，公共字段抽成 `XxxInfo` 基座内嵌；`Post` 这个名字不再使用。

| 类别 | 规则 | 例子 |
| --- | --- | --- |
| 入参 | `XxxParams`（body）/ `XxxQuery`（query） | `ArticleCreateParams`、`MessageListQuery` |
| 出参 | `XxxResponse`，只放**该接口实际返回**的字段 | `ArticleDetailResponse`、`ArticleLikeResponse`、`FollowResponse` |
| 公共字段 | `XxxInfo` 作为内嵌基座，不单独作为接口返回值 | `ArticleInfo`、`CommentInfo`、`UserInfo`、`ConversationInfo`、`MessageInfo` |
| 列表接口 | 每个列表**独立**结构：`List` + `Total` + 内嵌 `common.PagerForm` | `ArticleListResponse`、`UserArticlesResponse`、`ConversationListResponse` |
| 不用 `Post` | spec 的 `Post` 只是 schema 名，代码里拆成 `ArticleInfo` + 各种 `ArticleXxxResponse` | — |
| 不用泛型 | 删掉 `PageResult[T]`，每个列表结构自己写字段（以后能各自演进） | — |
| 计数字段 | `XxxCount`，不用复数名词 | `LikeCount`（JSON 仍是 `likes`） |
| 布尔状态 | 形容词或 `Is` 前缀 | `Liked`、`Following`、`IsFollowing` |
| 时间字段 | 统一 `CreateTime`/`UpdateTime`（string）；不用 `CreatedAt` | — |

```go
// ============ 认证 ============
type AuthRegisterParams struct {
	Name     string `json:"name" binding:"required,min=1,max=20"`
	Password string `json:"password" binding:"required,min=1,max=20"`
}

type AuthLoginParams struct {
	Name     string `json:"name" binding:"required,min=1,max=20"`
	Password string `json:"password" binding:"required,min=1,max=20"`
}

type AuthProfileResponse struct { // GET /auth/profile
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Slogan string `json:"slogan"`
}

// 登录只用于写会话，不进响应体
type AuthLoginResult struct {
	ID   int64
	Name string
}

// ============ 文章 ============
// ArticleInfo 文章公共字段：列表项 / 详情 / 发布返回共用
type ArticleInfo struct {
	ID           int64    `json:"id"`
	AuthorID     int64    `json:"authorId"`
	AuthorName   string   `json:"author"`   // spec 的 author 是名字，不是对象
	AuthorAvatar string   `json:"avatar"`
	Cover        string   `json:"cover"`    // = Images[0]
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Category     Category `json:"category"`
	Images       []string `json:"images"`
	Tags         []string `json:"tags"`
	LikeCount    int64    `json:"likes"`
	CreateTime   string   `json:"createTime"`
}

type ArticleListParams struct {
	common.PagerForm
	Category Category `form:"category" binding:"omitempty,oneof=音乐 美食 萌宠 旅行 游戏 影视 读书 摄影 编程 其他"`
	Keyword  string   `form:"keyword" binding:"max=40"`
	Sort     Sort     `form:"sort" binding:"omitempty,oneof=date-desc date-asc likes-desc"`
}

type ArticleListResponse struct { // GET /articles/list
	common.PagerForm
	List  []ArticleInfo `json:"list"`
	Total int64         `json:"total"`
}

type ArticleDetailResponse struct { // GET /articles/{id}
	ArticleInfo
}

type ArticleCreateParams struct {
	Title    string   `json:"title" binding:"required,max=40"`
	Content  string   `json:"content" binding:"required,max=500"`
	Category Category `json:"category" binding:"required,oneof=音乐 美食 萌宠 旅行 游戏 影视 读书 摄影 编程 其他"`
	Images   []string `json:"images" binding:"max=9"`
	Tags     []string `json:"tags" binding:"max=10,dive,max=40"`
}

type ArticleCreateResponse struct { // POST /articles
	ArticleInfo
}

type ArticleLikeResponse struct { // POST/DELETE /articles/{id}/like
	Liked     bool  `json:"liked"`
	LikeCount int64 `json:"likes"`
}

// ============ 评论 ============
type CommentInfo struct {
	ID         int64     `json:"id"`
	ArticleID  int64     `json:"articleId"`
	UserID     int64     `json:"userId"`
	User       UserBrief `json:"user"`
	Content    string    `json:"content"`
	LikeCount  int64     `json:"likes"`
	CreateTime string    `json:"createTime"`
}

type ArticleCommentListParams struct {
	Sort Sort `form:"sort" binding:"omitempty,oneof=likes-desc time-desc"`
}

type CommentListResponse struct { // GET /articles/{id}/comments
	List  []CommentInfo `json:"list"`
	Total int64         `json:"total"`
}

type CommentCreateParams struct {
	Content string `json:"content" binding:"required,max=80"`
}

type CommentCreateResponse struct { // POST /articles/{id}/comments
	CommentInfo
}

type CommentLikeResponse struct { // POST/DELETE /comments/{id}/like
	Liked     bool  `json:"liked"`
	LikeCount int64 `json:"likes"`
}

// ============ 用户 ============
type UserBrief struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// UserInfo 本人资料公共字段：/users/me 与 PUT /users/me 共用
type UserInfo struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	Avatar       string   `json:"avatar"`
	CreateTime   string   `json:"createTime"`
	Slogan       string   `json:"slogan"`
	Description  string   `json:"description"`
	Location     string   `json:"location"`
	Interests    []string `json:"interests"`
	Following    []int64  `json:"following"`
	LikeArticles []int64  `json:"likeArticles"`
	LikeComments []int64  `json:"likeComments"`
}

type UserMeResponse struct { // GET /users/me
	UserInfo
}

type UserUpdateResponse struct { // PUT /users/me
	UserInfo
}

type UserUpdateParams struct { // 全指针：区分"没传"与"传空"
	Name        *string  `json:"name" binding:"omitempty,min=1,max=20"`
	Avatar      *string  `json:"avatar"`
	Slogan      *string  `json:"slogan" binding:"omitempty,max=32"`
	Description *string  `json:"description" binding:"omitempty,max=120"`
	Location    *string  `json:"location" binding:"omitempty,max=40"`
	Interests   []string `json:"interests"`
	Password    *string  `json:"password" binding:"omitempty,min=1,max=20"`
}

type UserPublicResponse struct { // GET /users/{id}
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Avatar      string   `json:"avatar"`
	CreateTime  string   `json:"createTime"`
	Slogan      string   `json:"slogan"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	Interests   []string `json:"interests"`
	IsFollowing bool     `json:"isFollowing"`
}

type UserArticlesParams struct {
	common.PagerForm
	Sort Sort `form:"sort" binding:"omitempty,oneof=date-desc date-asc likes-desc"`
}

type UserArticlesResponse struct { // GET /users/{id}/articles
	common.PagerForm
	List  []ArticleInfo `json:"list"`
	Total int64         `json:"total"`
}

type FollowResponse struct { // POST/DELETE /users/{id}/follow
	Following bool `json:"following"`
}

// ============ 消息 ============
type MessageInfo struct {
	ID             int64  `json:"id"`
	ConversationID int64  `json:"conversationId"`
	SenderID       int64  `json:"senderId"`
	ReceiverID     int64  `json:"receiverId"`
	Content        string `json:"content"`
	CreateTime     string `json:"createTime"`
}

type ConversationInfo struct {
	ID          int64        `json:"id"`
	Peer        UserBrief    `json:"peer"`
	LastMessage *MessageInfo `json:"lastMessage,omitempty"`
	CreateTime  string       `json:"createTime"`
	UpdateTime  string       `json:"updateTime"`
}

type ConversationListResponse struct { // GET /conversations
	common.PagerForm
	List  []ConversationInfo `json:"list"`
	Total int64              `json:"total"`
}

type ConversationCreateParams struct {
	UserID int64 `json:"userId" binding:"required,min=1"`
}

type ConversationCreateResponse struct { // POST /conversations
	ConversationInfo
}

type MessageListQuery struct {
	AfterID  int64 `form:"afterId" binding:"omitempty,min=0"`
	BeforeID int64 `form:"beforeId" binding:"omitempty,min=1"`
	Limit    int   `form:"limit" binding:"omitempty,min=1,max=100"`
}

type MessageListResponse struct { // GET /conversations/{id}/messages
	List    []MessageInfo `json:"list"`
	HasMore bool          `json:"hasMore"`
	Peer    UserBrief     `json:"peer"`
}

type MessageSendParams struct {
	Content string `json:"content" binding:"required,max=1000"`
}

type MessageSendResponse struct { // POST /conversations/{id}/messages
	MessageInfo
}

// ============ 上传 ============
type ImageUploadResponse struct { // POST /image
	URL string `json:"url"`
}
```

- 时间统一 `2006-01-02T15:04:05.000Z`：`service.FormatTime(t)`，DTO 用 `string`，避开 `time.Time` 的 RFC3339Nano；
- 一律去掉 `omitempty`（仅 `ConversationInfo.LastMessage` 保留）：`images`/`tags`/`likes` 空时输出 `[]`/`0`，前端不用判空；
- **所有列表统一 `{list,total}`**：评论列表也用 `CommentListResponse`，与文章列表同构；评论不分页（spec 无 `page`/`size`），`total` 为全部条数、服务端上限 500，前端 `DetailsPage` 需从 `res.data` 改成 `res.data.list`。

**旧字段 → 新字段对照**（重写 `service/struct.go` 时逐条替换）：

| 现有 | 改后 | 原因 |
| --- | --- | --- |
| `Article`（DTO） | `ArticleInfo` + `ArticleDetailResponse`/`ArticleCreateResponse` | 不搞万能 DTO；基座与接口响应分开 |
| `Article.Author UserBrief` | `ArticleInfo.AuthorName` + `AuthorID` + `AuthorAvatar` | spec 的 `author` 是字符串 |
| `Article.Likes`（带 omitempty） | `ArticleInfo.LikeCount`（json 仍为 `likes`） | 计数用 Count；去掉 omitempty |
| `Article.Comments []Comment` | 删除 | 详情不返回评论 |
| `Comment`（DTO） | `CommentInfo` + `CommentCreateResponse` | 同上 |
| `Comment.Likes *int64` | `CommentInfo.LikeCount int64` | 去指针，0 也输出 |
| `Comment.User *UserBrief` | `CommentInfo.User UserBrief` | spec 里 user 一定存在 |
| `User`（DTO） | `UserInfo` + `UserMeResponse`/`UserUpdateResponse` | 按接口拆 |
| `User.Description *string` / `Slogan *string` | `string` | 去指针，空串照常输出 |
| `UserPublic` | `UserPublicResponse` | 出参统一 Response 后缀 |
| `Conversation` | `ConversationInfo` + `ConversationCreateResponse` | 同上 |
| `ConversationListResponse`（模板遗留） | 重写为 `{PagerForm, List, Total}` | 复用名字，字段按 spec |
| `LikeResult` | `ArticleLikeResponse` / `CommentLikeResponse` | 两个接口各自独立 |
| `FollowResult` | `FollowResponse` | 统一后缀 |
| `Message`（DTO） | `MessageInfo` + `MessageSendResponse` | 同上 |
| `MessageList` | `MessageListResponse` | 统一后缀 |
| `PageResult[T]` / `PostList` | 删除 | 每个列表接口用自己的 Response 结构 |

### 3.5 分页与校验（`PagerForm` 小幅修改）

`common/form.go` 只动 `PagerForm` 一个结构体：补 `json` tag、放宽 `binding`、上限改 100、加两个方法；字段名保持 `Page`/`Size`（现有调用点不用改名）。

```go
// common/form.go（改）
type PagerForm struct {
	Page int `form:"page" json:"page" binding:"omitempty,min=1"`
	Size int `form:"size" json:"size" binding:"omitempty,min=1,max=100"`
}

// Normalize 兜底：page=1、size=10、上限 100
func (p PagerForm) Normalize() PagerForm {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = 10
	}
	if p.Size > 100 {
		p.Size = 100
	}
	return p
}

// Offset 与 model.Paginate 口径一致，需要手工分页时用
func (p PagerForm) Offset() int { return (p.Page - 1) * p.Size }
```

改完后入参可以直接 `ShouldBindQuery`（不传分页不再 400），响应里 `page`/`size` 也自然输出小写（E2 解除）：

```go
type ArticleListParams struct {
	common.PagerForm // 直接绑定 page/size，由 Normalize 兜底
	Category Category `form:"category" binding:"omitempty,oneof=音乐 美食 萌宠 旅行 游戏 影视 读书 摄影 编程 其他"`
	Keyword  string   `form:"keyword" binding:"max=40"`
	Sort     Sort     `form:"sort" binding:"omitempty,oneof=date-desc date-asc likes-desc"`
}
```

controller 统一在绑定后 `params.PagerForm = params.PagerForm.Normalize()` 再传 service；`page=0`/`size=0` 视为未传（宽容策略）。

其余长度限制直接落 `binding`（与 spec 一致）：标题 40、正文 500、评论 80、消息 1000、slogan 32、description 120、location 40、用户名/密码 20。

### 3.6 路由：只补 handler，不改路径

`router/router.go` 的改动是给空路由补第二个参数，路径字符串保持原样；`router/init.go` 额外挂一个上传图片的静态目录：

```go
userRouter.GET("/me", ctr.User.Me)
userRouter.PUT("/me", ctr.User.UpdateMe)
userRouter.GET("/:id", ctr.User.Public)
userRouter.POST("/:id/follow", ctr.User.Follow)
userRouter.DELETE("/:id/follow", ctr.User.Unfollow)
userRouter.GET("/:id/articles", ctr.User.Articles)

articleRouter.GET("/list", ctr.Article.ArticlesList)
articleRouter.POST("/", ctr.Article.ArticleCreate)
articleRouter.GET("/:id", ctr.Article.ArticleDetail)
articleRouter.DELETE("/:id", ctr.Article.ArticleDelete)
articleRouter.POST("/:id/like", ctr.Article.ArticleLike)
articleRouter.DELETE("/:id/like", ctr.Article.ArticleUnlike)
articleRouter.GET("/:id/comments", ctr.Comment.CommentList)
articleRouter.POST("/:id/comments", ctr.Comment.CommentCreate)

commentRouter.DELETE("/:id", ctr.Comment.CommentDelete)
commentRouter.POST("/:id/like", ctr.Comment.CommentLike)
commentRouter.DELETE("/:id/like", ctr.Comment.CommentUnlike)

conversationRouter.GET("/", ctr.Conversation.Conversations)
conversationRouter.POST("/", ctr.Conversation.ConversationCreate)
conversationRouter.GET("/:id/messages", ctr.Conversation.Messages)
conversationRouter.POST("/:id/messages", ctr.Conversation.SendMessage)

imagesRouter.POST("/", ctr.Image.Upload)
```

**实际路径 vs spec**（不改路径的必然结果，前端 baseURL 需按实际路径配）：

| spec | 本项目实际 | 差异 |
| --- | --- | --- |
| `/articles/list` | `/api/articles/list` | 前缀 `/api`（全表同理） |
| `/articles` | `/api/articles/` | 尾斜杠；POST 会 307 重定向 |
| `/conversations` | `/api/conversations/` | 尾斜杠 |
| `/image` | `/api/image/` | 前缀 `/api` + 尾斜杠；静态目录已从 `/image` 改为 `/update` |

### 3.7 上传（落盘 + 相对路径）

- `POST /api/image/`，`multipart/form-data` 字段 `image`；限 5MB；白名单 jpg/png/webp/gif；文件名 `时间戳_随机串.ext`，落到 **`Config.AppUploadDir`**（`APP_UPLOAD_DIR`，默认 `./update`，启动时自动创建）；
- 响应返回**相对路径** `{"url":"/update/20261111211328_xxx.png"}`，`image_url.url` 也存这个相对路径；
- `router/init.go` 挂载静态目录：`r.Static(service.ImageURLPath, config.Config.AppUploadDir)`，所以 `GET /update/xxx.png` 可直接访问（E4 解除）；
- 文章 `images` 仍兼容前端提交的 base64 data URL（spec 现状），两者都能渲染；
- 失败返回 400/1002、超限 413（spec 写的 408 不复用）。

### 3.8 一致性：点赞、关注、会话

- **点赞**（文章/评论，事务）：
  1. 校验目标存在（3001/4002）；
  2. `SELECT status FROM like_record WHERE user_id=? AND target_id=? AND target_type=? FOR UPDATE`（`clause.Locking`）；
  3. 无记录 → `INSERT status=1` + `UPDATE article SET like_count=like_count+1`；记录存在且状态变化 → `UPDATE` + `±1`；状态相同 → 不动（**幂等**）；
  4. 回读 `like_count` 返回 `{liked, likes}`。
- **关注**：`INSERT IGNORE user_follows` / `DELETE`；`id==uid` → 6001，目标不存在 → 5001；返回 `{following}`。
- **会话**：始终 `lowID=min(a,b)`、`highID=max(a,b)`；先查后建，命中 `uk_pair` 唯一索引时回查返回已有；列表按 `updated_at DESC` 分页，`peer` 取另一方，`lastMessage` 用 `(conversation_id, MAX(id))` 子查询批量取（避免 N+1）。
- **消息**：
  - `afterId` → `id > afterId ORDER BY id ASC LIMIT ?`；
  - `beforeId` → `id < beforeId ORDER BY id DESC LIMIT ?` 后反转，`hasMore` 判断"是否还有更早一条"；
  - 发送：校验成员（7002）→ 插入 → `UPDATE conversation SET updated_at=NOW(3)`；`receiverId` 为对方。
- **登录限流**：`service/auth.go` 内进程内 `map[ip+name]` + 锁，5 次/15 分钟 → 2002。

### 3.9 model 小幅调整（允许范围内）

原则：**只加索引、改字段长度、加一张表**，不改字段语义、不改表名、不动关联关系。

| 文件 | 改动 | 理由 |
| --- | --- | --- |
| `model/image_url.go` | `URL` 的 tag 改 `type:longtext` | 前端提交 base64 data URL，`varchar(512)` 装不下 |
| `model/conversation.go` | `UserLowID`/`UserHighID` 加 `uniqueIndex:uk_pair` | 两人唯一会话，靠 DB 兜底并发重复创建 |
| `model/user.go` | `Name` 加 `uniqueIndex:uk_user_name`；`Slogan`/`Description`/`Location` 长度改 32/120/40 | 防并发重名；长度对齐 spec |
| `model/tag.go` | **补 `BaseModel`（主键）** | 原模型没有主键，GORM 建 many2many 关联表时 `tag_id` 缺失，AutoMigrate 直接 panic |
| `model/like_record.go` | **补 `BaseModel`（主键）** | 同上；清理点赞记录改用 `Unscoped()` 硬删除，避免软删除占唯一索引 |
| `model/scopes.go` | `Paginate` 上限 20 → 100 | 与 spec `size≤100`、`PagerForm` 上限保持一致 |
| `model/article.go`、`comment.go`、`message.go`（可选） | 复合索引：`(category,created_at)`、`(article_id,like_count)`、`(conversation_id,id)` | 列表/评论/消息查询走索引 |

三点注意：

- **唯一索引 + 软删除**：本项目没有删除用户的接口，`user.name` 唯一索引不会与 `deleted_at` 冲突；以后若加注销，改成 `uk(name, deleted_at)`（MySQL 允许多个 NULL）。
- **dev 与生产的差异**：dev 模式 `model.init()` 的 AutoMigrate 会按新 tag 自动改列/建表；生产模式跳过 AutoMigrate，需要手工执行附录 B 的 DDL（`sql/` 不在本次范围）。
- `user` 是 MySQL 关键字：GORM 生成 SQL 会自动加反引号，写裸 SQL 时自己加。

**顺带的命名修正（不改行为，同一批改动里做）**：

| 现有 | 改后 | 说明 |
| --- | --- | --- |
| `Conversation.Message []Message` | `Messages []Message` | 切片用复数，避免与单条 `Message` 混淆 |
| `LikeRecord.TargetType` 里的裸字符串 `"article"`/`"comment"` | service 侧常量 `LikeTargetArticle`/`LikeTargetComment` | 消魔法字符串；不改 model 字段类型 |
| `ImageUrl`（类型名） | `ImageURL` | Go 首字母缩写规范；同步改 `model/article.go` 与 service 引用 |
| `Article.Images []ImageUrl` | 保留 | 已是复数，语义清晰 |

---

## 4. 文件级落地清单

| 文件 | 动作 | 内容 |
| --- | --- | --- |
| `common/form.go` | 小改 | `PagerForm` 补 `json` tag + `omitempty` + 上限 100 + `Normalize()`/`Offset()` |
| `service/errs.go` | 新增 | `AppError`、15 个业务码、`With/WithMsg` |
| `model/image_url.go` | 小改 | `URL` → `longtext` |
| `model/conversation.go` | 小改 | `uk_pair(user_low_id, user_high_id)` |
| `model/user.go` | 小改 | `uk_user_name`；slogan/description/location 长度对齐 spec |
| `model/init.go` | 小改 | `AutoMigrate` 追加 `&Session{}` |
| `service/service.go` | 改 | 删 `Hello` 相关（无会话服务） |
| `service/struct.go` | 重写 | `XxxParams`/`XxxResponse` 各接口一套 + `ArticleInfo/CommentInfo/UserInfo/ConversationInfo/MessageInfo` 基座（3.4） |
| `service/auth.go` | 改 | 注册/登录（限流）/profile/logout，返回 `*AppError` |
| `service/article.go` | 改 | 列表（修 SQL、分类/关键字/排序/分页）、详情、发布、删除（级联）、点赞/取消 |
| `service/comment.go` | 重写 | 列表（likes-desc/time-desc）、发表、删除、点赞/取消 |
| `service/user.go` | 重写 | Me 聚合、UpdateMe、Public(isFollowing)、Follow/Unfollow、Articles |
| `service/conversation.go` | 重写 | 会话列表/创建、消息列表/发送 |
| `service/image.go` | 重写 | 落盘 `Config.AppUploadDir` + 返回相对路径 `/update/xxx` |
| `controller/init.go` | 改 | `ok/fail/paramErr` + `currentUID`/`bindID` 等 helper；**保留旧 `Response`**（middleware 引用它） |
| `controller/session.go` | **不改** | 保持初版：gin-contrib cookie store（`tz-sessions`） |
| `controller/auth.go` | 重写 | JSON 绑定、错误码、登录写会话、登出 `SessionClear` |
| `controller/articles.go` | 重写 | 列表/详情/发布/删除/点赞（`Article` 类型 8 个方法） |
| `controller/comment.go` | 重写 | `CommentList/CommentCreate/CommentDelete/CommentLike/CommentUnlike` |
| `controller/user.go` | 重写 | `Me/UpdateMe/Public/Follow/Unfollow/Articles` |
| `controller/conversation.go` | 重写 | `Conversations/ConversationCreate/Messages/SendMessage` |
| `controller/image.go` | 重写 | `Upload` |
| `router/router.go`、`router/init.go` | 改 | 补 handler 绑定；`r.Static(service.ImageURLPath, config.Config.AppUploadDir)` |

### 4.1 Service 方法签名

```go
// 认证
func (a *AuthSrv) Register(p AuthRegisterParams) error
func (a *AuthSrv) Login(p AuthLoginParams, ip string) (AuthLoginResult, error)
func (a *AuthSrv) Profile(uid int64) (AuthProfileResponse, error)

// 文章
func (a *ArticleSrv) List(p ArticleListParams) (ArticleListResponse, error)
func (a *ArticleSrv) Detail(id int64) (ArticleDetailResponse, error)
func (a *ArticleSrv) Create(uid int64, p ArticleCreateParams) (ArticleCreateResponse, error)
func (a *ArticleSrv) Delete(uid, id int64) error
func (a *ArticleSrv) Like(uid, id int64) (ArticleLikeResponse, error)
func (a *ArticleSrv) Unlike(uid, id int64) (ArticleLikeResponse, error)

// 评论
func (c *CommentSrv) List(articleID int64, p ArticleCommentListParams) (CommentListResponse, error)
func (c *CommentSrv) Create(uid, articleID int64, p CommentCreateParams) (CommentCreateResponse, error)
func (c *CommentSrv) Delete(uid, id int64) error
func (c *CommentSrv) Like(uid, id int64) (CommentLikeResponse, error)
func (c *CommentSrv) Unlike(uid, id int64) (CommentLikeResponse, error)

// 用户
func (u *UserSrv) Me(uid int64) (UserMeResponse, error)
func (u *UserSrv) UpdateMe(uid int64, p UserUpdateParams) (UserUpdateResponse, bool, error) // 第二返回值 = 是否改了密码
func (u *UserSrv) Public(uid, viewer int64) (UserPublicResponse, error)
func (u *UserSrv) Follow(uid, target int64) (FollowResponse, error)
func (u *UserSrv) Unfollow(uid, target int64) (FollowResponse, error)
func (u *UserSrv) Articles(target int64, p UserArticlesParams) (UserArticlesResponse, error)

// 消息
func (cv *ConversationSrv) List(uid int64, p common.PagerForm) (ConversationListResponse, error)
func (cv *ConversationSrv) Create(uid int64, p ConversationCreateParams) (ConversationCreateResponse, error)
func (cv *ConversationSrv) Messages(uid, cid int64, q MessageListQuery) (MessageListResponse, error)
func (cv *ConversationSrv) Send(uid, cid int64, p MessageSendParams) (MessageSendResponse, error)

// 上传
func (i *ImageSrv) Save(r io.Reader, filename string, size int64) (ImageUploadResponse, error)
```

### 4.2 实现注意点

1. `ArticleListResponse.List` 用 `pkg.Map(rows, toArticleInfo)` 生成（非 nil）；`Total` 的 count 查询**不能**带 `Offset/Limit`。
2. 文章列表 `category=all`/空 → 不加条件；`sort` 默认 `date-desc`（`created_at DESC, id DESC`）。
3. `cover`/`author`/`avatar` 在 service 内拼进 `ArticleInfo.Cover`/`AuthorName`/`AuthorAvatar`，controller 不做数据加工。
4. 详情不返回 `comments`（spec 可选，避免与评论接口重复加载）。
5. 所有写操作开事务；错误统一 `ErrInternal.With(err)` 并 `logger.Errorf`。
6. `UpdateMe` 用指针字段做部分更新；改 password 时 `SessionDestroyByUser`。
7. controller 每个需登录接口第一行 `uid, ok := currentUID(c); if !ok { return }`（helper 内部写 401/2003）。
8. 分页与业务参数一次 `ShouldBindQuery` 绑定，随后 `p.PagerForm = p.PagerForm.Normalize()` 再传 service（3.5）。
9. 上传的文件名、目录、体积校验都在 service，controller 只取 `c.FormFile("image")`。
10. 删除文章事务内清理 `like_record`（article + 其下 comment）、`comment`、`image_url`、`article_tags`。

---

## 5. 接口对照表（26 项，含实际路径）

| # | spec 路径 | 实际路径 | controller | service |
| --- | --- | --- | --- | --- |
| 1 | POST `/auth/register` | `/api/auth/register` | `Auth.AuthRegister` | `AuthSrv.Register` |
| 2 | POST `/auth/login` | `/api/auth/login` | `Auth.AuthLogin` | `AuthSrv.Login` |
| 3 | GET `/auth/profile` | `/api/auth/profile` | `Auth.AuthProfile` | `AuthSrv.Profile` |
| 4 | DELETE `/auth/logout` | `/api/auth/logout` | `Auth.AuthLogout` | `AuthSrv.Logout` |
| 5 | GET `/articles/list` | `/api/articles/list` | `Article.ArticlesList` | `ArticleSrv.List` |
| 6 | POST `/articles` | `/api/articles/` | `Article.ArticleCreate` | `ArticleSrv.Create` |
| 7 | GET `/articles/{id}` | `/api/articles/:id` | `Article.ArticleDetail` | `ArticleSrv.Detail` |
| 8 | DELETE `/articles/{id}` | `/api/articles/:id` | `Article.ArticleDelete` | `ArticleSrv.Delete` |
| 9 | POST `/articles/{id}/like` | 同 | `Article.ArticleLike` | `ArticleSrv.Like` |
| 10 | DELETE `/articles/{id}/like` | 同 | `Article.ArticleUnlike` | `ArticleSrv.Unlike` |
| 11 | GET `/articles/{id}/comments` | 同 | `Comment.CommentList` | `CommentSrv.List` |
| 12 | POST `/articles/{id}/comments` | 同 | `Comment.CommentCreate` | `CommentSrv.Create` |
| 13 | DELETE `/comments/{id}` | 同 | `Comment.CommentDelete` | `CommentSrv.Delete` |
| 14 | POST `/comments/{id}/like` | 同 | `Comment.CommentLike` | `CommentSrv.Like` |
| 15 | DELETE `/comments/{id}/like` | 同 | `Comment.CommentUnlike` | `CommentSrv.Unlike` |
| 16 | GET `/users/me` | 同 | `User.Me` | `UserSrv.Me` |
| 17 | PUT `/users/me` | 同 | `User.UpdateMe` | `UserSrv.UpdateMe` |
| 18 | GET `/users/{id}` | 同 | `User.Public` | `UserSrv.Public` |
| 19 | POST `/users/{id}/follow` | 同 | `User.Follow` | `UserSrv.Follow` |
| 20 | DELETE `/users/{id}/follow` | 同 | `User.Unfollow` | `UserSrv.Unfollow` |
| 21 | GET `/users/{id}/articles` | 同 | `User.Articles` | `UserSrv.Articles` |
| 22 | GET `/conversations` | `/api/conversations/` | `Conversation.Conversations` | `ConversationSrv.List` |
| 23 | POST `/conversations` | `/api/conversations/` | `Conversation.ConversationCreate` | `ConversationSrv.Create` |
| 24 | GET `/conversations/{id}/messages` | 同 | `Conversation.Messages` | `ConversationSrv.Messages` |
| 25 | POST `/conversations/{id}/messages` | 同 | `Conversation.SendMessage` | `ConversationSrv.Send` |
| 26 | POST `/image` | `/api/image/` | `Image.Upload` | `ImageSrv.Save` |

---

## 6. 分阶段实施计划

| 阶段 | 内容 | 验收 |
| --- | --- | --- |
| **P0 骨架** | model 小幅调整（3.9，含 `session.go`）、`service/errs.go`、`controller/init.go` 的 `ok/fail/paramErr/pager`；删 4 个 Hello 桩；`router` 补全部 handler；空方法先返回 `ErrInternal` | `go build ./...` 通过；AutoMigrate 后建出 `session` 表与唯一索引；所有 26 条路由不再 404/空响应 |
| **P1 认证** | JSON 绑定、限流、argon2id、登录写 cookie 会话、登出 `SessionClear` | 注册 409/1001；密码错 401/2001；6 次 429/2002；登出后 profile 401/2003 |
| **P2 文章** | 列表（修 SQL）、详情、发布（图片/标签）、删除（级联）、点赞 | `category=all` 正常；`total` 正确；非作者删 403/3002；点赞幂等 |
| **P3 评论** | 列表两种排序、发表、删除、点赞 | 文章不存在 404/3001；非本人删 403/4001 |
| **P4 用户** | Me 聚合、改资料/改密、他人资料、关注、用户文章 | 改密后当前 cookie 会话被清除（其它设备副本仍有效）；关注自己 400/6001 |
| **P5 私信** | 会话列表/创建、消息增量与历史、发送 | 重复创建同 id；非成员 403/7002；`afterId` 只返回新消息 |
| **P6 上传+联调** | `/api/image/` 上传、`./update` 落盘 + `/update` 静态目录、全流程回归 | 上传返回 `/update/xxx` 且能直接 GET 到；Apifox 按实际路径跑通 26 个接口 |

---

## 7. 验收

### 7.1 命令

```powershell
go build ./...
go vet ./...
go run .            # 或 tz-gin run
```

### 7.2 冒烟（注意 `/api` 前缀与尾斜杠）

```powershell
$b = "http://127.0.0.1:8088/api"
curl.exe -i -c c.txt -H "Content-Type: application/json" -d '{\"name\":\"小明\",\"password\":\"abc123456\"}' $b/auth/register
curl.exe -i -c c.txt -H "Content-Type: application/json" -d '{\"name\":\"小明\",\"password\":\"abc123456\"}' $b/auth/login
curl.exe -i -b c.txt $b/auth/profile
curl.exe -i -b c.txt -H "Content-Type: application/json" -d '{\"title\":\"第一次去成都\",\"content\":\"好吃\",\"category\":\"旅行\"}' $b/articles/
curl.exe -i "$b/articles/list?category=all&page=1&size=10"
curl.exe -i "$b/users/1/articles?sort=likes-desc&page=1&size=10"
curl.exe -i -b c.txt -X POST $b/articles/1/like
curl.exe -i -b c.txt -H "Content-Type: application/json" -d '{\"content\":\"写得好棒！\"}' $b/articles/1/comments
curl.exe -i -b c.txt -F "image=@test.jpg" $b/image/
```

### 7.3 边界用例

| 类别 | 用例 | 期望 |
| --- | --- | --- |
| 未登录 | 19 个需登录接口 | `401` + `2003`，`data:null` |
| 参数 | 缺 title / category 非枚举 / size=101 / id=abc | `400` + `1002` |
| 资源 | 文章/评论/用户/会话不存在 | `404` + 3001/4002/5001/7001 |
| 权限 | 删他人文章/评论、访问他人会话 | `403` + 3002/4001/7002 |
| 冲突 | 注册重名、改名重名 | `409` + `1001` |
| 限流 | 连续 6 次错密码 | `429` + `2002` |
| 幂等 | 并发点赞 20 次、重复关注 | `likes` 只 +1、`following` 稳定 |
| 软删除 | 删除后列表不含、详情/评论 404 | 一致 |
| 时间 | `createTime` 形如 `2026-11-01T08:30:00.000Z` | 可被 `new Date()` 解析 |
| 列表结构 | 文章/评论/用户文章列表 | 均为 `{list,total(,page,size)}`，空时 `list: []` |
| 分页 | 不传 page/size | 默认 1/10，不报 400（`PagerForm` 已放宽 tag） |

---

## 8. 出界项（需你或队友在上层修，本方案只保证不被卡死）

| # | 出界项 | 现状 | 本方案的绕过 | 建议 |
| --- | --- | --- | --- | --- |
| ~~E1~~ | ~~模型缺索引/长度~~ | 已解除：model 小幅调整（3.9），dev 由 AutoMigrate 生效 | — | 生产库执行附录 B DDL |
| ~~E2~~ | ~~`PagerForm` 无 json tag~~ | 已解除：`common/form.go` 小幅修改（3.5） | 入参直接绑定，响应输出 `page`/`size` | — |
| E3 | 路由前缀/尾斜杠 | `/api`、`/articles/`、`/api/image/` | 按实际路径实现并在文档标注 | 允许改 router 时对齐 spec |
| ~~E4~~ | ~~无静态目录~~ | 已解除：落盘 `Config.AppUploadDir` + `r.Static(service.ImageURLPath, config.Config.AppUploadDir)`，返回相对路径 | — | 生产环境把该目录换成持久化卷即可 |
| E5 | CORS 未开凭证 | `AllowCredentials=false` + `*` | 同域开发不受影响 | 改 `config/default.go` 回显 Origin |
| E6 | 错误中间件格式旧 | 恒 200 + 分类码 | controller 不调用 `c.Error` | 后续统一到 `common` |
| E7 | `go.sum` 被 gitignore | 依赖不可复现 | 不新增依赖即可 | 移出 .gitignore |

---

## 附录 A：主体改动一览（与原方案的取舍）

| 项 | 原方案（全层改造） | 本方案（controller+service+model 小改） |
| --- | --- | --- |
| 响应封装 | 新增 `common/response.go` | `controller/init.go` 内 `ok/fail` |
| 错误码 | 新增 `common/errcode.go` | `service/errs.go` |
| 错误中间件 | 重写 `middleware/error.go` | 不调用 `c.Error`，绕开它 |
| 会话 | 改 `config/default.go` + 自实现 store | 保留初版 `controller/session.go`（cookie store），不动 config |
| 分页 | 改 `common.PagerForm` tag | 同样小幅改 `PagerForm`（tag/上限/方法），入参直接绑定 |
| 数据表修补 | 改 model tag / 写 DDL | model 小幅调整（索引/长度/新表），dev 自动生效 |
| DTO 命名 | `Article` + `Author UserBrief` 万能结构 | 每接口 `XxxParams`/`XxxResponse` + `XxxInfo` 基座，不再有 `Post` |
| 路由 | 双前缀 + 静态目录 + 兜底 | 仅补 handler，路径不动 |
| 上传响应 | 相对路径 + `r.Static` | 落盘 `Config.AppUploadDir`、返回 `/update/xxx`、`r.Static(service.ImageURLPath, config.Config.AppUploadDir)` |

**契约变更（需前端同步，且 spec 建议一并更新）**：

| 接口 | spec | 本项目 | 前端改法 |
| --- | --- | --- | --- |
| `GET /users/{id}/articles` | `Post[]` | `{list,total,page,size}` | `data` → `data.list`，分页用 `data.total` |
| `GET /articles/{id}/comments` | `Comment[]` | `{list,total}` | `data` → `data.list`（评论不分页，`total` 为全部条数） |

## 附录 B：model 改动对应的 DDL（仅生产库手工执行）

dev 模式 `model/init.go` 的 AutoMigrate 会按 3.9 的新 tag 自动建表改列，**不需要**手工执行；生产模式跳过 AutoMigrate，需执行：

```sql
-- 1. base64 图片
ALTER TABLE `image_url` MODIFY COLUMN `url` LONGTEXT NOT NULL;

-- 2. 两人唯一会话
ALTER TABLE `conversation` ADD UNIQUE KEY `uk_pair` (`user_low_id`,`user_high_id`);

-- 3. 用户名唯一 + 长度对齐 spec
ALTER TABLE `user`
  MODIFY COLUMN `slogan` varchar(32) NOT NULL DEFAULT '',
  MODIFY COLUMN `description` varchar(120) NOT NULL DEFAULT '',
  MODIFY COLUMN `location` varchar(40) NOT NULL DEFAULT '',
  ADD UNIQUE KEY `uk_user_name` (`name`);
```

> 执行前先确认没有重名用户，否则第 3 条的唯一索引会失败。
> 会话改回 cookie store 后**不再需要** `session` 表；此前若已建过，可手工 `DROP TABLE session;`。

---

## 附录 C：实现状态（已落地）

代码已实现并通过 `go build ./...`、`go vet ./...`、`gofmt -l`（无输出）与端到端冒烟。

### C.1 已实现范围

| 模块 | 接口 | 状态 |
| --- | --- | --- |
| 认证 | register / login / profile / logout | ✅ 含登录限流（2002）、登录写 cookie 会话、登出清 cookie |
| 文章 | list / create / detail / delete / like / unlike | ✅ 分类+关键字+排序+分页、`cover` 取首图、软删除级联 |
| 评论 | list / create / delete / like / unlike | ✅ `{list,total}`、两种排序、权限校验 |
| 用户 | me / update me / public / follow / unfollow / articles | ✅ 关注聚合、`isFollowing`、改密后原会话失效 |
| 消息 | conversations / create / messages / send | ✅ 唯一会话、`afterId`/`beforeId`、`hasMore`、`peer` |
| 上传 | `POST /api/image/` | ✅ 落盘 `Config.AppUploadDir`（`./update`）+ 返回 `/update/xxx.png`，静态目录可直接访问 |

### C.2 实施中发现并修复的两个隐藏问题

1. **`Tag` / `LikeRecord` 没有主键**（未嵌 `BaseModel`）→ GORM 建 many2many 关联表时拿不到 `tag_id`，`AutoMigrate` 在 `ReorderModels` 里 nil 解引用 panic。补主键后解决，无需关闭关系迁移。
2. **`model.Paginate` 上限 20** → `size=50` 只返回 20 条。改为 100。

> 会话最终保留初版 cookie store，因此不再涉及 `session` 表与 `SID` 列名问题。

### C.3 冒烟结果（真实 MySQL `be_exam`）

| 场景 | 结果 |
| --- | --- |
| 注册 / 重名注册 | `200 / 409(1001)` |
| 登录 / 错误密码 | `200` + `Set-Cookie: tz-sessions=...` / `401(2001)` |
| 未登录访问 profile | `401` + `{"code":2003,...}` |
| 发文章（含标签/图片） | `200`，返回 `author/avatar/cover/tags/likes` |
| 列表（分类+关键字+排序+分页） | `200`，`{list,total,page,size}` |
| 详情 / 删除后详情 | `200` / `404(3001)` |
| 点赞×2 / 取消 | `likes=1`（幂等） / `likes=0` |
| 评论增删 / 评论点赞 / 评论列表 | `200`，`{list,total}` |
| 用户文章列表 | `200`，`{list,total,page,size}` |
| 关注 / 他人资料 `isFollowing` | `{following:true}` / `true` |
| 关注自己 | `400(6001)` |
| 会话创建两次（双向） | 返回同一 `id` |
| 发消息 / 拉消息 / 会话列表 | `200`，含 `lastMessage`、`hasMore`、`peer` |
| 上传 png | `200`，`url=data:image/png;base64,...` |
| 改资料（slogan/location/interests） | `200`，回显更新后的 `UserInfo` |
| 登出后再访问 | `401(2003)` |
| 参数缺失 | `400(1002)`，中文校验文案 |
| 和自己聊天 | `400(7003)` |

### C.4 仍未处理（受约束限制，见第 8 章）

- E3 路由前缀/尾斜杠：当前实际路径 `/api/...`、`/api/articles/`、`/api/image/`；
- E5 CORS 未开凭证：跨域带 cookie 需改 `config/default.go`；
- E6 `middleware/error.go` 仍是旧格式（已通过“controller 不调用 `c.Error`”绕开）；
- E7 `go.sum` 仍被 `.gitignore` 忽略。

