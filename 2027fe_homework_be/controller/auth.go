package controller

import (
	"fe_exam/service"

	"github.com/gin-gonic/gin"
)

type Auth struct {
}

// AuthRegister POST /auth/register
func (a *Auth) AuthRegister(c *gin.Context) {
	var p service.AuthRegisterParams
	if !bindJSON(c, &p) {
		return
	}
	if err := srv.AuthSrv.Register(p); err != nil {
		fail(c, err)
		return
	}
	ok(c, nil)
}

// AuthLogin POST /auth/login
func (a *Auth) AuthLogin(c *gin.Context) {
	var p service.AuthLoginParams
	if !bindJSON(c, &p) {
		return
	}
	res, err := srv.AuthSrv.Login(p, c.ClientIP())
	if err != nil {
		fail(c, err)
		return
	}
	SessionSet(c, sessionKey, UserSession{ID: int(res.ID), Username: res.Name})
	ok(c, nil)
}

// AuthProfile GET /auth/profile
func (a *Auth) AuthProfile(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	resp, err := srv.AuthSrv.Profile(uid)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// AuthLogout DELETE /auth/logout
func (a *Auth) AuthLogout(c *gin.Context) {
	if _, logged := currentUID(c); !logged {
		return
	}
	SessionClear(c)
	ok(c, nil)
}
