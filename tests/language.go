package main

import (
	"log"
	"fmt"
	"reflect"
)

func main() {
	log.SetFlags(log.Lshortfile)
	fmt.Println("abcd")
	var t float64 = 58500000000000000000
	log.Println(reflect.TypeOf(t).Name())
}
