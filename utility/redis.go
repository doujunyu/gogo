package utility

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
	"sync"
)

var sendRedisLine redis.Conn
var sendRedisLineOnce sync.Once

// RedisLine 建立链接,,实现单例模式
func RedisLine() redis.Conn {
	sendRedisLineOnce.Do(func() {
		REDIS_HOST := os.Getenv("REDIS_HOST")
		REDIS_PORT := os.Getenv("REDIS_PORT")
		REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
		REDIS_SELECT := os.Getenv("REDIS_SELECT")
		redisLine, err := redis.Dial("tcp", REDIS_HOST+":"+REDIS_PORT)
		if err != nil {
			return
		}
		_, err = redisLine.Do("auth", REDIS_PASSWORD)
		if err != nil {
			_ = redisLine.Close()
			return
		}
		_, err = redisLine.Do("select", REDIS_SELECT)
		if err != nil {
			_ = redisLine.Close()
			return
		}
		if err == nil {
			sendRedisLine = redisLine
		}
		fmt.Println("redis链接执行")
	})
	return sendRedisLine
}








