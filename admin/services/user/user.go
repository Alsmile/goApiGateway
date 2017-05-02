package user

import (
  "log"
  "fmt"
  "time"
  "crypto/hmac"
  "crypto/sha256"
  "encoding/hex"
  "github.com/dgrijalva/jwt-go"
  "github.com/alsmile/goMicroServer/admin/models"
  "github.com/alsmile/goMicroServer/db/pq"
  "github.com/alsmile/goMicroServer/utils"
)

func EncodePassword(pwd string) string {
  mac := hmac.New(sha256.New, []byte(utils.GlobalConfig.Secret))
  mac.Write([]byte(pwd))
  str := hex.EncodeToString(mac.Sum([]byte("goMicroServer")))

  return str
}

func AddUser(u *models.User) (err error) {
  _, err = pq.ConnPool.Exec("INSERT INTO users (id, email, username, password) VALUES ($1, $2, $3); ", u.Id, u.Email, u.Username, EncodePassword(u.Password))
  return
}

func GetUserByPassword(u *models.User) (err error) {
  err = pq.ConnPool.QueryRow("SELECT id, email,username FROM users WHERE email=$1 and password=$2", u.Email, EncodePassword(u.Password)).
    Scan(&u.Id, &u.Email, &u.Username)

  return
}

func GetUserById(u *models.User) (err error) {
  err = pq.ConnPool.QueryRow("SELECT id, email,username FROM users WHERE id=$1", u.Id).
    Scan(&u.Id, &u.Email, &u.Username)

  return
}

func GetToken(u *models.User, hours int) (data string) {
  if u == nil {
    return
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "uid": u.Id,
    "exp": time.Now().Add(time.Hour * time.Duration(hours)).Unix(),
  })
  data, err := token.SignedString([]byte(utils.GlobalConfig.Jwt))
  if err != nil {
    log.Printf("service.user.GetToken: error= %s", err)
  }
  return
}

func ValidToken(t string) (user map[string]interface{}) {
  token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("签名方法错误: %v", token.Header["alg"])
    }
    return []byte(utils.GlobalConfig.Jwt), nil
  })

  if err != nil {
    log.Printf("service.user.ValidToken.Parse: error= %s", err)
    return
  }

  claims, ok := token.Claims.(jwt.MapClaims)
  if !ok || !token.Valid {
    log.Printf("service.user.ValidToken.Claims error")
    return
  }
  user = claims
  return
}
