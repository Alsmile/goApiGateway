package models

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// User 用户数据结构
type User struct {
	ID       bson.ObjectId `json:"_id" bson:"_id"`
	Profile  UserProfile   `json:"profile" `
	Password string        `json:"password" `
	Captcha  string        `json:"captcha" bson:",omitempty"`
	Active   struct {
		Code string    `json:"code" `
		Time time.Time `json:"time" bson:",omitempty"`
	} `json:"active" `
	PasswordCode string    `json:"passwordCode" bson:"passwordCode"`
	RememberMe   bool      `json:"rememberMe" bson:"rememberMe,omitempty"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
	DeletedAt    time.Time `json:"deletedAt" bson:"deletedAt,omitempty"`
}

// UserProfile 用户基本信息数据结构
type UserProfile struct {
	Email    string `json:"email" `
	Phone    string `json:"phone" `
	Username string `json:"username" `
	// 下面时对接第三方用户系统时，存储的用户信息
	UserID   string `json:"userId" bson:"userId,omitempty"`
	Nickname string `json:"nickname" bson:"nickname,omitempty"`
}

// GetUsername 当username为空时，截取邮箱用户名
func (profile *UserProfile) GetUsername() {
	pos := strings.Index(profile.Email, "@")
	if pos < 0 {
		return
	}

	profile.Username = profile.Email[0:pos]
}

// QuotedUser 被其他数据引用时的精简用户信息
type QuotedUser struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email    string        `json:"email" bson:"email,omitempty"`
	Phone    string        `json:"phone" bson:"phone,omitempty"`
	Username string        `json:"username" bson:"username,omitempty"`
}
