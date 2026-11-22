package controller

import (
	"fe_exam/service"

	"github.com/gin-gonic/gin"
)

type Comment struct {
}

// CommentList GET /articles/:id/comments
func (cm *Comment) CommentList(c *gin.Context) {
	articleID, valid := bindID(c)
	if !valid {
		return
	}
	var p service.ArticleCommentListParams
	if !bindQuery(c, &p) {
		return
	}
	resp, err := srv.CommentSrv.List(articleID, p)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// CommentCreate POST /articles/:id/comments
func (cm *Comment) CommentCreate(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	articleID, valid := bindID(c)
	if !valid {
		return
	}
	var p service.CommentCreateParams
	if !bindJSON(c, &p) {
		return
	}
	resp, err := srv.CommentSrv.Create(uid, articleID, p)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// CommentDelete DELETE /comments/:id
func (cm *Comment) CommentDelete(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	id, valid := bindID(c)
	if !valid {
		return
	}
	if err := srv.CommentSrv.Delete(uid, id); err != nil {
		fail(c, err)
		return
	}
	ok(c, nil)
}

// CommentLike POST /comments/:id/like
func (cm *Comment) CommentLike(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	id, valid := bindID(c)
	if !valid {
		return
	}
	resp, err := srv.CommentSrv.Like(uid, id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// CommentUnlike DELETE /comments/:id/like
func (cm *Comment) CommentUnlike(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	id, valid := bindID(c)
	if !valid {
		return
	}
	resp, err := srv.CommentSrv.Unlike(uid, id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}
