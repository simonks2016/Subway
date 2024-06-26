package Core

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	//"github.com/simonks2016/Subway/DataAdapter"
	errors2 "github.com/simonks2016/Subway/errors"
	"strings"
)

type OperationLib struct {
	Wheel
	Fuel *redis.Pool
}

func (v OperationLib) SPop(string2 string, count int) (result []string, err error) {
	Redis := v.Fuel.Get()
	if Redis.Err() != nil {
		return nil, Redis.Err()
	}
	defer func() {
		_ = Redis.Close()
	}()

	result, err = redis.Strings(Redis.Do("SPop", string2, count))
	return
}

func (this OperationLib) SAdd(SetName string, key ...interface{}) (err error) {

	Redis := this.Fuel.Get()
	if Redis.Err() != nil {
		return Redis.Err()
	}
	defer func() {
		_ = Redis.Close()
	}()

	var args []interface{}
	var num int

	args = append(args, SetName)
	args = append(args, key...)

	if num, err = redis.Int(Redis.Do("SADD", args...)); err != nil {
		return err
	}
	if num == 0 {
		return errors2.ErrAddFailed
	}
	return nil
}
func (this OperationLib) SMembers(SetName string) (result []string, err error) {

	Redis := this.Fuel.Get()
	if Redis.Err() != nil {
		return nil, Redis.Err()
	}
	defer func() {
		_ = Redis.Close()
	}()
	result, err = redis.Strings(Redis.Do("SMEMBERS", SetName))
	return
}

func (this OperationLib) SCard(SetName string) (num int, err error) {

	Redis := this.Fuel.Get()
	if Redis.Err() != nil {
		return 0, Redis.Err()
	}
	defer func() {
		_ = Redis.Close()
	}()

	num, err = redis.Int(Redis.Do("SCARD", SetName))
	return
}

func (this OperationLib) SIsMember(SetName string, key string) (result bool, err error) {

	Redis := this.Fuel.Get()
	if Redis.Err() != nil {
		return false, Redis.Err()
	}
	defer func() {
		_ = Redis.Close()
	}()

	result, err = redis.Bool(Redis.Do("SISMEMBER", SetName, key))
	return
}

func (this OperationLib) SRemove(SetName string, key ...interface{}) (err error) {

	Redis := this.Fuel.Get()
	if Redis.Err() != nil {
		return Redis.Err()
	}
	defer func() {
		_ = Redis.Close()
	}()

	var args []interface{}

	args = append(args, SetName)
	args = append(args, key...)

	_, err = Redis.Do("SREM", args...)
	return
}

func (this OperationLib) SRandMember(SetName string, Count int) (result []string, err error) {

	Redis := this.Fuel.Get()
	if Redis.Err() != nil {
		return nil, Redis.Err()
	}
	defer func() {
		_ = Redis.Close()
	}()
	result, err = redis.Strings(Redis.Do("SRANDMEMBER", SetName, Count))
	return
}

func (this OperationLib) SIncr(sets ...interface{}) (result []string, err error) {

	Redis := this.Fuel.Get()
	if Redis.Err() != nil {
		return nil, Redis.Err()
	}
	defer func() {
		_ = Redis.Close()
	}()

	var args []interface{}

	args = append(args, sets...)

	result, err = redis.Strings(Redis.Do("SINTER", args...))
	return
}

func (this OperationLib) SUnion(sets ...interface{}) (result []string, err error) {

	Redis := this.Fuel.Get()
	if Redis.Err() != nil {
		return nil, Redis.Err()
	}
	defer func() {
		_ = Redis.Close()
	}()

	result, err = redis.Strings(Redis.Do("SUNION", sets...))
	return

}

func (this OperationLib) Exist(key string) (bool, error) {

	if this.Fuel != nil {
		Redis := this.Fuel.Get()
		if Redis.Err() != nil {
			return false, Redis.Err()
		}
		defer func() {
			if err := Redis.Close(); err != nil {
				return
			}
		}()
		ret, err := redis.Bool(Redis.Do("EXISTS", key))
		if err != nil {
			return ret, errors.New("query by redis failed")
		}
		return ret, err
	}
	return false, errors2.ErrUnable2ConnectRedis

}
func (this OperationLib) GetByte(key string) (error, []byte) {

	if this.Fuel != nil {
		Redis := this.Fuel.Get()
		if Redis.Err() != nil {
			return Redis.Err(), nil
		}
		defer func() {
			if err := Redis.Close(); err != nil {
				return
			}
		}()
		exists, err := redis.Bool(Redis.Do("exists", key))
		if err != nil {
			return err, nil
		}
		if !exists {
			return errors2.ErrNil, nil
		}
		ret, err := redis.Bytes(Redis.Do("get", key))
		if err != nil {
			if errors.Is(err, redis.ErrNil) {
				return errors2.ErrNil, nil
			}
			return err, nil
		}
		return nil, ret
	}
	return errors2.ErrUnable2ConnectRedis, nil

}

