package air_redisclient

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
)

//redis client, 会全局维护已初始化的redis client

//创建client
//config: redis配置
func NewRedisClient(addr string, password string) (*RedisClient, error) {
	options := &redis.Options{
		Addr:     addr,
		Password: password,
	}
	client := redis.NewClient(options)
	cli := &RedisClient{
		Conn:     client,
		Addr:     addr,
		Password: password,
	}

	_, err := cli.Conn.Ping().Result()
	if err != nil {
		return nil, err
	}

	return cli, nil
}

//redis client
type RedisClient struct {
	Conn     *redis.Client //redis
	Addr     string        //地址，ip:port
	Password string        //密码
}

//get conn
func (cli *RedisClient) GetConn() *redis.Client {
	return cli.Conn
}

//close
func (cli *RedisClient) Close() {
	if cli.Conn != nil {
		_ = cli.Conn.Close()
	}
}

//Set
//key: key
//value: 设置值
//expiredS: 过期时间，为0则不过期
func (cli *RedisClient) Set(key string, value string, expiredS int) error {
	err := cli.GetConn().Set(key, value, time.Duration(expiredS)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

//Get
//key: key
func (cli *RedisClient) Get(key string) (string, error) {
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
//key: key
func (cli *RedisClient) MGet(key ...string) ([]interface{}, error) {
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
//key: key
func (cli *RedisClient) Del(key string) (int64, error) {
	ret, err := cli.GetConn().Del(key).Result()
	if err != nil {
		return -1, err
	}

	return ret, nil
}

//Incr
//key: key
func (cli *RedisClient) Incr(key string) (int64, error) {
	ret, err := cli.GetConn().Incr(key).Result()
	if err != nil {
		return -1, err
	}

	return ret, nil
}

//Decr
//key: key
func (cli *RedisClient) Decr(key string) (int64, error) {
	ret, err := cli.GetConn().Decr(key).Result()
	if err != nil {
		return -1, err
	}

	return ret, nil
}

//SetNX-支持过期时间的NX，即不存在才返回成功，加上过期可实现分布式锁
//key: key
//value: 设置值
//expiredS: 过期时间，为0则不过期
func (cli *RedisClient) SetNX(key string, value string, expiredS int) error {
	err := cli.GetConn().SetNX(key, value, time.Duration(expiredS)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

//HSet
//key: key
//field: field
//value: 值，一般场景用数值类型
func (cli *RedisClient) HSet(key string, field string, value interface{}) (bool, error) {
	ret, err := cli.GetConn().HSet(key, field, value).Result()
	if err != nil {
		return false, err
	}

	return ret, nil
}

//HGet
//key: key
//field: field
func (cli *RedisClient) HGet(key string, field string) (string, error) {
	value, err := cli.GetConn().HGet(key, field).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

//HDel
//key: key
//field: field
func (cli *RedisClient) HDel(key string, field string) (int64, error) {
	value, err := cli.GetConn().HDel(key, field).Result()
	if err != nil {
		return -1, err
	}

	return value, nil
}

//HMGet
//key: key
//field: field
func (cli *RedisClient) HMGet(key string, field ...string) ([]interface{}, error) {
	values, err := cli.GetConn().HMGet(key, field...).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("empty")
		}
		return nil, err
	}

	return values, nil
}
