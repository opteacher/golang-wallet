package svc

import (
	"services"
	"log"
	"time"
	"testing"
)

func TestDepositService(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	//Test service
	depositService := services.GetDepositService()
	notifyService := services.GetStableService()
	depositService.Init()
	notifyService.Init()

	depositService.Start()
	notifyService.Start()

	time.Sleep(20 * time.Second)

	depositService.Stop()
	notifyService.Stop()
	for !depositService.IsDestroy() && !notifyService.IsDestroy() {}
}