func (this OperationLib) GetString(key string) (error, string) {

	if this.Fuel != nil {
		Redis := this.Fuel.Get()
		if Redis.Err() != nil {
			return Redis.Err(), ""
		}
		defer func() {
			if err := Redis.Close(); err != nil {
				return
			}
		}()
		ret, err := redis.String(Redis.Do("get", key))
		if err != nil {
			if errors.Is(err, redis.ErrNil) {
				return errors2.ErrNil, ""
			}
			return err, ""
		}
		if ret == "nil" || ret == "null" {
			return errors2.ErrNil, ""
		}
		return nil, ret
	}
	return errors2.ErrUnable2ConnectRedis, ""
}

func (this OperationLib) SetStringEx(key, value string, expireTime int64) error {
	if this.Fuel != nil {
		Redis := this.Fuel.Get()
		if Redis.Err() != nil {
			return Redis.Err()
		}
		defer func() {
			if err := Redis.Close(); err != nil {
				return
			}
		}()
		if _, err := Redis.Do("SETEX", key, expireTime, value); err != nil {
			return err
		}
		return nil
	}
	return errors2.ErrUnable2ConnectRedis
}

func (this OperationLib) SetStringNx(key, value string) error {
	if this.Fuel != nil {
		Redis := this.Fuel.Get()
		if Redis.Err() != nil {
			return Redis.Err()
		}
		defer func() {
			if err := Redis.Close(); err != nil {
				return
			}
		}()
		if _, err := Redis.Do("SETNX", key, value); err != nil {
			return err
		}
		return nil
	}
	return errors2.ErrUnable2ConnectRedis
}

func (this OperationLib) SetString(key string, value string) error {
	if this.Fuel != nil {
		Redis := this.Fuel.Get()
		if Redis.Err() != nil {
			return Redis.Err()
		}
		defer func() {
			if err := Redis.Close(); err != nil {
				return
			}
		}()
		if _, err := Redis.Do("set", key, value); err != nil {
			return err
		}
		return nil
	}
	return errors2.ErrUnable2ConnectRedis
}

func (this OperationLib) Set(key string, value []byte) error {
	if this.Fuel != nil {
		Redis := this.Fuel.Get()
		if Redis.Err() != nil {
			return Redis.Err()
		}
		defer func() {
			if err := Redis.Close(); err != nil {
				return
			}
		}()
		if _, err := Redis.Do("set", key, value); err != nil {
			return err
		}
		return nil
	}
	return errors2.ErrUnable2ConnectRedis
}

func (this OperationLib) BatchGetStrings(key ...interface{}) (err error, ret []string) {

	if this.Fuel != nil {
		Redis := this.Fuel.Get()
		if Redis.Err() != nil {
			return Redis.Err(), nil
		}
		defer func() {
			if err = Redis.Close(); err != nil {
				return
			}
		}()
		ret, err = redis.Strings(Redis.Do("mget", key...))
		if err != nil {
			if errors.Is(err, redis.ErrNil) {
				return errors2.ErrNil, nil
			}
			return err, nil
		}

		if len(ret) <= 0 {
			return errors2.ErrNil, nil
		}
		return nil, ret
	}
	return errors2.ErrUnable2ConnectRedis, nil
}

func (c OperationLib) Delete(Key string) error {

	Redis := c.Fuel.Get()
	if Redis.Err() != nil {
		return Redis.Err()
	}
	defer func() {
		_ = Redis.Close()
	}()

	if _, err := Redis.Do("del", Key); err != nil {
		return err
	}
	return nil
}

func (c OperationLib) BatchDelete(key ...interface{}) error {

	if c.Fuel != nil {
		Redis := c.Fuel.Get()
		if Redis.Err() != nil {
			return Redis.Err()
		}
		defer func() {
			_ = Redis.Close()
		}()
		var args []interface{}
		//append to key ele
		args = append(args, key...)
		if _, err := Redis.Do("del", args...); err != nil {
			return err
		}
		return nil
	}
	return errors2.ErrUnable2ConnectRedis
}

func (c OperationLib) Expire(Key string, expireTime int64) error {

	Redis := c.Fuel.Get()
	if Redis.Err() != nil {
		return Redis.Err()
	}
	defer func() {
		_ = Redis.Close()
	}()

	if _, err := Redis.Do("expire", Key, expireTime); err != nil {
		return err
	}
	return nil
}

