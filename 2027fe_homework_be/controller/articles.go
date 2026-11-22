package controller

import (
	"fe_exam/service"

	"github.com/gin-gonic/gin"
)

type Article struct {
}

// ArticlesList GET /articles/list
func (a *Article) ArticlesList(c *gin.Context) {
	var p service.ArticleListParams
	if !bindQuery(c, &p) {
		return
	}
	resp, err := srv.ArticleSrv.List(p)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// ArticleDetail GET /articles/:id
func (a *Article) ArticleDetail(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	resp, err := srv.ArticleSrv.Detail(id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// ArticleCreate POST /articles
func (a *Article) ArticleCreate(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	var p service.ArticleCreateParams
	if !bindJSON(c, &p) {
		return
	}
	resp, err := srv.ArticleSrv.Create(uid, p)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// ArticleDelete DELETE /articles/:id
func (a *Article) ArticleDelete(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	id, valid := bindID(c)
	if !valid {
		return
	}
	if err := srv.ArticleSrv.Delete(uid, id); err != nil {
		fail(c, err)
		return
	}
	ok(c, nil)
}

// ArticleLike POST /articles/:id/like
func (a *Article) ArticleLike(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	id, valid := bindID(c)
	if !valid {
		return
	}
	resp, err := srv.ArticleSrv.Like(uid, id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// ArticleUnlike DELETE /articles/:id/like
func (a *Article) ArticleUnlike(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	id, valid := bindID(c)
	if !valid {
		return
	}
	resp, err := srv.ArticleSrv.Unlike(uid, id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}
