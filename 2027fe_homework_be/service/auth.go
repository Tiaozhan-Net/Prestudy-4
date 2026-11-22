package service

import (
	"errors"
	"strings"
	"sync"
	"time"

	"fe_exam/model"

	"github.com/alexedwards/argon2id"
	"gorm.io/gorm"
)

type AuthSrv struct {
}

// ---------------- 登录限流 ----------------

const (
	loginMaxFails = 5
	loginWindow   = 15 * time.Minute
)

type loginAttempt struct {
	fails int
	first time.Time
}

var loginLimiter = struct {
	sync.Mutex
	m map[string]*loginAttempt
}{m: make(map[string]*loginAttempt)}

func loginKey(ip, name string) string {
	return ip + "|" + strings.ToLower(strings.TrimSpace(name))
}

func loginAllowed(key string) bool {
	loginLimiter.Lock()
	defer loginLimiter.Unlock()
	att, ok := loginLimiter.m[key]
	if !ok {
		return true
	}
	if time.Since(att.first) > loginWindow {
		delete(loginLimiter.m, key)
		return true
	}
	return att.fails < loginMaxFails
}

func loginFail(key string) {
	loginLimiter.Lock()
	defer loginLimiter.Unlock()
	att, ok := loginLimiter.m[key]
	if !ok || time.Since(att.first) > loginWindow {
		loginLimiter.m[key] = &loginAttempt{fails: 1, first: time.Now()}
		return
	}
	att.fails++
}

func loginOK(key string) {
	loginLimiter.Lock()
	defer loginLimiter.Unlock()
	delete(loginLimiter.m, key)
}

// ---------------- 业务方法 ----------------

func (a *AuthSrv) Register(p AuthRegisterParams) error {
	name := strings.TrimSpace(p.Name)
	if name == "" {
		return ErrParam.WithMsg("用户名不能为空")
	}

	var count int64
	if err := model.DB.Model(&model.User{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return ErrInternal.With(err)
	}
	if count > 0 {
		return ErrUserNameExists
	}

	hash, err := argon2id.CreateHash(p.Password, argon2id.DefaultParams)
	if err != nil {
		return ErrInternal.With(err)
	}

	user := model.User{
		Name:     name,
		Password: hash,
		Avatar:   "/avatar.png",
	}
	if err := model.DB.Create(&user).Error; err != nil {
		if isDuplicateKey(err) {
			return ErrUserNameExists
		}
		return ErrInternal.With(err)
	}
	return nil
}

func (a *AuthSrv) Login(p AuthLoginParams, ip string) (AuthLoginResult, error) {
	var resp AuthLoginResult
	name := strings.TrimSpace(p.Name)
	key := loginKey(ip, name)
	if !loginAllowed(key) {
		return resp, ErrTooManyAttempts
	}

	var user model.User
	err := model.DB.Where("name = ?", name).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		loginFail(key)
		return resp, ErrCredentials
	}
	if err != nil {
		return resp, ErrInternal.With(err)
	}

	match, err := argon2id.ComparePasswordAndHash(p.Password, user.Password)
	if err != nil {
		return resp, ErrInternal.With(err)
	}
	if !match {
		loginFail(key)
		return resp, ErrCredentials
	}

	loginOK(key)
	resp = AuthLoginResult{ID: user.ID, Name: user.Name}
	return resp, nil
}

func (a *AuthSrv) Profile(uid int64) (AuthProfileResponse, error) {
	var user model.User
	err := model.DB.Where("id = ?", uid).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return AuthProfileResponse{}, ErrUserNotFound
	}
	if err != nil {
		return AuthProfileResponse{}, ErrInternal.With(err)
	}
	return AuthProfileResponse{
		ID:     user.ID,
		Name:   user.Name,
		Avatar: user.Avatar,
		Slogan: user.Slogan,
	}, nil
}
