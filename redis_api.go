package air_redisclient

import (
	"errors"
	"github.com/airingone/config"
	"github.com/airingone/log"
	"sync"
)

//redis api是为了方便使用而进行封装，这里需要用到的操作可以随着使用继续增加

var AllRedisClients map[string]*RedisClient //全局的redis client
var AllRedisClientsRmu sync.RWMutex

//初始化全局redis对象，进程初始化的时候初始化一次
//configName: 配置文件redis配置名
func InitRedisClient(configName ...string) {
	if AllRedisClients == nil {
		AllRedisClients = make(map[string]*RedisClient)
	}

	for _, name := range configName {
		config := config.GetRedisConfig(name)
		cli, err := NewRedisClient(config.Addr, config.Password)
		if err != nil {
			log.Error("[REDIS]: InitRedisClient err, config name: %s, err: %+v", name, err)
			continue
		}

		AllRedisClientsRmu.Lock()
		if oldCli, ok := AllRedisClients[name]; ok { //	如果已存在则先关闭
			_ = oldCli.GetConn().Close()
		}
		AllRedisClients[name] = cli
		AllRedisClientsRmu.Unlock()
		log.Info("[REDIS]: InitRedisClient succ, config name: %s", name)
	}
}

//close all client
func CloseRedisClient() {
	if AllRedisClients == nil {
		return
	}
	AllRedisClientsRmu.RLock()
	defer AllRedisClientsRmu.RUnlock()
	for _, cli := range AllRedisClients {
		cli.Close()
	}
}

//get client对象
//configName: 配置文件redis配置名
func GetRedisClient(configName string) (*RedisClient, error) {
	AllRedisClientsRmu.RLock()
	defer AllRedisClientsRmu.RUnlock()
	if _, ok := AllRedisClients[configName]; !ok {
		return nil, errors.New("redis client not exist")
	}

	return AllRedisClients[configName], nil
}

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

	return cli.Set(key, value, expiredS)
}

//Get
//configName: redis配置名
//key: key
func RedisGet(configName string, key string) (string, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return "", err
	}

	return cli.Get(key)
}

//MGet
//configName: redis配置名
//key: key
func RedisMGet(configName string, key ...string) ([]interface{}, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return nil, err
	}

	return cli.MGet(key...)
}

//Del
//configName: redis配置名
//key: key
func RedisDel(configName string, key string) (int64, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return -1, err
	}

	return cli.Del(key)
}

//Incr
//configName: redis配置名
//key: key
func RedisIncr(configName string, key string) (int64, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return -1, err
	}

	return cli.Incr(key)
}

//Decr
//configName: redis配置名
//key: key
func RedisDecr(configName string, key string) (int64, error) {
	cli, err := GetRedisClient(configName)
	if err != nil {
		return -1, err
	}

	return cli.Decr(key)
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

	return cli.SetNX(key, value, expiredS)
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

	return cli.HSet(key, field, value)
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

	return cli.HGet(key, field)
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

	return cli.HDel(key, field)
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

	return cli.HMGet(key, field...)
}
