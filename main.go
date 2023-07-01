package subway

import (
	"github.com/gomodule/redigo/redis"
	"github.com/simonks2016/Subway/ViewModel"
	"strings"
	"time"
)

type Subway struct {
	redis *redis.Pool

	d map[string]ViewModel.ModelOperation[any]
}

func NewSubway(redisPool *redis.Pool) *Subway {

	return &Subway{redis: redisPool}
}

func NewRedisConnWithSubway(address, userName, password string) *redis.Pool {

	redisAuth := userName + ":" + password
	return &redis.Pool{
		Dial: func() (redisConn redis.Conn, err error) {
			redisConn, err = redis.Dial("tcp", address)
			if err != nil {
				if redisConn != nil {
					_ = redisConn.Close()
				}
				return
			}
			if _, err = redisConn.Do("config", "get", "requirepass"); err != nil {
				if strings.Compare(err.Error(), "NOAUTH Authentication required.") != 0 {
					_ = redisConn.Close()
					return
				}
				if _, err = redisConn.Do("auth", redisAuth); err != nil {
					_ = redisConn.Close()
					return
				}
			}

			return
		},
		MaxIdle:     0,
		MaxActive:   0,
		IdleTimeout: time.Minute,
	}
}