func (o OperationLib) ZRemoveAll(sets ...interface{}) (err error) {

	var (
		Redis redis.Conn
	)
	if o.Fuel != nil {
		Redis = o.Fuel.Get()
		if err = Redis.Err(); err != nil {
			return
		} else {
			defer func() {
				_ = Redis.Close()
			}()
			//查阅redis
			_, err = Redis.Do("zremrangebyrank", sets, 0, -1)
			//返回
			return
		}
	}
	return errors2.ErrUnable2ConnectRedis
}

func (o OperationLib) ZAdd(SetName string, Score float64, key ...interface{}) (err error) {

	if o.Fuel != nil {
		var Redis = o.Fuel.Get()
		if err = Redis.Err(); err != nil {
			return err
		}

		defer func() {
			_ = Redis.Close()
		}()

		var args []interface{}
		args = append(args, SetName, Score)
		args = append(args, key...)

		var num int
		if num, err = redis.Int(Redis.Do("ZADD", args...)); err != nil {
			return err
		}
		if num == 0 {
			return errors.New("failed to add cache")
		}
		return nil
	}
	return errors2.ErrUnable2ConnectRedis
}

func (o OperationLib) ZRange(SetName string, Start, End int64, Desc bool) (err error, result []string) {

	var (
		Redis redis.Conn
	)

	if o.Fuel != nil {
		Redis = o.Fuel.Get()
		if err = Redis.Err(); err != nil {
			return
		} else {
			defer func() {
				_ = Redis.Close()
			}()
			//查阅redis
			if result, err = redis.Strings(Redis.Do(func() string {
				if Desc {
					return "zrevrange"
				}
				return "zrange"
			}(), SetName, Start, End)); err != nil {
				return err, nil
			} else {
				return nil, result
			}
		}
	}
	return errors2.ErrUnable2ConnectRedis, nil
}

func (o OperationLib) ZCard(SetName string) (err error, total int64) {

	var (
		Redis redis.Conn
	)

	if o.Fuel != nil {
		Redis = o.Fuel.Get()
		if err = Redis.Err(); err != nil {
			return
		} else {
			defer func() {
				_ = Redis.Close()
			}()
			//查阅redis
			total, err = redis.Int64(Redis.Do("zcard", SetName))
			//返回
			return
		}
	}
	return errors2.ErrUnable2ConnectRedis, 0
}

func (o OperationLib) ZRangeBySore(SetName string, min, max, offset, limit int64) (err error, result []string) {

	var (
		Redis redis.Conn
	)

	if o.Fuel != nil {
		Redis = o.Fuel.Get()
		if err = Redis.Err(); err != nil {
			return
		} else {
			defer func() {
				_ = Redis.Close()
			}()
			//查阅redis
			result, err = redis.Strings(Redis.Do("ZRANGEBYSCORE", SetName,
				min,
				max,
				"LIMIT",
				offset,
				limit,
			))
			//返回
			return
		}
	}
	return errors2.ErrUnable2ConnectRedis, nil
}

func (o OperationLib) ZRemove(SetName string, key ...interface{}) (err error, num int) {

	var (
		Redis redis.Conn
	)

	if o.Fuel != nil {
		Redis = o.Fuel.Get()
		if err = Redis.Err(); err != nil {
			return
		} else {
			defer func() {
				_ = Redis.Close()
			}()

			var args []interface{}

			if len(SetName) > 0 {
				args = append(args, SetName)
			}
			args = append(args, key...)

			//查阅redis
			num, err = redis.Int(Redis.Do("ZREM", args...))
			//返回
			return
		}
	}
	return errors2.ErrUnable2ConnectRedis, 0
}

func (o OperationLib) ZIsMember(setName string, key string) (error, bool) {

	var (
		r redis.Conn
	)

	if o.Fuel != nil {
		r = o.Fuel.Get()
		if err := r.Err(); err != nil {
			return err, false
		} else {
			defer func() {
				_ = r.Close()
			}()

			var d []string
			//查阅redis
			d, err = redis.Strings(r.Do("ZRANGE", setName, 0, -1))
			if err != nil {
				return err, false
			}

			for _, s := range d {
				if strings.Compare(s, key) == 0 {
					return nil, true
				}
			}
			//返回
			return nil, false
		}
	}
	return errors2.ErrUnable2ConnectRedis, false
}

func (o OperationLib) Persist(key ...interface{}) error {

	if o.Fuel != nil {
		redisConn := o.Fuel.Get()
		if err := redisConn.Err(); err != nil {
			return err
		} else {
			defer func() {
				_ = redisConn.Close()
			}()
			//do command
			if _, err = redisConn.Do("persist", key...); err != nil {
				return err
			}
			return nil
		}
	}
	return errors2.ErrUnable2ConnectRedis
}

