package controller

import (
	"fe_exam/service"

	"github.com/gin-gonic/gin"
)

type User struct {
}

// Me GET /users/me
func (u *User) Me(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	resp, err := srv.UserSrv.Me(uid)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// UpdateMe PUT /users/me
func (u *User) UpdateMe(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	var p service.UserUpdateParams
	if !bindJSON(c, &p) {
		return
	}
	resp, passwordChanged, err := srv.UserSrv.UpdateMe(uid, p)
	if err != nil {
		fail(c, err)
		return
	}
	if passwordChanged {
		// 改密码后清空当前会话（cookie store 无法吊销其它设备）
		SessionClear(c)
	}
	ok(c, resp)
}

// Public GET /users/:id
func (u *User) Public(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	resp, err := srv.UserSrv.Public(id, optionalUID(c))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// Follow POST /users/:id/follow
func (u *User) Follow(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	id, valid := bindID(c)
	if !valid {
		return
	}
	resp, err := srv.UserSrv.Follow(uid, id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// Unfollow DELETE /users/:id/follow
func (u *User) Unfollow(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	id, valid := bindID(c)
	if !valid {
		return
	}
	resp, err := srv.UserSrv.Unfollow(uid, id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// Articles GET /users/:id/articles
func (u *User) Articles(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	var p service.UserArticlesParams
	if !bindQuery(c, &p) {
		return
	}
	resp, err := srv.UserSrv.Articles(id, p)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}
