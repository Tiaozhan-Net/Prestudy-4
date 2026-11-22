package router

import (
	"fe_exam/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Use(middleware.Error)
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	apiRouter := r.Group("/api")
	{
		imagesRouter := apiRouter.Group("/image")
		{
			imagesRouter.POST("/", ctr.Image.Upload)
		}
		authRouter := apiRouter.Group("/auth")
		{
			authRouter.POST("/register", ctr.Auth.AuthRegister)
			authRouter.POST("/login", ctr.Auth.AuthLogin)
			authRouter.GET("/profile", ctr.Auth.AuthProfile)
			authRouter.DELETE("/logout", ctr.Auth.AuthLogout)
		}
		userRouter := apiRouter.Group("/users")
		{
			userRouter.GET("/me", ctr.User.Me)
			userRouter.PUT("/me", ctr.User.UpdateMe)
			userRouter.GET("/:id", ctr.User.Public)
			userRouter.POST("/:id/follow", ctr.User.Follow)
			userRouter.DELETE("/:id/follow", ctr.User.Unfollow)
			userRouter.GET("/:id/articles", ctr.User.Articles)
		}
		articleRouter := apiRouter.Group("/articles")
		{
			articleRouter.GET("/list", ctr.Article.ArticlesList)
			articleRouter.POST("/", ctr.Article.ArticleCreate)
			articleRouter.GET("/:id", ctr.Article.ArticleDetail)
			articleRouter.DELETE("/:id", ctr.Article.ArticleDelete)

			// 笔记点赞
			articleRouter.POST("/:id/like", ctr.Article.ArticleLike)
			articleRouter.DELETE("/:id/like", ctr.Article.ArticleUnlike)

			// 笔记评论
			articleRouter.GET("/:id/comments", ctr.Comment.CommentList)
			articleRouter.POST("/:id/comments", ctr.Comment.CommentCreate)
		}
		commentRouter := apiRouter.Group("/comments")
		{
			commentRouter.DELETE("/:id", ctr.Comment.CommentDelete)
			// 评论点赞
			commentRouter.POST("/:id/like", ctr.Comment.CommentLike)
			commentRouter.DELETE("/:id/like", ctr.Comment.CommentUnlike)
		}
		conversationRouter := apiRouter.Group("/conversations")
		{
			conversationRouter.GET("/", ctr.Conversation.Conversations)
			conversationRouter.POST("/", ctr.Conversation.ConversationCreate)
			// 会话消息
			conversationRouter.GET("/:id/messages", ctr.Conversation.Messages)
			conversationRouter.POST("/:id/messages", ctr.Conversation.SendMessage)
		}
	}
}
