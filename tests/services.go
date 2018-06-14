package main

import (
	"services"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)

	//Test service
	service := services.GetDepositService()
	service.Init()
	service.Start()

	time.Sleep(20 * time.Second)

	service.Stop()
}
