package tests

import (
	"testing"
	"github.com/go-redis/redis"
	"utils"
	"log"
	"time"
)

func TestRedis(t *testing.T) {
	//var cli *redis.Cmdable
	//var err error
	//if cli, err = databases.ConnectRedis(); err != nil {
	//	log.Fatal(err)
	//}
	//
	//if err = (*cli).Set("abcd", "helloWorld", 20 * time.Second).Err(); err != nil {
	//	log.Fatal(err)
	//}

	redisCfg := utils.GetConfig().GetSubsSettings().Redis
	redisCli := redis.NewClient(&redis.Options {
		Addr: redisCfg.Clusters[0].Url,
		Password: redisCfg.Password,
		DB: 0,
	})
	err := redisCli.Ping().Err()
	if err != nil {
		log.Fatal(err)
	}
	if err = redisCli.Set("abcd", "helloWorld", 20 * time.Second).Err(); err != nil {
		log.Fatal(err)
	}
}
