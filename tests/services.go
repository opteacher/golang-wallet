package main

import (
	"services"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)

	//Test service
	depositService := services.GetDepositService()
	notifyService := services.GetNotifyService()
	depositService.Init()
	notifyService.Init()

	depositService.Start()
	notifyService.Start()

	time.Sleep(20 * time.Second)

	depositService.Stop()
	notifyService.Stop()
	for !depositService.IsDestroy() && !notifyService.IsDestroy() {}
}
