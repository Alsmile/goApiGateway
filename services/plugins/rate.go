package plugins

import (
	myRedis "github.com/alsmile/goApiGateway/db/redis"
	"github.com/garyburd/redigo/redis"
)

// RateLimit 速率限制
func RateLimit(host, method, url, ip string, rate uint64) bool {
	redisConn := myRedis.RedisPool.Get()
	defer redisConn.Close()

	val, _ := redis.Uint64(redisConn.Do("GET", "rateLimit:"+host+method+url))
	if val > rate {
		return true
	}

	val = val + 1
	redisConn.Do("SETEX", "rateLimit:"+host+method+url, 60, val)

	return false
}
