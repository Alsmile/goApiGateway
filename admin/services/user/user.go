package user

import (
  "log"
  "fmt"
  "time"
  "crypto/hmac"
  "crypto/sha256"
  "encoding/hex"
  "errors"
  "github.com/dgrijalva/jwt-go"
  "github.com/alsmile/goMicroServer/admin/models"
  "github.com/alsmile/goMicroServer/db/pq"
  "github.com/alsmile/goMicroServer/utils"
  "github.com/alsmile/goMicroServer/services"
  "strings"
  "github.com/alsmile/goMicroServer/services/email"
  "github.com/alsmile/goMicroServer/db/redis"
  "gopkg.in/kataras/iris.v6"
)

func EncodePassword(pwd string) string {
  mac := hmac.New(sha256.New, []byte(utils.GlobalConfig.Secret))
  mac.Write([]byte(pwd))
  str := hex.EncodeToString(mac.Sum([]byte("goMicroServer")))

  return str
}

func AddUser(u *models.User) (err error) {
  u.Id = utils.GetGuid()
  u.ActiveCode = utils.GetGuid()
  u.GetUsername()
  _, err = pq.ConnPool.Exec("INSERT INTO users (id, email, phone, username, password, \"activeCode\") VALUES ($1,$2,$3,$4,$5,$6); ",
    u.Id, u.Email, u.Phone, u.Username, EncodePassword(u.Password), u.ActiveCode)

  if err != nil {
    if strings.Contains(err.Error(), "unique_email") {
      err = errors.New(services.ErrorEmailExists)
    } else {
      log.Printf("admin.services.user.user.AddUser error: %v\r\n", err)
      err = errors.New(services.ErrorSave)
    }
  } else {
    email.SendSignUpEmail(u)
  }

  return
}

func GetUserByPassword(u *models.User) (err error) {
  err = pq.ConnPool.QueryRow("SELECT id, email,username, \"activeCode\" FROM users WHERE email=$1 and password=$2", u.Email, EncodePassword(u.Password)).
    Scan(&u.Id, &u.Email, &u.Username, &u.ActiveCode)

  if err != nil {
    err = errors.New("用户名或密码错误")
  } else if u.ActiveCode != "" {
    err = errors.New(services.ErrorNoActive)
  }
  return
}

func Active(u *models.User) (err error) {
  err = pq.ConnPool.QueryRow("SELECT id,email,username FROM users WHERE \"activeCode\"=$1 ", u.ActiveCode).
    Scan(&u.Id, &u.Email, &u.Username)
  if err != nil || u.Id == "" || u.ActiveCode == ""{
    err = errors.New(services.ErrorActiveCode)
    return
  }

  _, err = pq.ConnPool.Exec("UPDATE users SET \"activeCode\"=$1 ", "")
  if err != nil {
    log.Printf("admin.services.user.user.Active error: %v\r\n", err)
    err = errors.New(services.ErrorSave)
  }

  return
}

func ForgetPassword(u *models.User) (err error) {
  err = pq.ConnPool.QueryRow("SELECT id,username FROM users WHERE email=$1 ", u.Email).Scan(&u.Id, &u.Username)
  if err != nil || u.Id == "" {
    err = errors.New(services.ErrorUserNoExists)
    return
  }

  redisConn := redis.RedisPool.Get()
  defer redisConn.Close()

  u.PasswordCode = utils.GetGuid()
  _, err = redisConn.Do("SETEX", u.PasswordCode, services.TokenValidHours*3600, u.Id)
  if err != nil {
    log.Printf("admin.services.user.user.ForgetPassword: redis error=%v\r\n", err)
    err = errors.New(services.ErrorSave)
    return
  }
  email.SendForgetPasswordEmail(u)

  return
}

func NewPassword(u *models.User) (err error) {
  redisConn := redis.RedisPool.Get()
  defer redisConn.Close()
  val, err := redisConn.Do("GET", u.PasswordCode)
  if err != nil {
    err = errors.New(services.ErrorCaptchaCode)
    return
  }

  u.Id = utils.String(val)
  err = pq.ConnPool.QueryRow("SELECT email,username FROM users WHERE id=$1 ", u.Id).Scan(&u.Email, &u.Username)
  if err != nil || u.Email == "" {
    err = errors.New(services.ErrorUserNoExists)
    return
  }

  _, err = pq.ConnPool.Exec("UPDATE users SET password=$1 where id=$2", EncodePassword(u.Password), u.Id)
  if err != nil {
    log.Printf("admin.services.user.user.NewPassword error: %v\r\n", err)
    err = errors.New(services.ErrorSave)
  }

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

func ValidToken(ctx *iris.Context, user *models.User) {
  data := ctx.RequestHeader("Authorization")
  if data == "" || !strings.HasPrefix(data, services.HeaderTrim) {
    ctx.SetStatusCode(401)
    return
  }

  t := data[len(services.HeaderTrim):]
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
    return
  }
  user.Id = utils.String(claims["uid"])
  return
}
