package captcha

import (
  "github.com/dchest/captcha"
  "idengta.xyz/dengta/db/redis"
  goRedis "github.com/garyburd/redigo/redis"
  "github.com/alsmile/goMicroServer/session"
)

const (
  CaptchaSessionName = "captcha"
)

var smsTemplates map[string]string

func init() {
  captcha.SetCustomStore(&captchaRedisStore{})
}

func IsNeedSignCaptcha(id string) (b bool) {
  redisConn := redis.RedisPool.Get()
  defer redisConn.Close()

  val, err := goRedis.Int(redisConn.Do("GET", "Captcha.Sign.Error." + id))
  if err != nil || val < 3 {
    return
  }

  b = true
  return
}

func SignError(id string) (need bool) {
  redisConn := redis.RedisPool.Get()
  defer redisConn.Close()

  var num int
  val, err := goRedis.Int(redisConn.Do("GET", "Captcha.Sign.Error." + id))
  if err == nil{
    num = val + 1
    need = num > 2
  } else {
    num = 1
  }

  redisConn.Do("SETEX", "Captcha.Sign.Error." + id, session.SessionMaxAge, num)

  return
}

func ClearSignError(id string){
  redisConn := redis.RedisPool.Get()
  defer redisConn.Close()

  redisConn.Do("DEL", "Captcha.Sign.Error." + id)

  return
}
