package session
import (
  "net/http"
  "github.com/kataras/iris/context"
  "github.com/alsmile/goApiGateway/utils"
  "github.com/alsmile/goApiGateway/db/redis"
)

const (
  SessionId  = "sid"
  // 10分钟，只用于短期存储，比如验证码
  SessionMaxAge = 10*60
)

func GetSessionId(ctx context.Context) (sessionId string) {
  sessionId = ctx.GetCookie(SessionId)
  if sessionId == "" {
    sessionId = utils.GetGuid()
    cookie := &http.Cookie{}
    cookie.Path = "/"
    cookie.HttpOnly = true
    cookie.Name = SessionId
    cookie.Value = sessionId
    ctx.SetCookie(cookie)
  }

  return
}

func SetSession(ctx context.Context, name string, val interface{}) error {
  sessionId := GetSessionId(ctx)

  if sessionId == "" {
    sessionId = GetSessionId(ctx)
  }

  redisConn := redis.RedisPool.Get()
  defer redisConn.Close()
  _, err := redisConn.Do("SETEX", sessionId + "." + name, SessionMaxAge, val)
  return err
}

func GetSession(ctx context.Context, name string) (val interface{}, err error) {
  sessionId := GetSessionId(ctx)
  if sessionId == "" {
    return
  }

  redisConn := redis.RedisPool.Get()
  defer redisConn.Close()

  val, err = redisConn.Do("GET", sessionId + "." + name)
  return
}
