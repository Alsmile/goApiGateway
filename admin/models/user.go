package models

import (
  "strings"
)

type User struct {
  Id           string `json:"id"`
  Email        string `json:"email"`
  Phone        string `json:"phone"`
  Username     string `json:"username"`
  Password     string `json:"password"`
  Captcha      string `json:"captcha"`
  ActiveCode   string `json:"activeCode"`
  PasswordCode string `json:"passwordCode"`
  RememberMe   bool `json:"rememberMe"`
}

func (user *User) GetUsername() {
  pos := strings.Index(user.Email, "@")
  if pos < 0 {
    return
  }

  user.Username = user.Email[0: pos]
}
