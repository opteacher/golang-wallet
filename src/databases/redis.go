package databases

import (
	"github.com/go-redis/redis"
	"sync"
	"utils"
)

var __redisClis redis.Cmdable
var __once sync.Once

func ConnectRedis() (redis.Cmdable, error) {
	var err error
	if __redisClis == nil {
		__once = sync.Once {}
		__once.Do(func() {
			err = createClients()
		})
	}
	return __redisClis, err
}

func createClients() error {
	redisCfg := utils.GetConfig().GetSubsSettings().Redis
	if len(redisCfg.Clusters) > 1 {
		var redisAddr []string
		for _, rds := range redisCfg.Clusters {
			redisAddr = append(redisAddr, rds.Url)
		}
		redisClis := redis.NewClusterClient(&redis.ClusterOptions {
			Addrs: redisAddr,
			Password: redisCfg.Password,
		})
		err := redisClis.Ping().Err()
		__redisClis = redisClis
		return err
	} else {
		redisCli := redis.NewClient(&redis.Options {
			Addr: redisCfg.Clusters[0].Url,
			Password: redisCfg.Password,
			DB: 0,
		})
		err := redisCli.Ping().Err()
		__redisClis = redisCli
		return err
	}
}