func (o OperationLib) SetHashMap(key interface{}, field interface{}, value interface{}) error {

	if o.Fuel != nil {
		redisConn := o.Fuel.Get()
		if err := redisConn.Err(); err != nil {
			return err
		} else {
			defer func() {
				_ = redisConn.Close()
			}()
			//do command
			if _, err = redisConn.Do("HSET", key, field, value); err != nil {
				return err
			}
			return nil
		}
	}
	return errors2.ErrUnable2ConnectRedis

}
func (o OperationLib) GetHashMap(key, field interface{}) (any, error) {

	if o.Fuel != nil {
		redisConn := o.Fuel.Get()
		if err := redisConn.Err(); err != nil {
			return nil, err
		} else {
			defer func() {
				_ = redisConn.Close()
			}()
			//do command
			if result, err := redisConn.Do("HGET", key, field); err != nil {
				if errors.Is(err, redis.ErrNil) {
					return nil, errors2.ErrNil
				}
				return nil, err
			} else {
				return result, nil
			}
		}
	}
	return nil, errors2.ErrUnable2ConnectRedis
}

func (o OperationLib) GetALLHashMap(key interface{}) ([]any, error) {

	if o.Fuel != nil {
		redisConn := o.Fuel.Get()
		if err := redisConn.Err(); err != nil {
			return nil, err
		} else {
			defer func() {
				_ = redisConn.Close()
			}()
			//do command
			if result, err := redis.Values(redisConn.Do("HGETALL", key)); err != nil {
				return nil, err
			} else {
				return result, nil
			}
		}
	}
	return nil, errors2.ErrUnable2ConnectRedis

}

func (o OperationLib) MSetHashMap(key interface{}, data map[any]any) error {

	var args []interface{}

	args = append(args, key)

	for field, value := range data {
		args = append(args, field, value)
	}

	if o.Fuel != nil {
		redisConn := o.Fuel.Get()
		if err := redisConn.Err(); err != nil {
			return err
		} else {
			defer func() {
				_ = redisConn.Close()
			}()
			//do command
			if _, err = redisConn.Do("HMSET", args...); err != nil {
				return err
			}
			return nil
		}
	}
	return errors2.ErrUnable2ConnectRedis
}

func (o OperationLib) MGetHashMap(key string, fields ...interface{}) ([]any, error) {

	if o.Fuel != nil {
		redisConn := o.Fuel.Get()
		if err := redisConn.Err(); err != nil {
			return nil, err
		} else {
			defer func() {
				_ = redisConn.Close()
			}()
			//args
			var args []interface{}
			//set the args
			args = append(args, key)
			args = append(args, fields...)

			//do command
			if result, err := redis.Values(redisConn.Do("HMGET", args...)); err != nil {
				return nil, err
			} else {
				return result, nil
			}
		}
	}
	return nil, errors2.ErrUnable2ConnectRedis

}

func (o OperationLib) DelHashMap(key, field interface{}) error {

	if o.Fuel != nil {
		redisConn := o.Fuel.Get()
		if err := redisConn.Err(); err != nil {
			return err
		} else {
			defer func() {
				_ = redisConn.Close()
			}()
			//do command
			if _, err = redisConn.Do("HDEL", key, field); err != nil {
				return err
			}
			return nil
		}
	}
	return errors2.ErrUnable2ConnectRedis
}
func (o OperationLib) GetFieldsHashMap(key interface{}) ([]any, error) {
	if o.Fuel != nil {
		redisConn := o.Fuel.Get()
		if err := redisConn.Err(); err != nil {
			return nil, err
		} else {
			defer func() {
				_ = redisConn.Close()
			}()
			//do command
			if result, err := redis.Values(redisConn.Do("HKEYS", key)); err != nil {
				return nil, err
			} else {
				return result, nil
			}
		}
	}
	return nil, errors2.ErrUnable2ConnectRedis
}

func (o OperationLib) ExistsHashMap(key interface{}, field interface{}) (bool, error) {

	if o.Fuel != nil {
		redisConn := o.Fuel.Get()
		if err := redisConn.Err(); err != nil {
			return false, err
		} else {
			defer func() {
				_ = redisConn.Close()
			}()
			//do command
			if result, err := redis.Bool(redisConn.Do("HEXISTS", key, field)); err != nil {
				return false, err
			} else {
				return result, nil
			}
		}
	}
	return false, errors2.ErrUnable2ConnectRedis
}

func (o OperationLib) NewDocumentId(topic, id string) string {
	return fmt.Sprintf("subway-document-%s-%s", topic, id)
}

func NewOperationLib(r *redis.Pool) *OperationLib {

	return &OperationLib{
		Fuel: r,
	}
}
