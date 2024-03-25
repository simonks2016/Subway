package subway

import (
	"github.com/gomodule/redigo/redis"
	"strings"
	"time"
)

type _subwayConnection struct {
	pool *redis.Pool
}

var Subway *_subwayConnection = nil

func NewRedisConnWithSubway(address, userName, password string) *_subwayConnection {

	redisAuth := userName + ":" + password

	pool := &redis.Pool{
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

	Subway = &_subwayConnection{pool: pool}
	//return Subway
	return Subway
}
