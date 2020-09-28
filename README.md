# redis client组建
## 1.组件描述
redis client用于redis客户端

## 2.如何使用
```
import (
    "github.com/airingone/config"
    "github.com/airingone/log"
    redisclient "github.com/airingone/air-redisclient"
)

func main() {
    config.InitConfig()                        //进程启动时调用一次初始化配置文件，配置文件名为config.yml，目录路径为../conf/或./
    log.InitLog(config.GetLogConfig("log"))    //进程启动时调用一次初始化日志
    redisclient.InitRedisClient("redis_test1") //初始化创建redis client
    defer redisclient.CloseRedisClient()       //进程退出时关闭client

    #打印日志
    err := redisclient.RedisSet("redis_test1", "key_test01", "value_test01", 0)
    if err != nil {
        log.Error("RedisSet: err: %+v", err)
    } else {
        log.Error("RedisSet: succ")
    }
  
    value0, err := RedisGet("redis_test1", "key_test01")
    if err != nil {
        log.Error("RedisGet: err: %+v", err)
    } else {
        log.Error("RedisGet: value: %s", value01)
    }  
}
```
更多使用请参考[测试例子](https://github.com/airingone/air-redisclient/blob/master/redis_test.go)
