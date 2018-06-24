package tests

import (
	"log"
	"fmt"
	"rpcs"
	"testing"
)

func TestRPC(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	fmt.Println(rpcs.GetRPC("ETH").GetCurrentHeight())
}