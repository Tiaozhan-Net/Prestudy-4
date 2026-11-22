package controller

import (
	"fe_exam/common"
	"fe_exam/service"

	"github.com/gin-gonic/gin"
)

type Conversation struct {
}

// Conversations GET /conversations
func (cv *Conversation) Conversations(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	var p common.PagerForm
	if !bindQuery(c, &p) {
		return
	}
	resp, err := srv.ConversationSrv.List(uid, p)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// ConversationCreate POST /conversations
func (cv *Conversation) ConversationCreate(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	var p service.ConversationCreateParams
	if !bindJSON(c, &p) {
		return
	}
	resp, err := srv.ConversationSrv.Create(uid, p)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// Messages GET /conversations/:id/messages
func (cv *Conversation) Messages(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	id, valid := bindID(c)
	if !valid {
		return
	}
	var q service.MessageListQuery
	if !bindQuery(c, &q) {
		return
	}
	resp, err := srv.ConversationSrv.Messages(uid, id, q)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}

// SendMessage POST /conversations/:id/messages
func (cv *Conversation) SendMessage(c *gin.Context) {
	uid, logged := currentUID(c)
	if !logged {
		return
	}
	id, valid := bindID(c)
	if !valid {
		return
	}
	var p service.MessageSendParams
	if !bindJSON(c, &p) {
		return
	}
	resp, err := srv.ConversationSrv.Send(uid, id, p)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}
