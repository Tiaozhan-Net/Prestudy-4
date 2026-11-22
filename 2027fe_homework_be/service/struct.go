package service

import (
	"fe_exam/common"
	"time"
)

// ============================== 基础类型 ==============================

// Category 文章分类（与前端/接口文档一致）
type Category string

const (
	CategoryMusic      Category = "音乐"
	CategoryFood       Category = "美食"
	CategoryPet        Category = "萌宠"
	CategoryTravel     Category = "旅行"
	CategoryGame       Category = "游戏"
	CategoryMovie      Category = "影视"
	CategoryBook       Category = "读书"
	CategoryPhotograph Category = "摄影"
	CategoryCoding     Category = "编程"
	CategoryOther      Category = "其他"
)

// Sort 列表排序
type Sort string

const (
	SortDateDesc  Sort = "date-desc"
	SortDateAsc   Sort = "date-asc"
	SortLikesDesc Sort = "likes-desc"
)

// 点赞对象类型（like_record.target_type）
const (
	LikeTargetArticle = "article"
	LikeTargetComment = "comment"
)

// FormatTime 统一时间格式：2026-11-01T08:30:00.000Z
func FormatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}

// ============================== 认证 ==============================

type AuthRegisterParams struct {
	Name     string `json:"name" binding:"required,min=1,max=20"`
	Password string `json:"password" binding:"required,min=1,max=20"`
}

type AuthLoginParams struct {
	Name     string `json:"name" binding:"required,min=1,max=20"`
	Password string `json:"password" binding:"required,min=1,max=20"`
}

// AuthProfileResponse GET /auth/profile
type AuthProfileResponse struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Slogan string `json:"slogan"`
}

// AuthLoginResult 登录结果：只用于写会话，不进响应体
type AuthLoginResult struct {
	ID   int64
	Name string
}

// ============================== 文章 ==============================

// ArticleInfo 文章公共字段：列表项 / 详情 / 发布返回共用
type ArticleInfo struct {
	ID           int64    `json:"id"`
	AuthorID     int64    `json:"authorId"`
	AuthorName   string   `json:"author"` // spec 的 author 是名字，不是对象
	AuthorAvatar string   `json:"avatar"`
	Cover        string   `json:"cover"` // = Images[0]
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
	Category Category `form:"category" binding:"omitempty,oneof=all 音乐 美食 萌宠 旅行 游戏 影视 读书 摄影 编程 其他"`
	Keyword  string   `form:"keyword" binding:"max=40"`
	Sort     Sort     `form:"sort" binding:"omitempty,oneof=date-desc date-asc likes-desc"`
}

// ArticleListResponse GET /articles/list
type ArticleListResponse struct {
	common.PagerForm
	List  []ArticleInfo `json:"list"`
	Total int64         `json:"total"`
}

// ArticleDetailResponse GET /articles/{id}
type ArticleDetailResponse struct {
	ArticleInfo
}

type ArticleCreateParams struct {
	Title    string   `json:"title" binding:"required,max=40"`
	Content  string   `json:"content" binding:"required,max=500"`
	Category Category `json:"category" binding:"required,oneof=音乐 美食 萌宠 旅行 游戏 影视 读书 摄影 编程 其他"`
	Images   []string `json:"images" binding:"max=9"`
	Tags     []string `json:"tags" binding:"max=10,dive,max=40"`
}

// ArticleCreateResponse POST /articles
type ArticleCreateResponse struct {
	ArticleInfo
}

// ArticleLikeResponse POST/DELETE /articles/{id}/like
type ArticleLikeResponse struct {
	Liked     bool  `json:"liked"`
	LikeCount int64 `json:"likes"`
}

// ============================== 评论 ==============================

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

// CommentListResponse GET /articles/{id}/comments
type CommentListResponse struct {
	List  []CommentInfo `json:"list"`
	Total int64         `json:"total"`
}

type CommentCreateParams struct {
	Content string `json:"content" binding:"required,max=80"`
}

// CommentCreateResponse POST /articles/{id}/comments
type CommentCreateResponse struct {
	CommentInfo
}

// CommentLikeResponse POST/DELETE /comments/{id}/like
type CommentLikeResponse struct {
	Liked     bool  `json:"liked"`
	LikeCount int64 `json:"likes"`
}

// ============================== 用户 ==============================

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

// UserMeResponse GET /users/me
type UserMeResponse struct {
	UserInfo
}

// UserUpdateResponse PUT /users/me
type UserUpdateResponse struct {
	UserInfo
}

// UserUpdateParams 全指针：区分"没传"与"传空"
type UserUpdateParams struct {
	Name        *string  `json:"name" binding:"omitempty,min=1,max=20"`
	Avatar      *string  `json:"avatar"`
	Slogan      *string  `json:"slogan" binding:"omitempty,max=32"`
	Description *string  `json:"description" binding:"omitempty,max=120"`
	Location    *string  `json:"location" binding:"omitempty,max=40"`
	Interests   []string `json:"interests"`
	Password    *string  `json:"password" binding:"omitempty,min=1,max=20"`
}

// UserPublicResponse GET /users/{id}
type UserPublicResponse struct {
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

// UserArticlesResponse GET /users/{id}/articles
type UserArticlesResponse struct {
	common.PagerForm
	List  []ArticleInfo `json:"list"`
	Total int64         `json:"total"`
}

// FollowResponse POST/DELETE /users/{id}/follow
type FollowResponse struct {
	Following bool `json:"following"`
}

// ============================== 消息 ==============================

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

// ConversationListResponse GET /conversations
type ConversationListResponse struct {
	common.PagerForm
	List  []ConversationInfo `json:"list"`
	Total int64              `json:"total"`
}

type ConversationCreateParams struct {
	UserID int64 `json:"userId" binding:"required,min=1"`
}

// ConversationCreateResponse POST /conversations
type ConversationCreateResponse struct {
	ConversationInfo
}

type MessageListQuery struct {
	AfterID  int64 `form:"afterId" binding:"omitempty,min=0"`
	BeforeID int64 `form:"beforeId" binding:"omitempty,min=1"`
	Limit    int   `form:"limit" binding:"omitempty,min=1,max=100"`
}

// MessageListResponse GET /conversations/{id}/messages
type MessageListResponse struct {
	List    []MessageInfo `json:"list"`
	HasMore bool          `json:"hasMore"`
	Peer    UserBrief     `json:"peer"`
}

type MessageSendParams struct {
	Content string `json:"content" binding:"required,max=1000"`
}

// MessageSendResponse POST /conversations/{id}/messages
type MessageSendResponse struct {
	MessageInfo
}

// ============================== 上传 ==============================

// ImageUploadResponse POST /image
type ImageUploadResponse struct {
	URL string `json:"url"`
}
