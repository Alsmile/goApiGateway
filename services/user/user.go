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
  "github.com/alsmile/goMicroServer/models"
  "github.com/alsmile/goMicroServer/db/mongo"
  "github.com/alsmile/goMicroServer/utils"
  "github.com/alsmile/goMicroServer/services"
  "strings"
  "github.com/alsmile/goMicroServer/services/email"
  "github.com/alsmile/goMicroServer/db/redis"
  "gopkg.in/kataras/iris.v6"
  "gopkg.in/mgo.v2/bson"
)

func EncodePassword(pwd string) string {
  mac := hmac.New(sha256.New, []byte(utils.GlobalConfig.Secret))
  mac.Write([]byte(pwd))
  str := hex.EncodeToString(mac.Sum([]byte("goMicroServer")))

  return str
}

func AddUser(u *models.User) (err error) {
  if u.Profile.Email == "" && u.Profile.Phone == "" {
    err = errors.New(services.ErrorEmailEmpty)
    return
  }

  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Find(
    bson.M{"profile.email": u.Profile.Email}).Select(bson.M{"_id":1}).One(&u)

  if u.Id != "" {
    err = errors.New(services.ErrorEmailExists)
    return
  }

  u.Id = bson.NewObjectId()
  u.Active.Code = bson.NewObjectId().Hex()
  u.Profile.GetUsername()
  u.Password = EncodePassword(u.Password)
  u.CreatedAt = time.Now().UTC()

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Insert(u)
  u.Password = ""
  if err != nil {
    if strings.Contains(err.Error(), "unique_email") {
      err = errors.New(services.ErrorEmailExists)
    } else {
      log.Printf("services.user.user.AddUser error: %v\r\n", err)
      err = errors.New(services.ErrorSave)
    }
  } else {
    email.SendSignUpEmail(u)
  }

  return
}

func GetUserByPassword(u *models.User) (err error) {
  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Find(
    bson.M{
      "profile.email": u.Profile.Email,
      "password": EncodePassword(u.Password),
    }).Select(bson.M{"_id":1, "profile": 1, "roles": 1, "active": 1}).One(&u)
  u.Password = ""
  if err != nil {
    log.Printf("[error]services.user.GetUserByPassword: err=%v\r\n", err)
    err = errors.New("用户名或密码错误")
  } else if u.Active.Code != "" {
    err = errors.New(services.ErrorNoActive)
  }

  return
}

func Active(u *models.User) (err error) {
  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err =  mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Find(
    bson.M{"active.code": u.Active.Code}).Select(bson.M{"_id":1, "profile": 1, "roles": 1, "active": 1}).One(&u)
  if err != nil {
    err = errors.New(services.ErrorActiveCode)
    return
  }

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).
    Update(bson.M{"active.code": u.Active.Code}, bson.M{"$set": bson.M{"active": bson.M{"time": time.Now().UTC()}}})
  if err != nil {
    err = errors.New(services.ErrorSave)
    return
  }

  return
}

func ForgetPassword(u *models.User) (err error) {
  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err =  mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Find(
    bson.M{
      "profile.email": u.Profile.Email,
    }).Select(bson.M{"_id":1}).One(&u)
  if err != nil || u.Id == "" {
    err = errors.New(services.ErrorUserNoExists)
    return
  }

  redisConn := redis.RedisPool.Get()
  defer redisConn.Close()

  u.PasswordCode = utils.GetGuid()
  _, err = redisConn.Do("SETEX", u.PasswordCode, services.TokenValidHours*3600, u.Id)
  if err != nil {
    log.Printf("services.user.user.ForgetPassword: redis error=%v\r\n", err)
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

  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  u.Id = bson.ObjectIdHex(utils.String(val))
  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).
    Update(bson.M{"_id": u.Id}, bson.M{"$set": bson.M{"password": EncodePassword(u.Password)}})
  u.Password = ""
  if err != nil {
    err = errors.New(services.ErrorUserNoExists)
  }

  return
}

func GetUserById(u *models.User) (err error) {
  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err =  mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Find(bson.M{"_id": u.Id}).
    Select(bson.M{"_id": 1, "profile": 1, "roles": 1}).One(&u)
  if err != nil {
    err = errors.New(services.ErrorUserNoExists)
  }

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
  user.Id = bson.ObjectIdHex(utils.String(claims["uid"]))
  return
}
