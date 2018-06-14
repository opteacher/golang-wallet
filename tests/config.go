package main

import (
	"utils"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)

	//Test config
	config := utils.GetConfig()
	log.Println(config.GetBaseSettings())
	log.Println(config.GetSubsSettings())
}