package tests

import (
	"utils"
	"log"
	"testing"
)

func TestConfig(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	//Test config
	config := utils.GetConfig()
	log.Println(config.GetBaseSettings())
	log.Println(config.GetSubsSettings())
}