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

func (v OperationLib) SPop(string2 string) (error, string) {
	return nil, ""
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
			return errors.New("An error occurred! error message:" + err.Error()), ""
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
			return fmt.Errorf("对不起客官,查阅Redis失败!错误信息:%s", err.Error()), nil
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
