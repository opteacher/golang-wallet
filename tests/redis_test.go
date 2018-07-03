package tests

import (
	"testing"
	"databases"
	"github.com/go-redis/redis"
	"log"
	"fmt"
)

func TestRedis(t *testing.T) {
	var cli *redis.Cmdable
	var err error
	if cli, err = databases.ConnectRedis(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cli)

	//if err = (*cli).Set("abcd", "helloWorld", 20 * time.Second).Err(); err != nil {
	//	log.Fatal(err)
	//}
}
