package main

import (
	"log"
	"fmt"
	"rpcs"
)

func main() {
	log.SetFlags(log.Lshortfile)

	fmt.Println(rpcs.GetRPC("ETH").GetCurrentHeight())
}