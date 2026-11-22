package controller

import (
	"encoding/gob"
	"errors"
	"net/http"
	"strings"

	"fe_exam/common"
	"fe_exam/logger"
	"fe_exam/service"
	vl "fe_exam/service/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response 旧响应结构：middleware/error.go 仍引用，保留勿删
type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Code    uint64 `json:"code,omitempty"`
}

// ResponseNew 保留兼容
func ResponseNew(_ *gin.Context, obj any) *Response {
	return &Response{Success: true, Data: obj}
}

var srv = service.New()

// sessionKey 会话中存放登录用户的键（controller/session.go 保持初版，不在其中定义常量）
const sessionKey = "user-session"

func init() {
	gob.Register(UserSession{})
}

// successMessage 成功响应的固定文案（想换英文/换语气只改这一处）
const successMessage = "success"

// ok 成功响应：{code:0, message:<successMessage>, data:...}
func ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": successMessage, "data": data})
}

// fail 失败响应：把 *service.AppError 映射为 HTTP 状态码 + 业务码
func fail(c *gin.Context, err error) {
	var appErr *service.AppError
	if errors.As(err, &appErr) {
		if appErr.Err != nil {
			logger.Errorf("controller: %v", appErr.Err)
		}
		c.JSON(appErr.HTTP, gin.H{"code": appErr.Code, "message": appErr.Msg, "data": nil})
		return
	}
	logger.Errorf("controller: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"code": 9001, "message": "服务器内部错误", "data": nil})
}

// paramErr 参数校验失败 -> 1002/400
func paramErr(c *gin.Context, err error) {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) && vl.Trans != nil {
		msgs := make([]string, 0, len(verrs))
		for _, v := range verrs.Translate(vl.Trans) {
			msgs = append(msgs, v)
		}
		fail(c, service.ErrParam.WithMsg(strings.Join(msgs, "；")))
		return
	}
	fail(c, service.ErrParam.With(err))
}

// currentUID 当前登录用户；未登录写 401/2003 并返回 false
func currentUID(c *gin.Context) (int64, bool) {
	us, typeOK := SessionGet(c, sessionKey).(UserSession)
	if !typeOK || us.ID <= 0 {
		fail(c, service.ErrUnauthenticated)
		return 0, false
	}
	return int64(us.ID), true
}

// optionalUID 未登录返回 0（浏览类接口用）
func optionalUID(c *gin.Context) int64 {
	us, typeOK := SessionGet(c, sessionKey).(UserSession)
	if !typeOK || us.ID <= 0 {
		return 0
	}
	return int64(us.ID)
}

// bindID 绑定路径参数 :id
func bindID(c *gin.Context) (int64, bool) {
	var form common.IDUriForm
	if err := c.ShouldBindUri(&form); err != nil || form.ID <= 0 {
		fail(c, service.ErrParam.WithMsg("id 不合法"))
		return 0, false
	}
	return int64(form.ID), true
}

// bindJSON 绑定请求体
func bindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		paramErr(c, err)
		return false
	}
	return true
}

// bindQuery 绑定查询参数
func bindQuery(c *gin.Context, obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		paramErr(c, err)
		return false
	}
	return true
}
