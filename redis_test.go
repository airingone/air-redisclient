package air_redisclient

import (
	"github.com/airingone/config"
	"github.com/airingone/log"
	"testing"
	"time"
)

//redis client测试
func TestRedisClient(t *testing.T) {
	config.InitConfig()                     //配置文件初始化
	log.InitLog(config.GetLogConfig("log")) //日志初始化
	InitRedisClient("redis_test1")          //初始化创建redis client
	defer CloseRedisClient()

	err := RedisSet("redis_test1", "key_test01", "value_test01", 0)
	log.Error("RedisSet: err: %+v", err)
	err = RedisSet("redis_test1", "key_test02", "value_test02", 10)
	log.Error("RedisSet: err: %+v", err)

	value01, err := RedisGet("redis_test1", "key_test01")
	log.Error("RedisGet: err: %+v, value: %s", err, value01)

	value02, err := RedisMGet("redis_test1", "key_test01", "key_test02")
	log.Error("RedisMGet: err: %+v, value: %+v", err, value02)

	time.Sleep(10 * time.Second)

	value03, err := RedisGet("redis_test1", "key_test02")
	log.Error("RedisGet: err: %+v, value: %s", err, value03)

	ret, err := RedisDel("redis_test1", "key_test01")
	log.Error("RedisDel: err: %+v, ret: %d", err, ret)
	ret, err = RedisDel("redis_test1", "key_test02")
	log.Error("RedisDel: err: %+v, ret: %d", err, ret)

	value04, err := RedisMGet("redis_test1", "key_test01", "key_test02")
	log.Error("RedisMGet: err: %+v, value: %+v", err, value04)

	//or
	cli, err := GetRedisClient("redis_test1")
	if err != nil {
		log.Error("GetRedisClient: err: %+v", err)
		return
	}
	err = cli.Set("key_test05", "value05", 0)
	log.Error("Set: err: %+v, value: %+v", err)
	value05, err := cli.Get("key_test05")
	log.Error("Get: err: %+v, value: %+v", err, value05)
}

//brew install redis
//redis-server
//redis-cli -h 127.0.0.1 -p 6379
//config set requirepass 123456 设置密码
//127.0.0.1:6379> AUTH 123456
