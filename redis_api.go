package air_redisclient

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
)

//redis api是为了方便使用而进行封装，这里需要用到的操作可以随着使用继续增加

//Set
//configName: redis配置名
//key: key
//value: 设置值
//expiredS: 过期时间，为0则不过期
func RedisSet(configName string, key string, value string, expiredS int) error {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return err
	}

	err = cli.GetConn().Set(key, value, time.Duration(expiredS)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

//Get
//configName: redis配置名
//key: key
func RedisGet(configName string, key string) (string, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return "", err
	}

	value, err := cli.GetConn().Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("empty")
		}
		return "", err
	}

	return value, nil
}

//MGet
//configName: redis配置名
//key: key
func RedisMGet(configName string, key ...string) ([]interface{}, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return nil, err
	}

	//cli.GetConn().Do("MGET", key...)
	values, err := cli.GetConn().MGet(key...).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("empty")
		}
		return nil, err
	}

	return values, nil
}

//Del
//configName: redis配置名
//key: key
func RedisDel(configName string, key string) (int64, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return -1, err
	}

	ret, err := cli.GetConn().Del(key).Result()
	if err != nil {
		return -1, err
	}

	return ret, nil
}

//Incr
//configName: redis配置名
//key: key
func RedisIncr(configName string, key string) (int64, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return -1, err
	}

	ret, err := cli.GetConn().Incr(key).Result()
	if err != nil {
		return -1, err
	}

	return ret, nil
}

//Decr
//configName: redis配置名
//key: key
func RedisDecr(configName string, key string) (int64, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return -1, err
	}

	ret, err := cli.GetConn().Decr(key).Result()
	if err != nil {
		return -1, err
	}

	return ret, nil
}

//SetNX-支持过期时间的NX，即不存在才返回成功，加上过期可实现分布式锁
//configName: redis配置名
//key: key
//value: 设置值
//expiredS: 过期时间，为0则不过期
func RedisSetNX(configName string, key string, value string, expiredS int) error {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return err
	}

	err = cli.GetConn().SetNX(key, value, time.Duration(expiredS)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

//HSet
//configName: redis配置名
//key: key
//field: field
//value: 值，一般场景用数值类型
func RedisHSet(configName string, key string, field string, value interface{}) (bool, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return false, err
	}

	ret, err := cli.GetConn().HSet(key, field, value).Result()
	if err != nil {
		return false, err
	}

	return ret, nil
}

//HGet
//configName: redis配置名
//key: key
//field: field
func RedisHGet(configName string, key string, field string) (string, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return "", err
	}

	value, err := cli.GetConn().HGet(key, field).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

//HDel
//configName: redis配置名
//key: key
//field: field
func RedisHDel(configName string, key string, field string) (int64, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return -1, err
	}

	value, err := cli.GetConn().HDel(key, field).Result()
	if err != nil {
		return -1, err
	}

	return value, nil
}

//HMGet
//configName: redis配置名
//key: key
//field: field
func RedisHMGet(configName string, key string, field ...string) ([]interface{}, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return nil, err
	}

	values, err := cli.GetConn().HMGet(key, field...).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("empty")
		}
		return nil, err
	}

	return values, nil
}
