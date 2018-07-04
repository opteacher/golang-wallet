package tests

import (
	"testing"
	"github.com/go-redis/redis"
	"databases"
	"log"
	"time"
)

func TestRedis(t *testing.T) {
	var cli redis.Cmdable
	var err error
	if cli, err = databases.ConnectRedis(); err != nil {
		log.Fatal(err)
	}
	err = cli.Set("abcd", "helloWorld", 20 * time.Second).Err()
	if err != nil {
		log.Fatal(err)
	}
}
