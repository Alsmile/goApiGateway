package models

type User struct {
  Id       string `json:"id"`
  Email    string `json:"email"`
  Phone    string `json:"phone"`
  Username string `json:"username"`
  Password string `json:"password"`
}
