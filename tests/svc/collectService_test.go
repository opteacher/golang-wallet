package svc

import (
	"services"
	"time"
	"testing"
)

func TestCollectService(t *testing.T) {
	services.GetCollectService().Start()
	time.Sleep(20 * time.Second)
}
