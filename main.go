package main

import (
	"services"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	//Test service
	services.GetDepositService().Init()
}