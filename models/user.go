package models

import (
  "strings"
  "time"
  "gopkg.in/mgo.v2/bson"
)

type User struct {
  Id       bson.ObjectId `json:"_id" bson:"_id"`
  Profile  UserProfile `json:"profile" `
  Password string `json:"password" `
  Captcha  string `json:"captcha" bson:",omitempty"`
  Active struct {
    Code string `json:"code" `
    Time time.Time `json:"time" bson:",omitempty"`
  } `json:"active" `
  PasswordCode string `json:"passwordCode" bson:"passwordCode"`
  RememberMe   bool `json:"rememberMe" bson:"rememberMe,omitempty"`
  CreatedAt    time.Time `json:"createdAt" bson:"createdAt,omitempty"`
  UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
  DeletedAt    time.Time `json:"deletedAt" bson:"deletedAt,omitempty"`
}

type UserProfile struct {
  Email    string `json:"email" `
  Phone    string `json:"phone" `
  Username string `json:"username" `
}

func (profile *UserProfile) GetUsername() {
  pos := strings.Index(profile.Email, "@")
  if pos < 0 {
    return
  }

  profile.Username = profile.Email[0: pos]
}

type QuotedUser struct {
  Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
  Email    string `json:"email" bson:"email,omitempty"`
  Phone    string `json:"phone" bson:"phone,omitempty"`
  Username string `json:"username" bson:"username,omitempty"`
}
