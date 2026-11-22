package service

import "fmt"

// AppError 业务错误：Code 进响应体，HTTP 决定状态码，Msg 为默认中文文案
type AppError struct {
	Code int
	HTTP int
	Msg  string
	Err  error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

func (e *AppError) Unwrap() error { return e.Err }

// With 附加内层错误（返回副本，避免污染全局变量）
func (e *AppError) With(err error) *AppError {
	cp := *e
	cp.Err = err
	return &cp
}

// WithMsg 覆盖文案（返回副本）
func (e *AppError) WithMsg(msg string) *AppError {
	cp := *e
	cp.Msg = msg
	return &cp
}

// 业务错误码表（与 Tenzors Blog API 一致）
var (
	ErrParam            = &AppError{Code: 1002, HTTP: 400, Msg: "参数不合法"}
	ErrUserNameExists   = &AppError{Code: 1001, HTTP: 409, Msg: "用户名已存在"}
	ErrCredentials      = &AppError{Code: 2001, HTTP: 401, Msg: "用户名或密码错误"}
	ErrTooManyAttempts  = &AppError{Code: 2002, HTTP: 429, Msg: "尝试次数过多，请稍后再试"}
	ErrUnauthenticated  = &AppError{Code: 2003, HTTP: 401, Msg: "未登录或登录已过期"}
	ErrArticleNotFound  = &AppError{Code: 3001, HTTP: 404, Msg: "文章不存在"}
	ErrArticleForbidden = &AppError{Code: 3002, HTTP: 403, Msg: "无权限删除该文章"}
	ErrCommentForbidden = &AppError{Code: 4001, HTTP: 403, Msg: "无权限删除该评论"}
	ErrCommentNotFound  = &AppError{Code: 4002, HTTP: 404, Msg: "评论不存在"}
	ErrUserNotFound     = &AppError{Code: 5001, HTTP: 404, Msg: "用户不存在"}
	ErrFollowSelf       = &AppError{Code: 6001, HTTP: 400, Msg: "不能关注自己"}

	ErrConversationNotFound  = &AppError{Code: 7001, HTTP: 404, Msg: "会话不存在"}
	ErrConversationForbidden = &AppError{Code: 7002, HTTP: 403, Msg: "无权访问该会话"}
	ErrChatSelf              = &AppError{Code: 7003, HTTP: 400, Msg: "不能和自己聊天"}

	ErrInternal = &AppError{Code: 9001, HTTP: 500, Msg: "服务器内部错误"}
)
