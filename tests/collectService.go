package main

import (
	"services"
	"time"
)

func main() {
	services.GetCollectService().Start()
	time.Sleep(20 * time.Second)
}
