package captcha

import (
	"github.com/dchest/captcha"
	"github.com/garyburd/redigo/redis"
	"github.com/kataras/iris"

	myRedis "github.com/alsmile/goApiGateway/db/redis"
	"github.com/alsmile/goApiGateway/session"
)

const (
	CaptchaSessionName = "captcha"
)

var smsTemplates map[string]string

func init() {
	captcha.SetCustomStore(&captchaRedisStore{})
}

func IsNeedSignCaptcha(id string) (b bool) {
	redisConn := myRedis.RedisPool.Get()
	defer redisConn.Close()

	val, err := redis.Int(redisConn.Do("GET", "Captcha.Sign.Error."+id))
	if err != nil || val < 3 {
		return
	}

	b = true
	return
}

func SignError(id string) (need bool) {
	redisConn := myRedis.RedisPool.Get()
	defer redisConn.Close()

	var num int
	val, err := redis.Int(redisConn.Do("GET", "Captcha.Sign.Error."+id))
	if err == nil {
		num = val + 1
		need = num > 2
	} else {
		num = 1
	}

	redisConn.Do("SETEX", "Captcha.Sign.Error."+id, session.SessionMaxAge, num)

	return
}

func ClearSignError(id string) {
	redisConn := myRedis.RedisPool.Get()
	defer redisConn.Close()

	redisConn.Do("DEL", "Captcha.Sign.Error."+id)

	return
}

func VerifyImage(ctx iris.Context, code string) bool {
	captchaId, _ := redis.String(session.GetSession(ctx, CaptchaSessionName))
	if captchaId == "" {
		return false
	}
	return captcha.VerifyString(captchaId, code)
}
