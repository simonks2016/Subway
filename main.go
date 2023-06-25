package Subway

import (
	"Subway/Basic"
	"github.com/gomodule/redigo/redis"
)

type Subway struct {
	basicOperation Basic.Wheel
}

func NewSubway(redisPool *redis.Pool) *Subway {

	var v = new(Basic.OperationLib)
	v.Fuel = redisPool

	return &Subway{basicOperation: v}
}
