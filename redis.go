package air_redisclient

import (
	"errors"
	"github.com/airingone/config"
	"github.com/airingone/log"
	"github.com/go-redis/redis"
	"sync"
)

//redis client, 会全局维护已初始化的redis client

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
		cli, err := NewRedisClient(config)
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

//创建client
//config: redis配置
func NewRedisClient(config config.ConfigRedis) (*RedisClient, error) {
	options := &redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
	}
	client := redis.NewClient(options)
	cli := &RedisClient{
		Conn:   client,
		config: config,
	}

	_, err := cli.Conn.Ping().Result()
	if err != nil {
		return nil, err
	}

	return cli, nil
}

//redis client
type RedisClient struct {
	Conn   *redis.Client      //redis
	config config.ConfigRedis //redis配置
}

//get conn
func (cli *RedisClient) GetConn() *redis.Client {
	return cli.Conn
}

//close
func (cli *RedisClient) Close() {
	_ = cli.Conn.Close()
}
