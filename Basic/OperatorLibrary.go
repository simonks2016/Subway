package Basic

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
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
		return errors.New("add cache failed")
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
	return false, errors.New("unable to connect redis")

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
			return err, ""
		}
		return nil, ret
	}
	return errors.New("unable to connect to Redis"), ""
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
			return fmt.Errorf("对不起客官,查阅Redis失败!错误信息:%s", err.Error())
		}
		return nil
	}
	return fmt.Errorf("%s\n", "连接Redis数据库失败!")
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
			return err, nil
		}
		return nil, ret
	}
	return fmt.Errorf("%s\n", "连接Redis数据库失败!"), nil
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
	return errors.New("failed to connect to redis")
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
	return errors.New("failed to connect to redis")
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
	return errors.New("failed to connect to redis"), nil
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
	return errors.New("failed to connect to redis"), 0
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
	return errors.New("failed to connect to redis"), nil
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
	return errors.New("failed to connect to redis"), 0
}

func (o OperationLib) NewDocumentId(topic, id string) string {
	return fmt.Sprintf("subway-document-%s-%s", topic, id)
}
